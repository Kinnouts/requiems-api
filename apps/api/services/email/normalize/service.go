package normalize

import (
	"strings"

	normalizer "github.com/bobadilla-tech/go-email-normalizer"
)

type Service struct {
	n normalizer.Normalizer
}

func NewService() *Service {
	return &Service{
		n: *normalizer.NewNormalizer(),
	}
}

func (s *Service) Normalize(email string) EmailNormalization {
	result := s.n.Normalize2(email)

	parts := strings.Split(result.Normalized, "@")

	return EmailNormalization{
		Original:   email,
		Normalized: result.Normalized,
		Local:      parts[0],
		Domain:     parts[1],
		Changes:    result.Changes,
	}
}
