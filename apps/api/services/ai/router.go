package ai

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/ai/similarity"
	"requiems-api/services/ai/detectlanguage"
)


// RegisterRoutes mounts all AI service handlers on the given router.
func RegisterRoutes(r chi.Router) {
	similaritySvc := similarity.NewService()
	similarity.RegisterRoutes(r, similaritySvc)
  
  detectlanguageSvc := detectlanguage.NewService()
	detectlanguage.RegisterRoutes(r, detectlanguageSvc)
}
