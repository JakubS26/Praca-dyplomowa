package parser

// Zwraca stan, w którym znajdziemy się, gdy wczytamy dany ciąg symboli z obecnego stanu
// (Jeśli brak takiej ścieżki w automacie, zwraca -1)
func (p *Parser) readSymbolsFromState(state int, symbols []int) int {

	for _, symbol := range symbols {

		transitionsFromCurrentState := p.transitions[state]

		found := false

		for _, transition := range transitionsFromCurrentState {
			if transition.symbol == symbol {
				state = transition.destState
				found = true
			}
		}

		if !found {
			return -1
		}

	}

	return state
}

func (p *Parser) generateIncludesRelation() map[stateSymbolPair][]stateSymbolPair {

	result := make(map[stateSymbolPair][]stateSymbolPair)

	// Przeglądamy po kolei wszystkie produkcje
	for _, rule := range p.rules {

		leftSymbol := rule.getLeftHandSideSymbol()
		rightSymbols := rule.getRightHandSide()

		// Przeglądamy wszystkie symbole po prawej stronie produkcji i szukamy nieterminali
		// Jeśli trafimy na nieterminal, sprawdzamy, czy są spełnione warunki relacji includes

		for index, symbol := range rightSymbols {
			if symbol >= p.getMinimalNonTerminalIndex() {

				beta := rightSymbols[0:index]
				gamma := rightSymbols[index+1:]

				// Sprawdzamy warunek 1: łańcuch znaków gamma musi składać się z samych symboli, z których
				// możemy wyprowadzić epsilon

				isGammaNullable := true

				for _, g := range gamma {
					_, ok := p.nullableSymbols[g]
					isGammaNullable = isGammaNullable && ok
				}

				if !isGammaNullable {
					continue
				}

				// Sprawdzamy warunek 2: między którymi stanami możemy przejść, wczytując ciąg symboli beta
				// W tym celu sprawdzamy po kolei wszystkie stany

				numberOfStates := len(p.transitions)

				for state := 0; state < numberOfStates; state++ {

					finalState := p.readSymbolsFromState(state, beta)

					// Jeśli ścieżka odpowiadająca danemu ciągowi przejść istnieje w automacie
					if finalState != -1 {

						// Sprawdzamy czy dana para ma już przypisany swój wycinek
						if result[stateSymbolPair{finalState, symbol}] == nil {
							result[stateSymbolPair{finalState, symbol}] = make([]stateSymbolPair, 0)
						}

						// Do wycinka odpowiadającego danej parze dodajemy parę, z którą jest ona w relacji includes
						// (p, A) includes (p', B)
						result[stateSymbolPair{finalState, symbol}] = append(result[stateSymbolPair{finalState, symbol}], stateSymbolPair{state, leftSymbol})
					}

				}

			}
		}

	}

	return result

}
