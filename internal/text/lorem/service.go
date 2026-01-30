package lorem

import (
	"strings"

	lorelai "github.com/bobadilla-tech/lorelai/pkg"
)

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
			b.WriteString(lorelai.ClassicSentence())
		}
	}

	text := b.String()

	return Lorem{
		Text:       text,
		Paragraphs: paragraphs,
		WordCount:  paragraphs * sentences * lorelai.ClassicWordsPerSentence(),
	}
}
