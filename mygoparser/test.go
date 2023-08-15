package main

import (
	"fmt"
	"goparser/lexer"
)

func main() {

	lexer.AddTokenDefinition("NUMBER", "[0-9]+")
	lexer.AddTokenDefinition("WORD", "[a-z]+")
	lexer.AddTokenDefinition("KEYWORD_INT", "int")

	err := lexer.Init()

	if err == nil {
		lexer.PrintTokens()
	} else {
		fmt.Println(err)
	}

}
