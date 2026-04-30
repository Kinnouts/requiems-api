package similarity

// Request is the input for the text similarity endpoint.
type Request struct {
	Text1 string `json:"text1" validate:"required"`
	Text2 string `json:"text2" validate:"required"`
}

// Result is the response payload for the text similarity endpoint.
type Result struct {
	Similarity float64 `json:"similarity"`
	Method     string  `json:"method"`
}

func (Result) IsData() {}
