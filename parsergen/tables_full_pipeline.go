package parsergen

import "goparser/parser"

func GenerateParser() {

	C := CreateLr0ItemSets()
	_ = C

	transitions := GetTransitions()

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

	// Za pomocą zbiorów podglądów (LA) wyznaczamy tabele parsowania

	result, _ := GenerateLalrParseTables(transitions, lookaheadSets, rules, C,
		parser.GetEndOfInputSymbolId(), parser.GetMinimalNonTerminalIndex(), parser.GetNumberOfGrammarSymbols())

	parser.SetParseTable(result)

}
