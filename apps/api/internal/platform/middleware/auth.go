package middleware

import (
	"net/http"

	"requiems-api/internal/platform/httpx"
)

// Validates the X-Backend-Secret header
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
				httpx.Error(w, http.StatusUnauthorized, "Missing X-Backend-Secret header")
				return
			}

			if providedSecret != expectedSecret {
				httpx.Error(w, http.StatusForbidden, "Invalid X-Backend-Secret")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
