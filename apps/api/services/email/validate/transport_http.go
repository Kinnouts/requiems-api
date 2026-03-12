package validate

import (
	"context"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes registers HTTP routes for the validate package on r.
// It registers a POST "/validate" endpoint that accepts a Request containing an email
// and responds with an EmailValidation produced by svc for the provided email.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Post("/validate", httpx.Handle(func(ctx context.Context, req Request) (EmailValidation, error) {
		return svc.ValidateEmail(ctx, req.Email), nil
	}))
}
