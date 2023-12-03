package parser

import (
	"fmt"
	"goparser/lexer"
	"testing"
)

func TestEps(t *testing.T) {

	lexer.AddTokenDefinition("a", `a`)
	lexer.AddTokenDefinition("b", `b`)
	lexer.AddTokenDefinition("c", `c`)

	lexer.Init()

	AddParserRule("S -> A B C", nil)
	AddParserRule("A -> a A", nil)
	AddParserRule("A -> epsilon", nil)
	AddParserRule("B -> b B", nil)
	AddParserRule("B -> epsilon", nil)
	AddParserRule("C -> c C", nil)
	AddParserRule("C -> epsilon", nil)

	C := CreateLr0ItemSets()

	transitions := GetTransitions()

	// Wyznaczamy zbiory DR

	drSets := GenerateDrSets(GetMinimalNonTerminalIndex())

	// Wyznaczamy zbiór terminali, z których można wyprowadzić słowo puste

	nullableSymbols := FindNullable(GetParserRules())

	// Wyznaczamy relację reads

	readsRelation := generateReadsRelation(transitions, nullableSymbols, GetMinimalNonTerminalIndex())

	// Za pomocą relacji reads i zbiorów DR wyznaczamy zbiory Read

	readSets := digraphAlgorithm(drSets, readsRelation, GetMinimalNonTerminalIndex(), GetNumberOfGrammarSymbols()-1, len(transitions))

	// Wyznaczamy relację includes

	nonterminalCheck := func(id int) bool {
		if id >= GetMinimalNonTerminalIndex() && id <= GetNumberOfGrammarSymbols()-1 {
			return true
		}
		return false
	}

	includesRelation := generateIncludesRelation(transitions, nullableSymbols, GetParserRules(), nonterminalCheck)

	// Za pomocą relacji includes i zbiorów Read wyznaczamy zbiory Follow

	followSets := digraphAlgorithm(readSets, includesRelation, GetMinimalNonTerminalIndex(), GetNumberOfGrammarSymbols()-1, len(transitions))

	// Wyznaczamy relację lookback

	lookbackRelation := generateLookbackRelation(transitions, GetParserRules())

	// Za pomocą realcji lookback oraz zbiorów Follow wyznaczamy zbiory LA

	lookaheadSets := generateLookaheadSets(lookbackRelation, followSets)

	// Za pomocą zbiorów podglądów (LA) wyznaczamy tabele parsowania

	result, _ := GenerateLalrParseTables(transitions, lookaheadSets, rules, C,
		GetEndOfInputSymbolId(), GetMinimalNonTerminalIndex(), GetNumberOfGrammarSymbols())

	fmt.Print(" ")
	for i := 0; i < GetNumberOfGrammarSymbols(); i++ {
		fmt.Printf("%6.6s", GetSymbolName(i))
	}
	fmt.Println()

	for index, row := range result {
		fmt.Print(index)
		for _, action := range row {
			fmt.Printf("%6.6s", action)
		}
		fmt.Println()
	}

	SetParseTable(result)

}
