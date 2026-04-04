package places

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"

	"requiems-api/platform/config"
	"requiems-api/services/places/cities"
	"requiems-api/services/places/geocode"
	"requiems-api/services/places/holidays"
	"requiems-api/services/places/postal"
	"requiems-api/services/places/timezone"
	workingdays "requiems-api/services/places/working-days"
)

func RegisterRoutes(r chi.Router, cfg config.Config, rdb *redis.Client) {
	workingDaysSvc := workingdays.NewService()
	workingdays.RegisterRoutes(r, workingDaysSvc)

	timezoneSvc, err := timezone.NewService()
	if err != nil {
		log.Printf("places: failed to initialize timezone service: %v", err)
	}
	timezone.RegisterRoutes(r, timezoneSvc)

	holidaysSvc := holidays.NewService()
	holidays.RegisterRoutes(r, holidaysSvc)

	postalSvc := postal.NewService(cfg.PostalCodeDBPath)
	postal.RegisterRoutes(r, postalSvc)

	geocodeSvc := geocode.NewService(cfg.NominatimURL, &http.Client{Timeout: 10 * time.Second}, rdb)
	geocode.RegisterRoutes(r, geocodeSvc)

	citiesSvc := cities.NewService(cfg.CitiesDBPath)
	cities.RegisterRoutes(r, citiesSvc)
}
