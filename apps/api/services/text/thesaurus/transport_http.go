package thesaurus

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts thesaurus handlers on the given router.
// Paths are relative to the parent mount point.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/thesaurus/{word}", func(w http.ResponseWriter, r *http.Request) {
		word := chi.URLParam(r, "word")
		if word == "" {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "word is required")
			return
		}

		result, err := svc.Lookup(word)
		if err != nil {
			httpx.Error(w, http.StatusNotFound, "not_found", "word not found in thesaurus")
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	})
}
