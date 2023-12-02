package main

import (
	"bufio"
	"fmt"
	"goparser/lexer"
	"goparser/parser"
	"goparser/parsergen"
	"os"
	"strconv"
)

func main() {

	lexer.AddTokenDefinition("NEWLINE", `\n`)
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

	parser.AddParserRule("S -> E NEWLINE", func(p []any) { fmt.Printf("Wynik: %d\n\n", p[1]) })
	parser.AddParserRule("E -> E PLUS T", func(p []any) { p[0] = p[1].(int) + p[3].(int) })
	parser.AddParserRule("E -> E MINUS T", func(p []any) { p[0] = p[1].(int) - p[3].(int) })
	parser.AddParserRule("E -> T", func(p []any) { p[0] = p[1].(int) })
	parser.AddParserRule("T -> T TIMES F", func(p []any) { p[0] = p[1].(int) * p[3].(int) })
	parser.AddParserRule("T -> T DIV F", func(p []any) {
		if p[3].(int) != 0 {
			p[0] = p[1].(int) / p[3].(int)
		} else {
			p[0] = 0
		}
	})
	parser.AddParserRule("T -> F", func(p []any) { p[0] = p[1].(int) })
	parser.AddParserRule("F -> L_PAR E R_PAR", func(p []any) { p[0] = p[2].(int) })
	parser.AddParserRule("F -> NUM", func(p []any) { p[0], _ = strconv.Atoi(p[1].(string)) })

	parsergen.GenerateParser()

	for true {
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		if len(line) == 1 {
			break
		}
		lexer.SetInputString(line)
		parser.ParseWithSemanticActions()
	}

}
