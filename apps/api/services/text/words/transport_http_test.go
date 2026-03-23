package words

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
	RegisterRoutes(r, &Service{db: &mockQuerier{}})
	return r
}

func TestDictionary_KnownWord(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dictionary/ephemeral", http.NoBody)
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[DictionaryEntry]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Word != "ephemeral" {
		t.Errorf("expected word %q, got %q", "ephemeral", resp.Data.Word)
	}

	if resp.Data.Phonetic == "" {
		t.Error("expected non-empty phonetic")
	}

	if len(resp.Data.Definitions) == 0 {
		t.Error("expected at least one definition")
	}

	if resp.Data.Definitions[0].PartOfSpeech == "" {
		t.Error("expected non-empty partOfSpeech")
	}

	if resp.Data.Definitions[0].Definition == "" {
		t.Error("expected non-empty definition text")
	}

	if len(resp.Data.Synonyms) == 0 {
		t.Error("expected at least one synonym")
	}
}

func TestDictionary_CaseInsensitive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dictionary/EPHEMERAL", http.NoBody)
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[DictionaryEntry]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Word != "ephemeral" {
		t.Errorf("expected normalized word %q, got %q", "ephemeral", resp.Data.Word)
	}
}

func TestDictionary_UnknownWord(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dictionary/zzyzx", http.NoBody)
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDictionary_ExampleField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dictionary/ephemeral", http.NoBody)
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[DictionaryEntry]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Definitions[0].Example == "" {
		t.Error("expected non-empty example")
	}
}

func TestDictionary_MultipleDefinitions(t *testing.T) {
	// melancholy has two definitions (noun and adjective)
	req := httptest.NewRequest(http.MethodGet, "/dictionary/melancholy", http.NoBody)
	w := httptest.NewRecorder()

	setupRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[DictionaryEntry]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Data.Definitions) < 2 {
		t.Errorf("expected at least 2 definitions, got %d", len(resp.Data.Definitions))
	}
}
