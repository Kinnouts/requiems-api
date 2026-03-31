package inflation

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts inflation handlers on the given router.
// Paths are relative to the parent mount point (/v1/finance).
func RegisterRoutes(r chi.Router, svc *Service) {
	registerInflationRoutes(r, svc)
}

// registerInflationRoutes wires the Getter interface to the router. Kept
// unexported so tests can inject a stub without going through the concrete
// *Service type.
func registerInflationRoutes(r chi.Router, g Getter) {
	// GET /inflation?country=US — return latest and historical CPI inflation rate
	r.Get("/inflation", func(w http.ResponseWriter, r *http.Request) {
		// Uppercase the country param before binding — iso3166_1_alpha2 is case-sensitive.
		if country := r.URL.Query().Get("country"); country != "" {
			q := r.URL.Query()
			q.Set("country", strings.ToUpper(country))
			r.URL.RawQuery = q.Encode()
		}

		var req Request
		if err := httpx.BindQuery(r, &req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		resp, err := g.GetInflation(r.Context(), req.Country)
		if err != nil {
			if ae, ok := err.(*httpx.AppError); ok {
				httpx.Error(w, ae.Status, ae.Code, ae.Message)
				return
			}
			httpx.Error(w, http.StatusInternalServerError, "internal_error", "internal server error")
			return
		}

		httpx.JSON(w, http.StatusOK, resp)
	})
}
