package facts

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
	svc := NewService()
	RegisterRoutes(r, svc)
	return r
}

func TestFacts_RandomFact(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/facts", http.NoBody)
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
	if f.Category == "" {
		t.Error("expected non-empty category")
	}
	if f.Source == "" {
		t.Error("expected non-empty source")
	}
}

func TestFacts_CategoryFilter(t *testing.T) {
	r := setupRouter()

	categories := []string{"science", "history", "technology", "nature", "space", "food"}
	for _, cat := range categories {
		t.Run(cat, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/facts?category="+cat, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200 for category %q, got %d", cat, w.Code)
			}

			var resp httpx.Response[Fact]
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Data.Category != cat {
				t.Errorf("expected category %q, got %q", cat, resp.Data.Category)
			}
		})
	}
}

func TestFacts_InvalidCategory(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/facts?category=invalid", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestFacts_CategoryCaseInsensitive(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/facts?category=SCIENCE", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for uppercase category, got %d", w.Code)
	}

	var resp httpx.Response[Fact]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Category != "science" {
		t.Errorf("expected category 'science', got %q", resp.Data.Category)
	}
}

func TestFacts_ServiceRandom(t *testing.T) {
	svc := NewService()

	f, err := svc.Random("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Fact == "" {
		t.Error("expected non-empty fact")
	}
}

func TestFacts_ServiceRandomByCategory(t *testing.T) {
	svc := NewService()

	f, err := svc.Random("science")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Category != "science" {
		t.Errorf("expected category 'science', got %q", f.Category)
	}
}
