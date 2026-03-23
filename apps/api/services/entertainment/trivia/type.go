package trivia

// Request holds the optional query parameters for the trivia endpoint.
type Request struct {
	Category   string `query:"category"   validate:"omitempty,oneof=science history geography sports music movies literature math technology nature"`
	Difficulty string `query:"difficulty" validate:"omitempty,oneof=easy medium hard"`
}

// Question is a trivia question with multiple-choice answers.
type Question struct {
	Question   string   `json:"question"`
	Options    []string `json:"options"`
	Answer     string   `json:"answer"`
	Category   string   `json:"category"`
	Difficulty string   `json:"difficulty"`
}

func (Question) IsData() {}
