package detectlanguage

import (
	"strings"

	lingua "github.com/pemistahl/lingua-go"
)

// detector is the default language detector shared across all requests.
var detector = lingua.NewLanguageDetectorBuilder().
	FromAllLanguages().
	Build()

// Service performs language detection.
type Service struct{}

// NewService returns a new language detection Service.
func NewService() *Service { return &Service{} }

// Detect identifies the language of the given text.
func (s *Service) Detect(text string) Result {
	language, reliable := detector.DetectLanguageOf(text)
	if !reliable {
		return Result{
			Language:   "Unknown",
			Code:       "",
			Confidence: 0,
		}
	}

	confidence := detector.ComputeLanguageConfidence(text, language)

	return Result{
		Language:   language.String(),
		Code:       strings.ToLower(language.IsoCode639_1().String()),
		Confidence: confidence,
	}
}
