package advice

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/httpx"
)

// RegisterRoutes mounts advice handlers on the given router.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/v1/advice", func(w http.ResponseWriter, r *http.Request) {
		a, err := svc.Random(r.Context())
		if err != nil {
			httpx.Error(w, http.StatusServiceUnavailable, "no advice available")
			return
		}

		httpx.JSON(w, http.StatusOK, a)
	})
}


