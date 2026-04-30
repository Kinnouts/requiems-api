package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestValidateURLParam(t *testing.T) {
	alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	// okHandler records that it was reached.
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	setupRouter := func() chi.Router {
		r := chi.NewRouter()
		r.With(ValidateURLParam("id", alphanumeric, "id must be alphanumeric")).
			Get("/{id}", okHandler)
		return r
	}

	t.Run("valid param passes through", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/abc123", http.NoBody)
		w := httptest.NewRecorder()

		setupRouter().ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status: want 200, got %d", w.Code)
		}
	})

	t.Run("invalid param returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bad-param!", http.NoBody)
		w := httptest.NewRecorder()

		setupRouter().ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status: want 400, got %d", w.Code)
		}
	})
}
