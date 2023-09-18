package main

import (
	"fmt"
	"goparser/lexer"
	"goparser/parser"
)

func main() {

	// lexer.AddTokenDefinition("KEYWORD_INT", `int`)
	// lexer.AddTokenDefinition("KEYWORD_IF", `if`)
	// lexer.AddTokenDefinition("KEYWORD_ELSE", `else`)
	// lexer.AddTokenDefinition("KEYWORD_RETURN", `return`)
	// lexer.AddTokenDefinition("WHITESPACE", `[ \t]`)
	// lexer.AddTokenDefinition("NEWLINE", `\n`)
	// lexer.AddTokenDefinition("IDENTIFIER", `[a-zA-Z][a-zA-Z0-9]*`)
	// lexer.AddTokenDefinition("NUMBER", `[0-9]+`)
	// lexer.AddTokenDefinition("LEFT_PAR", `\(`)
	// lexer.AddTokenDefinition("RIGHT_PAR", `\)`)
	// lexer.AddTokenDefinition("SEMICOLON", `;`)
	// lexer.AddTokenDefinition("EQUALS", `=`)
	// lexer.AddTokenDefinition("INCREMENT", `\+\+`)
	// lexer.AddTokenDefinition("DECREMENT", `\-\-`)
	// lexer.AddTokenDefinition("PLUS_EQUALS", `\+=`)
	// lexer.AddTokenDefinition("PLUS", `\+`)
	// lexer.AddTokenDefinition("MINUS", `\-`)
	// lexer.AddTokenDefinition("TIMES", `\*`)
	// lexer.AddTokenDefinition("GREATER", `>`)
	// lexer.AddTokenDefinition("LEFT_BR", `\{`)
	// lexer.AddTokenDefinition("RIGHT_BR", `}`)

	// lexer.Init()

	// lexer.OpenFile("expr_test.txt")
	// lexer.OpenFile("c_test.txt")

	// var err error = nil
	// var tok lexer.Token

	// for err == nil {
	// 	tok, err = lexer.NextToken()
	// 	lexer.PrintToken(tok)
	// }

	/***********************************************************/

	lexer.AddTokenDefinition("ID", `[a-zA-Z][a-zA-Z0-9]*`)
	lexer.AddTokenDefinition("PLUS", `\+`)
	lexer.AddTokenDefinition("TIMES", `\*`)
	lexer.AddTokenDefinition("L_PAR", `\(`)
	lexer.AddTokenDefinition("R_PAR", `\)`)

	lexer.Init()

	lexer.OpenFile("expr_test2.txt")

	parser.Parse()

	testSliceToken := make(map[string]int)
	testSliceNonTerminal := make(map[string]int)

	testSliceToken["ID"] = 0
	testSliceToken["PLUS"] = 1
	testSliceToken["TIMES"] = 2
	testSliceToken["L_PAR"] = 3
	testSliceToken["R_PAR"] = 4

	result, _ := parser.ToParserRule("E -> E PLUS T", testSliceToken, testSliceNonTerminal)
	fmt.Println(result)

	result, _ = parser.ToParserRule("E -> T", testSliceToken, testSliceNonTerminal)
	fmt.Println(result)

	result, _ = parser.ToParserRule("T -> T TIMES F", testSliceToken, testSliceNonTerminal)
	fmt.Println(result)

	result, _ = parser.ToParserRule("T -> F", testSliceToken, testSliceNonTerminal)
	fmt.Println(result)

	result, _ = parser.ToParserRule("F -> L_PAR E R_PAR", testSliceToken, testSliceNonTerminal)
	fmt.Println(result)

	result, _ = parser.ToParserRule("F -> ID", testSliceToken, testSliceNonTerminal)
	fmt.Println(result)

}
