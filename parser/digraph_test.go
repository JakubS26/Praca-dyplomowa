package parser

import (
	"testing"
)

func TestSimpleSetUnion(t *testing.T) {

	set1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	set2 := []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

	set3 := simpleSetUnion(set1, set2)

	if len(set3) != 13 {
		t.Fatalf("Wrong number of elements in result set")
	}

	check := checkElements(set3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13})

	if !check {
		t.Fatalf("Not all elements are present in result set")
	}

}

func TestDigraph1(t *testing.T) {

	predefinedSets := map[stateSymbolPair][]int{
		{0, 0}: {0},
		{1, 1}: {1},
		{2, 2}: {2},
		{3, 3}: {3},
		{4, 4}: {3, 4},
		{5, 5}: {3, 5},
	}

	relation := map[stateSymbolPair][]stateSymbolPair{
		{0, 0}: {{1, 1}, {2, 2}},
		{1, 1}: {{3, 3}},
		{2, 2}: {{4, 4}, {5, 5}},
	}

	minNonTerminalIndex := 0
	maxNonterminalIndex := 5
	numberOfStates := 6

	result := digraphAlgorithm(predefinedSets, relation,
		minNonTerminalIndex, maxNonterminalIndex, numberOfStates)

	_ = result

	// for key, value := range result {
	// 	if len(value) > 0 {
	// 		fmt.Println(key)
	// 		fmt.Print(value, "\n\n")
	// 	}
	// }

}
