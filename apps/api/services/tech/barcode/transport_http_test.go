package barcode

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

// ── GET /barcode (PNG) ─────────────────────────────────────────────────────

func TestBarcode_PNG_Code128(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/barcode?data=123456789&type=code128", http.NoBody)
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

func TestBarcode_PNG_AllTypes(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		name    string
		query   string
	}{
		{"code128", "/barcode?data=HELLO123&type=code128"},
		{"code93", "/barcode?data=HELLO&type=code93"},
		{"code39", "/barcode?data=HELLO&type=code39"},
		{"ean8", "/barcode?data=1234567&type=ean8"},
		{"ean13", "/barcode?data=123456789012&type=ean13"},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, tc.query, http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("type=%q: expected status 200, got %d", tc.name, w.Code)
		}
	}
}

func TestBarcode_PNG_MissingData(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/barcode?type=code128", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestBarcode_PNG_MissingType(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/barcode?data=123456789", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestBarcode_PNG_InvalidType(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/barcode?data=123456789&type=invalid", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid type, got %d", w.Code)
	}
}

func TestBarcode_PNG_InvalidEAN8Data(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/barcode?data=123&type=ean8", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", w.Code)
	}
}

// ── GET /barcode/base64 (JSON) ─────────────────────────────────────────────

func TestBarcode_Base64_Code128(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/barcode/base64?data=123456789&type=code128", http.NoBody)
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

	if resp.Data.Type != "code128" {
		t.Errorf("expected type code128, got %q", resp.Data.Type)
	}

	if resp.Data.Width != defaultWidth {
		t.Errorf("expected width %d, got %d", defaultWidth, resp.Data.Width)
	}

	if resp.Data.Height != defaultHeight {
		t.Errorf("expected height %d, got %d", defaultHeight, resp.Data.Height)
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

func TestBarcode_Base64_MissingData(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/barcode/base64?type=code128", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestBarcode_Base64_InvalidType(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/barcode/base64?data=123456789&type=invalid", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid type, got %d", w.Code)
	}
}
