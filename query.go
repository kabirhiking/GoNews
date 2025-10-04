package main

import (
	"strings"
)

// QueryToRPN: parse a user query into RPN tokens supporting:
// - quoted phrases: "small cat" -> token PHRASE:small cat
// - operators: AND, OR, NOT (case-insensitive)
// - parentheses ( )
func QueryToRPN(q string) []string {
	// tokenize: keep quoted phrases together
	var toks []string
	q = strings.TrimSpace(q)
	if q == "" {
		return nil
	}
	// parse tokens
	cur := ""
	inQuote := false
	for i := 0; i < len(q); i++ {
		c := q[i]
		if c == '"' {
			if inQuote {
				// end quote
				if cur != "" {
					toks = append(toks, "PHRASE:"+cur)
				}
				cur = ""
				inQuote = false
			} else {
				inQuote = true
				cur = ""
			}
			continue
		}
		if inQuote {
			cur += string(c)
			continue
		}
		// outside quote: split on spaces and parentheses
		if c == ' ' {
			if cur != "" {
				toks = append(toks, cur)
				cur = ""
			}
			continue
		}
		if c == '(' || c == ')' {
			if cur != "" {
				toks = append(toks, cur)
				cur = ""
			}
			toks = append(toks, string(c))
			continue
		}
		cur += string(c)
	}
	if cur != "" {
		toks = append(toks, cur)
	}

	// normalize operators
	for i, t := range toks {
		t := strings.ToUpper(t)
		if t == "AND" || t == "OR" || t == "NOT" || t == "(" || t == ")" || strings.HasPrefix(t, "PHRASE:") {
			// keep as-is (phrase keeps case inside)
		} else {
			// normal token -> lowercase + tokenization step
			t = strings.ToLower(t)
			// break token into word tokens if it contains non-word chars
			sub := Tokenize(t)
			if len(sub) == 0 {
				// keep original token
				toks[i] = t
			} else if len(sub) == 1 {
				toks[i] = sub[0]
			} else {
				// if tokenization produced multiple tokens, join with _
				toks[i] = strings.Join(sub, "_")
			}
		}
	}

	// shunting-yard to convert to RPN
	prec := map[string]int{"OR": 1, "AND": 2, "NOT": 3}
	var out []string
	var opstack []string
	pushOp := func(op string) { opstack = append(opstack, op) }
	popOp := func() string {
		op := opstack[len(opstack)-1]
		opstack = opstack[:len(opstack)-1]
		return op
	}
	for _, tk := range toks {
		if tk == "(" {
			pushOp(tk)
			continue
		}
		if tk == ")" {
			for len(opstack) > 0 && opstack[len(opstack)-1] != "(" {
				out = append(out, popOp())
			}
			if len(opstack) > 0 && opstack[len(opstack)-1] == "(" {
				popOp()
			}
			continue
		}
		u := strings.ToUpper(tk)
		if u == "AND" || u == "OR" || u == "NOT" {
			for len(opstack) > 0 {
				op := opstack[len(opstack)-1]
				if op == "(" {
					break
				}
				if prec[strings.ToUpper(op)] >= prec[u] {
					out = append(out, popOp())
				} else {
					break
				}
			}
			pushOp(u)
			continue
		}
		// term or phrase
		if strings.HasPrefix(tk, "PHRASE:") {
			// normalize phrase content
			ph := strings.TrimPrefix(tk, "PHRASE:")
			out = append(out, "PHRASE:"+ph)
		} else {
			out = append(out, tk)
		}
	}
	for len(opstack) > 0 {
		out = append(out, popOp())
	}
	return out
}

// isOperator helper
func isOperator(t string) bool {
	u := strings.ToUpper(t)
	return u == "AND" || u == "OR" || u == "NOT"
}

// MakeSnippet returns a small preview around first matched term(s)
func MakeSnippet(content string, terms []string) string {
	if len(content) == 0 {
		return ""
	}
	// tokenize content (lowercase tokens)
	toks := Tokenize(content)
	first := -1
	for i, w := range toks {
		for _, t := range terms {
			// if phrase term, check first token
			if strings.HasPrefix(t, "PHRASE:") {
				ph := strings.TrimPrefix(t, "PHRASE:")
				phToks := Tokenize(ph)
				if len(phToks) > 0 && w == phToks[0] {
					first = i
					break
				}
			} else {
				if w == t {
					first = i
					break
				}
			}
		}
		if first != -1 {
			break
		}
	}
	if first == -1 {
		// fallback: return start of doc up to 30 tokens
		end := 30
		if end > len(toks) {
			end = len(toks)
		}
		return strings.Join(toks[:end], " ") + "..."
	}
	start := first - 8
	if start < 0 {
		start = 0
	}
	end := first + 12
	if end > len(toks) {
		end = len(toks)
	}
	snippet := strings.Join(toks[start:end], " ")
	return "..." + snippet + "..."
}