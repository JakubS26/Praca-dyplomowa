package parser

import (
	"goparser/lexer"
	"testing"
)

func TestEps(t *testing.T) {

	lexer.AddTokenDefinition("a", `a`)
	lexer.AddTokenDefinition("b", `b`)
	lexer.AddTokenDefinition("c", `c`)

	lexer.Init()

	AddParserRule("S -> A B C", nil)
	AddParserRule("A -> a A", nil)
	AddParserRule("A -> epsilon", nil)
	AddParserRule("B -> b B", nil)
	AddParserRule("B -> epsilon", nil)
	AddParserRule("C -> c C", nil)
	AddParserRule("C -> epsilon", nil)

	properStrings := []string{"", "aabbcc", "abc", "aaaa", "ab", "a", "b", "bc", "c"}

	for _, s := range properStrings {
		lexer.SetInputString(s)
		err := Parse()
		if err != nil {
			t.Fatalf("Parsing failed for string: " + s)
		}
	}

	improperStrings := []string{"cba", "aba", "ccccccb", "ababab"}

	for _, s := range improperStrings {
		lexer.SetInputString(s)
		err := Parse()
		if err == nil {
			t.Fatalf("Parsing should have failed for string: " + s)
		}
	}

}
