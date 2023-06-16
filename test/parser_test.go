package test

import (
	"fmt"
	"github.com/highgrav/scribble"
	"testing"
)

/*
Test with a simple state machine that can handle trivial select SQL.
*/
func TestParser(t *testing.T) {
	sm := scribble.NewStateMachine()
	sm.Terminal("__EPSILON__")
	sm.State("", "SELECT", "SELECT")
	sm.State("SELECT", "ID", "SELECT_COL")
	sm.State("SELECT_COL", "COMMA", "SELECT_COL_SEPARATOR")
	sm.State("SELECT_COL_SEPARATOR", "ID", "SELECT_COL")
	sm.State("SELECT_COL", "FROM", "SELECT_FROM")
	sm.State("SELECT_FROM", "ID", "SELECT_FROM_ID")
	sm.State("SELECT_FROM_ID", "SEMICOLON", "__EPSILON__")
	sm.State("SELECT_FROM_ID", "WHERE", "SELECT_WHERE")
	sm.State("SELECT_WHERE", "ID", "SELECT_LVAL")
	sm.State("SELECT_WHERE", "INTEGER", "SELECT_LVAL")
	sm.State("SELECT_WHERE", "DECIMAL", "SELECT_LVAL")
	sm.State("SELECT_LVAL", "EQ", "SELECT_COMPARATOR")
	sm.State("SELECT_COMPARATOR", "ID", "SELECT_RVAL")
	sm.State("SELECT_COMPARATOR", "INTEGER", "SELECT_RVAL")
	sm.State("SELECT_COMPARATOR", "DECIMAL", "SELECT_RVAL")
	sm.State("SELECT_RVAL", "AND", "SELECT_WHERE")
	sm.State("SELECT_RVAL", "SEMICOLON", "__EPSILON__")

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

	toks, err := tok.Tokenize("select foo, bar, baz from foobar where foo = 1 and bar = baz and 1.23 = 5;")
	if err != nil {
		t.Error(err)
		return
	}

	states, err := sm.Parse("", toks)
	if err != nil {
		fmt.Printf("Processed %d tokens\n", len(states))
		t.Error(err)
		return
	}
	tgtStates := []string{
		"SELECT",
		"SELECT_COL", "SELECT_COL_SEPARATOR",
		"SELECT_COL", "SELECT_COL_SEPARATOR",
		"SELECT_COL",
		"SELECT_FROM", "SELECT_FROM_ID",
		"SELECT_WHERE",
		"SELECT_LVAL", "SELECT_COMPARATOR", "SELECT_RVAL",
		"SELECT_WHERE", "SELECT_LVAL", "SELECT_COMPARATOR", "SELECT_RVAL",
		"SELECT_WHERE", "SELECT_LVAL", "SELECT_COMPARATOR", "SELECT_RVAL",
		"__EPSILON__",
	}
	for i, s := range states {
		if s.State != tgtStates[i] {
			t.Errorf("at token %d expected state %q got %q", i, tgtStates[i], s.State)
		}
	}
}
