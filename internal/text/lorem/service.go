package lorem

import "strings"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Generate(paragraphs int, sentences int) Lorem {
	var b strings.Builder

	// Heuristic to reduce reallocations; exact sizing not required
	b.Grow(paragraphs * sentences * 64)

	for p := range paragraphs {
		if p > 0 {
			b.WriteString("\n\n")
		}
		for s := range sentences {
			if s > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(GenerateSentence(10))
		}
	}

	text := b.String()

	return Lorem{
		Text:       text,
		Paragraphs: paragraphs,
		WordCount:  paragraphs * sentences * 10,
	}
}
