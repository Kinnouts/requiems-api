package numbase

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts the base conversion handler on the given router.
// Paths are relative to the parent mount point.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/base", func(w http.ResponseWriter, r *http.Request) {
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		value := r.URL.Query().Get("value")

		if fromStr == "" || toStr == "" || value == "" {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "from, to, and value parameters are required")
			return
		}

		fromBase, err := strconv.Atoi(fromStr)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "from must be a valid integer")
			return
		}

		toBase, err := strconv.Atoi(toStr)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "to must be a valid integer")
			return
		}

		result, err := svc.Convert(value, fromBase, toBase)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	})
}
