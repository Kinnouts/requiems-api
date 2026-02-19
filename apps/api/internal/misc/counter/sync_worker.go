package counter

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const syncInterval = 60 * time.Second

// StartSyncWorker runs a background goroutine that periodically flushes Redis
// counter values to PostgreSQL. It stops when ctx is cancelled.
func StartSyncWorker(ctx context.Context, rdb *redis.Client, repo Repository) {
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := syncToPostgres(ctx, rdb, repo); err != nil {
				log.Printf("counter sync error: %v", err)
			}
		}
	}
}

func syncToPostgres(ctx context.Context, rdb *redis.Client, repo Repository) error {
	var cursor uint64
	var keys []string

	for {
		var batch []string
		var err error

		batch, cursor, err = rdb.Scan(ctx, cursor, "counter:*", 100).Result()
		if err != nil {
			return err
		}

		keys = append(keys, batch...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) == 0 {
		return nil
	}

	// Fetch all values in a single pipeline round-trip.
	pipe := rdb.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))

	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	for i, cmd := range cmds {
		val, err := cmd.Int64()
		if err != nil {
			continue
		}

		namespace := keys[i][len("counter:"):]

		if err := repo.Upsert(ctx, namespace, val); err != nil {
			log.Printf("counter upsert error for %q: %v", namespace, err)
		}
	}

	return nil
}
