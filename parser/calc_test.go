package parser

import (
	"goparser/lexer"
	"testing"
)

func TestCalc(t *testing.T) {

	lexer.AddTokenDefinition("NL", `\n`)
	lexer.AddTokenDefinition("NUM", `[0-9]+`)
	lexer.AddTokenDefinition("PLUS", `\+`)
	lexer.AddTokenDefinition("TIMES", `\*`)
	lexer.AddTokenDefinition("L_PAR", `\(`)
	lexer.AddTokenDefinition("R_PAR", `\)`)

	lexer.Init()

	AddParserRule("S -> E NL", nil)
	AddParserRule("E -> E PLUS T", nil)
	AddParserRule("E -> T", nil)
	AddParserRule("T -> T TIMES F", nil)
	AddParserRule("T -> F", nil)
	AddParserRule("F -> L_PAR E R_PAR", nil)
	AddParserRule("F -> NUM", nil)

	generateParser()

}
