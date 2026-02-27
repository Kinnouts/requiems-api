package entertainment

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/entertainment/horoscope"
)

func RegisterRoutes(r chi.Router) {
	horoscopeSvc := horoscope.NewService()
	horoscope.RegisterRoutes(r, horoscopeSvc)
}
