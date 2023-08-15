package main

import (
	"goparser/lexer"
)

func main() {

	lexer.AddTokenDefinition("KEYWORD_INT", `int`)
	lexer.AddTokenDefinition("WHITESPACE", ` `)
	lexer.AddTokenDefinition("NEWLINE", `\n`)
	lexer.AddTokenDefinition("NUMBER", `[0-9]+`)
	lexer.AddTokenDefinition("WORD", `[a-zA-z]+`)
	//lexer.AddTokenDefinition("WRONG", `][`)

	err := lexer.Init()

	// if err == nil {
	// 	lexer.PrintTokens()
	// } else {
	// 	fmt.Println(err)
	// }

	err = lexer.OpenFile("ala.txt")

	// if err == nil {
	// 	lexer.TestPrintFile()
	// }

	err = nil
	var tok lexer.Token

	for err == nil {
		tok, err = lexer.NextToken()
		lexer.PrintToken(tok)
	}

}
