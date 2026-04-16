package health

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"requiems-api/services/health/exercises"
)

// RegisterRoutes wires all health sub-services onto the given router.
func RegisterRoutes(r chi.Router, pool *pgxpool.Pool) {
	svc := exercises.NewService(pool)
	exercises.RegisterRoutes(r, svc)
}
