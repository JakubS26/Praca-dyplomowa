package parser

import (
	"goparser/lexer"
	"testing"
)

func TestFull(t *testing.T) {

	lexer.AddTokenDefinition("c", `c`)
	lexer.AddTokenDefinition("d", `d`)

	lexer.Init()

	AddParserRule("S -> C C", nil)
	AddParserRule("C -> c C", nil)
	AddParserRule("C -> d", nil)

	properStrings := []string{"dd", "cdd", "cccccccdd", "cdcd", "cccdcccd", "cdccccd"}

	for _, s := range properStrings {
		lexer.SetInputString(s)
		err := Parse()
		if err != nil {
			t.Fatalf("Parsing failed for string: " + s)
		}
	}

	improperStrings := []string{"c", "cd", "cdcdc", "ddd"}

	for _, s := range improperStrings {
		lexer.SetInputString(s)
		err := Parse()
		if err == nil {
			t.Fatalf("Parsing should have failed for string: " + s)
		}
	}

}
