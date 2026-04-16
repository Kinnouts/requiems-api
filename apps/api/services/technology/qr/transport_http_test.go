package qr

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	svc := NewService()
	RegisterRoutes(r, svc)
	return r
}

// ── GET /qr (PNG) ──────────────────────────────────────────────────────────

func TestQR_PNG_DefaultSize(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/qr?data=https://example.com", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "image/png" {
		t.Errorf("expected Content-Type image/png, got %q", ct)
	}

	body := w.Body.Bytes()
	if len(body) == 0 {
		t.Error("expected non-empty PNG response body")
	}

	if string(body[:4]) != "\x89PNG" {
		t.Error("expected valid PNG signature in response body")
	}
}

func TestQR_PNG_CustomSize(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/qr?data=https://example.com&size=200", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestQR_PNG_RecoveryLevels(t *testing.T) {
	r := setupRouter()

	levels := []string{"low", "medium", "high", "highest"}
	for _, level := range levels {
		req := httptest.NewRequest(http.MethodGet, "/qr?data=test&recovery="+level, http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("recovery=%q: expected status 200, got %d", level, w.Code)
		}
	}
}

func TestQR_PNG_MissingData(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/qr", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestQR_PNG_SizeTooSmall(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/qr?data=test&size=10", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestQR_PNG_SizeTooLarge(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/qr?data=test&size=2000", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestQR_PNG_InvalidRecovery(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/qr?data=test&recovery=ultra", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid recovery level, got %d", w.Code)
	}
}

// ── GET /qr/base64 (JSON) ──────────────────────────────────────────────────

func TestQR_Base64_DefaultSize(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/qr/base64?data=https://example.com", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var resp httpx.Response[Base64Response]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Image == "" {
		t.Error("expected non-empty base64 image string")
	}

	if resp.Data.Width != defaultSize {
		t.Errorf("expected width %d, got %d", defaultSize, resp.Data.Width)
	}

	if resp.Data.Height != defaultSize {
		t.Errorf("expected height %d, got %d", defaultSize, resp.Data.Height)
	}

	// Verify the base64 string decodes to valid PNG bytes.
	decoded, err := base64.StdEncoding.DecodeString(resp.Data.Image)
	if err != nil {
		t.Fatalf("failed to decode base64 image: %v", err)
	}

	if string(decoded[:4]) != "\x89PNG" {
		t.Error("expected valid PNG signature in decoded base64 data")
	}
}

func TestQR_Base64_CustomSize(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/qr/base64?data=hello&size=300", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Base64Response]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Width != 300 {
		t.Errorf("expected width 300, got %d", resp.Data.Width)
	}

	if resp.Data.Height != 300 {
		t.Errorf("expected height 300, got %d", resp.Data.Height)
	}
}

func TestQR_Base64_RecoveryLevels(t *testing.T) {
	r := setupRouter()

	levels := []string{"low", "medium", "high", "highest"}
	for _, level := range levels {
		req := httptest.NewRequest(http.MethodGet, "/qr/base64?data=test&recovery="+level, http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("recovery=%q: expected status 200, got %d", level, w.Code)
		}
	}
}

func TestQR_Base64_MissingData(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/qr/base64", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestQR_Base64_InvalidRecovery(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/qr/base64?data=test&recovery=ultra", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid recovery level, got %d", w.Code)
	}
}
