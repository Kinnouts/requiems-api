package sentiment

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

func TestSentimentHandler_Positive(t *testing.T) {
	r := setupRouter()

	body := `{"text":"I love this product! It's amazing."}`
	req := httptest.NewRequest(http.MethodPost, "/sentiment", strings.NewReader(body))
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

	if resp.Data.Sentiment != "positive" {
		t.Errorf("expected sentiment positive, got %q", resp.Data.Sentiment)
	}
	if resp.Data.Score <= 0 || resp.Data.Score > 1 {
		t.Errorf("score out of range [0,1]: %.2f", resp.Data.Score)
	}
}

func TestSentimentHandler_Negative(t *testing.T) {
	r := setupRouter()

	body := `{"text":"This is terrible and I hate it."}`
	req := httptest.NewRequest(http.MethodPost, "/sentiment", strings.NewReader(body))
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

	if resp.Data.Sentiment != "negative" {
		t.Errorf("expected sentiment negative, got %q", resp.Data.Sentiment)
	}
}

func TestSentimentHandler_MissingTextField(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/sentiment", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSentimentHandler_MissingBody(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/sentiment", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
