package parser

type stateSymbolPair struct {
	state  int
	symbol int
}

func (p *Parser) generateDrSets() map[stateSymbolPair][]int {
	result := make(map[stateSymbolPair][]int)

	//Przeszukujemy wszystkie możliwe przejścia z danego stanu
	for state, edges := range p.transitions {
		for _, edge := range edges {

			drSet := make([]int, 0)

			//Napotkano przejście z aktualnego stanu do innego stanu z symbolem nieterminalnym
			if edge.symbol >= p.getMinimalNonTerminalIndex() {

				for _, nextEdge := range p.transitions[edge.destState] {
					if !isNonTerminal(nextEdge.symbol, p.getMinimalNonTerminalIndex()) {
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
