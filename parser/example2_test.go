package parser

import (
	"fmt"
	"goparser/lexer"
	"testing"
)

func TestExample2(t *testing.T) {

	p := NewParser(lexer.NewLexer())

	id := map[string]int{
		"a": 0,
		"b": 1,
		"$": 2,
		"S": 3,
		"A": 4,
		"B": 5,
	}

	name := map[int]string{
		0: "a",
		1: "b",
		2: "$",
		3: "S",
		4: "A",
		5: "B",
	}

	p.endOfInputSymbolId = 2
	p.minimalNonTerminalIndex = 3
	p.numberOfGrammarSymbols = 6

	p.rules = []parserRule{
		createParserRule(id["S"], []int{id["A"], id["B"]}, nil),
		createParserRule(id["A"], []int{id["a"], id["A"]}, nil),
		createParserRule(id["A"], []int{}, nil),
		createParserRule(id["B"], []int{id["b"], id["B"]}, nil),
		createParserRule(id["B"], []int{}, nil),
	}

	p.findNullable()

	p.transitions = [][]automatonTransition{
		{
			createAutomatonTransition(0, 1, id["S"]),
			createAutomatonTransition(0, 2, id["A"]),
			createAutomatonTransition(0, 3, id["a"]),
		},
		{
			createAutomatonTransition(1, 11, id["$"]),
		},
		{
			createAutomatonTransition(2, 5, id["b"]),
			createAutomatonTransition(2, 4, id["B"]),
		},
		{
			createAutomatonTransition(3, 3, id["a"]),
			createAutomatonTransition(3, 6, id["A"]),
		},
		{},
		{
			createAutomatonTransition(5, 5, id["b"]),
			createAutomatonTransition(5, 7, id["B"]),
		},
		{},
		{},
		{},
	}

	dr := p.generateDrSets()

	fmt.Println("DR sets:")
	for pair, set := range dr {
		fmt.Println(pair.state, name[pair.symbol], ":")
		for _, elem := range set {
			fmt.Print(name[elem])
		}
		fmt.Print("\n\n")
	}

	reads := p.generateReadsRelation()

	fmt.Println("reads relation:")
	for pair, set := range reads {
		fmt.Println(pair.state, name[pair.symbol], ":")
		for _, elem := range set {
			fmt.Print(elem.state, " ", name[elem.symbol], "\n")
		}
		fmt.Print("\n\n")
	}

	read := digraphAlgorithm(dr, reads, p.minimalNonTerminalIndex, p.numberOfGrammarSymbols-1, len(p.transitions))

	fmt.Println("Read sets:")
	for pair, set := range read {
		fmt.Println(pair.state, name[pair.symbol], ":")
		for _, elem := range set {
			fmt.Print(name[elem])
		}
		fmt.Print("\n\n")
	}

	includes := p.generateIncludesRelation()

	fmt.Println("includes relation:")
	for pair, set := range includes {
		fmt.Println(pair.state, name[pair.symbol], ":")
		for _, elem := range set {
			fmt.Print(elem.state, " ", name[elem.symbol], "\n")
		}
		fmt.Print("\n\n")
	}

}
