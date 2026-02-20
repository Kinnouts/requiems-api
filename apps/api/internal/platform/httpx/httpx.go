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

// ErrorResponse is the standard error envelope.
// Fields is populated only for validation errors (error: "validation_failed").
type ErrorResponse struct {
	Error    string       `json:"error"`
	Message  string       `json:"message,omitempty"`
	Fields   []FieldError `json:"fields,omitempty"`
	Metadata Metadata     `json:"metadata"`
}

// JSON writes a 200-class success response wrapped in {"data": ..., "metadata": ...}.
func JSON[T Data](w http.ResponseWriter, status int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response[T]{
		Data:     v,
		Metadata: Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

// Error writes a JSON error response with a machine-readable code and a
// human-readable message.
//
//	httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid email format")
func Error(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:    code,
		Message:  message,
		Metadata: Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

// writeValidationError writes a 422 Unprocessable Entity with a structured
// list of field-level constraint violations.
func writeValidationError(w http.ResponseWriter, fields []FieldError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:    "validation_failed",
		Fields:   fields,
		Metadata: Metadata{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}
