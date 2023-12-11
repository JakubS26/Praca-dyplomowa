package parser

import (
	"goparser/lexer"
	"testing"
)

func isElementStateSymbolPair(s stateSymbolPair, setOfPairs []stateSymbolPair) bool {
	for _, elem := range setOfPairs {
		if elem == s {
			return true
		}
	}
	return false
}

func TestIncludes(t *testing.T) {

	id := map[byte]int{
		'a': 0,
		'b': 1,
		'c': 2,
		'A': 4,
		'B': 5,
		'C': 6,
		'E': 7,
	}

	automatonTransitions := [][]automatonTransition{
		{
			{0, 1, id['a']},
			{0, 2, id['b']},
			{0, 3, id['c']},
		},
		{
			{1, 4, id['b']},
			{1, 3, id['c']},
		},
		{
			{2, 3, id['c']},
		},
		{
			{3, 5, id['B']},
		},
		{
			{4, 5, id['C']},
		},
		{},
	}

	productions := []parserRule{
		createParserRule(id['A'], []int{id['a'], id['b'], id['C'], id['E']}, nil),
		createParserRule(id['A'], []int{id['c'], id['B'], id['C']}, nil),
		createParserRule(id['E'], []int{}, nil),
	}

	nullableSymbols := map[int]struct{}{
		id['E']: {},
	}

	p := NewParser(&lexer.Lexer{})

	p.rules = productions
	p.transitions = automatonTransitions
	p.nullableSymbols = nullableSymbols

	result := p.generateIncludesRelation()

	//Sprawdzenie wyniku dla pary (4, C)

	if len(result[stateSymbolPair{4, id['C']}]) != 1 {
		t.Fatalf("Includes relation is not correctly determined for pair %v!", stateSymbolPair{4, id['C']})
	}

	if !isElementStateSymbolPair(stateSymbolPair{0, id['A']}, result[stateSymbolPair{4, id['C']}]) {
		t.Fatalf("Includes relation is not correctly determined for pair %v!", stateSymbolPair{4, id['C']})
	}

	//Sprawdzenie wyniku dla pary (5, C)

	if len(result[stateSymbolPair{5, id['C']}]) != 3 {
		t.Fatalf("Includes relation is not correctly determined for pair %v!", stateSymbolPair{5, id['C']})
	}

	if !isElementStateSymbolPair(stateSymbolPair{0, id['A']}, result[stateSymbolPair{5, id['C']}]) {
		t.Fatalf("Includes relation is not correctly determined for pair %v!", stateSymbolPair{5, id['C']})
	}

	if !isElementStateSymbolPair(stateSymbolPair{1, id['A']}, result[stateSymbolPair{5, id['C']}]) {
		t.Fatalf("Includes relation is not correctly determined for pair %v!", stateSymbolPair{5, id['C']})
	}

	if !isElementStateSymbolPair(stateSymbolPair{2, id['A']}, result[stateSymbolPair{5, id['C']}]) {
		t.Fatalf("Includes relation is not correctly determined for pair %v!", stateSymbolPair{5, id['C']})
	}

	//Sprawdzenie wyniku dla pary (5, E)

	if len(result[stateSymbolPair{5, 7}]) != 1 {
		t.Fatalf("Includes relation is not correctly determined for pair %v!", stateSymbolPair{5, 7})
	}

	if !isElementStateSymbolPair(stateSymbolPair{0, id['A']}, result[stateSymbolPair{5, id['E']}]) {
		t.Fatalf("Includes relation is not correctly determined for pair %v!", stateSymbolPair{5, id['E']})
	}

	if len(result) != 3 {
		t.Fatalf("Includes relation is not correctly determined!")
	}

}
