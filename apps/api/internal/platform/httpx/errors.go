package httpx

// FieldError describes a single field-level validation failure.
type FieldError struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
}

// ValidationFailure is returned by BindAndValidate when struct validation fails.
// Handlers can check for this type to distinguish validation errors from other
// decode failures.
type ValidationFailure struct {
	Fields []FieldError
}

func (e *ValidationFailure) Error() string { return "validation_failed" }

// AppError is a typed error the service layer can return to signal a specific
// HTTP status code and error code to the client. The Handle wrapper maps it
// automatically.
//
// Example:
//
//	return nil, &httpx.AppError{Status: 404, Code: "not_found", Message: "resource not found"}
type AppError struct {
	Code    string
	Message string
	Status  int
}

func (e *AppError) Error() string { return e.Message }
