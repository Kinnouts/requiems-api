package validation

import (
	"github.com/go-chi/chi/v5"

	emailvalidate "requiems-api/services/validation/email"
	"requiems-api/services/validation/phone"
	"requiems-api/services/validation/profanity"
)

// RegisterRoutes wires all validation sub-services onto the given router.
func RegisterRoutes(r chi.Router) {
	emailSvc := emailvalidate.NewService()
	emailvalidate.RegisterRoutes(r, emailSvc)

	phoneSvc := phone.NewService()
	phone.RegisterRoutes(r, phoneSvc)

	profanitySvc := profanity.NewService()
	profanity.RegisterRoutes(r, profanitySvc)
}
