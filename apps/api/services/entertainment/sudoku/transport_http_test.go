package sudoku

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

func TestSudoku_DefaultDifficulty(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/sudoku", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Puzzle]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Difficulty != "medium" {
		t.Errorf("expected default difficulty 'medium', got %q", resp.Data.Difficulty)
	}
}

func TestSudoku_AllDifficulties(t *testing.T) {
	r := setupRouter()

	difficulties := []string{"easy", "medium", "hard"}

	for _, d := range difficulties {
		t.Run(d, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/sudoku?difficulty="+d, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200 for difficulty %q, got %d", d, w.Code)
			}

			var resp httpx.Response[Puzzle]
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Data.Difficulty != d {
				t.Errorf("expected difficulty %q, got %q", d, resp.Data.Difficulty)
			}
		})
	}
}

func TestSudoku_InvalidDifficulty(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/sudoku?difficulty=impossible", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSudoku_PuzzleHasEmptyCells(t *testing.T) {
	svc := NewService()

	for _, d := range []string{"easy", "medium", "hard"} {
		t.Run(d, func(t *testing.T) {
			p := svc.Generate(d)

			empty := 0
			for r := range 9 {
				for c := range 9 {
					if p.Puzzle[r][c] == 0 {
						empty++
					}
				}
			}

			expected := cellsToRemove[d]
			if empty != expected {
				t.Errorf("difficulty %q: expected %d empty cells, got %d", d, expected, empty)
			}
		})
	}
}

func TestSudoku_SolutionIsComplete(t *testing.T) {
	svc := NewService()
	p := svc.Generate("hard")

	for r := range 9 {
		for c := range 9 {
			if p.Solution[r][c] < 1 || p.Solution[r][c] > 9 {
				t.Errorf("solution[%d][%d] = %d, want 1-9", r, c, p.Solution[r][c])
			}
		}
	}
}

func TestSudoku_SolutionIsValid(t *testing.T) {
	svc := NewService()
	p := svc.Generate("medium")

	// Check each row contains 1-9.
	for r := range 9 {
		if !hasAllDigits(p.Solution[r][:]) {
			t.Errorf("row %d does not contain all digits 1-9", r)
		}
	}

	// Check each column contains 1-9.
	for c := range 9 {
		col := make([]int, 9)
		for r := range 9 {
			col[r] = p.Solution[r][c]
		}
		if !hasAllDigits(col) {
			t.Errorf("column %d does not contain all digits 1-9", c)
		}
	}

	// Check each 3×3 box contains 1-9.
	for br := range 3 {
		for bc := range 3 {
			box := make([]int, 0, 9)
			for r := range 3 {
				for c := range 3 {
					box = append(box, p.Solution[br*3+r][bc*3+c])
				}
			}
			if !hasAllDigits(box) {
				t.Errorf("box [%d,%d] does not contain all digits 1-9", br, bc)
			}
		}
	}
}

func TestSudoku_PuzzleMatchesSolution(t *testing.T) {
	svc := NewService()
	p := svc.Generate("easy")

	for r := range 9 {
		for c := range 9 {
			if p.Puzzle[r][c] != 0 && p.Puzzle[r][c] != p.Solution[r][c] {
				t.Errorf("puzzle[%d][%d]=%d differs from solution[%d][%d]=%d",
					r, c, p.Puzzle[r][c], r, c, p.Solution[r][c])
			}
		}
	}
}

// hasAllDigits returns true when values contains each of 1-9 exactly once.
func hasAllDigits(values []int) bool {
	seen := make(map[int]bool, 9)
	for _, v := range values {
		if v < 1 || v > 9 || seen[v] {
			return false
		}
		seen[v] = true
	}
	return len(seen) == 9
}
