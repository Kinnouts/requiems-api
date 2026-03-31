package finance

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"requiems-api/services/finance/bin"
	"requiems-api/services/finance/commodities"
	"requiems-api/services/finance/exchange"
	"requiems-api/services/finance/iban"
	"requiems-api/services/finance/inflation"
	"requiems-api/services/finance/mortgage"
)

// RegisterRoutes mounts all finance domain handlers on the given router.
// The router is expected to be mounted at /v1/finance by the caller.
func RegisterRoutes(r chi.Router, pool *pgxpool.Pool, rdb *redis.Client) {
	binSvc := bin.NewService(pool)
	bin.RegisterRoutes(r, binSvc)

	inflationSvc := inflation.NewService(pool)
	inflation.RegisterRoutes(r, inflationSvc)

	mortgageSvc := mortgage.NewService()
	mortgage.RegisterRoutes(r, mortgageSvc)

	ibanSvc := iban.NewService(pool)
	iban.RegisterRoutes(r, ibanSvc)

	commoditiesSvc := commodities.NewService(pool)
	commodities.RegisterRoutes(r, commoditiesSvc)

	exchangeSvc := exchange.NewService(rdb)
	exchange.RegisterRoutes(r, exchangeSvc)
}
