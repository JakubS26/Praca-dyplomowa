package parser

func generateParser() {

	C := createLr0ItemSets()
	_ = C

	transitions := getTransitions()

	// Wyznaczamy zbiory DR

	drSets := generateDrSets(getMinimalNonTerminalIndex(), transitions)

	// Wyznaczamy zbiór terminali, z których można wyprowadzić słowo puste

	nullableSymbols := findNullable(getParserRules())

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

	// Za pomocą zbiorów podglądów (LA) wyznaczamy tabele parsowania

	result, _ := generateLalrParseTables(transitions, lookaheadSets, rules, C,
		getEndOfInputSymbolId(), getMinimalNonTerminalIndex(), getNumberOfGrammarSymbols())

	tablesGenerated = true
	setParseTable(result)

}
