package detectlanguage

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

func TestDetectLanguage_French(t *testing.T) {
	r := setupRouter()

	body := `{"text":"Bonjour, comment ça va?"}`
	req := httptest.NewRequest(http.MethodPost, "/detect-language", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Result]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Language != "French" {
		t.Errorf("expected language French, got %q", resp.Data.Language)
	}
	if resp.Data.Code != "fr" {
		t.Errorf("expected code fr, got %q", resp.Data.Code)
	}
	if resp.Data.Confidence <= 0 {
		t.Error("expected confidence to be greater than 0")
	}
}

func TestDetectLanguage_English(t *testing.T) {
	r := setupRouter()

	body := `{"text":"The quick brown fox jumps over the lazy dog"}`
	req := httptest.NewRequest(http.MethodPost, "/detect-language", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Result]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Language != "English" {
		t.Errorf("expected language English, got %q", resp.Data.Language)
	}
	if resp.Data.Code != "en" {
		t.Errorf("expected code en, got %q", resp.Data.Code)
	}
}

func TestDetectLanguage_MissingTextField(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/detect-language", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDetectLanguage_MissingBody(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/detect-language", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
