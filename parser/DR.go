package parser

type stateSymbolPair struct {
	state  int
	symbol int
}

func GenerateDrSets(minNonterminalSymbol int) map[stateSymbolPair][]int {
	result := make(map[stateSymbolPair][]int)

	//Przeszukujemy wszystkie możliwe przejścia z danego stanu
	for state, edges := range transitions {
		for _, edge := range edges {

			drSet := make([]int, 0)

			//fmt.Println("MTI: ", parser.GetMinimalNonTerminalIndex())
			//fmt.Println("SYMBOL: ", parser.GetSymbolName(edge.symbol), "IS NONTERMINAL: ", isNonTerminal(edge.symbol))

			//Napotkano przejście z aktualnego stanu do innego stanu z symbolem nieterminalnym
			if edge.symbol >= minNonterminalSymbol {

				for _, nextEdge := range transitions[edge.destState] {
					if !isNonTerminal(nextEdge.symbol) {
						drSet = append(drSet, nextEdge.symbol)
					}
				}

			}

			if len(drSet) != 0 {
				result[stateSymbolPair{state, edge.symbol}] = drSet

				//* To do wywalenia po zakończeniu testów
				// fmt.Println()
				// fmt.Println("STATE: ", state)
				// fmt.Println("SYMBOL: ", parser.GetSymbolName(edge.symbol))
				// for _, i := range drSet {
				// 	fmt.Println(parser.GetSymbolName(i))
				// }
				// fmt.Println()
				//*************************
			}

		}
	}

	return result
}
