package randomuser

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts the random-user handler on the given router.
// Paths are relative to the parent mount point (e.g. /v1/misc).
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/random-user", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, http.StatusOK, svc.Generate())
	})
}
