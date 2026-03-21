package convert

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/convert/base64"
)

// RegisterRoutes wires all convert-domain services onto r.
// Routes are relative to the mount point (e.g. /v1/convert).
func RegisterRoutes(r chi.Router) {
	base64Svc := base64.NewService()
	base64.RegisterRoutes(r, base64Svc)
}
