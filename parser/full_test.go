package parser

import (
	"fmt"
	"goparser/lexer"
	"testing"
)

func TestFull(t *testing.T) {

	lexer.AddTokenDefinition("c", `c`)
	lexer.AddTokenDefinition("d", `d`)

	lexer.Init()

	AddParserRule("S -> C C", nil)
	AddParserRule("C -> c C", nil)
	AddParserRule("C -> d", nil)

	C := CreateLr0ItemSets()
	_ = C

	for index, set := range C {
		fmt.Println(index)
		Print(set)
		fmt.Printf("\n")
	}

	transitions := GetTransitions()

	for _, x := range transitions {
		for _, y := range x {
			fmt.Println(y.GetSourceState(), "  ", GetSymbolName(y.GetSymbol()), "  ", y.GetDestState())
		}
	}

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

	_ = lookaheadSets

	for key, value := range lookaheadSets {
		fmt.Println("State:", key.state, "Rule number:", key.productionId)
		for _, symbol := range value {
			fmt.Println(GetSymbolName(symbol))
		}
	}

	// Za pomocą zbiorów podglądów (LA) wyznaczamy tabele parsowania

	result, _ := GenerateLalrParseTables(transitions, lookaheadSets, rules, C,
		GetEndOfInputSymbolId(), GetMinimalNonTerminalIndex(), GetNumberOfGrammarSymbols())

	fmt.Print(" ")
	for i := 0; i < GetNumberOfGrammarSymbols(); i++ {
		fmt.Printf("%5.5s", GetSymbolName(i))
	}
	fmt.Println()

	for index, row := range result {
		fmt.Print(index)
		for _, action := range row {
			fmt.Printf("%5.5s", action)
		}
		fmt.Println()
	}

	SetParseTable(result)

	sampleInput := "ccdcd"

	lexer.SetInputString(sampleInput)
	Parse()

}
