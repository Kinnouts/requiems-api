package convert

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/services/convert/markdown"
)

func RegisterRoutes(r chi.Router) {
	markdownSvc := markdown.NewService()
	markdown.RegisterRoutes(r, markdownSvc)
}
