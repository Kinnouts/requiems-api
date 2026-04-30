package jokes

// DadJoke represents a single dad joke with its identifier.
type DadJoke struct {
	ID   string `json:"id"`
	Joke string `json:"joke"`
}

func (DadJoke) IsData() {}
