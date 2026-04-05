package convert

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/convert/base64"
	"requiems-api/services/convert/markdown"
	"requiems-api/services/convert/numbase"
)

func RegisterRoutes(r chi.Router) {
	markdownSvc := markdown.NewService()
	markdown.RegisterRoutes(r, markdownSvc)

	base64Svc := base64.NewService()
	base64.RegisterRoutes(r, base64Svc)

	numbaseSvc := numbase.NewService()
	numbase.RegisterRoutes(r, numbaseSvc)
}
