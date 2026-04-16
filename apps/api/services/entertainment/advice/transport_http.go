package advice

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts advice handlers on the given router.
// Paths are relative to the parent mount point.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/advice", func(w http.ResponseWriter, r *http.Request) {
		a, err := svc.Random(r.Context())
		if err != nil {
			httpx.Error(w, http.StatusServiceUnavailable, "service_unavailable", "no advice available")
			return
		}

		httpx.JSON(w, http.StatusOK, a)
	})
}
