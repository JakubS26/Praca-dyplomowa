package parser

type stateProductionPair struct {
	state        int
	productionId int
}

func generateLookbackRelation(automatonTransitions [][]automatonTransition,
	rules []ParserRule) map[stateProductionPair][]stateSymbolPair {

	result := make(map[stateProductionPair][]stateSymbolPair)

	numberOfStates := len(automatonTransitions)

	for ruleIndex, rule := range rules {
		for stateP := 0; stateP < numberOfStates; stateP++ {

			A := rule.getLeftHandSideSymbol()
			omega := rule.getRightHandSide()
			stateQ := readSymbolsFromState(automatonTransitions, stateP, omega)

			if stateQ != -1 {

				if result[stateProductionPair{stateQ, ruleIndex}] == nil {
					result[stateProductionPair{stateQ, ruleIndex}] = make([]stateSymbolPair, 0)
				}

				result[stateProductionPair{stateQ, ruleIndex}] = append(result[stateProductionPair{stateQ, ruleIndex}], stateSymbolPair{stateP, A})

			}

		}
	}

	return result
}
