package counter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"requiems-api/platform/httpx"
	"testing"

	"github.com/go-chi/chi/v5"
)

var errRedisDown = errors.New("connection refused")

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

		var resp httpx.Response[Counter]

		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		got := resp.Data

		if got.Namespace != "hits" {
			t.Errorf("namespace: want %q got %q", "hits", got.Namespace)
		}

		if got.Value != 5 {
			t.Errorf("value: want 5 got %d", got.Value)
		}
	})

	t.Run("returns 400 for invalid namespace from URL param validation", func(t *testing.T) {
		svc := &mockService{
			incrementFn: func(_ context.Context, ns string) (int64, error) {
				return 1, nil
			},
		}

		req := httptest.NewRequest(http.MethodPost, "/counter/!!!invalid!!!", http.NoBody)
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

		var resp httpx.Response[Counter]
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		got := resp.Data

		if got.Namespace != "page-views" {
			t.Errorf("namespace: want %q got %q", "page-views", got.Namespace)
		}

		if got.Value != 42 {
			t.Errorf("value: want 42 got %d", got.Value)
		}
	})

	t.Run("returns 400 for invalid namespace from URL param validation", func(t *testing.T) {
		svc := &mockService{
			getFn: func(_ context.Context, ns string) (int64, error) {
				return 42, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/counter/!!!invalid!!!", http.NoBody)
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
