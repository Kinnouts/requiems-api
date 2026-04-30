package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	captureHandler := func(gotID *string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*gotID = GetRequestID(r.Context())
			w.WriteHeader(http.StatusOK)
		})
	}

	t.Run("propagates existing X-Request-ID", func(t *testing.T) {
		var ctxID string
		handler := RequestID(captureHandler(&ctxID))

		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		req.Header.Set("X-Request-ID", "trace-abc-123")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if ctxID != "trace-abc-123" {
			t.Errorf("context id = %q, want %q", ctxID, "trace-abc-123")
		}
		if got := w.Header().Get("X-Request-ID"); got != "trace-abc-123" {
			t.Errorf("response header = %q, want %q", got, "trace-abc-123")
		}
	})

	t.Run("generates ID when header is absent", func(t *testing.T) {
		var ctxID string
		handler := RequestID(captureHandler(&ctxID))

		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if ctxID == "" {
			t.Error("expected a generated request ID in context")
		}
		if got := w.Header().Get("X-Request-ID"); got != ctxID {
			t.Errorf("response header %q does not match context ID %q", got, ctxID)
		}
	})

	t.Run("generated IDs are unique", func(t *testing.T) {
		ids := make(map[string]struct{}, 100)
		for i := 0; i < 100; i++ {
			var ctxID string
			handler := RequestID(captureHandler(&ctxID))
			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			handler.ServeHTTP(httptest.NewRecorder(), req)
			ids[ctxID] = struct{}{}
		}
		if len(ids) != 100 {
			t.Errorf("expected 100 unique IDs, got %d", len(ids))
		}
	})
}

func TestGetRequestID_missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	if id := GetRequestID(req.Context()); id != "" {
		t.Errorf("expected empty string, got %q", id)
	}
}
