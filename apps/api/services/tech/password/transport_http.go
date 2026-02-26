package password

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/password", func(w http.ResponseWriter, r *http.Request) {
		req := Request{Length: 16}

		if err := httpx.BindQuery(r, &req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		result, err := svc.Generate(req.Length, req.Uppercase, req.Numbers, req.Symbols)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "internal_error", "failed to generate password")
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	})
}
