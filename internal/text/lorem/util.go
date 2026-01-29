package lorem

import (
	"math/rand"
	"strings"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func formatSentence(sentence string) string {
	if len(sentence) == 0 {
		return sentence
	}

	// Took me 30 minutes to do this, don't judge me
	sentence = strings.TrimSpace(sentence)
	sentence = strings.ToUpper(sentence[:1]) + sentence[1:] + "."
	return sentence
}

func GenerateSentence(wordCount int) string {
	var b strings.Builder

	b.Grow(wordCount * 8)

	for i := range wordCount {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(Words[random.Intn(len(Words))])
	}

	return formatSentence(b.String())
}
