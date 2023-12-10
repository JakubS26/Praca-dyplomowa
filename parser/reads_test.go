package parser

import (
	"testing"
)

func TestIsNullable1(t *testing.T) {

	rules := make([]parserRule, 0)

	// A -> B C D
	rules = append(rules, createParserRule(0, []int{1, 2, 3}, nil))
	// B -> E F
	rules = append(rules, createParserRule(1, []int{4, 5}, nil))
	// E -> D x
	rules = append(rules, createParserRule(4, []int{3, 6}, nil))
	// F -> epsilon
	rules = append(rules, createParserRule(5, []int{}, nil))
	// C -> D D D
	rules = append(rules, createParserRule(2, []int{3, 3, 3}, nil))
	// D -> epsilon
	rules = append(rules, createParserRule(3, []int{}, nil))
	// G -> D D x D D
	rules = append(rules, createParserRule(7, []int{3, 3, 6, 3, 3}, nil))

	p := Parser{rules: rules}

	p.findNullable()

	result := p.nullableSymbols

	_, ok := result[2]
	if !ok {
		t.Fatalf("Symbol number 2 should be classified as nullable")
	}
	_, ok = result[3]
	if !ok {
		t.Fatalf("Symbol number 3 should be classified as nullable")
	}
	_, ok = result[5]
	if !ok {
		t.Fatalf("Symbol number 5 should be classified as nullable")
	}

	if len(result) > 3 {
		t.Fatalf("Too many symbols were classified as nullable")
	}

}

func TestIsNullable2(t *testing.T) {

	rules := make([]parserRule, 0)

	// A -> B B B B B
	rules = append(rules, createParserRule(0, []int{1, 1, 1, 1, 1}, nil))
	// B -> C C
	rules = append(rules, createParserRule(1, []int{2, 2}, nil))
	// C -> epsilon
	rules = append(rules, createParserRule(2, []int{}, nil))

	p := Parser{rules: rules}

	p.findNullable()

	result := p.nullableSymbols

	_, ok := result[0]
	if !ok {
		t.Fatalf("Symbol number 0 should be classified as nullable")
	}
	_, ok = result[1]
	if !ok {
		t.Fatalf("Symbol number 1 should be classified as nullable")
	}
	_, ok = result[2]
	if !ok {
		t.Fatalf("Symbol number 2 should be classified as nullable")
	}

}

func TestReadsRelation(t *testing.T) {

	testTransitions := make([][]automatonTransition, 7)

	testTransitions[0] = append(testTransitions[0], automatonTransition{0, 1, 0})
	testTransitions[0] = append(testTransitions[0], automatonTransition{0, 2, 1})
	testTransitions[0] = append(testTransitions[0], automatonTransition{0, 3, 2})

	testTransitions[1] = append(testTransitions[1], automatonTransition{1, 4, 3})

	testTransitions[2] = append(testTransitions[2], automatonTransition{2, 5, 4})

	testTransitions[3] = append(testTransitions[3], automatonTransition{3, 6, 5})

	testTransitions[4] = append(testTransitions[4], automatonTransition{4, 5, 6})

	nullableSymbols := map[int]struct{}{
		4: {},
		5: {},
		6: {},
	}

	p := Parser{transitions: testTransitions, nullableSymbols: nullableSymbols}

	readsRelation := p.generateReadsRelation()

	if len(readsRelation) > 3 {
		t.Fatalf("Too many pairs in reads relation were found")
	}

	value, ok := readsRelation[stateSymbolPair{0, 1}]
	if !ok || len(value) != 1 || value[0].state != 2 || value[0].symbol != 4 {
		t.Fatalf("The reads relation for pair (0, 1) is not correctly determined")
	}

	value, ok = readsRelation[stateSymbolPair{1, 3}]
	if !ok || len(value) != 1 || value[0].state != 4 || value[0].symbol != 6 {
		t.Fatalf("The reads relation for pair (1, 3) has not been determined correctly")
	}

	value, ok = readsRelation[stateSymbolPair{0, 2}]
	if !ok || len(value) != 1 || value[0].state != 3 || value[0].symbol != 5 {
		t.Fatalf("The reads relation for pair (0, 2) is not correctly determined")
	}

}
