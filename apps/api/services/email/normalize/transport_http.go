package normalize

import (
	"context"
	"requiems-api/platform/httpx"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router, svc *Service) {
	router.Post("/normalize", httpx.Handle(func(_ context.Context, req EmailNormalizationRequest) (EmailNormalization, error) {
		return svc.Normalize(req.Email)
	}))
}
