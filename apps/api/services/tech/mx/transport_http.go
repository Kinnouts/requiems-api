package mx

import (
	"net"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// domainRe matches valid fully-qualified domain names.
var domainRe = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)

func RegisterRoutes(r chi.Router, svc *Service) {
	r.Get("/mx/{domain}", func(w http.ResponseWriter, r *http.Request) {
		domain := chi.URLParam(r, "domain")

		if !domainRe.MatchString(domain) {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid domain name")
			return
		}

		result, err := svc.Lookup(r.Context(), domain)
		if err != nil {
			if isDNSNotFound(err) {
				httpx.Error(w, http.StatusNotFound, "not_found", "no MX records found for domain")
				return
			}
			httpx.Error(w, http.StatusInternalServerError, "internal_error", "internal error")
			return
		}

		httpx.JSON(w, http.StatusOK, result)
	})
}

// isDNSNotFound reports whether the error is a DNS "no such host" / NXDOMAIN error.
func isDNSNotFound(err error) bool {
	if dnsErr, ok := err.(*net.DNSError); ok {
		return dnsErr.IsNotFound
	}
	return false
}
