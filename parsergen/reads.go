package parsergen

import (
	"goparser/parser"
)

func FindNullable(rules []parser.ParserRule) map[int]struct{} {
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

func generateReadsRelation(automatonTransitions [][]automatonTransition, nullableSymbols map[int]struct{}) map[stateSymbolPair][]stateSymbolPair {

	result := make(map[stateSymbolPair][]stateSymbolPair)

	//Przeszukujemy wszystkie możliwe przejścia z kolejnych stanów
	for state, edges := range automatonTransitions {
		for _, edge := range edges {

			readsRelation := make([]stateSymbolPair, 0)

			//Napotkano przejście z aktualnego stanu do innego stanu z symbolem nieterminalnym
			if isNonTerminal(edge.symbol) {

				for _, nextEdge := range automatonTransitions[edge.destState] {
					_, isNullable := nullableSymbols[nextEdge.symbol]
					if isNullable {
						readsRelation = append(readsRelation, stateSymbolPair{edge.destState, nextEdge.symbol})
					}
				}

			}

			if len(readsRelation) != 0 {
				result[stateSymbolPair{state, edge.symbol}] = readsRelation

				// fmt.Println()
				// fmt.Println("STATE: ", state)
				// fmt.Println("SYMBOL: ", parser.GetSymbolName(edge.symbol))
				// for _, i := range readsRelation {
				// 	fmt.Println(i.state, "   ", parser.GetSymbolName(i.symbol))
				// }
				// fmt.Println()
			}

		}
	}

	return result
}
