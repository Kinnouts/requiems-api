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

func (s *Service) Normalize(email string) (EmailNormalization, error) {
	result, err := s.n.Normalize2(email)
	if err != nil {
		return EmailNormalization{}, err
	}

	changes := result.Changes
	if changes == nil {
		changes = []normalizer.Change{}
	}

	local, domain, ok := strings.Cut(result.Normalized, "@")
	if !ok {
		return EmailNormalization{
			Original:   email,
			Normalized: result.Normalized,
			Changes:    changes,
		}, nil
	}

	return EmailNormalization{
		Original:   email,
		Normalized: result.Normalized,
		Local:      local,
		Domain:     domain,
		Changes:    changes,
	}, nil
}
