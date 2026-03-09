package barcode

import (
	"encoding/base64"
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	// GET /barcode — returns a raw PNG image.
	r.Get("/barcode", func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if err := httpx.BindQuery(r, &req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		png, _, _, err := svc.Generate(req.Data, req.Type)
		if err != nil {
			httpx.Error(w, http.StatusUnprocessableEntity, "unprocessable_entity", err.Error())
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(png)
	})

	// GET /barcode/base64 — returns a JSON envelope with a base64-encoded PNG.
	r.Get("/barcode/base64", func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if err := httpx.BindQuery(r, &req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		png, width, height, err := svc.Generate(req.Data, req.Type)
		if err != nil {
			httpx.Error(w, http.StatusUnprocessableEntity, "unprocessable_entity", err.Error())
			return
		}

		httpx.JSON(w, http.StatusOK, Base64Response{
			Image:  base64.StdEncoding.EncodeToString(png),
			Type:   req.Type,
			Width:  width,
			Height: height,
		})
	})
}
