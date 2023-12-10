package parser

import (
	"errors"
	"fmt"
	"goparser/lexer"
	"strconv"
	"testing"
)

func TestClc2(t *testing.T) {

	lexer := lexer.NewLexer()

	lexer.AddTokenDefinition("NL", `\n`)
	lexer.AddTokenDefinition("NUM", `[0-9]+`)
	lexer.AddTokenDefinition("PLUS", `\+`)
	lexer.AddTokenDefinition("MINUS", `\-`)
	lexer.AddTokenDefinition("TIMES", `\*`)
	lexer.AddTokenDefinition("DIV", `\/`)
	lexer.AddTokenDefinition("L_PAR", `\(`)
	lexer.AddTokenDefinition("R_PAR", `\)`)

	lexer.Ignore(` `)
	lexer.Ignore(`\t`)

	lexer.Init()

	parser := NewParser(lexer)

	parser.AddParserRule("S -> E NL", func(p []any) { fmt.Printf("Wynik: %d\n\n", p[1]) })
	parser.AddParserRule("E -> E PLUS T", func(p []any) { p[0] = p[1].(int) + p[3].(int) })
	parser.AddParserRule("E -> E MINUS T", func(p []any) { p[0] = p[1].(int) - p[3].(int) })
	parser.AddParserRule("E -> T", func(p []any) { p[0] = p[1].(int) })
	parser.AddParserRule("T -> T TIMES F", func(p []any) { p[0] = p[1].(int) * p[3].(int) })
	parser.AddParserRule("T -> T DIV F", func(p []any) {
		if p[3].(int) != 0 {
			p[0] = p[1].(int) / p[3].(int)
		} else {
			parser.RaiseError(errors.New("Error: division by 0"))
		}
	})
	parser.AddParserRule("T -> F", func(p []any) { p[0] = p[1].(int) })
	parser.AddParserRule("F -> L_PAR E R_PAR", func(p []any) { p[0] = p[2].(int) })
	parser.AddParserRule("F -> NUM", func(p []any) { p[0], _ = strconv.Atoi(p[1].(string)) })
	parser.AddParserRule("F -> MINUS NUM", func(p []any) { p[0], _ = strconv.Atoi(p[2].(string)); p[0] = p[0].(int) * (-1) })

	fmt.Print("  ")
	for i := 0; i < parser.numberOfGrammarSymbols; i++ {
		fmt.Printf("%6.6s", parser.getSymbolName(i))
	}
	fmt.Println()

	parser.generateParser()

	for index, row := range parser.parsingTable {
		fmt.Printf("%2d", index)
		for _, action := range row {
			fmt.Printf("%6.6s", action)
		}
		fmt.Println()
	}

}
