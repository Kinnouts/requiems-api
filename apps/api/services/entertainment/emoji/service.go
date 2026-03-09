package emoji

import (
	"math/rand"
	"strings"
)

// Service provides emoji lookup and search operations.
type Service struct{}

// NewService returns a new Service.
func NewService() *Service { return &Service{} }

// Random returns a randomly selected emoji.
func (s *Service) Random() Emoji {
	return emojis[rand.Intn(len(emojis))]
}

// GetByName returns the emoji matching the given name (snake_case).
// Returns the emoji and true if found, or a zero value and false if not.
func (s *Service) GetByName(name string) (Emoji, bool) {
	name = strings.ToLower(name)
	for _, e := range emojis {
		if e.Name == name {
			return e, true
		}
	}
	return Emoji{}, false
}

// Search returns all emojis whose name contains the query string
// (case-insensitive). Returns an EmojiList with matching results.
func (s *Service) Search(query string) EmojiList {
	query = strings.ToLower(query)
	var matches []Emoji
	for _, e := range emojis {
		if strings.Contains(e.Name, query) || strings.Contains(strings.ToLower(e.Category), query) {
			matches = append(matches, e)
		}
	}
	if matches == nil {
		matches = []Emoji{}
	}
	return EmojiList{Items: matches, Total: len(matches)}
}
