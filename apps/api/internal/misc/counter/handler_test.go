package counter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

var errRedisDown = errors.New("connection refused")

// mockService is a test double for Service.
type mockService struct {
	incrementFn func(ctx context.Context, namespace string) (int64, error)
	getFn       func(ctx context.Context, namespace string) (int64, error)
}

func (m *mockService) Increment(ctx context.Context, namespace string) (int64, error) {
	return m.incrementFn(ctx, namespace)
}

func (m *mockService) Get(ctx context.Context, namespace string) (int64, error) {
	return m.getFn(ctx, namespace)
}

func newTestRouter(svc Service) http.Handler {
	r := chi.NewRouter()
	RegisterRoutes(r, svc)
	return r
}

func TestIncrementHandler(t *testing.T) {
	t.Run("returns updated counter value", func(t *testing.T) {
		svc := &mockService{
			incrementFn: func(_ context.Context, ns string) (int64, error) {
				return 5, nil
			},
		}

		req := httptest.NewRequest(http.MethodPost, "/counter/hits", http.NoBody)
		w := httptest.NewRecorder()

		newTestRouter(svc).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		var got Counter
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if got.Namespace != "hits" {
			t.Errorf("namespace: want %q got %q", "hits", got.Namespace)
		}

		if got.Value != 5 {
			t.Errorf("value: want 5 got %d", got.Value)
		}
	})

	t.Run("returns 400 for invalid namespace via service error", func(t *testing.T) {
		svc := &mockService{
			incrementFn: func(_ context.Context, ns string) (int64, error) {
				return 0, validateNamespace("!!!invalid!!!")
			},
		}

		req := httptest.NewRequest(http.MethodPost, "/counter/hits", http.NoBody)
		w := httptest.NewRecorder()

		newTestRouter(svc).ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 500 for internal server error", func(t *testing.T) {
		svc := &mockService{
			incrementFn: func(_ context.Context, ns string) (int64, error) {
				return 0, errRedisDown
			},
		}

		req := httptest.NewRequest(http.MethodPost, "/counter/hits", http.NoBody)
		w := httptest.NewRecorder()

		newTestRouter(svc).ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}

func TestGetHandler(t *testing.T) {
	t.Run("returns counter value", func(t *testing.T) {
		svc := &mockService{
			getFn: func(_ context.Context, ns string) (int64, error) {
				return 42, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/counter/page-views", http.NoBody)
		w := httptest.NewRecorder()

		newTestRouter(svc).ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		var got Counter
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if got.Namespace != "page-views" {
			t.Errorf("namespace: want %q got %q", "page-views", got.Namespace)
		}

		if got.Value != 42 {
			t.Errorf("value: want 42 got %d", got.Value)
		}
	})

	t.Run("returns 400 on validation error", func(t *testing.T) {
		svc := &mockService{
			getFn: func(_ context.Context, ns string) (int64, error) {
				return 0, validateNamespace("!!!invalid!!!")
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/counter/page-views", http.NoBody)
		w := httptest.NewRecorder()

		newTestRouter(svc).ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 500 for internal server error", func(t *testing.T) {
		svc := &mockService{
			getFn: func(_ context.Context, ns string) (int64, error) {
				return 0, errRedisDown
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/counter/page-views", http.NoBody)
		w := httptest.NewRecorder()

		newTestRouter(svc).ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}

func TestValidateNamespace(t *testing.T) {
	valid := []string{"hits", "page-views", "my_counter", "a", "A1b2-c3_d4"}
	for _, ns := range valid {
		if err := validateNamespace(ns); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", ns, err)
		}
	}

	invalid := []string{"", "has space", "has/slash", "has.dot", "toolong_toolong_toolong_toolong_toolong_toolong_toolong_toolong_toolong"}
	for _, ns := range invalid {
		if err := validateNamespace(ns); err == nil {
			t.Errorf("expected %q to be invalid", ns)
		}
	}
}
