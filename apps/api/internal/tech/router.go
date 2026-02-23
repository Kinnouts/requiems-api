package tech

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/internal/tech/phone"
)

func RegisterRoutes(r chi.Router) {
	phoneSvc := phone.NewService()
	phone.RegisterRoutes(r, phoneSvc)
  
	passwordSvc := password.NewService()
	password.RegisterRoutes(r, passwordSvc)
)