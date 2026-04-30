package facts

// Request holds optional query parameters for the facts endpoint.
type Request struct {
	Category string `query:"category"`
}

// Fact is the response payload for a single random fact.
type Fact struct {
	Fact     string `json:"fact"`
	Category string `json:"category"`
	Source   string `json:"source"`
}

func (Fact) IsData() {}
