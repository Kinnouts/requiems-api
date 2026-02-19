package misc

import (
	"github.com/go-chi/chi/v5"

	"requiems-api/internal/misc/convert"
)

func RegisterRoutes(r chi.Router) {
	convertSvc := convert.NewService()
	convert.RegisterRoutes(r, convertSvc)
}
