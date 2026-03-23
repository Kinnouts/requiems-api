package chucknorris

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

func TestChuckNorris_Random(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/chuck-norris", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Fact]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	f := resp.Data

	if f.Fact == "" {
		t.Error("expected non-empty fact")
	}

	if !strings.HasPrefix(f.ID, "cn_") {
		t.Errorf("expected ID to start with 'cn_', got %q", f.ID)
	}
}

func TestChuckNorris_Randomness(t *testing.T) {
	svc := NewService()

	seen := make(map[string]bool)
	for range 50 {
		f := svc.Random()
		seen[f.ID] = true
	}

	// With 30 facts and 50 draws, expect at least 5 distinct facts.
	if len(seen) < 5 {
		t.Errorf("expected variety in random facts, got only %d distinct IDs in 50 draws", len(seen))
	}
}

func TestChuckNorris_FactsNonEmpty(t *testing.T) {
	svc := NewService()
	for i := range 10 {
		f := svc.Random()
		if f.Fact == "" {
			t.Errorf("call %d: expected non-empty fact", i)
		}
		if f.ID == "" {
			t.Errorf("call %d: expected non-empty ID", i)
		}
	}
}
