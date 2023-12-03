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
			fmt.Println(y.GetSourceState(), "  ", getSymbolName(y.GetSymbol()), "  ", y.GetDestState())
		}
	}

	// Wyznaczamy zbiory DR

	drSets := GenerateDrSets(getMinimalNonTerminalIndex(), transitions)

	// Wyznaczamy zbiór terminali, z których można wyprowadzić słowo puste

	nullableSymbols := FindNullable(getParserRules())

	// Wyznaczamy relację reads

	readsRelation := generateReadsRelation(transitions, nullableSymbols, getMinimalNonTerminalIndex())

	// Za pomocą relacji reads i zbiorów DR wyznaczamy zbiory Read

	readSets := digraphAlgorithm(drSets, readsRelation, getMinimalNonTerminalIndex(), getNumberOfGrammarSymbols()-1, len(transitions))

	// Wyznaczamy relację includes

	nonterminalCheck := func(id int) bool {
		if id >= getMinimalNonTerminalIndex() && id <= getNumberOfGrammarSymbols()-1 {
			return true
		}
		return false
	}

	includesRelation := generateIncludesRelation(transitions, nullableSymbols, getParserRules(), nonterminalCheck)

	// Za pomocą relacji includes i zbiorów Read wyznaczamy zbiory Follow

	followSets := digraphAlgorithm(readSets, includesRelation, getMinimalNonTerminalIndex(), getNumberOfGrammarSymbols()-1, len(transitions))

	// Wyznaczamy relację lookback

	lookbackRelation := generateLookbackRelation(transitions, getParserRules())

	// Za pomocą realcji lookback oraz zbiorów Follow wyznaczamy zbiory LA

	lookaheadSets := generateLookaheadSets(lookbackRelation, followSets)

	_ = lookaheadSets

	for key, value := range lookaheadSets {
		fmt.Println("State:", key.state, "Rule number:", key.productionId)
		for _, symbol := range value {
			fmt.Println(getSymbolName(symbol))
		}
	}

	// Za pomocą zbiorów podglądów (LA) wyznaczamy tabele parsowania

	result, _ := GenerateLalrParseTables(transitions, lookaheadSets, rules, C,
		getEndOfInputSymbolId(), getMinimalNonTerminalIndex(), getNumberOfGrammarSymbols())

	fmt.Print(" ")
	for i := 0; i < getNumberOfGrammarSymbols(); i++ {
		fmt.Printf("%5.5s", getSymbolName(i))
	}
	fmt.Println()

	for index, row := range result {
		fmt.Print(index)
		for _, action := range row {
			fmt.Printf("%5.5s", action)
		}
		fmt.Println()
	}

	setParseTable(result)

	sampleInput := "ccdcd"

	lexer.SetInputString(sampleInput)
	Parse()

}
