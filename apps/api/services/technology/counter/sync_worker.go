package counter

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// The Dirty-Set Swap Pattern for Race-Safe Counter Syncing

const syncInterval = 60 * time.Second
const dirtySetKey = "counter:dirty"
const processingSetKey = "counter:dirty:processing"

// StartSyncWorker runs a background goroutine that periodically syncs dirty
// redis counters to PostgreSQL using the atomic dirty-set swap pattern.
// It stops when ctx is cancelled and isolates panics to a single sync cycle.
func StartSyncWorker(ctx context.Context, rdb *redis.Client, repo Repository) {
	ticker := time.NewTicker(syncInterval)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := runSyncCycle(ctx, rdb, repo); err != nil {
				log.Printf("counter sync error: %v", err)
			}
		}
	}
}

func runSyncCycle(ctx context.Context, rdb *redis.Client, repo Repository) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("counter sync worker panicked during cycle")
			log.Printf("counter sync worker panicked: %v — worker will continue on next cycle", r)
		}
	}()

	return syncDirtyCounters(ctx, rdb, repo)
}

func syncDirtyCounters(ctx context.Context, rdb *redis.Client, repo Repository) error {
	// Step 1: Reuse any existing processing snapshot from a prior failed cycle.
	// This preserves at-least-once delivery by retrying the same dirty set
	// before moving new dirty namespaces into processing.
	acquired, err := acquireProcessingSet(ctx, rdb)
	if err != nil {
		return err
	}
	if !acquired {
		return nil
	}

	// Step 2: Get all namespaces from the processing set.
	// Use SMEMBERS to fetch the snapshot. For very large sets (10K+),
	// consider SSCAN for incremental retrieval.
	namespaces, err := rdb.SMembers(ctx, processingSetKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	if len(namespaces) == 0 {
		return rdb.Del(ctx, processingSetKey).Err()
	}

	// Step 3: Fetch all counter values in a single MGET operation.
	// MGET is O(N) for N counters and highly optimized in Redis.
	counterKeys := make([]string, len(namespaces))
	for i, ns := range namespaces {
		counterKeys[i] = redisKey(ns)
	}

	vals, err := rdb.MGet(ctx, counterKeys...).Result()
	if err != nil {
		return err
	}

	// Step 4: Build the map of namespaces to counter values for batch upsert.
	// Skip counters that have been deleted (vals[i] == nil).
	counters := parseCounterValues(namespaces, vals)
	if len(counters) == 0 {
		return rdb.Del(ctx, processingSetKey).Err()
	}

	// Step 5: Perform a single batch UPSERT to PostgreSQL.
	if err := repo.UpsertBatch(ctx, counters); err != nil {
		return err
	}

	// Step 6: Delete the processing set after successful sync.
	// If this fails, the processing set remains and will be retried next cycle.
	return rdb.Del(ctx, processingSetKey).Err()
}

// acquireProcessingSet atomically moves the dirty set into the processing set.
// Returns true if a processing set is ready to be consumed, false if there is
// nothing to sync (dirty set is empty). Reuses an existing processing set when
// a prior cycle failed mid-flight (at-least-once guarantee).
func acquireProcessingSet(ctx context.Context, rdb *redis.Client) (bool, error) {
	processingExists, err := rdb.Exists(ctx, processingSetKey).Result()
	if err != nil {
		return false, err
	}
	if processingExists > 0 {
		return true, nil
	}

	// No in-flight snapshot exists — atomically rename the dirty set.
	// Check existence first to avoid the ERR no such key error from RENAMENX.
	dirtyExists, err := rdb.Exists(ctx, dirtySetKey).Result()
	if err != nil {
		return false, err
	}
	if dirtyExists == 0 {
		return false, nil
	}

	_, err = rdb.RenameNX(ctx, dirtySetKey, processingSetKey).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

// parseCounterValues builds a namespace→value map from parallel MGET results.
// Entries where the Redis key no longer exists (nil) are silently skipped.
func parseCounterValues(namespaces []string, vals []any) map[string]int64 {
	counters := make(map[string]int64, len(namespaces))
	for i, ns := range namespaces {
		if vals[i] == nil {
			continue
		}

		counterVal, ok := vals[i].(string)
		if !ok {
			log.Printf("counter value type assertion failed for %q", ns)
			continue
		}

		val, err := strconv.ParseInt(counterVal, 10, 64)
		if err != nil {
			log.Printf("counter parse error for %q: %v", ns, err)
			continue
		}
		counters[ns] = val
	}
	return counters
}
