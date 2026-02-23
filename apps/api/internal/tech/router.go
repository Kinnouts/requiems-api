package tech

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/internal/tech/password"
)

func RegisterRoutes(r chi.Router) {
	passwordSvc := password.NewService()
	password.RegisterRoutes(r, passwordSvc)
}
