package spellcheck

// Request is the input for the spell check endpoint.
type Request struct {
	Text string `json:"text" validate:"required"`
}

// Correction describes a single spelling mistake and its suggested fix.
type Correction struct {
	Original  string `json:"original"`
	Suggested string `json:"suggested"`
	Position  int    `json:"position"`
}

// Result is the response payload for the spell check endpoint.
type Result struct {
	Corrected   string       `json:"corrected"`
	Corrections []Correction `json:"corrections"`
}

func (Result) IsData() {}
