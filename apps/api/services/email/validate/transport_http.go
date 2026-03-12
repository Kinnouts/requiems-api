package validate

import (
	"context"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Post("/validate", httpx.Handle(func(ctx context.Context, req Request) (EmailValidation, error) {
		return svc.ValidateEmail(ctx, req.Email), nil
	}))
}
