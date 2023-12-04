package parser

import (
	"errors"
	"strconv"
)

func generateLalrParseTables(automatonTransitions [][]automatonTransition,
	lookaheadSets map[stateProductionPair][]int, rules []parserRule, lr0ItemSetCollection [][]lr0Item,
	endOfInputSymbolId int, startingSymbolIndex int, numberOfGrammarSymbols int) ([][]string, error) {

	//augmentedStartingSymbolIndex := -1

	parseTable := make([][]string, len(automatonTransitions))

	for i := range parseTable {
		parseTable[i] = make([]string, numberOfGrammarSymbols)
	}

	// W opdowiednim stanie gdy widzimy symbol końca inputu, ustawiamy akcję akceptuj

	for index, itemSet := range lr0ItemSetCollection {
		for _, item := range itemSet {
			if rules[item.ruleNumber].getLeftHandSideSymbol() == -1 && item.markerLocation == 1 {
				parseTable[index][endOfInputSymbolId] = "a"
			}
		}
	}

	// W tabelach parsowania ustawiamy akcję "shift" przy odpowiednich przejściach pomiędzy stanami automatu dla terminali

	maxTerminalIndex := endOfInputSymbolId - 1

	for _, transitionsFromState := range automatonTransitions {
		for _, transition := range transitionsFromState {
			i := transition.sourceState
			j := transition.destState

			a := transition.symbol

			if a == endOfInputSymbolId {
				continue
			}

			//Sprawdzamy, czy w tym miejscu tabeli nie ma już przypisanej jakiejś akcji

			if parseTable[i][a] != "" {
				return nil, errors.New("This grammar is not LALR(1)!")
			}

			//Sprawdzamy, czy symbol a jest terminalem

			if a <= maxTerminalIndex {
				parseTable[i][a] = "s" + strconv.Itoa(j)
			}
		}

	}

	// Ustawiamy redukcje zgodnie z produkcjami oraz zbiorami podglądów

	for setIndex, lrOItemSet := range lr0ItemSetCollection {
		for _, lr0Item := range lrOItemSet {

			// Sprawdzamy, czy dana sytuacja LR(0) ma znacznik (kropkę) na końcu
			// oraz czy po lewej stronie produkcji nie mamy "dodatkowego" symbolu startowego S'

			currentRule := rules[lr0Item.ruleNumber]

			if lr0Item.markerLocation == currentRule.getRightHandSideLength() && currentRule.getLeftHandSideSymbol() != -1 {

				// Sprawdzamy, jaki mamy zbiór podglądów dla danej produkcji w obecnym stanie (numerem obecnego stanu jest setIndex)
				// Dla symboli a, które należą do zbioru podglądów ustawiamy akcja[setIndex][a] = redukuj zgodnie z regułą lr0Item.ruleNumber

				lookaheadSymbols := lookaheadSets[stateProductionPair{setIndex, lr0Item.ruleNumber}]

				for _, a := range lookaheadSymbols {

					if parseTable[setIndex][a] != "" {
						return nil, errors.New("This grammar is not LALR(1)!")
					}

					parseTable[setIndex][a] = "r" + strconv.Itoa(lr0Item.ruleNumber)
				}

			}

		}
	}

	// Ustawiamy "przejście" dla odpowiednich nieterminali

	for _, transitionsFromState := range automatonTransitions {
		for _, transition := range transitionsFromState {
			i := transition.sourceState
			j := transition.destState

			A := transition.symbol

			//Sprawdzamy, czy symbol A jest nieterminalem

			if A >= endOfInputSymbolId+1 {
				parseTable[i][A] = strconv.Itoa(j)
			}
		}

	}

	return parseTable, nil
}
