package exercises

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"requiems-api/platform/httpx"
)

// stubQuerier implements exerciseQuerier for HTTP handler tests.
type stubQuerier struct {
	listResult   ExerciseList
	getResult    Exercise
	randomResult Exercise
	stringResult StringList
	err          error
}

func (s *stubQuerier) List(_ context.Context, _ ListParams) (ExerciseList, error) {
	return s.listResult, s.err
}

func (s *stubQuerier) Get(_ context.Context, _ int) (Exercise, error) {
	if s.err != nil {
		return Exercise{}, s.err
	}
	return s.getResult, nil
}

func (s *stubQuerier) Random(_ context.Context, _ ListParams) (Exercise, error) {
	if s.err != nil {
		return Exercise{}, s.err
	}
	return s.randomResult, nil
}

func (s *stubQuerier) BodyParts(_ context.Context) (StringList, error) {
	return s.stringResult, s.err
}

func (s *stubQuerier) Equipment(_ context.Context) (StringList, error) {
	return s.stringResult, s.err
}

func (s *stubQuerier) Muscles(_ context.Context) (StringList, error) {
	return s.stringResult, s.err
}

func setupTestRouter(q exerciseQuerier) chi.Router {
	r := chi.NewRouter()
	registerExerciseRoutes(r, q)
	return r
}

// ---- GET /exercises ----

func TestListExercises_Returns200(t *testing.T) {
	stub := &stubQuerier{
		listResult: ExerciseList{
			Items:   []Exercise{{ID: 1, Name: "squat"}},
			Total:   1,
			Page:    1,
			PerPage: 20,
		},
	}

	r := setupTestRouter(stub)
	req := httptest.NewRequest(http.MethodGet, "/exercises", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[ExerciseList]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Data.Total != 1 {
		t.Errorf("expected total 1, got %d", resp.Data.Total)
	}
	if len(resp.Data.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Data.Items))
	}
}

func TestListExercises_ResponseEnvelope(t *testing.T) {
	stub := &stubQuerier{listResult: ExerciseList{Items: []Exercise{}}}
	r := setupTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/exercises", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if _, ok := raw["data"]; !ok {
		t.Error("response must have a 'data' key")
	}
	if _, ok := raw["metadata"]; !ok {
		t.Error("response must have a 'metadata' key")
	}
}

func TestListExercises_ServiceError_Returns500(t *testing.T) {
	stub := &stubQuerier{err: &httpx.AppError{
		Status: http.StatusInternalServerError, Code: "internal_error", Message: "db down",
	}}

	r := setupTestRouter(stub)
	req := httptest.NewRequest(http.MethodGet, "/exercises", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestListExercises_InvalidPerPage_Returns400(t *testing.T) {
	stub := &stubQuerier{}
	r := setupTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/exercises?per_page=0", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// ---- GET /exercises/random ----

func TestRandomExercise_Returns200(t *testing.T) {
	stub := &stubQuerier{randomResult: Exercise{ID: 7, Name: "deadlift"}}
	r := setupTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/exercises/random", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Exercise]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Data.Name != "deadlift" {
		t.Errorf("expected 'deadlift', got %q", resp.Data.Name)
	}
}

func TestRandomExercise_NoMatch_Returns404(t *testing.T) {
	stub := &stubQuerier{err: &httpx.AppError{
		Status: http.StatusNotFound, Code: "not_found", Message: "no exercises found",
	}}

	r := setupTestRouter(stub)
	req := httptest.NewRequest(http.MethodGet, "/exercises/random?body_part=unknown", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// ---- GET /exercises/{id} ----

func TestGetExercise_ValidID_Returns200(t *testing.T) {
	stub := &stubQuerier{getResult: Exercise{ID: 3, Name: "bench press"}}
	r := setupTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/exercises/3", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[Exercise]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Data.Name != "bench press" {
		t.Errorf("expected 'bench press', got %q", resp.Data.Name)
	}
}

func TestGetExercise_NonNumericID_Returns400(t *testing.T) {
	stub := &stubQuerier{}
	r := setupTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/exercises/abc", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetExercise_ZeroID_Returns400(t *testing.T) {
	stub := &stubQuerier{}
	r := setupTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/exercises/0", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetExercise_NotFound_Returns404(t *testing.T) {
	stub := &stubQuerier{err: &httpx.AppError{
		Status: http.StatusNotFound, Code: "not_found", Message: "exercise not found",
	}}

	r := setupTestRouter(stub)
	req := httptest.NewRequest(http.MethodGet, "/exercises/9999", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// ---- GET /body-parts, /equipment, /muscles ----

func TestBodyParts_Returns200(t *testing.T) {
	stub := &stubQuerier{stringResult: StringList{Items: []string{"chest", "back"}, Total: 2}}
	r := setupTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/body-parts", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp httpx.Response[StringList]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Data.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Data.Total)
	}
}

func TestEquipment_Returns200(t *testing.T) {
	stub := &stubQuerier{stringResult: StringList{Items: []string{"barbell", "dumbbell"}, Total: 2}}
	r := setupTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/equipment", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMuscles_Returns200(t *testing.T) {
	stub := &stubQuerier{stringResult: StringList{Items: []string{"biceps", "triceps"}, Total: 2}}
	r := setupTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/muscles", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMetadata_ServiceError_Returns500(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"body-parts", "/body-parts"},
		{"equipment", "/equipment"},
		{"muscles", "/muscles"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &stubQuerier{err: &httpx.AppError{
				Status: http.StatusInternalServerError, Code: "internal_error", Message: "db error",
			}}

			r := setupTestRouter(stub)
			req := httptest.NewRequest(http.MethodGet, tt.path, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusInternalServerError {
				t.Errorf("expected 500, got %d", w.Code)
			}
		})
	}
}
