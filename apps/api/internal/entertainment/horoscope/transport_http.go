package horoscope

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

// RegisterRoutes mounts horoscope handlers on the given router.
// Paths are relative to the parent mount point.
func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/horoscope/{sign}", func(w http.ResponseWriter, r *http.Request) {
		sign := strings.ToLower(chi.URLParam(r, "sign"))
		if !IsValidSign(sign) {
			httpx.Error(w, http.StatusBadRequest, "invalid zodiac sign")
			return
		}

		h, err := svc.Daily(sign)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "failed to generate horoscope")
			return
		}

		httpx.JSON(w, http.StatusOK, h)
	})
}
