//nolint:revive // misc is an acceptable package name for miscellaneous utilities
package misc

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"requiems-api/services/misc/convert"
	"requiems-api/services/misc/counter"
)

func RegisterRoutes(ctx context.Context, r chi.Router, pool *pgxpool.Pool, rdb *redis.Client) {
	convertSvc := convert.NewService()
	convert.RegisterRoutes(r, convertSvc)

	counterRepo := counter.NewRepository(pool)
	counterSvc := counter.NewService(rdb, counterRepo)
	go counter.StartSyncWorker(ctx, rdb, counterRepo)
	counter.RegisterRoutes(r, counterSvc)
}
