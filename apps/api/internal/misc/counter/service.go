package counter

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// Service defines the counter business logic interface.
type Service interface {
	Increment(ctx context.Context, namespace string) (int64, error)
	Get(ctx context.Context, namespace string) (int64, error)
}

type service struct {
	rdb  *redis.Client
	repo Repository
}

func NewService(rdb *redis.Client, repo Repository) Service {
	return &service{rdb: rdb, repo: repo}
}

// Increment atomically increments the Redis counter and returns the new value.
func (s *service) Increment(ctx context.Context, namespace string) (int64, error) {
	if err := validateNamespace(namespace); err != nil {
		return 0, err
	}
	return s.rdb.Incr(ctx, redisKey(namespace)).Result()
}

// Get returns the counter value from Redis, falling back to PostgreSQL when the
// key is absent. If a PostgreSQL value is found it is hydrated back into Redis.
func (s *service) Get(ctx context.Context, namespace string) (int64, error) {
	if err := validateNamespace(namespace); err != nil {
		return 0, err
	}

	val, err := s.rdb.Get(ctx, redisKey(namespace)).Int64()
	if err == nil {
		return val, nil
	}

	if !errors.Is(err, redis.Nil) {
		return 0, err
	}

	// Fallback to PostgreSQL
	total, err := s.repo.Get(ctx, namespace)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	// Hydrate Redis
	if err := s.rdb.Set(ctx, redisKey(namespace), total, 0).Err(); err != nil {
		log.Printf("counter redis hydrate error for %q: %v", namespace, err)
	}

	return total, nil
}
