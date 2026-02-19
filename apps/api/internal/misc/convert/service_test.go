package convert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
)

func newRouter() chi.Router {
	r := chi.NewRouter()
	svc := NewService()
	RegisterRoutes(r, svc)
	return r
}

func TestConvert_Service(t *testing.T) {
	svc := NewService()

	tests := []struct {
		name    string
		from    string
		to      string
		value   float64
		want    float64
		formula string
		wantErr bool
	}{
		// Length
		{"miles to km", "miles", "km", 10, 16.09344, "miles × 1.609344", false},
		{"km to miles", "km", "miles", 10, 6.2137119224, "km × 0.621371", false},
		{"m to cm", "m", "cm", 1, 100, "m × 100", false},
		{"ft to m", "ft", "m", 1, 0.3048, "ft × 0.3048", false},
		{"inch to cm", "inch", "cm", 1, 2.54, "inch × 2.54", false},

		// Weight
		{"kg to lb", "kg", "lb", 1, 2.2046226218, "kg × 2.204623", false},
		{"lb to kg", "lb", "kg", 1, 0.45359237, "lb × 0.453592", false},
		{"g to oz", "g", "oz", 28.349523, 1, "g × 0.035274", false},

		// Volume
		{"l to gallon", "l", "gallon", 1, 0.2641720512, "l × 0.264172", false},
		{"gallon to l", "gallon", "l", 1, 3.7854118, "gallon × 3.785412", false},

		// Temperature
		{"celsius to fahrenheit", "celsius", "fahrenheit", 100, 212, "°C × 9/5 + 32", false},
		{"fahrenheit to celsius", "fahrenheit", "celsius", 32, 0, "(°F − 32) × 5/9", false},
		{"celsius to kelvin", "celsius", "kelvin", 0, 273.15, "°C + 273.15", false},
		{"kelvin to celsius", "kelvin", "celsius", 273.15, 0, "°K − 273.15", false},

		// Area
		{"hectare to acre", "hectare", "acre", 1, 2.4710538283, "hectare × 2.471054", false},

		// Speed
		{"mph to km/h", "mph", "km/h", 60, 96.560640, "mph × 1.609344", false},

		// Data
		{"gb to mb", "gb", "mb", 1, 1000, "gb × 1000", false},

		// Time
		{"hour to min", "hour", "min", 2, 120, "hour × 60", false},
		{"day to hours", "day", "hour", 1, 24, "day × 24", false},

		// Errors
		{"unknown from", "lightyear", "km", 1, 0, "", true},
		{"unknown to", "km", "lightyear", 1, 0, "", true},
		{"incompatible units", "km", "kg", 1, 0, "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := svc.Convert(tc.from, tc.to, tc.value)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Result != tc.want {
				t.Errorf("result: got %v, want %v", resp.Result, tc.want)
			}
			if resp.Formula != tc.formula {
				t.Errorf("formula: got %q, want %q", resp.Formula, tc.formula)
			}
		})
	}
}

func TestConvert_HTTP(t *testing.T) {
	r := newRouter()

	t.Run("valid conversion returns 200", func(t *testing.T) {
		params := url.Values{"from": {"miles"}, "to": {"km"}, "value": {"10"}}
		req := httptest.NewRequest(http.MethodGet, "/convert?"+params.Encode(), http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		var resp Response
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if resp.From != "miles" || resp.To != "km" || resp.Input != 10 {
			t.Errorf("unexpected response fields: %+v", resp)
		}
		if resp.Result != 16.09344 {
			t.Errorf("result: got %v, want 16.09344", resp.Result)
		}
	})

	t.Run("missing parameters returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/convert?from=miles&to=km", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("invalid value returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/convert?from=miles&to=km&value=abc", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("unknown unit returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/convert?from=lightyear&to=km&value=1", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("incompatible units returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/convert?from=km&to=kg&value=1", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}
