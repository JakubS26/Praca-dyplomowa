package parsergen_test

import (
	"goparser/parser"
	"goparser/parsergen"
	"testing"
)

func TestIsNullable1(t *testing.T) {

	rules := make([]parser.ParserRule, 0)

	// A -> B C D
	rules = append(rules, parser.CreateParserRule(0, []int{1, 2, 3}, nil))
	// B -> E F
	rules = append(rules, parser.CreateParserRule(1, []int{4, 5}, nil))
	// E -> D x
	rules = append(rules, parser.CreateParserRule(4, []int{3, 6}, nil))
	// F -> epsilon
	rules = append(rules, parser.CreateParserRule(5, []int{}, nil))
	// C -> D D D
	rules = append(rules, parser.CreateParserRule(2, []int{3, 3, 3}, nil))
	// D -> epsilon
	rules = append(rules, parser.CreateParserRule(3, []int{}, nil))
	// G -> D D x D D
	rules = append(rules, parser.CreateParserRule(7, []int{3, 3, 6, 3, 3}, nil))

	result := parsergen.FindNullable(rules)

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

	rules := make([]parser.ParserRule, 0)

	// A -> B B B B B
	rules = append(rules, parser.CreateParserRule(0, []int{1, 1, 1, 1, 1}, nil))
	// B -> C C
	rules = append(rules, parser.CreateParserRule(1, []int{2, 2}, nil))
	// C -> epsilon
	rules = append(rules, parser.CreateParserRule(2, []int{}, nil))

	result := parsergen.FindNullable(rules)

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
