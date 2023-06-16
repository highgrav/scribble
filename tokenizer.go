package scribble

import (
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
)

type Tokenizer struct {
	ready           bool
	tokens          map[string]*regexp2.Regexp
	bounds          []*regexp2.Regexp
	tokenPriorities map[string]int
	sortedTokens    []string
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		ready:           false,
		tokens:          make(map[string]*regexp2.Regexp),
		bounds:          make([]*regexp2.Regexp, 0),
		tokenPriorities: make(map[string]int),
	}
}

func (tok *Tokenizer) Token(priority int, tokenName, tokenRegex string) error {
	rxp := regexp2.MustCompile("^"+tokenRegex+"$", 0)
	tok.tokens[tokenName] = rxp
	tok.tokenPriorities[tokenName] = priority
	return nil
}

func (tok *Tokenizer) Boundary(boundsRegex string) error {
	rxp := regexp2.MustCompile("^"+boundsRegex+"$", 0)
	tok.bounds = append(tok.bounds, rxp)
	return nil
}

func (tok *Tokenizer) Compile() {
	tok.sortTokens()
	tok.ready = true
}

func (tok *Tokenizer) Tokenize(stream string) ([]Token, error) {
	if !tok.ready {
		return nil, errors.New("tokenizer not compiled")
	}
	tokens := make([]Token, 0)
	loc := 0
	for loc < len(stream) {
		match := MATCH_NONE
		name := ""

		// We try to parse from loc to the end of the stream, then keep going backwards until we
		// get the longest possible match
		var lookahead int = len(stream) + 1

		// keep going until we have a match
		for match == MATCH_NONE {

			lookahead = lookahead - 1
			if lookahead > (len(stream) - loc) {
				lookahead = len(stream) - loc
			}
			if lookahead < 1 {
				return tokens, fmt.Errorf("(%d:%d) unresolvable token: %q\n", loc, lookahead, Substring(stream, loc, len(stream)-loc))
			}
			match, name = tok.MatchTypes(Substring(stream, loc, lookahead))
		}

		switch match {
		case MATCH_BOUNDARY:
			loc = loc + lookahead
			break
		case MATCH_TOKEN:
			tokens = append(tokens, Token{
				Name:    name,
				Literal: Substring(stream, loc, lookahead),
			})
			loc = loc + lookahead
			break
		case MATCH_NONE:
			// fall through
		}
	}

	return tokens, nil
}

func (tok *Tokenizer) TokenizeOneAhead(stream string) ([]Token, error) {
	if !tok.ready {
		return nil, errors.New("tokenizer not compiled")
	}
	return nil, nil
}

func (tok *Tokenizer) sortTokens() {
	min := 9999999999
	max := -1
	sortables := make(map[int][]string)
	sorted := make([]string, 0)
	for k, v := range tok.tokenPriorities {
		if sortables[v] == nil {
			sortables[v] = make([]string, 0)
		}
		sortables[v] = append(sortables[v], k)
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	for i := min; i <= max; i++ {
		if sortables[i] != nil {
			sorted = append(sorted, sortables[i]...)
		}
	}
	tok.sortedTokens = sorted
}

func (tok *Tokenizer) MatchToken(buffer string) string {
	if !tok.ready {
		return ""
	}
	for _, k := range tok.sortedTokens {
		if match, _ := tok.tokens[k].MatchString(buffer); match {
			return k
		}
	}
	return ""
}

func (tok *Tokenizer) MatchTypes(buffer string) (MatchTypes, string) {
	if !tok.ready {
		return MATCH_NONE, ""
	}
	for _, rx := range tok.bounds {
		if match, _ := rx.MatchString(buffer); match {
			return MATCH_BOUNDARY, ""
		}
	}
	for _, k := range tok.sortedTokens {
		if match, _ := tok.tokens[k].MatchString(buffer); match {
			return MATCH_TOKEN, k
		}
	}
	return MATCH_NONE, ""
}

func Substring(buf string, start, length int) string {
	if start > len(buf)-1 || length <= 0 {
		return ""
	}
	if length > len(buf)-start {
		return buf[start:]
	}
	return buf[start:(start + length)]
}
