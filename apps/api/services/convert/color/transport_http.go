package color //nolint:revive // package name matches the service domain it implements

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts the color conversion handler on the given router.
// Path is relative to the parent mount point.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/color", func(w http.ResponseWriter, r *http.Request) {
		req := Request{}

		if err := httpx.BindQuery(r, &req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		result, err := svc.Convert(req.From, req.To, req.Value)
		if err != nil {
			var appErr *httpx.AppError
			if errors.As(err, &appErr) {
				httpx.Error(w, appErr.Status, appErr.Code, appErr.Message)
				return
			}
			httpx.Error(w, http.StatusInternalServerError, "internal_error", "failed to convert color")
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	})
}
