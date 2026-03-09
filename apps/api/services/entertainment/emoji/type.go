package emoji

// Emoji represents a single emoji with its metadata.
type Emoji struct {
	Emoji    string `json:"emoji"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Unicode  string `json:"unicode"`
}

func (Emoji) IsData() {}

// EmojiList represents a collection of emoji search results.
type EmojiList struct {
	Items []Emoji `json:"items"`
	Total int     `json:"total"`
}

func (EmojiList) IsData() {}

// SearchRequest holds the query parameter for emoji search.
type SearchRequest struct {
	Query string `query:"q" validate:"required,min=1,max=100"`
}
