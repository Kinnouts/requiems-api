package ai

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/ai/detectlanguage"
)

// RegisterRoutes wires all AI domain services onto the given router.
func RegisterRoutes(r chi.Router) {
	detectlanguageSvc := detectlanguage.NewService()
	detectlanguage.RegisterRoutes(r, detectlanguageSvc)
}
