package swift

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts SWIFT/BIC lookup handlers on the given router.
// Paths are relative to the parent mount point (/v1/finance).
func RegisterRoutes(r chi.Router, svc *Service) {
	registerSWIFTRoutes(r, svc)
}

// registerSWIFTRoutes wires the Looker interface to the router. Kept unexported
// so tests can inject a stub without going through the concrete *Service type.
func registerSWIFTRoutes(r chi.Router, l Looker) {
	// GET /swift/{code} — look up bank metadata for a SWIFT/BIC code
	r.Get("/swift/{code}", func(w http.ResponseWriter, r *http.Request) {
		rawCode := chi.URLParam(r, "code")

		result, err := l.Lookup(r.Context(), rawCode)
		if err != nil {
			if ae, ok := err.(*httpx.AppError); ok {
				httpx.Error(w, ae.Status, ae.Code, ae.Message)
				return
			}
			httpx.Error(w, http.StatusInternalServerError, "internal_error", "internal server error")
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	})
}
