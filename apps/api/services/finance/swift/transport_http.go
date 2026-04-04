package swift

import (
	"net/http"
	"strconv"

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
	// GET /swift — list SWIFT records with optional filters.
	r.Get("/swift", func(w http.ResponseWriter, r *http.Request) {
		filter, err := buildListFilter(r)
		if err != nil {
			if ae, ok := err.(*httpx.AppError); ok {
				httpx.Error(w, ae.Status, ae.Code, ae.Message)
				return
			}
			httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid query parameters")
			return
		}

		result, err := l.List(r.Context(), filter)
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

	// GET /swift/country/{country_code} — list SWIFT records for a country.
	r.Get("/swift/country/{country_code}", func(w http.ResponseWriter, r *http.Request) {
		filter, err := buildListFilter(r)
		if err != nil {
			if ae, ok := err.(*httpx.AppError); ok {
				httpx.Error(w, ae.Status, ae.Code, ae.Message)
				return
			}
			httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid query parameters")
			return
		}

		filter.CountryCode = chi.URLParam(r, "country_code")

		result, err := l.List(r.Context(), filter)
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

func buildListFilter(r *http.Request) (ListFilter, error) {
	q := r.URL.Query()
	filter := ListFilter{
		CountryCode: q.Get("country_code"),
		BankCode:    q.Get("bank_code"),
		Query:       q.Get("q"),
		Limit:       50,
		Offset:      0,
	}

	if raw := q.Get("limit"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return ListFilter{}, &httpx.AppError{Status: http.StatusBadRequest, Code: "bad_request", Message: "limit must be an integer"}
		}
		filter.Limit = v
	}

	if raw := q.Get("offset"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return ListFilter{}, &httpx.AppError{Status: http.StatusBadRequest, Code: "bad_request", Message: "offset must be an integer"}
		}
		filter.Offset = v
	}

	return filter, nil
}
