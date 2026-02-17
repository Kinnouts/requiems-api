package workingdays

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/working-days", func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		country := r.URL.Query().Get("country")

		// Validate required parameters
		if fromStr == "" {
			httpx.Error(w, http.StatusBadRequest, "from parameter is required (format: YYYY-MM-DD)")
			return
		}

		if toStr == "" {
			httpx.Error(w, http.StatusBadRequest, "to parameter is required (format: YYYY-MM-DD)")
			return
		}

		// Parse dates
		from, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid from date format, expected YYYY-MM-DD")
			return
		}

		to, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid to date format, expected YYYY-MM-DD")
			return
		}

		// Validate date range
		if to.Before(from) {
			httpx.Error(w, http.StatusBadRequest, "to date must be on or after from date")
			return
		}

		// Calculate working days
		workingDays := svc.GetWorkingDays(from, to, country)

		// Build response
		response := WorkingDays{
			WorkingDays: workingDays,
			From:        fromStr,
			To:          toStr,
			Country:     country,
		}

		httpx.JSON(w, http.StatusOK, response)
	})
}
