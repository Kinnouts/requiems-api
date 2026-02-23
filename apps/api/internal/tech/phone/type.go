package phone

// ValidateRequest holds query parameters for the phone validation endpoint.
type ValidateRequest struct {
	Number string `query:"number" validate:"required"`
}

// ValidateResponse is the response for a phone number validation request.
type ValidateResponse struct {
	Number    string `json:"number"`
	Valid     bool   `json:"valid"`
	Country   string `json:"country,omitempty"`
	Type      string `json:"type,omitempty"`
	Formatted string `json:"formatted,omitempty"`
}

func (ValidateResponse) IsData() {}
