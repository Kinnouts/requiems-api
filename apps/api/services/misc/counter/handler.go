package counter

import (
	"net/http"
	"regexp"
	"requiems-api/platform/httpx"
	"requiems-api/platform/middleware"

	"github.com/go-chi/chi/v5"
)

var namespaceRe = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

const namespaceValidationErrorMessage = "invalid namespace: must be 1-64 chars, alphanumeric, hyphen or underscore only"

func RegisterRoutes(r chi.Router, svc Service) {
	r.Group(func(validated chi.Router) {
		validated.Use(middleware.ValidateURLParam("namespace", namespaceRe, namespaceValidationErrorMessage))

		validated.Post("/counter/{namespace}", incrementHandler(svc))
		validated.Get("/counter/{namespace}", getHandler(svc))
	})
}

func incrementHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ns := chi.URLParam(r, "namespace")

		val, err := svc.Increment(r.Context(), ns)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "internal_error", "Failed to increment counter")
			return
		}

		httpx.JSON(w, http.StatusOK, Counter{Namespace: ns, Value: val})
	}
}

func getHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ns := chi.URLParam(r, "namespace")

		val, err := svc.Get(r.Context(), ns)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "internal_error", "Failed to get counter")
			return
		}

		httpx.JSON(w, http.StatusOK, Counter{Namespace: ns, Value: val})
	}
}