package useragent

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/useragent", func(w http.ResponseWriter, r *http.Request) {
		ua := r.URL.Query().Get("ua")
		if ua == "" {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "ua parameter is required")
			return
		}

		httpx.JSON(w, http.StatusOK, svc.Parse(ua))
	})
}
