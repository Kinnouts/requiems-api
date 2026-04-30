package postal

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/postal/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		if code == "" {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "postal code is required")
			return
		}

		country := strings.ToUpper(r.URL.Query().Get("country"))
		if country == "" {
			country = "US"
		}

		result, ok := svc.Lookup(code, country)
		if !ok {
			httpx.Error(w, http.StatusNotFound, "not_found", "postal code not found")
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	})
}
