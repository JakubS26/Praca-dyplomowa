package parsergen

import (
	"fmt"
	"goparser/lexer"
	"goparser/parser"
	"strconv"
	"testing"
)

type T = []parser.Object

func TestCalc(t *testing.T) {

	lexer.AddTokenDefinition("NL", `\n`)
	lexer.AddTokenDefinition("NUM", `[0-9]+`)
	lexer.AddTokenDefinition("PLUS", `\+`)
	lexer.AddTokenDefinition("TIMES", `\*`)
	lexer.AddTokenDefinition("L_PAR", `\(`)
	lexer.AddTokenDefinition("R_PAR", `\)`)

	lexer.Init()

	parser.AddParserRule("S -> E NL", func(p T) { fmt.Printf("Wynik: %d\n", p[1].IntegerValue) })
	parser.AddParserRule("E -> E PLUS T", func(p T) { p[0].IntegerValue = p[1].IntegerValue + p[3].IntegerValue })
	parser.AddParserRule("E -> T", func(p T) { p[0].IntegerValue = p[1].IntegerValue })
	parser.AddParserRule("T -> T TIMES F", func(p T) { p[0].IntegerValue = p[1].IntegerValue * p[3].IntegerValue })
	parser.AddParserRule("T -> F", func(p T) { p[0].IntegerValue = p[1].IntegerValue })
	parser.AddParserRule("F -> L_PAR E R_PAR", func(p T) { p[0].IntegerValue = p[2].IntegerValue })
	parser.AddParserRule("F -> NUM", func(p T) { p[0].IntegerValue, _ = strconv.Atoi(p[1].GetStringValue()) })

	C := CreateLr0ItemSets()
	_ = C

	// for index, set := range C {
	// 	fmt.Println(index)
	// 	Print(set)
	// 	fmt.Printf("\n")
	// }

	transitions := GetTransitions()

	// for _, x := range transitions {
	// 	for _, y := range x {
	// 		fmt.Println(y.GetSourceState(), "  ", parser.GetSymbolName(y.GetSymbol()), "  ", y.GetDestState())
	// 	}
	// }

	// Wyznaczamy zbiory DR

	drSets := GenerateDrSets(parser.GetMinimalNonTerminalIndex())

	// Wyznaczamy zbiór terminali, z których można wyprowadzić słowo puste

	nullableSymbols := FindNullable(parser.GetParserRules())

	// Wyznaczamy relację reads

	readsRelation := generateReadsRelation(transitions, nullableSymbols, parser.GetMinimalNonTerminalIndex())

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

	_ = lookaheadSets

	// for key, value := range lookaheadSets {
	// 	fmt.Println("State:", key.state, "Rule number:", key.productionId)
	// 	for _, symbol := range value {
	// 		fmt.Println(parser.GetSymbolName(symbol))
	// 	}
	// }

	// Za pomocą zbiorów podglądów (LA) wyznaczamy tabele parsowania

	result, _ := GenerateLalrParseTables(transitions, lookaheadSets, rules, C,
		parser.GetEndOfInputSymbolId(), parser.GetMinimalNonTerminalIndex(), parser.GetNumberOfGrammarSymbols())

	fmt.Print(" ")
	for i := 0; i < parser.GetNumberOfGrammarSymbols(); i++ {
		fmt.Printf("%6.6s", parser.GetSymbolName(i))
	}
	fmt.Println()

	for index, row := range result {
		fmt.Print(index)
		for _, action := range row {
			fmt.Printf("%6.6s", action)
		}
		fmt.Println()
	}

	parser.SetParseTable(result)

	sampleInput := "333*(123+456)\n"
	lexer.SetInputString(sampleInput)
	parser.ParseWithSemanticActions()

	sampleInput = "(1+2)*(3+4)\n"
	lexer.SetInputString(sampleInput)
	parser.ParseWithSemanticActions()

	sampleInput = "1+0+1+0+1+0\n"
	lexer.SetInputString(sampleInput)
	parser.ParseWithSemanticActions()

}
