package tech

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/tech/password"
	"requiems-api/services/tech/phone"
	"requiems-api/services/tech/qr"
	"requiems-api/services/tech/useragent"
)

func RegisterRoutes(r chi.Router) {
	phoneSvc := phone.NewService()
	phone.RegisterRoutes(r, phoneSvc)

	passwordSvc := password.NewService()
	password.RegisterRoutes(r, passwordSvc)

	uaSvc := useragent.NewService()
	useragent.RegisterRoutes(r, uaSvc)

	qrSvc := qr.NewService()
	qr.RegisterRoutes(r, qrSvc)
}
