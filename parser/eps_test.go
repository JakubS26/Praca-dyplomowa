package parser

import (
	"goparser/lexer"
	"testing"
)

func TestEps(t *testing.T) {

	lexer := lexer.NewLexer()

	lexer.AddTokenDefinition("a", `a`)
	lexer.AddTokenDefinition("b", `b`)
	lexer.AddTokenDefinition("c", `c`)

	lexer.Init()

	parser := NewParser(lexer)

	parser.AddParserRule("S -> A B C", nil)
	parser.AddParserRule("A -> a A", nil)
	parser.AddParserRule("A -> epsilon", nil)
	parser.AddParserRule("B -> b B", nil)
	parser.AddParserRule("B -> epsilon", nil)
	parser.AddParserRule("C -> c C", nil)
	parser.AddParserRule("C -> epsilon", nil)

	parser.Init()

	properStrings := []string{"", "aabbcc", "abc", "aaaa", "ab", "a", "b", "bc", "c"}

	for _, s := range properStrings {
		lexer.SetInputString(s)
		err := parser.Parse()
		if err != nil {
			t.Fatalf("Parsing failed for string: " + s)
		}
	}

	improperStrings := []string{"cba", "aba", "ccccccb", "ababab"}

	for _, s := range improperStrings {
		lexer.SetInputString(s)
		err := parser.Parse()
		if err == nil {
			t.Fatalf("Parsing should have failed for string: " + s)
		}
	}

}
