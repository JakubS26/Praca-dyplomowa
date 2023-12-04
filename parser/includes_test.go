package parser

import (
	"testing"
)

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

	checkNonterminal := func(id int) bool {
		if id <= 7 && id >= 4 {
			return true
		}
		return false
	}

	result := generateIncludesRelation(automatonTransitions, nullableSymbols, productions, checkNonterminal)

	_ = result
}
