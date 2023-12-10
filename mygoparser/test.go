package main

import (
	"bufio"
	"errors"
	"fmt"
	"goparser/lexer"
	"goparser/parser"
	"os"
	"strconv"
)

func main() {

	lexer := lexer.NewLexer()

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

	parser := parser.NewParser(lexer)

	parser.AddParserRule("S -> E NEWLINE", func(p []any) { fmt.Printf("Wynik: %d\n\n", p[1]) })
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

	parser.Init()

	for true {
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		if len(line) == 1 {
			break
		}
		lexer.SetInputString(line)
		err := parser.Parse()

		if err != nil {
			fmt.Print(err, "\n\n")
		}
	}

	parser.ExportParseTablesToFile("parse_table.txt")

}
