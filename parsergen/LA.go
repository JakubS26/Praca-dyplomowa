package parsergen

func generateLookaheadSets(lookbackRelation map[stateProductionPair][]stateSymbolPair,
	followSets map[stateSymbolPair][]int) map[stateProductionPair][]int {

	result := make(map[stateProductionPair][]int)

	for key, value := range lookbackRelation {

		lookaheadSet := []int{}

		for _, r := range value {
			lookaheadSet = simpleSetUnion(lookaheadSet, followSets[r])
		}

		result[key] = lookaheadSet

	}

	return result
}
