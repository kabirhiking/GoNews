package main

import (
	"math"
	"sort"
	"strings"
)

// Posting: map of docID to positions
type Posting map[int][]int

// Index structure
type Index struct {
	Terms        map[string]Posting
	Docs         map[int]Document
	DocTokCounts map[int]int // number of tokens in each doc (for TF normalization)
	N            int         // number of documents
}

func NewIndex() *Index {
	return &Index{Terms: make(map[string]Posting), Docs: make(map[int]Document), DocTokCounts: make(map[int]int)}
}

// AddDocument tokenizes and adds to the inverted index
func (idx *Index) AddDocument(d Document) {
	idx.Docs[d.ID] = d
	tokens := Tokenize(d.Title + " " + d.Content)
	idx.DocTokCounts[d.ID] = len(tokens)
	for pos, tok := range tokens {
		if _, ok := idx.Terms[tok]; !ok {
			idx.Terms[tok] = make(Posting)
		}
		idx.Terms[tok][d.ID] = append(idx.Terms[tok][d.ID], pos)
	}
	idx.N = len(idx.Docs)
}

// helper: convert posting map to sorted slice of ids
func postingIDs(post Posting) []int {
	var ids []int
	for id := range post {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	return ids
}

// SearchResult holds docID and score/matches
type SearchResult struct {
	DocID        int
	Score        float64
	MatchedTerms []string
}

// Search is a full query processor: supports AND/OR/NOT and quoted phrases
func (idx *Index) Search(query string) []SearchResult {
	if len(query) == 0 {
		return nil
	}
	// parse query -> RPN tokens
	rpn := QueryToRPN(query)
	// evaluate RPN to get set of matching docIDs
	resSet := idx.EvaluateRPN(rpn)
	// convert set to scored results
	var results []SearchResult
	for doc := range resSet {
		// gather matched terms: any query term present in doc
		matched := idx.matchedTermsInDoc(doc, rpn)
		score := idx.scoreDoc(doc, matched)
		results = append(results, SearchResult{DocID: doc, Score: score, MatchedTerms: matched})
	}
	// sort by score desc
	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	return results
}

// matchedTermsInDoc extracts which query terms (non-operators) appear in the doc
func (idx *Index) matchedTermsInDoc(doc int, rpn []string) []string {
	set := map[string]bool{}
	for _, tok := range rpn {
		if isOperator(tok) { // skip
			continue
		}
		if strings.HasPrefix(tok, "PHRASE:") {
			phrase := strings.TrimPrefix(tok, "PHRASE:")
			tokens := Tokenize(phrase)
			if idx.checkPhraseInDoc(doc, tokens) {
				set[phrase] = true
			}
		} else {
			// normal token
			if posting, ok := idx.Terms[tok]; ok {
				if len(posting[doc]) > 0 {
					set[tok] = true
				}
			}
		}
	}
	var out []string
	for t := range set {
		out = append(out, t)
	}
	return out
}

// scoreDoc: TF-IDF style scoring using matched terms
func (idx *Index) scoreDoc(doc int, matched []string) float64 {
	score := 0.0
	for _, t := range matched {
		if strings.HasPrefix(t, "PHRASE:") {
			// give a boost for phrase matches
			score += 2.0
			continue
		}
		posting := idx.Terms[t]
		if posting == nil {
			continue
		}
		tf := float64(len(posting[doc]))
		df := float64(len(posting))
		if df == 0 || idx.DocTokCounts[doc] == 0 {
			continue
		}
		// normalize tf by doc length
		tfNorm := tf / float64(idx.DocTokCounts[doc])
		idf := math.Log(1 + float64(idx.N)/df)
		score += tfNorm * idf
	}
	return score
}

// EvaluateRPN evaluates RPN query tokens and returns a set (map[int]struct{}) of matching docs
func (idx *Index) EvaluateRPN(rpn []string) map[int]struct{} {
	stack := []map[int]struct{}{}
	universe := idx.allDocsSet()
	for _, tok := range rpn {
		if tok == "AND" || tok == "OR" {
			// binary
			if len(stack) < 2 {
				continue
			}
			r := stack[len(stack)-1]
			l := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if tok == "AND" {
				stack = append(stack, setIntersect(l, r))
			} else {
				stack = append(stack, setUnion(l, r))
			}
		} else if tok == "NOT" {
			// unary: pop one
			if len(stack) < 1 {
				continue
			}
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, setDiff(universe, a))
		} else {
			// term or phrase
			var s map[int]struct{}
			if strings.HasPrefix(tok, "PHRASE:") {
				phrase := strings.TrimPrefix(tok, "PHRASE:")
				toks := Tokenize(phrase)
				s = idx.docsWithPhrase(toks)
			} else {
				if posting, ok := idx.Terms[tok]; ok {
					s = make(map[int]struct{})
					for id := range posting {
						s[id] = struct{}{}
					}
				} else {
					s = map[int]struct{}{} // empty set
				}
			}
			stack = append(stack, s)
		}
	}
	if len(stack) == 0 {
		return map[int]struct{}{}
	}
	return stack[len(stack)-1]
}

// helpers to work with sets
func (idx *Index) allDocsSet() map[int]struct{} {
	out := make(map[int]struct{})
	for id := range idx.Docs {
		out[id] = struct{}{}
	}
	return out
}

func setIntersect(a, b map[int]struct{}) map[int]struct{} {
	res := make(map[int]struct{})
	// iterate smaller
	if len(a) > len(b) {
		a, b = b, a
	}
	for k := range a {
		if _, ok := b[k]; ok {
			res[k] = struct{}{}
		}
	}
	return res
}

func setUnion(a, b map[int]struct{}) map[int]struct{} {
	res := make(map[int]struct{})
	for k := range a {
		res[k] = struct{}{}
	}
	for k := range b {
		res[k] = struct{}{}
	}
	return res
}

func setDiff(a, b map[int]struct{}) map[int]struct{} {
	res := make(map[int]struct{})
	for k := range a {
		if _, ok := b[k]; !ok {
			res[k] = struct{}{}
		}
	}
	return res
}

// docsWithPhrase: return docs where tokens appear consecutively
func (idx *Index) docsWithPhrase(tokens []string) map[int]struct{} {
	res := make(map[int]struct{})
	if len(tokens) == 0 {
		return res
	}
	// get candidate docs by intersecting postings for each token
	var candidate []int
	for i, t := range tokens {
		posting, ok := idx.Terms[t]
		if !ok {
			return res
		}
		ids := postingIDs(posting)
		if i == 0 {
			candidate = ids
		} else {
			candidate = intersectSorted(candidate, ids)
		}
		if len(candidate) == 0 {
			return res
		}
	}
	for _, doc := range candidate {
		if idx.checkPhraseInDoc(doc, tokens) {
			res[doc] = struct{}{}
		}
	}
	return res
}

// checkPhraseInDoc: naive consecutive position check
func (idx *Index) checkPhraseInDoc(doc int, tokens []string) bool {
	posLists := make([][]int, len(tokens))
	for i, t := range tokens {
		posLists[i] = idx.Terms[t][doc]
		if len(posLists[i]) == 0 {
			return false
		}
	}
	for _, p := range posLists[0] {
		ok := true
		for i := 1; i < len(tokens); i++ {
			need := p + i
			if !contains(posLists[i], need) {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func contains(arr []int, x int) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func intersectSorted(a, b []int) []int {
	i, j := 0, 0
	var res []int
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			res = append(res, a[i])
			i++; j++
		} else if a[i] < b[j] {
			i++
		} else {
			j++
		}
	}
	return res
}