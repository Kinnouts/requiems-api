package misc

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"requiems-api/internal/misc/counter"
)

func RegisterRoutes(ctx context.Context, r chi.Router, pool *pgxpool.Pool, rdb *redis.Client) {
	repo := counter.NewRepository(pool)
	svc := counter.NewService(rdb, repo)

	go counter.StartSyncWorker(ctx, rdb, repo)

	counter.RegisterRoutes(r, svc)
}
