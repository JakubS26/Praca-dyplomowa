package parser

import (
	"goparser/lexer"
	"testing"
)

func TestParseTables(t *testing.T) {

	// Gramatyka dla tego przykładu:
	//[0] S' -> S		-1 -> 3
	//[1] S -> CC		3 -> 4, 4
	//[2] C -> cC		4 -> 0, 4
	//[3] C -> d		4 -> 1

	productions := []parserRule{
		createParserRule(-1, []int{3, 2}, nil),
		createParserRule(3, []int{4, 4}, nil),
		createParserRule(4, []int{0, 4}, nil),
		createParserRule(4, []int{1}, nil),
	}

	I0 := []lr0Item{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
	}

	I1 := []lr0Item{
		{0, 1},
	}

	I2 := []lr0Item{
		{1, 1},
		{2, 0},
		{3, 0},
	}

	I3 := []lr0Item{
		{2, 1},
		{2, 0},
		{3, 0},
	}

	I4 := []lr0Item{
		{3, 1},
	}

	I5 := []lr0Item{
		{1, 2},
	}

	I6 := []lr0Item{
		{2, 2},
	}

	I7 := []lr0Item{
		{0, 2},
	}

	lr0SetCollection := []lr0ItemSet{I0, I1, I2, I3, I4, I5, I6, I7}

	endOfInputSymbolIndex := 2
	startingSymbolIndex := 3
	numberOfSymbols := 5

	transitions := [][]automatonTransition{
		{
			createAutomatonTransition(0, 1, 3),
			createAutomatonTransition(0, 2, 4),
			createAutomatonTransition(0, 3, 0),
			createAutomatonTransition(0, 4, 1),
		},
		{
			createAutomatonTransition(1, 7, 2),
		},
		{
			createAutomatonTransition(2, 3, 0),
			createAutomatonTransition(2, 4, 1),
			createAutomatonTransition(2, 5, 4),
		},
		{
			createAutomatonTransition(3, 3, 0),
			createAutomatonTransition(3, 4, 1),
			createAutomatonTransition(3, 6, 4),
		},
		{},
		{},
		{},
		{},
	}

	lookaheadSets := map[stateProductionPair][]int{
		{4, 3}: {0, 1, 2}, //c, d, $
		{6, 2}: {0, 1, 2}, //c, d, $
		{5, 1}: {2},       //$

	}

	p := NewParser(&lexer.Lexer{})

	p.transitions = transitions
	p.rules = productions
	p.lr0Sets = lr0SetCollection

	p.endOfInputSymbolId = endOfInputSymbolIndex
	p.minimalNonTerminalIndex = startingSymbolIndex
	p.numberOfGrammarSymbols = numberOfSymbols

	result, _ := p.generateLalrParseTables(lookaheadSets)

	ok := true

	ok = ok && result[0][0] == "s3"
	ok = ok && result[0][1] == "s4"
	ok = ok && result[0][2] == ""
	ok = ok && result[0][3] == "1"
	ok = ok && result[0][4] == "2"

	ok = ok && result[1][0] == ""
	ok = ok && result[1][1] == ""
	ok = ok && result[1][2] == "a"
	ok = ok && result[1][3] == ""
	ok = ok && result[1][4] == ""

	ok = ok && result[2][0] == "s3"
	ok = ok && result[2][1] == "s4"
	ok = ok && result[2][2] == ""
	ok = ok && result[2][3] == ""
	ok = ok && result[2][4] == "5"

	ok = ok && result[3][0] == "s3"
	ok = ok && result[3][1] == "s4"
	ok = ok && result[3][2] == ""
	ok = ok && result[3][3] == ""
	ok = ok && result[3][4] == "6"

	ok = ok && result[4][0] == "r3"
	ok = ok && result[4][1] == "r3"
	ok = ok && result[4][2] == "r3"
	ok = ok && result[4][3] == ""
	ok = ok && result[4][4] == ""

	ok = ok && result[5][0] == ""
	ok = ok && result[5][1] == ""
	ok = ok && result[5][2] == "r1"
	ok = ok && result[5][3] == ""
	ok = ok && result[5][4] == ""

	ok = ok && result[6][0] == "r2"
	ok = ok && result[6][1] == "r2"
	ok = ok && result[6][2] == "r2"
	ok = ok && result[6][3] == ""
	ok = ok && result[6][4] == ""

	ok = ok && result[7][0] == ""
	ok = ok && result[7][1] == ""
	ok = ok && result[7][2] == ""
	ok = ok && result[7][3] == ""
	ok = ok && result[7][4] == ""

	if !ok {
		t.Fatalf("Parse table is not correctly determined!")
	}

}
