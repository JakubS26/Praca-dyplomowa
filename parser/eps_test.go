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

	generateParser()

}
