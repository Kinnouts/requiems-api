package quotes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/quotes/random", func(w http.ResponseWriter, r *http.Request) {
		q, err := svc.Random(r.Context())

		if err != nil {
			httpx.Error(w, http.StatusServiceUnavailable, "no quotes available")
			return
		}

		httpx.JSON(w, http.StatusOK, q)
	})
}
