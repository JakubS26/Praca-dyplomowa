package parser

import (
	"goparser/lexer"
	"testing"
)

func TestFull(t *testing.T) {

	lexer := lexer.NewLexer()

	lexer.AddTokenDefinition("c", `c`)
	lexer.AddTokenDefinition("d", `d`)

	lexer.Init()

	parser := NewParser(lexer)

	parser.AddParserRule("S -> C C", nil)
	parser.AddParserRule("C -> c C", nil)
	parser.AddParserRule("C -> d", nil)

	// parser.generateParser()

	// for _, row := range parser.parsingTable {
	// 	for _, elem := range row {
	// 		fmt.Print(elem)
	// 		for i := 1; i <= 10-len(elem); i++ {
	// 			fmt.Print(" ")
	// 		}
	// 	}
	// 	fmt.Println()
	// }

	properStrings := []string{"dd", "cdd", "cccccccdd", "cdcd", "cccdcccd", "cdccccd"}

	for _, s := range properStrings {
		lexer.SetInputString(s)
		err := parser.Parse()
		if err != nil {
			t.Fatalf("Parsing failed for string: " + s)
		}
	}

	improperStrings := []string{"c", "cd", "cdcdc", "ddd"}

	for _, s := range improperStrings {
		lexer.SetInputString(s)
		err := parser.Parse()
		if err == nil {
			t.Fatalf("Parsing should have failed for string: " + s)
		}
	}

}
