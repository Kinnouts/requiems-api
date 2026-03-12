package phone

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/validate/phone", func(w http.ResponseWriter, r *http.Request) {
		req := ValidateRequest{}

		if err := httpx.BindQuery(r, &req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		httpx.JSON(w, http.StatusOK, svc.Validate(req.Number))
	})

	r.Post("/validate/phone/batch", httpx.Handle(
		func(_ context.Context, req BatchValidateRequest) (BatchValidateResponse, error) {
			return svc.ValidateBatch(req.Numbers), nil
		},
	))
}
