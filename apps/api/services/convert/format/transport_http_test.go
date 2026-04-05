package convformat_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
	convformat "requiems-api/services/convert/format"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	convformat.RegisterRoutes(r, convformat.NewService())
	return r
}

func TestFormat_HappyPath_JSONToYAML(t *testing.T) {
	r := setupRouter()

	body := `{"from":"json","to":"yaml","content":"{\"name\":\"Alice\",\"age\":30}"}`
	req := httptest.NewRequest(http.MethodPost, "/convert/format", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[convformat.Response]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !strings.Contains(resp.Data.Result, "Alice") {
		t.Errorf("expected YAML with 'Alice', got %q", resp.Data.Result)
	}
}

func TestFormat_HappyPath_CSVToJSON(t *testing.T) {
	r := setupRouter()

	body := `{"from":"csv","to":"json","content":"name,age\nAlice,30\n"}`
	req := httptest.NewRequest(http.MethodPost, "/convert/format", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[convformat.Response]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !strings.Contains(resp.Data.Result, "Alice") {
		t.Errorf("expected JSON with 'Alice', got %q", resp.Data.Result)
	}
}

func TestFormat_InvalidFromFormat(t *testing.T) {
	r := setupRouter()

	body := `{"from":"txt","to":"json","content":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/convert/format", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for unsupported format, got %d: %s", w.Code, w.Body.String())
	}
}

func TestFormat_MissingBody(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/convert/format", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestFormat_MalformedInput(t *testing.T) {
	r := setupRouter()

	body := `{"from":"json","to":"yaml","content":"{invalid json"}`
	req := httptest.NewRequest(http.MethodPost, "/convert/format", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for malformed input, got %d: %s", w.Code, w.Body.String())
	}
}

func TestFormat_MissingFields(t *testing.T) {
	r := setupRouter()

	body := `{"from":"json","to":"yaml"}`
	req := httptest.NewRequest(http.MethodPost, "/convert/format", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for missing content, got %d: %s", w.Code, w.Body.String())
	}
}
