package counter

import (
	"fmt"
	"regexp"
)

// Counter is the response model returned by both endpoints.
type Counter struct {
	Namespace string `json:"namespace"`
	Value     int64  `json:"value"`
}

var namespaceRe = regexp.MustCompile(`^[a-zA-Z0-9_\-]{1,64}$`)

func validateNamespace(ns string) error {
	if !namespaceRe.MatchString(ns) {
		return fmt.Errorf("namespace must be 1–64 chars: alphanumeric, hyphen or underscore only")
	}
	return nil
}

func redisKey(namespace string) string {
	return "counter:" + namespace
}
