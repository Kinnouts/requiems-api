package email

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/email/disposable"
	"requiems-api/services/email/normalize"
	"requiems-api/services/email/validate"
)

func RegisterRoutes(r chi.Router) {
	disposableSvc := disposable.NewService()
	disposable.RegisterRoutes(r, disposableSvc)

	normalizeSvc := normalize.NewService()
	normalize.RegisterRoutes(r, normalizeSvc)

	validateSvc := validate.NewService()
	validate.RegisterRoutes(r, validateSvc)
}
