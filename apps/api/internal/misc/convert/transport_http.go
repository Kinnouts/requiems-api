package convert

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

// RegisterRoutes mounts the unit-conversion handler on the given router.
func RegisterRoutes(r chi.Router, svc *Service) {
	// GET /convert?from=miles&to=km&value=10
	r.Get("/convert", func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		valueStr := r.URL.Query().Get("value")

		if from == "" || to == "" || valueStr == "" {
			httpx.Error(w, http.StatusBadRequest, "from, to, and value query parameters are required")
			return
		}

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "value must be a valid number")
			return
		}

		result, err := svc.Convert(from, to, value)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	})
}
