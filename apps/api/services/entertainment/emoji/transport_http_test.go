package emoji

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

func TestEmoji_Random(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/emoji/random", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Emoji]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	e := resp.Data
	if e.Emoji == "" {
		t.Error("expected non-empty emoji")
	}
	if e.Name == "" {
		t.Error("expected non-empty name")
	}
	if e.Category == "" {
		t.Error("expected non-empty category")
	}
	if e.Unicode == "" {
		t.Error("expected non-empty unicode")
	}
}

func TestEmoji_GetByName_Found(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/emoji/grinning_face", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Emoji]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	e := resp.Data
	if e.Name != "grinning_face" {
		t.Errorf("expected name 'grinning_face', got %q", e.Name)
	}
	if e.Emoji != "😀" {
		t.Errorf("expected emoji '😀', got %q", e.Emoji)
	}
	if e.Category != "Smileys & Emotion" {
		t.Errorf("expected category 'Smileys & Emotion', got %q", e.Category)
	}
	if e.Unicode != "U+1F600" {
		t.Errorf("expected unicode 'U+1F600', got %q", e.Unicode)
	}
}

func TestEmoji_GetByName_CaseInsensitive(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/emoji/GRINNING_FACE", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for uppercase name, got %d", w.Code)
	}

	var resp httpx.Response[Emoji]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Name != "grinning_face" {
		t.Errorf("expected name 'grinning_face', got %q", resp.Data.Name)
	}
}

func TestEmoji_GetByName_NotFound(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/emoji/does_not_exist", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestEmoji_Search_WithResults(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/emoji/search?q=happy", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[EmojiList]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Total < 0 {
		t.Error("expected non-negative total")
	}
}

func TestEmoji_Search_NoQuery(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/emoji/search", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for missing query, got %d", w.Code)
	}
}

func TestEmoji_Search_ReturnsMatchingEmojis(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/emoji/search?q=smile", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[EmojiList]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Total == 0 {
		t.Error("expected at least one result for 'smile'")
	}

	for _, e := range resp.Data.Items {
		if e.Emoji == "" || e.Name == "" || e.Category == "" || e.Unicode == "" {
			t.Errorf("expected all fields to be non-empty for emoji: %+v", e)
		}
	}
}

func TestEmoji_Search_EmptyResults(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/emoji/search?q=zzzyyyxxx", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[EmojiList]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.Total != 0 {
		t.Errorf("expected 0 results for nonsense query, got %d", resp.Data.Total)
	}

	if resp.Data.Items == nil {
		t.Error("expected non-nil items slice for empty results")
	}
}

func TestService_Random_ReturnsValidEmoji(t *testing.T) {
	svc := NewService()

	e := svc.Random()
	if e.Emoji == "" {
		t.Error("expected non-empty emoji from Random()")
	}
	if e.Name == "" {
		t.Error("expected non-empty name from Random()")
	}
}

func TestService_GetByName(t *testing.T) {
	svc := NewService()

	tests := []struct {
		name      string
		input     string
		wantFound bool
	}{
		{name: "valid name", input: "grinning_face", wantFound: true},
		{name: "valid name uppercase", input: "GRINNING_FACE", wantFound: true},
		{name: "unknown name", input: "not_a_real_emoji", wantFound: false},
		{name: "empty string", input: "", wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, found := svc.GetByName(tt.input)
			if found != tt.wantFound {
				t.Errorf("GetByName(%q): got found=%v, want %v", tt.input, found, tt.wantFound)
			}
		})
	}
}

func TestService_Search(t *testing.T) {
	svc := NewService()

	tests := []struct {
		name      string
		query     string
		wantEmpty bool
	}{
		{name: "smile query returns results", query: "smile", wantEmpty: false},
		{name: "heart query returns results", query: "heart", wantEmpty: false},
		{name: "nonsense returns empty", query: "zzzyyyxxx", wantEmpty: true},
		{name: "category search", query: "food", wantEmpty: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.Search(tt.query)
			isEmpty := result.Total == 0
			if isEmpty != tt.wantEmpty {
				t.Errorf("Search(%q): got empty=%v, want empty=%v (total=%d)", tt.query, isEmpty, tt.wantEmpty, result.Total)
			}
			if result.Total != len(result.Items) {
				t.Errorf("Search(%q): Total=%d does not match len(Items)=%d", tt.query, result.Total, len(result.Items))
			}
		})
	}
}
