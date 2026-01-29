package lorem

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/lorem", func(w http.ResponseWriter, r *http.Request) {
		// TODO(crydafan): implement lorem endpoint
		httpx.Error(w, http.StatusNotImplemented, "lorem endpoint not implemented")
	})
}
