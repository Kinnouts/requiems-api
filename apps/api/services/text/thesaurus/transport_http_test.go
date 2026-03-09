package thesaurus

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	RegisterRoutes(r, NewService())
	return r
}

func TestThesaurus_KnownWord(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/thesaurus/happy", http.NoBody)
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Result]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Word != "happy" {
		t.Errorf("expected word %q, got %q", "happy", resp.Data.Word)
	}

	if len(resp.Data.Synonyms) == 0 {
		t.Error("expected at least one synonym")
	}

	if len(resp.Data.Antonyms) == 0 {
		t.Error("expected at least one antonym")
	}
}

func TestThesaurus_CaseInsensitive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/thesaurus/HAPPY", http.NoBody)
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Result]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Word != "happy" {
		t.Errorf("expected normalized word %q, got %q", "happy", resp.Data.Word)
	}
}

func TestThesaurus_UnknownWord(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/thesaurus/zzyzx", http.NoBody)
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
