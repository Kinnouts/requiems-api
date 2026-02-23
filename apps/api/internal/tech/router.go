package tech

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/internal/tech/useragent"
)

func RegisterRoutes(r chi.Router) {
	uaSvc := useragent.NewService()
	useragent.RegisterRoutes(r, uaSvc)
}
