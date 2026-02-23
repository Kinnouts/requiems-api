package profanity

// Request is the input for the profanity check endpoint.
type Request struct {
	Text string `json:"text" validate:"required"`
}

// Result is the response payload for the profanity check endpoint.
type Result struct {
	HasProfanity bool     `json:"hasProfanity"`
	Censored     string   `json:"censored"`
	FlaggedWords []string `json:"flaggedWords"`
}

func (Result) IsData() {}
