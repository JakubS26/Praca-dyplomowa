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

	generateParser()

	sampleInput := "ccdcd"

	lexer.SetInputString(sampleInput)
	Parse()

}
