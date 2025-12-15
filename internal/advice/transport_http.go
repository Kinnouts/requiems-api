package advice

import (
	"encoding/json"
	"net/http"
)

// RegisterHTTP mounts advice handlers on the given mux.
func RegisterHTTP(mux *http.ServeMux, svc *Service) {
	mux.Handle("/v1/advice", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		a, err := svc.Random(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "no advice available"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(a)
	}))
}


