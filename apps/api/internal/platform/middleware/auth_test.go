package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBackendSecretAuth(t *testing.T) {
	validSecret := "this_is_a_valid_secret_with_32_chars_minimum"

	// Test handler that just returns 200 OK
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authorized"))
	})

	t.Run("allows request with valid secret", func(t *testing.T) {
		middleware := BackendSecretAuth(validSecret)
		handler := middleware(testHandler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Backend-Secret", validSecret)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("rejects request with missing header", func(t *testing.T) {
		middleware := BackendSecretAuth(validSecret)
		handler := middleware(testHandler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("rejects request with invalid secret", func(t *testing.T) {
		middleware := BackendSecretAuth(validSecret)
		handler := middleware(testHandler)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Backend-Secret", "wrong_secret")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected status 403, got %d", w.Code)
		}
	})

	t.Run("panics if secret is empty", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when secret is empty")
			}
		}()

		BackendSecretAuth("")
	})

	t.Run("panics if secret is too short", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when secret is too short")
			}
		}()

		BackendSecretAuth("short")
	})
}
