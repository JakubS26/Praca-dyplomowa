package parser

import (
	"fmt"
	"goparser/lexer"
	"testing"
)

func TestExample1(t *testing.T) {

	p := NewParser(lexer.NewLexer())

	id := map[string]int{
		"S'": -1,
		"a":  0,
		"b":  1,
		"c":  2,
		"$":  3,
		"S":  4,
		"A":  5,
		"B":  6,
		"C":  7,
	}

	name := map[int]string{
		-1: "S'",
		0:  "a",
		1:  "b",
		2:  "c",
		3:  "$",
		4:  "S",
		5:  "A",
		6:  "B",
		7:  "C",
	}

	p.endOfInputSymbolId = 3
	p.minimalNonTerminalIndex = 4
	p.numberOfGrammarSymbols = 8

	p.rules = []parserRule{
		createParserRule(id["S"], []int{id["A"], id["B"], id["C"]}, nil),
		createParserRule(id["A"], []int{id["a"], id["A"]}, nil),
		createParserRule(id["A"], []int{}, nil),
		createParserRule(id["B"], []int{id["b"], id["B"]}, nil),
		createParserRule(id["B"], []int{}, nil),
		createParserRule(id["C"], []int{id["c"], id["C"]}, nil),
		createParserRule(id["C"], []int{id["c"]}, nil),
		createParserRule(id["S'"], []int{id["S"], id["$"]}, nil),
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
		{
			createAutomatonTransition(4, 8, id["c"]),
			createAutomatonTransition(4, 7, id["C"]),
		},
		{
			createAutomatonTransition(5, 5, id["b"]),
			createAutomatonTransition(5, 9, id["B"]),
		},
		{},
		{},
		{
			createAutomatonTransition(8, 8, id["c"]),
			createAutomatonTransition(8, 10, id["C"]),
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

	follow := digraphAlgorithm(read, includes, p.minimalNonTerminalIndex, p.numberOfGrammarSymbols-1, len(p.transitions))

	fmt.Println("Follow sets:")
	for pair, set := range follow {
		fmt.Println(pair.state, name[pair.symbol], ":")
		for _, elem := range set {
			fmt.Print(name[elem])
		}
		fmt.Print("\n\n")
	}

	lookback := p.generateLookbackRelation()

	fmt.Println("lookback relation:")
	for pair, set := range lookback {
		fmt.Println("stan:", pair.state, "produkcja:", pair.productionId, ":")
		for _, elem := range set {
			fmt.Print(elem.state, " ", name[elem.symbol], "\n")
		}
		fmt.Print("\n\n")
	}

	lookahead := generateLookaheadSets(lookback, follow)

	fmt.Println("LA sets:")
	for pair, set := range lookahead {
		fmt.Println("stan:", pair.state, "produkcja:", pair.productionId, ":")
		for _, elem := range set {
			fmt.Print(name[elem], "\n")
		}
		fmt.Print("\n\n")
	}

}
