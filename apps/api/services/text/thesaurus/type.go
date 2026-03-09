package thesaurus

// Result is the response payload for the thesaurus endpoint.
type Result struct {
	Word     string   `json:"word"`
	Synonyms []string `json:"synonyms"`
	Antonyms []string `json:"antonyms"`
}

func (Result) IsData() {}
