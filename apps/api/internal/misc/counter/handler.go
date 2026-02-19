package counter

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

// RegisterRoutes mounts counter handlers on the given router.
// Paths are relative to the parent mount point (e.g. /v1/misc).
func RegisterRoutes(r chi.Router, svc Service) {
	r.Post("/counter/{namespace}", incrementHandler(svc))
	r.Get("/counter/{namespace}", getHandler(svc))
}

func counterError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrInvalidNamespace) {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	httpx.Error(w, http.StatusInternalServerError, "internal server error")
}

func incrementHandler(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ns := chi.URLParam(r, "namespace")

		val, err := svc.Increment(r.Context(), ns)
		if err != nil {
			counterError(w, err)
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
			counterError(w, err)
			return
		}

		httpx.JSON(w, http.StatusOK, Counter{Namespace: ns, Value: val})
	}
}
