package parser

type stateSymbolPair struct {
	state  int
	symbol int
}

func generateDrSets(minNonterminalSymbol int, autTransitions [][]automatonTransition) map[stateSymbolPair][]int {
	result := make(map[stateSymbolPair][]int)

	//Przeszukujemy wszystkie możliwe przejścia z danego stanu
	for state, edges := range autTransitions {
		for _, edge := range edges {

			drSet := make([]int, 0)

			//Napotkano przejście z aktualnego stanu do innego stanu z symbolem nieterminalnym
			if edge.symbol >= minNonterminalSymbol {

				for _, nextEdge := range autTransitions[edge.destState] {
					if !isNonTerminal(nextEdge.symbol) {
						drSet = append(drSet, nextEdge.symbol)
					}
				}

			}

			if len(drSet) != 0 {
				result[stateSymbolPair{state, edge.symbol}] = drSet
			}

		}
	}

	return result
}
