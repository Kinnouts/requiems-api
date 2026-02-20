package disposable

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func RegisterRoutes(router chi.Router, svc *Service) {
	// POST /disposable/check - Check single email
	router.Post("/disposable/check", func(w http.ResponseWriter, r *http.Request) {
		var req CheckEmailRequest
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Email == "" {
			httpx.Error(w, http.StatusBadRequest, "email is required")
			return
		}

		result := svc.CheckEmail(req.Email)
		httpx.JSON(w, http.StatusOK, result)
	})

	// POST /disposable/check-batch - Check multiple emails
	router.Post("/disposable/check-batch", func(w http.ResponseWriter, r *http.Request) {
		var req BatchCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if len(req.Emails) == 0 {
			httpx.Error(w, http.StatusBadRequest, "emails array is required and cannot be empty")
			return
		}

		if len(req.Emails) > 100 {
			httpx.Error(w, http.StatusBadRequest, "maximum 100 emails allowed per batch request")
			return
		}

		result := svc.CheckBatch(req.Emails)
		httpx.JSON(w, http.StatusOK, result)
	})

	// GET /disposable/domain/{domain} - Check if a specific domain is disposable
	router.Get("/disposable/domain/{domain}", func(w http.ResponseWriter, r *http.Request) {
		domain := chi.URLParam(r, "domain")
		if domain == "" {
			httpx.Error(w, http.StatusBadRequest, "domain is required")
			return
		}

		result := svc.CheckDomain(domain)
		httpx.JSON(w, http.StatusOK, result)
	})

	// GET /disposable/domains - Get paginated list of all disposable domains
	router.Get("/disposable/domains", func(w http.ResponseWriter, r *http.Request) {
		page := 1
		perPage := 100

		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
			if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
				perPage = pp
			}
		}

		result := svc.GetDomains(page, perPage)
		httpx.JSON(w, http.StatusOK, result)
	})

	// GET /disposable/stats - Get statistics
	router.Get("/disposable/stats", func(w http.ResponseWriter, r *http.Request) {
		result := svc.GetStats()
		httpx.JSON(w, http.StatusOK, result)
	})
}
