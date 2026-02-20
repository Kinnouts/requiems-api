package httpx

import (
	"encoding/json"
	"net/http"
	"time"
)

// Data is a marker interface for types that can be used as API response payloads.
// Add IsData() to your response struct to use it with httpx.JSON.
type Data interface {
	IsData()
}

// Metadata is included in every response.
type Metadata struct {
	Timestamp string `json:"timestamp"`
	TraceID   string `json:"trace_id,omitempty"`
}

// Response is the standard success envelope: {"data": ..., "metadata": ...}
type Response[T Data] struct {
	Data     T        `json:"data"`
	Metadata Metadata `json:"metadata"`
}

// ErrorResponse is the standard error envelope: {"error": ..., "metadata": ...}
type ErrorResponse struct {
	Error    string   `json:"error"`
	Metadata Metadata `json:"metadata"`
}

func JSON[T Data](w http.ResponseWriter, status int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response[T]{
		Data:     v,
		Metadata: Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func Error(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:    msg,
		Metadata: Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}
