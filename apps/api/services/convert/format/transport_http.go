package convformat

import (
	"context"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts the format conversion handler on the given router.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Post("/format", httpx.Handle(
		func(_ context.Context, req Request) (Response, error) {
			return svc.Convert(req)
		},
	))
}
