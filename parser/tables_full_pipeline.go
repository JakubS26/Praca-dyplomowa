package parser

func (p *Parser) generateParser() {

	p.createLr0ItemSets()

	// Wyznaczamy zbiory DR

	drSets := p.generateDrSets()

	// Wyznaczamy zbiór terminali, z których można wyprowadzić słowo puste

	p.findNullable()

	// Wyznaczamy relację reads

	readsRelation := p.generateReadsRelation()

	// Za pomocą relacji reads i zbiorów DR wyznaczamy zbiory Read

	readSets := digraphAlgorithm(drSets, readsRelation, p.getMinimalNonTerminalIndex(), p.getNumberOfGrammarSymbols()-1, len(p.transitions))

	// Wyznaczamy relację includes

	includesRelation := p.generateIncludesRelation()

	// Za pomocą relacji includes i zbiorów Read wyznaczamy zbiory Follow

	followSets := digraphAlgorithm(readSets, includesRelation, p.getMinimalNonTerminalIndex(), p.getNumberOfGrammarSymbols()-1, len(p.transitions))

	// Wyznaczamy relację lookback

	lookbackRelation := p.generateLookbackRelation()

	// Za pomocą realcji lookback oraz zbiorów Follow wyznaczamy zbiory LA

	lookaheadSets := generateLookaheadSets(lookbackRelation, followSets)

	// Za pomocą zbiorów podglądów (LA) wyznaczamy tabele parsowania

	result, _ := p.generateLalrParseTables(lookaheadSets)

	//fmt.Println(result)

	p.parsingTable = result

}
