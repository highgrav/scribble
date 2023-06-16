package test

import (
	"errors"
	"fmt"
	"github.com/highgrav/scribble"
	"testing"
)

func TestSubstring(t *testing.T) {
	str := "hello, world!"
	if scribble.Substring(str, 0, 4) != "hell" {
		t.Error("failed to get substring 0,4")
		return
	}
	if scribble.Substring(str, 1, 4) != "ello" {
		t.Error("failed to get substring 1,4")
		return
	}
	if scribble.Substring(str, 7, 6) != "world!" {
		t.Error("failed to get substring 7,6")
		return
	}

}

func TestTokenize(t *testing.T) {
	tok := scribble.NewTokenizer()
	tok.Boundary(`[\s]+`)
	tok.Token(0, "SEMICOLON", `\;`)
	tok.Token(0, "COMMA", `,`)
	tok.Token(0, "EQ", `\=`)

	tok.Token(1, "SELECT", `(?i)select`)
	tok.Token(1, "FROM", `(?i)from`)
	tok.Token(1, "WHERE", `(?i)where`)
	tok.Token(1, "AND", `(?i)and`)

	tok.Token(2, "DECIMAL", `[-]?[0-9]+\.[0-9]+`)

	tok.Token(3, "INTEGER", `[-]?[0-9]+`)
	tok.Token(3, "ID", `[A-Za-z_][A-Za-z0-9_]*`)
	tok.Compile()
	res, val := tok.MatchTypes("       ")
	if res != scribble.MATCH_BOUNDARY {
		t.Error(errors.New("failed to match whitespace boundary"))
		return
	}
	res, val = tok.MatchTypes("TESTS")
	if res != scribble.MATCH_TOKEN || val != "ID" {
		t.Error(fmt.Errorf("failed to match ID (got %s for %q)", res, val))
		return
	}
	res, val = tok.MatchTypes("SELECT")
	if res != scribble.MATCH_TOKEN || val != "SELECT" {
		t.Error(fmt.Errorf("failed to match SELECT (got %s for %q)", res, val))
		return
	}

	toks, err := tok.Tokenize("SELECT foo FROM bar where foobar = 1 and foobaz = 1.23;")
	if err != nil {
		t.Error(err)
		return
	}
	if len(toks) != 13 {
		t.Errorf("expected 13 tokens, got %d", len(toks))
		return
	}
	vals := []string{
		"SELECT", "ID", "FROM", "ID", "WHERE", "ID", "EQ", "INTEGER", "AND", "ID", "EQ", "DECIMAL", "SEMICOLON",
	}
	for i, v := range vals {
		if toks[i].Name != v {
			t.Errorf("got unexpected token at position %d: expected %q, got %q on literal %q", i, v, toks[i].Name, toks[i].Literal)
			return
		}
	}
}
