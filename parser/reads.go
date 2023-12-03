package parser

func FindNullable(rules []ParserRule) map[int]struct{} {
	result := make(map[int]struct{})

	change := true

	for change {

		change = false

		for _, rule := range rules {

			_, alreadyChecked := result[rule.GetLeftHandSideSymbol()]

			if alreadyChecked {
				continue
			}

			if rule.GetRightHandSideLength() == 0 {
				result[rule.GetLeftHandSideSymbol()] = struct{}{}
				change = true
			} else {
				rightHandSideSymbols := rule.GetRightHandSide()
				checkAll := true

				for _, symbol := range rightHandSideSymbols {
					_, ok := result[symbol]
					checkAll = checkAll && ok
				}

				if checkAll {
					result[rule.GetLeftHandSideSymbol()] = struct{}{}
					change = true
				}
			}

		}

	}

	return result
}

func generateReadsRelation(automatonTransitions [][]automatonTransition,
	nullableSymbols map[int]struct{}, minNonterminalId int) map[stateSymbolPair][]stateSymbolPair {

	result := make(map[stateSymbolPair][]stateSymbolPair)

	checkNonterminal := func(id int) bool {
		if id >= minNonterminalId {
			return true
		}
		return false
	}

	//Przeszukujemy wszystkie możliwe przejścia z kolejnych stanów
	for state, edges := range automatonTransitions {
		for _, edge := range edges {

			readsRelation := make([]stateSymbolPair, 0)

			//Napotkano przejście z aktualnego stanu do innego stanu z symbolem nieterminalnym
			if checkNonterminal(edge.symbol) {

				for _, nextEdge := range automatonTransitions[edge.destState] {
					_, isNullable := nullableSymbols[nextEdge.symbol]
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
