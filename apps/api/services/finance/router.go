package finance

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"requiems-api/services/finance/bin"
)

// RegisterRoutes mounts all finance domain handlers on the given router.
// The router is expected to be mounted at /v1/finance by the caller.
func RegisterRoutes(r chi.Router, pool *pgxpool.Pool) {
	binSvc := bin.NewService(pool)
	bin.RegisterRoutes(r, binSvc)
}
