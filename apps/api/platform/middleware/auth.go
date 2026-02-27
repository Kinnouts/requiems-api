package middleware

import (
	"net/http"

	"requiems-api/platform/httpx"
)

// BackendSecretAuth validates the X-Backend-Secret header, ensuring only the
// Cloudflare Worker gateway can reach protected routes.
func BackendSecretAuth(expectedSecret string) func(http.Handler) http.Handler {
	if expectedSecret == "" {
		panic("BACKEND_SECRET environment variable is required but not set")
	}

	if len(expectedSecret) < 32 {
		panic("BACKEND_SECRET must be at least 32 characters long")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			providedSecret := r.Header.Get("X-Backend-Secret")

			if providedSecret == "" {
				httpx.Error(w, http.StatusUnauthorized, "unauthorized", "Missing X-Backend-Secret header")
				return
			}

			if providedSecret != expectedSecret {
				httpx.Error(w, http.StatusForbidden, "forbidden", "Invalid X-Backend-Secret")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
