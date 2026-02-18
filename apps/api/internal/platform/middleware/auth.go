package middleware

import (
	"net/http"

	"requiems-api/internal/platform/httpx"
)

// BackendSecretAuth validates the X-Backend-Secret header
// This ensures only trusted sources (Cloudflare Worker) can access the API
func BackendSecretAuth(expectedSecret string) func(http.Handler) http.Handler {
	// Panic if secret is not configured (fail fast on startup)
	if expectedSecret == "" {
		panic("BACKEND_SECRET environment variable is required but not set")
	}

	// Require minimum length for security
	if len(expectedSecret) < 32 {
		panic("BACKEND_SECRET must be at least 32 characters long")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the X-Backend-Secret header
			providedSecret := r.Header.Get("X-Backend-Secret")

			// Reject if header is missing
			if providedSecret == "" {
				httpx.Error(w, http.StatusUnauthorized, "Missing X-Backend-Secret header")
				return
			}

			// Reject if secret doesn't match
			if providedSecret != expectedSecret {
				httpx.Error(w, http.StatusForbidden, "Invalid X-Backend-Secret")
				return
			}

			// Secret is valid, continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}
