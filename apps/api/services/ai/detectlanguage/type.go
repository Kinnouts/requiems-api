package detectlanguage

// Request is the input for the detect-language endpoint.
type Request struct {
	Text string `json:"text" validate:"required"`
}

// Result is the response payload for the detect-language endpoint.
type Result struct {
	Language   string  `json:"language"`
	Code       string  `json:"code"`
	Confidence float64 `json:"confidence"`
}

func (Result) IsData() {}
