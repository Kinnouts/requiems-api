package vpn

import (
	"net"
	"net/http"
	"requiems-api/platform/httpx"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/ip/vpn/{ip}", httpx.Guard(svc, func(w http.ResponseWriter, r *http.Request) {
		ip := net.ParseIP(getIP(r))
		if ip == nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid IP address")
			return
		}

		result, err := svc.CheckIP(r.Context(), ip)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "internal server error", "failed to check IP address")
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	}))
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
