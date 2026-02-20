package places

import (
	"github.com/go-chi/chi/v5"

	workingdays "requiems-api/internal/places/working-days"
)

func RegisterRoutes(r chi.Router) {
	workingDaysSvc := workingdays.NewService()
	workingdays.RegisterRoutes(r, workingDaysSvc)
}
