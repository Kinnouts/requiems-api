package ai

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/ai/similarity"
)

// RegisterRoutes mounts all AI service handlers on the given router.
func RegisterRoutes(r chi.Router) {
	similaritySvc := similarity.NewService()
	similarity.RegisterRoutes(r, similaritySvc)
}
