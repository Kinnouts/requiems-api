package profanity

import (
	"context"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

// RegisterRoutes mounts the profanity check handler on the given router.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Post("/profanity", httpx.Handle(
		func(ctx context.Context, req Request) (Result, error) {
			return svc.Check(req.Text), nil
		},
	))
}
