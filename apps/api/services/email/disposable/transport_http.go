package disposable

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func RegisterRoutes(router chi.Router, svc *Service) {
	// POST /disposable/check — single email check
	router.Post("/disposable/check", httpx.Handle(
		func(_ context.Context, req CheckEmailRequest) (CheckEmailResponse, error) {
			return svc.CheckEmail(req.Email), nil
		},
	))

	// POST /disposable/check-batch — batch email check (max 100)
	router.Post("/disposable/check-batch", httpx.Handle(
		func(_ context.Context, req BatchCheckRequest) (BatchCheckResponse, error) {
			return svc.CheckBatch(req.Emails), nil
		},
	))

	// GET /disposable/domain/{domain} — check a specific domain
	router.Get("/disposable/domain/{domain}", func(w http.ResponseWriter, r *http.Request) {
		domain := chi.URLParam(r, "domain")

		if domain == "" {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "domain is required")
			return
		}

		httpx.JSON(w, http.StatusOK, svc.CheckDomain(domain))
	})

	// GET /disposable/domains — paginated list of all disposable domains
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

		httpx.JSON(w, http.StatusOK, svc.GetDomains(page, perPage))
	})

	// GET /disposable/stats — blocklist statistics
	router.Get("/disposable/stats", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, http.StatusOK, svc.GetStats())
	})
}
