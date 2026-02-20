package horoscope

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"requiems-api/internal/platform/httpx"
)

func setupRouter() chi.Router {
	r := chi.NewRouter()
	svc := NewService()
	RegisterRoutes(r, svc)
	return r
}

func TestHoroscope_ValidSign(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/horoscope/aries", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp httpx.Response[Horoscope]
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	h := resp.Data

	if h.Sign != "aries" {
		t.Errorf("expected sign 'aries', got %q", h.Sign)
	}

	if h.Date != time.Now().UTC().Format("2006-01-02") {
		t.Errorf("expected today's date, got %q", h.Date)
	}

	if h.Horoscope == "" {
		t.Error("expected non-empty horoscope text")
	}

	if h.LuckyNumber < 1 || h.LuckyNumber > 99 {
		t.Errorf("expected lucky_number between 1 and 99, got %d", h.LuckyNumber)
	}

	if h.Mood == "" {
		t.Error("expected non-empty mood")
	}
}

func TestHoroscope_AllSigns(t *testing.T) {
	r := setupRouter()

	validSigns := []string{
		"aries", "taurus", "gemini", "cancer", "leo", "virgo",
		"libra", "scorpio", "sagittarius", "capricorn", "aquarius", "pisces",
	}

	for _, sign := range validSigns {
		t.Run(sign, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/horoscope/"+sign, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200 for sign %q, got %d", sign, w.Code)
			}
		})
	}
}

func TestHoroscope_InvalidSign(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/horoscope/invalid", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHoroscope_CaseInsensitive(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/horoscope/ARIES", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for uppercase sign, got %d", w.Code)
	}

	var resp2 httpx.Response[Horoscope]
	if err := json.NewDecoder(w.Body).Decode(&resp2); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp2.Data.Sign != "aries" {
		t.Errorf("expected sign normalized to 'aries', got %q", resp2.Data.Sign)
	}
}

func TestHoroscope_DailyConsistency(t *testing.T) {
	svc := NewService()

	h1, err := svc.Daily("leo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	h2, err := svc.Daily("leo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h1.Horoscope != h2.Horoscope {
		t.Error("expected same horoscope for same sign on same day")
	}

	if h1.LuckyNumber != h2.LuckyNumber {
		t.Error("expected same lucky_number for same sign on same day")
	}

	if h1.Mood != h2.Mood {
		t.Error("expected same mood for same sign on same day")
	}
}
