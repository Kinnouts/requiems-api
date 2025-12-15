package words

import (
	"encoding/json"
	"net/http"
)

func RegisterHTTP(mux *http.ServeMux, svc *Service) {
	mux.Handle("/v1/words/random", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		wrd, err := svc.Random(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "no words available"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(wrd)
	}))
}


