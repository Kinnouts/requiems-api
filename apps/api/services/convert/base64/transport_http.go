package base64

import (
	"context"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts Base64 encode and decode handlers on the given router.
// Paths are relative to the parent mount point.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Post("/base64/encode", httpx.Handle(
		func(_ context.Context, req EncodeRequest) (Result, error) {
			return svc.Encode(req.Value, req.Variant), nil
		},
	))

	r.Post("/base64/decode", httpx.Handle(
		func(_ context.Context, req DecodeRequest) (Result, error) {
			return svc.Decode(req.Value, req.Variant)
		},
	))
}
