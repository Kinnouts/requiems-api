package profanity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	RegisterRoutes(r, NewService())
	return r
}

func TestProfanity_CleanText(t *testing.T) {
	r := setupRouter()

	body := `{"text":"Hello, world!"}`
	req := httptest.NewRequest(http.MethodPost, "/profanity", strings.NewReader(body))
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

	if resp.Data.HasProfanity {
		t.Error("expected HasProfanity to be false")
	}
	if resp.Data.Censored != "Hello, world!" {
		t.Errorf("expected censored to equal input, got %q", resp.Data.Censored)
	}
	if len(resp.Data.FlaggedWords) != 0 {
		t.Errorf("expected no flagged words, got %v", resp.Data.FlaggedWords)
	}
}

func TestProfanity_ProfaneText(t *testing.T) {
	r := setupRouter()

	body := `{"text":"That is such bullshit"}`
	req := httptest.NewRequest(http.MethodPost, "/profanity", strings.NewReader(body))
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

	if !resp.Data.HasProfanity {
		t.Error("expected HasProfanity to be true")
	}
	if resp.Data.Censored != "That is such ********" {
		t.Errorf("unexpected censored output: %q", resp.Data.Censored)
	}
	if len(resp.Data.FlaggedWords) != 1 || resp.Data.FlaggedWords[0] != "bullshit" {
		t.Errorf("expected flagged word 'bullshit', got %v", resp.Data.FlaggedWords)
	}
}

func TestProfanity_MissingTextField(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/profanity", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProfanity_MissingBody(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/profanity", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
