package jokes

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
	svc := NewService()
	RegisterRoutes(r, svc)
	return r
}

func TestDadJoke_Random(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/jokes/dad", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[DadJoke]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	j := resp.Data
	if j.Joke == "" {
		t.Error("expected non-empty joke")
	}
	if j.ID == "" {
		t.Error("expected non-empty id")
	}
	if !strings.HasPrefix(j.ID, "joke_") {
		t.Errorf("expected id to start with 'joke_', got %q", j.ID)
	}
}

func TestDadJoke_Random_MultipleCallsReturnValidJokes(t *testing.T) {
	svc := NewService()

	for range 10 {
		j := svc.Random()
		if j.Joke == "" {
			t.Error("expected non-empty joke")
		}
		if j.ID == "" {
			t.Error("expected non-empty id")
		}
		if !strings.HasPrefix(j.ID, "joke_") {
			t.Errorf("expected id to start with 'joke_', got %q", j.ID)
		}
	}
}
