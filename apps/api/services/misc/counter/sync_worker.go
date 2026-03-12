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
	processingExists, err := rdb.Exists(ctx, processingSetKey).Result()

	if err != nil {
		return err
	}

	if processingExists == 0 {
		// No in-flight snapshot exists, so atomically move the current dirty set
		// into processing. If counter:dirty doesn't exist, Rename returns redis.Nil.
		err = rdb.Rename(ctx, dirtySetKey, processingSetKey).Err()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}

		if errors.Is(err, redis.Nil) {
			return nil
		}
	}

	// Step 2: Get all namespaces from the processing set.
	// Use SMEMBERS to fetch the snapshot. For very large sets (10K+),
	// consider SSCAN for incremental retrieval.

	namespaces, err := rdb.SMembers(ctx, processingSetKey).Result()
	
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if len(namespaces) == 0 {
		// No counters in the processing set; clean it up
		return rdb.Del(ctx, processingSetKey).Err()
	}

	// Step 3: Fetch all counter values in a single MGET operation.
	// Build the full counter keys and fetch them atomically.
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
	// This handles edge cases where a counter is incremented then deleted before sync.

	counters := make(map[string]int64, len(namespaces))
	
	for i, ns := range namespaces {
		if vals[i] == nil {
			// Counter no longer exists in Redis; skip it
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

	if len(counters) == 0 {
		// No valid counters to sync; clean up the processing set
		return rdb.Del(ctx, processingSetKey).Err()
	}

	// Step 5: Perform a single batch UPSERT to PostgreSQL.
	// This is the most efficient database operation possible:
	// - One TCP round-trip
	// - One transaction
	// - PostgreSQL planner optimizes the bulk insert
	// For 10K counters, this typically <100ms depending on network latency.
	if err := repo.UpsertBatch(ctx, counters); err != nil {
		return err
	}

	// Step 6: Delete the processing set after successful sync.
	// This marks the sync as complete and allows the pattern to repeat.
	// If this fails (network error), the processing set remains and will
	// be retried on the next worker cycle.
	return rdb.Del(ctx, processingSetKey).Err()
}