package trivia

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

func TestTrivia_NoFilters(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/trivia", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Question]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Question == "" {
		t.Error("expected non-empty question")
	}
	if len(resp.Data.Options) == 0 {
		t.Error("expected at least one option")
	}
	if resp.Data.Answer == "" {
		t.Error("expected non-empty answer")
	}
	if resp.Data.Category == "" {
		t.Error("expected non-empty category")
	}
	if resp.Data.Difficulty == "" {
		t.Error("expected non-empty difficulty")
	}
}

func TestTrivia_FilterByCategory(t *testing.T) {
	r := setupRouter()

	categories := []string{"science", "history", "geography", "sports", "music", "movies", "literature", "math", "technology", "nature"}
	for _, cat := range categories {
		t.Run(cat, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/trivia?category="+cat, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200 for category %q, got %d", cat, w.Code)
			}

			var resp httpx.Response[Question]
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Data.Category != cat {
				t.Errorf("expected category %q, got %q", cat, resp.Data.Category)
			}
		})
	}
}

func TestTrivia_FilterByDifficulty(t *testing.T) {
	r := setupRouter()

	difficulties := []string{"easy", "medium", "hard"}
	for _, d := range difficulties {
		t.Run(d, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/trivia?difficulty="+d, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200 for difficulty %q, got %d", d, w.Code)
			}

			var resp httpx.Response[Question]
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Data.Difficulty != d {
				t.Errorf("expected difficulty %q, got %q", d, resp.Data.Difficulty)
			}
		})
	}
}

func TestTrivia_FilterByCategoryAndDifficulty(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/trivia?category=science&difficulty=easy", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Question]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Category != "science" {
		t.Errorf("expected category 'science', got %q", resp.Data.Category)
	}
	if resp.Data.Difficulty != "easy" {
		t.Errorf("expected difficulty 'easy', got %q", resp.Data.Difficulty)
	}
}

func TestTrivia_InvalidCategory(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/trivia?category=invalid", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestTrivia_InvalidDifficulty(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/trivia?difficulty=impossible", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestTrivia_AnswerIsInOptions(t *testing.T) {
	for _, q := range questions {
		found := false
		for _, opt := range q.Options {
			if opt == q.Answer {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("question %q: answer %q is not in options %v", q.Question, q.Answer, q.Options)
		}
	}
}

func TestTrivia_AllQuestionsHaveFourOptions(t *testing.T) {
	for _, q := range questions {
		if len(q.Options) != 4 {
			t.Errorf("question %q: expected 4 options, got %d", q.Question, len(q.Options))
		}
	}
}

func TestService_Random_NoMatch(t *testing.T) {
	svc := NewService()
	_, err := svc.Random("science", "impossible")
	if err == nil {
		t.Error("expected error for filters with no matching questions, got nil")
	}
}
