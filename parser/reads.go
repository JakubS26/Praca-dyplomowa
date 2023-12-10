package parser

func (p *Parser) findNullable() {
	result := make(map[int]struct{})

	change := true

	for change {

		change = false

		for _, rule := range p.rules {

			_, alreadyChecked := result[rule.getLeftHandSideSymbol()]

			if alreadyChecked {
				continue
			}

			if rule.getRightHandSideLength() == 0 {
				result[rule.getLeftHandSideSymbol()] = struct{}{}
				change = true
			} else {
				rightHandSideSymbols := rule.getRightHandSide()
				checkAll := true

				for _, symbol := range rightHandSideSymbols {
					_, ok := result[symbol]
					checkAll = checkAll && ok
				}

				if checkAll {
					result[rule.getLeftHandSideSymbol()] = struct{}{}
					change = true
				}
			}

		}

	}

	p.nullableSymbols = result
}

func (p *Parser) generateReadsRelation() map[stateSymbolPair][]stateSymbolPair {

	result := make(map[stateSymbolPair][]stateSymbolPair)

	checkNonterminal := func(id int) bool {
		if id >= p.getMinimalNonTerminalIndex() {
			return true
		}
		return false
	}

	//Przeszukujemy wszystkie możliwe przejścia z kolejnych stanów
	for state, edges := range p.transitions {
		for _, edge := range edges {

			readsRelation := make([]stateSymbolPair, 0)

			//Napotkano przejście z aktualnego stanu do innego stanu z symbolem nieterminalnym
			if checkNonterminal(edge.symbol) {

				for _, nextEdge := range p.transitions[edge.destState] {
					_, isNullable := p.nullableSymbols[nextEdge.symbol]
					if isNullable {
						readsRelation = append(readsRelation, stateSymbolPair{edge.destState, nextEdge.symbol})
					}
				}

			}

			if len(readsRelation) != 0 {
				result[stateSymbolPair{state, edge.symbol}] = readsRelation
			}

		}
	}

	return result
}
