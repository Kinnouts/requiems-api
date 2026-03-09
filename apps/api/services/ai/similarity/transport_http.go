package similarity

import (
	"context"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts the text similarity handler on the given router.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Post("/similarity", httpx.Handle(
		func(_ context.Context, req Request) (Result, error) {
			return svc.Cosine(req.Text1, req.Text2), nil
		},
	))
}
