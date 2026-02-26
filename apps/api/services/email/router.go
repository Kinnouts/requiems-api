package email

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/email/disposable"
)

func RegisterRoutes(r chi.Router) {
	disposableSvc := disposable.NewService()
	disposable.RegisterRoutes(r, disposableSvc)
}
