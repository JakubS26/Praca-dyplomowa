package parser

import (
	"fmt"
	"sort"
	"testing"
)

func TestSetUnion(t *testing.T) {

	set1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	set2 := []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

	set3 := setUnion(set1, set2)

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

	expected := map[stateSymbolPair][]int{
		{0, 0}: {0, 1, 2, 3, 4, 5},
		{1, 1}: {1, 3},
		{2, 2}: {2, 3, 4, 5},
		{3, 3}: {3},
		{4, 4}: {3, 4},
		{5, 5}: {3, 5},
	}

	minNonTerminalIndex := 0
	maxNonterminalIndex := 5
	numberOfStates := 6

	result := digraphAlgorithm(predefinedSets, relation,
		minNonTerminalIndex, maxNonterminalIndex, numberOfStates)

	for key, value := range result {
		if key.state != key.symbol && len(value) != 0 {
			t.Fatalf(fmt.Sprintf("Value for key %v should be empty!", key))
		}
	}

	for i := 0; i <= 5; i++ {
		resultSlice := result[stateSymbolPair{i, i}]
		sort.Slice(resultSlice, func(i, j int) bool {
			return resultSlice[i] < resultSlice[j]
		})
		expectedSlice := expected[stateSymbolPair{i, i}]
		sort.Slice(expectedSlice, func(i, j int) bool {
			return expectedSlice[i] < expectedSlice[j]
		})

		if len(resultSlice) != len(expectedSlice) {
			t.Fatalf(fmt.Sprintf("Result set is not correct for %v", stateProductionPair{i, i}))
		}

		for j := 0; j < len(resultSlice); j++ {
			if resultSlice[j] != expectedSlice[j] {
				t.Fatalf(fmt.Sprintf("Result set is not correct for %v", stateProductionPair{i, i}))
			}
		}
	}

}
