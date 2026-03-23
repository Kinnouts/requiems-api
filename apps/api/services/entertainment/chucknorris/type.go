package chucknorris

// Fact represents a single Chuck Norris fact/joke.
type Fact struct {
	ID   string `json:"id"`
	Fact string `json:"fact"`
}

func (Fact) IsData() {}
