package lorem

type Lorem struct {
	Text       string `json:"text"`
	Paragraphs int    `json:"paragraphs"`
	WordCount  int    `json:"wordCount"`
}

func (Lorem) IsData() {}
