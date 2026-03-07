package qr

import (
	"encoding/base64"
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

const defaultSize = 256

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/qr", func(w http.ResponseWriter, r *http.Request) {
		req := Request{Size: defaultSize}

		if err := httpx.BindQuery(r, &req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		png, err := svc.Generate(req.Data, req.Size)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "internal_error", "failed to generate QR code")
			return
		}

		if req.Format == "base64" {
			httpx.JSON(w, http.StatusOK, Base64Response{
				Image:  base64.StdEncoding.EncodeToString(png),
				Width:  req.Size,
				Height: req.Size,
			})
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(png)
	})
}
