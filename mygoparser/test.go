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

type T = []parser.Object

func main() {

	// lexer.AddTokenDefinition("NUM", `[0-9]+`)
	// lexer.AddTokenDefinition("PLUS", `\+`)
	// lexer.AddTokenDefinition("TIMES", `\*`)
	// lexer.AddTokenDefinition("L_PAR", `\(`)
	// lexer.AddTokenDefinition("R_PAR", `\)`)

	// lexer.Init()
	// lexer.OpenFile("expr_test4.txt")

	// parser.AddParserRule("E -> E PLUS T", func(p T) { p[0].IntegerValue = p[1].IntegerValue + p[3].IntegerValue; fmt.Println(p[0].IntegerValue) })
	// parser.AddParserRule("E -> T", func(p T) { p[0].IntegerValue = p[1].IntegerValue })
	// parser.AddParserRule("T -> T TIMES F", func(p T) { p[0].IntegerValue = p[1].IntegerValue * p[3].IntegerValue })
	// parser.AddParserRule("T -> F", func(p T) { p[0].IntegerValue = p[1].IntegerValue })
	// parser.AddParserRule("F -> L_PAR E R_PAR", func(p T) { p[0].IntegerValue = p[2].IntegerValue })
	// parser.AddParserRule("F -> NUM", func(p T) { p[0].IntegerValue, _ = strconv.Atoi(p[1].GetStringValue()) })

	// //parser.ParseWithSemanticActions()

	// C := parsergen.CreateLr0ItemSets()

	// for _, set := range C {
	// 	parsergen.Print(set)
	// 	fmt.Printf("\n")
	// }

	// transitions := parsergen.GetTransitions()

	// for _, x := range transitions {
	// 	for _, y := range x {
	// 		fmt.Println(y.GetSourceState(), "  ", parser.GetSymbolName(y.GetSymbol()), "  ", y.GetDestState())
	// 	}
	// }

	// parsergen.GenerateDrSets(parser.GetMinimalNonTerminalIndex())

	lexer.AddTokenDefinition("NEWLINE", `\n`)
	lexer.AddTokenDefinition("NUM", `[0-9]+`)
	lexer.AddTokenDefinition("PLUS", `\+`)
	lexer.AddTokenDefinition("MINUS", `\-`)
	lexer.AddTokenDefinition("TIMES", `\*`)
	lexer.AddTokenDefinition("DIV", `\/`)
	lexer.AddTokenDefinition("L_PAR", `\(`)
	lexer.AddTokenDefinition("R_PAR", `\)`)

	lexer.Init()

	parser.AddParserRule("S -> E NEWLINE", func(p T) { fmt.Printf("Wynik: %d\n\n", p[1].IntegerValue) })
	parser.AddParserRule("E -> E PLUS T", func(p T) { p[0].IntegerValue = p[1].IntegerValue + p[3].IntegerValue })
	parser.AddParserRule("E -> E MINUS T", func(p T) { p[0].IntegerValue = p[1].IntegerValue - p[3].IntegerValue })
	parser.AddParserRule("E -> T", func(p T) { p[0].IntegerValue = p[1].IntegerValue })
	parser.AddParserRule("T -> T TIMES F", func(p T) { p[0].IntegerValue = p[1].IntegerValue * p[3].IntegerValue })
	parser.AddParserRule("T -> T DIV F", func(p T) {
		if p[3].IntegerValue != 0 {
			p[0].IntegerValue = p[1].IntegerValue / p[3].IntegerValue
		}
	})
	parser.AddParserRule("T -> F", func(p T) { p[0].IntegerValue = p[1].IntegerValue })
	parser.AddParserRule("F -> L_PAR E R_PAR", func(p T) { p[0].IntegerValue = p[2].IntegerValue })
	parser.AddParserRule("F -> NUM", func(p T) { p[0].IntegerValue, _ = strconv.Atoi(p[1].GetStringValue()) })

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
