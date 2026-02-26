package lorem

// Request holds the optional query parameters for the lorem endpoint.
// Defaults should be set before calling httpx.BindQuery.
type Request struct {
	Paragraphs int `query:"paragraphs" validate:"min=1,max=20"`
	Sentences  int `query:"sentences"  validate:"min=1,max=20"`
}

// Lorem is the response payload for the lorem generator.
type Lorem struct {
	Text       string `json:"text"`
	Paragraphs int    `json:"paragraphs"`
	WordCount  int    `json:"wordCount"`
}

func (Lorem) IsData() {}
