package vpn

import (
	"net"
	"net/http"
	"requiems-api/platform/httpx"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/ip/vpn/{ip}", func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)

		result, err := svc.CheckIP(r.Context(), ip)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	})
}

func getIP(r *http.Request) string {
	if ip := chi.URLParam(r, "ip"); ip != "" {
		return ip
	}

	addr := r.RemoteAddr
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}

	return addr
}
