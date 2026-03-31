package crypto

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// RegisterRoutes mounts crypto price handlers on the given router.
// Paths are relative to the parent mount point (/v1/finance).
func RegisterRoutes(r chi.Router, svc *Service) {
	registerCryptoRoutes(r, svc)
}

// registerCryptoRoutes wires the Getter interface to the router. Kept
// unexported so tests can inject a stub without going through the concrete
// *Service type.
func registerCryptoRoutes(r chi.Router, g Getter) {
	// GET /crypto/{symbol}
	r.Get("/crypto/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		symbol := strings.ToUpper(chi.URLParam(r, "symbol"))

		price, err := g.GetPrice(r.Context(), symbol)
		if err != nil {
			var ae *httpx.AppError
			if errors.As(err, &ae) {
				httpx.Error(w, ae.Status, ae.Code, ae.Message)
				return
			}
			httpx.Error(w, http.StatusServiceUnavailable, "upstream_error", "crypto price service unavailable")
			return
		}

		httpx.JSON(w, http.StatusOK, price)
	})
}
