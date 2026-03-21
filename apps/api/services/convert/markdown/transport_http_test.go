package markdown

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	RegisterRoutes(r, NewService())
	return r
}

func TestMarkdown_HappyPath(t *testing.T) {
	r := setupRouter()

	body := `{"markdown":"# Hello\n\nThis is **bold** text."}`
	req := httptest.NewRequest(http.MethodPost, "/convert/markdown", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Response]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	want := "<h1>Hello</h1>\n<p>This is <strong>bold</strong> text.</p>"
	if resp.Data.HTML != want {
		t.Errorf("html mismatch\ngot:  %q\nwant: %q", resp.Data.HTML, want)
	}
}

func TestMarkdown_Sanitize(t *testing.T) {
	r := setupRouter()

	body := `{"markdown":"Hello <script>alert('xss')</script>","sanitize":true}`
	req := httptest.NewRequest(http.MethodPost, "/convert/markdown", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Response]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if strings.Contains(resp.Data.HTML, "<script>") {
		t.Errorf("expected script tag to be stripped, got: %q", resp.Data.HTML)
	}
}

func TestMarkdown_MissingBody(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/convert/markdown", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestMarkdown_EmptyMarkdown(t *testing.T) {
	r := setupRouter()

	body := `{"markdown":""}`
	req := httptest.NewRequest(http.MethodPost, "/convert/markdown", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Empty markdown triggers validation failure (required field)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
}
