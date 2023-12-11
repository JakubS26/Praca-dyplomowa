package parser

import (
	"goparser/lexer"
	"testing"
)

func TestCalc(t *testing.T) {

	lexer := lexer.NewLexer()

	lexer.AddTokenDefinition("NUM", `[0-9]+`)
	lexer.AddTokenDefinition("PLUS", `\+`)
	lexer.AddTokenDefinition("TIMES", `\*`)
	lexer.AddTokenDefinition("L_PAR", `\(`)
	lexer.AddTokenDefinition("R_PAR", `\)`)

	lexer.Init()

	parser := NewParser(lexer)

	parser.AddParserRule("E -> E PLUS T", nil)
	parser.AddParserRule("E -> T", nil)
	parser.AddParserRule("T -> T TIMES F", nil)
	parser.AddParserRule("T -> F", nil)
	parser.AddParserRule("F -> L_PAR E R_PAR", nil)
	parser.AddParserRule("F -> NUM", nil)

	parser.Init()

	properStrings := []string{"3", "3+3", "3+3*3", "(3+3)*3", "4*4*4*4*4*4", "(5)"}

	for _, s := range properStrings {
		lexer.SetInputString(s)
		err := parser.Parse()
		if err != nil {
			t.Fatalf("Parsing failed for string: " + s)
		}
	}

	improperStrings := []string{"3++", "*", "()", "1*(2+5", ""}

	for _, s := range improperStrings {
		lexer.SetInputString(s)
		err := parser.Parse()
		if err == nil {
			t.Fatalf("Parsing should have failed for string: " + s)
		}
	}

}
