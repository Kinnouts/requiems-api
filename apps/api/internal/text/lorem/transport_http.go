package lorem

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/lorem", func(w http.ResponseWriter, r *http.Request) {
		// Set defaults before binding so unset params keep their default value.
		req := LoremRequest{Paragraphs: 1, Sentences: 5}

		if err := httpx.BindQuery(r, &req); err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		httpx.JSON(w, http.StatusOK, svc.Generate(req.Paragraphs, req.Sentences))
	})
}
