package spellcheck

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

func TestSpellcheck_CleanText(t *testing.T) {
	r := setupRouter()

	body := `{"text":"Hello world"}`
	req := httptest.NewRequest(http.MethodPost, "/spellcheck", strings.NewReader(body))
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

	if len(resp.Data.Corrections) != 0 {
		t.Errorf("expected no corrections, got %v", resp.Data.Corrections)
	}
}

func TestSpellcheck_MisspelledText(t *testing.T) {
	r := setupRouter()

	body := `{"text":"Ths is a tset"}`
	req := httptest.NewRequest(http.MethodPost, "/spellcheck", strings.NewReader(body))
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

	if len(resp.Data.Corrections) == 0 {
		t.Error("expected at least one correction for misspelled input")
	}
	if resp.Data.Corrected == "Ths is a tset" {
		t.Error("expected corrected text to differ from input")
	}
}

func TestSpellcheck_MissingTextField(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/spellcheck", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSpellcheck_MissingBody(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/spellcheck", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
