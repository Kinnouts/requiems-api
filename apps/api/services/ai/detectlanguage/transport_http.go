package detectlanguage

import (
	"context"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts the detect-language handler on the given router.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Post("/detect-language", httpx.Handle(
		func(_ context.Context, req Request) (Result, error) {
			return svc.Detect(req.Text), nil
		},
	))
}
