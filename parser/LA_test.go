package parser

import (
	"testing"
)

func checkElement(slice []int, elem int) bool {
	for _, value := range slice {
		if value == elem {
			return true
		}
	}
	return false
}

func checkElements(slice []int, elems []int) bool {
	result := true

	for _, elem := range elems {
		result = result && checkElement(slice, elem)
	}

	return result
}

func TestLookaheads1(t *testing.T) {

	lookbackRelation := map[stateProductionPair][]stateSymbolPair{
		{0, 0}: {{1, 1}, {1, 2}, {1, 3}, {1, 4}},
		{0, 1}: {{1, 1}, {1, 5}},
		{0, 2}: {{1, 2}, {1, 3}, {1, 4}, {1, 5}},
	}

	followSets := map[stateSymbolPair][]int{
		{1, 1}: {0, 1, 11},
		{1, 2}: {0, 2, 22},
		{1, 3}: {0, 3, 33},
		{1, 4}: {0, 4, 44},
		{1, 5}: {0, 5, 55},
	}

	result := generateLookaheadSets(lookbackRelation, followSets)

	if len(result[stateProductionPair{0, 0}]) != 9 {
		t.Fatalf("The lookahead relation for pair (0, 0) is not correctly determined")
	}

	if len(result[stateProductionPair{0, 1}]) != 5 {
		t.Fatalf("The lookahead relation for pair (0, 1) is not correctly determined")
	}

	if len(result[stateProductionPair{0, 2}]) != 9 {
		t.Fatalf("The lookahead relation for pair (0, 2) is not correctly determined")
	}

	if !checkElements(result[stateProductionPair{0, 0}], []int{0, 1, 2, 3, 4, 11, 22, 33, 44}) {
		t.Fatalf("The lookahead relation for pair (0, 0) is not correctly determined")
	}

	if !checkElements(result[stateProductionPair{0, 1}], []int{0, 1, 5, 11, 55}) {
		t.Fatalf("The lookahead relation for pair (0, 1) is not correctly determined")
	}

	if !checkElements(result[stateProductionPair{0, 2}], []int{0, 2, 3, 4, 5, 22, 33, 44, 55}) {
		t.Fatalf("The lookahead relation for pair (0, 2) is not correctly determined")
	}

}
