package parsergen

import (
	"fmt"
	"goparser/parser"
	"testing"
)

func TestParseTables(t *testing.T) {

	// Gramatyka dla tego przykładu:
	//[0] S' -> S		-1 -> 3
	//[1] S -> CC		3 -> 4, 4
	//[2] C -> cC		4 -> 0, 4
	//[3] C -> d		4 -> 1

	productions := []parser.ParserRule{
		parser.CreateParserRule(-1, []int{3, 2}, nil),
		parser.CreateParserRule(3, []int{4, 4}, nil),
		parser.CreateParserRule(4, []int{0, 4}, nil),
		parser.CreateParserRule(4, []int{1}, nil),
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
			CreateAutomatonTransition(0, 1, 3),
			CreateAutomatonTransition(0, 2, 4),
			CreateAutomatonTransition(0, 3, 0),
			CreateAutomatonTransition(0, 4, 1),
		},
		{
			CreateAutomatonTransition(1, 7, 2),
		},
		{
			CreateAutomatonTransition(2, 3, 0),
			CreateAutomatonTransition(2, 4, 1),
			CreateAutomatonTransition(2, 5, 4),
		},
		{
			CreateAutomatonTransition(3, 3, 0),
			CreateAutomatonTransition(3, 4, 1),
			CreateAutomatonTransition(3, 6, 4),
		},
		{},
		{},
		{},
		{},
	}

	lookaheadSets := map[stateProductionPair][]int{
		//{1, 0}: {2},       //$
		{4, 3}: {0, 1, 2}, //c, d, $
		{6, 2}: {0, 1, 2}, //c, d, $
		{5, 1}: {2},       //$
	}

	result, _ := GenerateLalrParseTables(transitions, lookaheadSets, productions, lr0SetCollection, endOfInputSymbolIndex, startingSymbolIndex, numberOfSymbols)

	fmt.Print(" ")

	symNames := []string{"c", "d", "$", "S", "C"}

	_ = result
	_ = symNames

	for i := 0; i < 5; i++ {
		fmt.Printf("%5.5s", symNames[i])
	}
	fmt.Println()

	for index, row := range result {
		fmt.Print(index)
		for _, action := range row {
			fmt.Printf("%5.5s", action)
		}
		fmt.Println()
	}

}

func TestHello(t *testing.T) {
	fmt.Println("Hello, world!")
}
