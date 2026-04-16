package similarity

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

func TestSimilarity_IdenticalTexts(t *testing.T) {
	r := setupRouter()

	body := `{"text1":"The cat sat on the mat","text2":"The cat sat on the mat"}`
	req := httptest.NewRequest(http.MethodPost, "/similarity", strings.NewReader(body))
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

	if resp.Data.Similarity != 1.0 {
		t.Errorf("identical texts should have similarity 1.0, got %f", resp.Data.Similarity)
	}
	if resp.Data.Method != "cosine" {
		t.Errorf("expected method 'cosine', got %q", resp.Data.Method)
	}
}

func TestSimilarity_UnrelatedTexts(t *testing.T) {
	r := setupRouter()

	body := `{"text1":"The cat sat on the mat","text2":"quantum physics nuclear reactor"}`
	req := httptest.NewRequest(http.MethodPost, "/similarity", strings.NewReader(body))
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

	if resp.Data.Similarity != 0.0 {
		t.Errorf("unrelated texts should have similarity 0.0, got %f", resp.Data.Similarity)
	}
}

func TestSimilarity_SimilarTexts(t *testing.T) {
	r := setupRouter()

	body := `{"text1":"The cat sat on the mat","text2":"A cat was sitting on a mat"}`
	req := httptest.NewRequest(http.MethodPost, "/similarity", strings.NewReader(body))
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

	// Texts share words (cat, on, mat), expect non-zero similarity.
	if resp.Data.Similarity <= 0 || resp.Data.Similarity >= 1 {
		t.Errorf("expected similarity between 0 and 1 (exclusive), got %f", resp.Data.Similarity)
	}
}

func TestSimilarity_MissingText1(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/similarity", strings.NewReader(`{"text2":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSimilarity_MissingText2(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/similarity", strings.NewReader(`{"text1":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSimilarity_MissingBody(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/similarity", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
