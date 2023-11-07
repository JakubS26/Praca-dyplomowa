package parsergen

import "goparser/parser"

// Zwraca stan, w którym znajdziemy się, gdy wczytamy dany ciąg symboli z obecnego stanu
// (Jeśli brak takiej ścieżki w automacie, zwraca -1)
func readSymbolsFromState(automatonTransitions [][]automatonTransition, state int, symbols []int) int {

	for _, symbol := range symbols {

		transitionsFromCurrentState := automatonTransitions[state]

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

func generateIncludesRelation(automatonTransitions [][]automatonTransition, nullableSymbols map[int]struct{},
	rules []parser.ParserRule, isNonterminalCheck func(int) bool) map[stateSymbolPair][]stateSymbolPair {

	result := make(map[stateSymbolPair][]stateSymbolPair)

	// Przeglądamy po kolei wszystkie produkcje
	for _, rule := range rules {

		leftSymbol := rule.GetLeftHandSideSymbol()
		rightSymbols := rule.GetRightHandSide()

		// Przeglądamy wszystkie symbole po prawej stronie produkcji i szukamy nieterminali
		// Jeśli trafimy na nieterminal, sprawdzamy, czy są spełnione warunki relacji includes

		for index, symbol := range rightSymbols {
			if isNonterminalCheck(symbol) {

				beta := rightSymbols[0:index]
				gamma := rightSymbols[index+1:]

				// Sprawdzamy warunek 1: łańcuch znaków gamma musi składać się z samych symboli, z których
				// możemy wyprowadzić epsilon

				isGammaNullable := true

				for _, g := range gamma {
					_, ok := nullableSymbols[g]
					isGammaNullable = isGammaNullable && ok
				}

				if !isGammaNullable {
					continue
				}

				// Sprawdzamy warunek 2: między którymi stanami możemy przejść, wczytując ciąg symboli beta
				// W tym celu sprawdzamy po kolei wszystkie stany

				numberOfStates := len(automatonTransitions)

				for state := 0; state < numberOfStates; state++ {

					finalState := readSymbolsFromState(automatonTransitions, state, beta)

					// Jeśli ścieżka odpowiadająca danemu ciągowi przejść istnieje w automacie
					if finalState != -1 {

						// Sprawdzamy czy dana para ma już przypisany swój wycinek
						if result[stateSymbolPair{symbol, state}] == nil {
							result[stateSymbolPair{symbol, state}] = make([]stateSymbolPair, 0)
						}

						// Do wycinka odpowiadającego danej parze dodajemy parę, z którą jest ona w relacji includes
						// (p, A) includes (p', B)
						result[stateSymbolPair{symbol, state}] = append(result[stateSymbolPair{symbol, state}], stateSymbolPair{finalState, leftSymbol})
					}

				}

			}
		}

	}

	return result

}
