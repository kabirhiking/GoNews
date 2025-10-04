package main

import (
	"regexp"
	"strings"
)

var wordRE = regexp.MustCompile(`[a-zA-Z0-9]+`)

// toggle for stemming
var EnableStemming = false

// compact stopword list; extend as needed
var stopwords = map[string]bool{
	"the": true, "is": true, "and": true, "a": true, "an": true, "of": true, "to": true, "in": true,
	"for": true, "on": true, "with": true, "by": true, "that": true, "this": true, "it": true, "as": true,
	"are": true, "was": true, "at": true, "from": true, "be": true, "has": true, "have": true,
}

// Tokenize returns lowercase tokens from text, filtering stopwords
func Tokenize(text string) []string {
	text = strings.ToLower(text)
	matches := wordRE.FindAllString(text, -1)
	var tokens []string
	for _, m := range matches {
		if stopwords[m] {
			continue
		}
		if EnableStemming {
			m = Stem(m)
		}
		tokens = append(tokens, m)
	}
	return tokens
}

// Stem is placeholder for a stemming function. To enable real stemming:
//    go get github.com/reiver/go-porterstemmer
// and replace this implementation with call to that package.
func Stem(w string) string {
	// placeholder: return as-is. If you want stemming, uncomment and use a porter stemmer.
	// return porterstemmer.StemString(w)
	return w
}