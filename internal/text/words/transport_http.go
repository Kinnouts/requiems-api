package words

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

// RegisterRoutes mounts words handlers on the given router.
// Paths are relative to the parent mount point.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/words/random", func(w http.ResponseWriter, r *http.Request) {
		wrd, err := svc.Random(r.Context())
		if err != nil {
			httpx.Error(w, http.StatusServiceUnavailable, "no words available")
			return
		}

		httpx.JSON(w, http.StatusOK, wrd)
	})
}
