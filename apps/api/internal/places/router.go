package places

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	workingdays "requiems-api/internal/places/working-days"
)

func RegisterRoutes(r chi.Router, pool *pgxpool.Pool) {
	workingDaysSvc := workingdays.NewService()
	workingdays.RegisterRoutes(r, workingDaysSvc)
}
