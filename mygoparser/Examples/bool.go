package main

import (
	"bufio"
	"fmt"
	"goparser/lexer"
	"goparser/parser"
	"goparser/parsergen"
	"os"
)

func main() {

	lexer.AddTokenDefinition("NEWLINE", `\n`)
	lexer.AddTokenDefinition("TRUE", `true`)
	lexer.AddTokenDefinition("FALSE", `false`)
	lexer.AddTokenDefinition("OR", `or`)
	lexer.AddTokenDefinition("AND", `and`)
	lexer.AddTokenDefinition("NOT", `not`)
	lexer.AddTokenDefinition("L_PAR", `\(`)
	lexer.AddTokenDefinition("R_PAR", `\)`)

	lexer.Ignore(` `)
	lexer.Ignore(`\t`)

	lexer.Init()

	parser.AddParserRule("S -> E NEWLINE", func(p []any) { fmt.Printf("Wynik: %v\n\n", p[1]) })
	parser.AddParserRule("E -> E OR T", func(p []any) { p[0] = p[1].(bool) || p[3].(bool) })
	parser.AddParserRule("E -> T", func(p []any) { p[0] = p[1].(bool) })
	parser.AddParserRule("T -> T AND F", func(p []any) { p[0] = p[1].(bool) && p[3].(bool) })
	parser.AddParserRule("T -> F", func(p []any) { p[0] = p[1].(bool) })
	parser.AddParserRule("F -> L_PAR E R_PAR", func(p []any) { p[0] = p[2].(bool) })
	parser.AddParserRule("F -> NOT F", func(p []any) { p[0] = !p[2].(bool) })
	parser.AddParserRule("F -> TRUE", func(p []any) { p[0] = true })
	parser.AddParserRule("F -> FALSE", func(p []any) { p[0] = false })

	parsergen.GenerateParser()

	for true {
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		if len(line) == 1 {
			break
		}
		lexer.SetInputString(line)
		parser.Parse()
	}

}
