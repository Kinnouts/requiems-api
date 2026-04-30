package info

import (
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func RegisterRoutes(r chi.Router, svc *Service) {
	handler := httpx.Guard(svc, func(w http.ResponseWriter, r *http.Request) {
		ipStr := chi.URLParam(r, "ip")
		if ipStr == "" {
			ipStr = callerIP(r)
		}

		if net.ParseIP(ipStr) == nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid IP address")
			return
		}

		result, err := svc.CheckInfo(r.Context(), ipStr)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "internal_error", "internal error")
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	})

	r.Get("/ip/{ip}", handler)
	r.Get("/ip", handler)
}

// callerIP extracts the real client IP from the request, checking
// X-Forwarded-For, X-Real-IP, and RemoteAddr in that order.
func callerIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if before, _, ok := strings.Cut(xff, ","); ok {
			return strings.TrimSpace(before)
		}
		return strings.TrimSpace(xff)
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	addr := r.RemoteAddr
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}
