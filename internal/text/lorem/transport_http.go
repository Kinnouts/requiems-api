package lorem

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/lorem", func(w http.ResponseWriter, r *http.Request) {
		paragraphs := 1
		sentences := 5

		if p := r.URL.Query().Get("paragraphs"); p != "" {
			val, err := strconv.Atoi(p)
			if err != nil || val <= 0 {
				httpx.Error(w, http.StatusBadRequest, "invalid paragraphs parameter")
				return
			}
			paragraphs = val
		}

		if wc := r.URL.Query().Get("sentences"); wc != "" {
			val, err := strconv.Atoi(wc)
			if err != nil || val <= 0 {
				httpx.Error(w, http.StatusBadRequest, "invalid sentences parameter")
				return
			}
			sentences = val
		}

		lorem := svc.Generate(paragraphs, sentences)
		httpx.JSON(w, http.StatusOK, lorem)
	})
}
