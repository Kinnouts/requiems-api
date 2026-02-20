package workingdays

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/working-days", func(w http.ResponseWriter, r *http.Request) {
		req := WorkingDaysRequest{}

		if err := httpx.BindQuery(r, &req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		// Calculate working days
		workingDays := svc.GetWorkingDays(req.From, req.To, req.Country)

		// Build response
		response := WorkingDays{
			WorkingDays: workingDays,
			From:        req.From.Format("2006-01-02"),
			To:          req.To.Format("2006-01-02"),
			Country:     req.Country,
		}

		httpx.JSON(w, http.StatusOK, response)
	})
}
