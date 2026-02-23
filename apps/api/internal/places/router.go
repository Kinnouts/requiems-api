package places

import (
	"log"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/places/timezone"
	workingdays "requiems-api/internal/places/working-days"
)

func RegisterRoutes(r chi.Router) {
	workingDaysSvc := workingdays.NewService()
	workingdays.RegisterRoutes(r, workingDaysSvc)

	timezoneSvc, err := timezone.NewService()
	if err != nil {
		log.Fatalf("places: failed to initialize timezone service: %v", err)
	}
	timezone.RegisterRoutes(r, timezoneSvc)
}
