package parsergen

import (
	"fmt"
	"goparser/lexer"
	"goparser/parser"
	"testing"
)

func TestFull(t *testing.T) {

	lexer.AddTokenDefinition("C_T", `C`)
	lexer.AddTokenDefinition("D_T", `D`)

	lexer.Init()

	parser.AddParserRule("START -> C_NT C_NT", nil)
	parser.AddParserRule("C_NT -> C_T C_NT", nil)
	parser.AddParserRule("C_NT -> D_T", nil)

	C := CreateLr0ItemSets()

	for _, set := range C {
		Print(set)
		fmt.Printf("\n")
	}

	transitions := GetTransitions()

	for _, x := range transitions {
		for _, y := range x {
			fmt.Println(y.GetSourceState(), "  ", parser.GetSymbolName(y.GetSymbol()), "  ", y.GetDestState())
		}
	}

	// Wyznaczamy zbiory DR

	drSets := GenerateDrSets()

	// Wyznaczamy zbiór terminali, z których można wyprowadzić słowo puste

	nullableSymbols := FindNullable(parser.GetParserRules())

	// Wyznaczamy relację reads

	readsRelation := generateReadsRelation(transitions, nullableSymbols)

	// Za pomocą relacji reads i zbiorów DR wyznaczamy zbiory Read

	readSets := digraphAlgorithm(drSets, readsRelation, parser.GetMinimalNonTerminalIndex(), parser.GetNumberOfGrammarSymbols()-1, len(transitions))

	// Wyznaczamy relację includes

	nonterminalCheck := func(id int) bool {
		if id >= parser.GetMinimalNonTerminalIndex() && id <= parser.GetNumberOfGrammarSymbols()-1 {
			return true
		}
		return false
	}

	includesRelation := generateIncludesRelation(transitions, nullableSymbols, parser.GetParserRules(), nonterminalCheck)

	// Za pomocą relacji includes i zbiorów Read wyznaczamy zbiory Follow

	followSets := digraphAlgorithm(readSets, includesRelation, parser.GetMinimalNonTerminalIndex(), parser.GetNumberOfGrammarSymbols()-1, len(transitions))

	// Wyznaczamy relację lookback

	lookbackRelation := generateLookbackRelation(transitions, parser.GetParserRules())

	// Za pomocą realcji lookback oraz zbiorów Follow wyznaczamy zbiory LA

	lookaheadSets := generateLookaheadSets(lookbackRelation, followSets)

	for key, value := range lookaheadSets {
		fmt.Println("State:", key.state, "Rule number:", key.productionId)
		for _, symbol := range value {
			fmt.Print(parser.GetSymbolName(symbol))
		}
	}

	// Za pomocą zbiorów podglądów (LA) wyznaczamy tabele parsowania

}
