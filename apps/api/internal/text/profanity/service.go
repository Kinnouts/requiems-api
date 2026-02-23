package profanity

import (
	"regexp"
	"strings"
)

// wordRe matches sequences of alphabetic characters (words).
var wordRe = regexp.MustCompile(`[a-zA-Z]+`)

// profanitySet is the curated list of flagged words (all lowercase).
// Words are matched case-insensitively against the input text.
var profanitySet = map[string]struct{}{
	"ass":          {},
	"asshole":      {},
	"bastard":      {},
	"bitch":        {},
	"bullshit":     {},
	"cock":         {},
	"crap":         {},
	"cunt":         {},
	"damn":         {},
	"dick":         {},
	"douche":       {},
	"douchebag":    {},
	"dumbass":      {},
	"fag":          {},
	"faggot":       {},
	"fuck":         {},
	"fucker":       {},
	"fucking":      {},
	"goddamn":      {},
	"jackass":      {},
	"jerk":         {},
	"moron":        {},
	"motherfucker": {},
	"piss":         {},
	"prick":        {},
	"pussy":        {},
	"shit":         {},
	"shithead":     {},
	"slut":         {},
	"twat":         {},
	"whore":        {},
	"wanker":       {},
}

// Service performs profanity detection and censoring.
type Service struct{}

// NewService returns a new profanity Service.
func NewService() *Service { return &Service{} }

// Check inspects text for profanity, returning a censored copy of the text
// and the deduplicated list of flagged words found.
func (s *Service) Check(text string) Result {
	flaggedSet := map[string]bool{}
	var flaggedWords []string

	censored := wordRe.ReplaceAllStringFunc(text, func(word string) string {
		lower := strings.ToLower(word)
		if _, ok := profanitySet[lower]; ok {
			if !flaggedSet[lower] {
				flaggedSet[lower] = true
				flaggedWords = append(flaggedWords, lower)
			}
			return strings.Repeat("*", len(word))
		}
		return word
	})

	if flaggedWords == nil {
		flaggedWords = []string{}
	}

	return Result{
		HasProfanity: len(flaggedWords) > 0,
		Censored:     censored,
		FlaggedWords: flaggedWords,
	}
}
