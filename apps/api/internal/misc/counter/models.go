package counter

import (
	"errors"
	"fmt"
	"regexp"
)

// Counter is the response model returned by both endpoints.
type Counter struct {
	Namespace string `json:"namespace"`
	Value     int64  `json:"value"`
}

func (Counter) IsData() {}

// ErrInvalidNamespace is returned when the namespace fails validation.
// Handlers use this to distinguish client errors (400) from server errors (500).
var ErrInvalidNamespace = errors.New("invalid namespace")

var namespaceRe = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

func validateNamespace(ns string) error {
	if !namespaceRe.MatchString(ns) {
		return fmt.Errorf("%w: must be 1–64 chars, alphanumeric, hyphen or underscore only", ErrInvalidNamespace)
	}
	return nil
}

func redisKey(namespace string) string {
	return "counter:" + namespace
}
