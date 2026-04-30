package randomuser

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func newTestRouter(svc *Service) http.Handler {
	r := chi.NewRouter()
	RegisterRoutes(r, svc)
	return r
}

func TestRandomUserHandler(t *testing.T) {
	t.Run("returns 200 with valid user fields", func(t *testing.T) {
		svc := NewService()

		req := httptest.NewRequest(http.MethodGet, "/random-user", http.NoBody)
		w := httptest.NewRecorder()

		newTestRouter(svc).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		var resp httpx.Response[User]
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		u := resp.Data
		if u.Name == "" {
			t.Error("Name should not be empty")
		}
		if u.Email == "" {
			t.Error("Email should not be empty")
		}
		if u.Phone == "" {
			t.Error("Phone should not be empty")
		}
		if u.Address.Street == "" {
			t.Error("Address.Street should not be empty")
		}
		if u.Avatar == "" {
			t.Error("Avatar should not be empty")
		}
	})

	t.Run("content-type is application/json", func(t *testing.T) {
		svc := NewService()

		req := httptest.NewRequest(http.MethodGet, "/random-user", http.NoBody)
		w := httptest.NewRecorder()

		newTestRouter(svc).ServeHTTP(w, req)

		ct := w.Header().Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}
	})

	t.Run("returns different users on successive calls", func(t *testing.T) {
		svc := NewService()
		router := newTestRouter(svc)

		names := make(map[string]struct{})
		for range 10 {
			req := httptest.NewRequest(http.MethodGet, "/random-user", http.NoBody)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var resp httpx.Response[User]
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			names[resp.Data.Name] = struct{}{}
		}

		// With 30 first × 30 last = 900 combinations, 10 calls should yield > 1 unique name.
		if len(names) <= 1 {
			t.Error("expected varied output across multiple calls")
		}
	})
}
