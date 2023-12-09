package parser

import (
	"math"
	"sort"
)

func minimum(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func setUnion(set1 []int, set2 []int) []int {

	sort.IntSlice.Sort(set1)

	result := make([]int, len(set1))
	copy(result, set1)

	for _, elem := range set2 {
		index := sort.SearchInts(set1, elem)
		if index >= len(set1) || set1[index] != elem {
			result = append(result, elem)
		}
	}

	return result

}

func digraphAlgorithm(predefinedSets map[stateSymbolPair][]int, relation map[stateSymbolPair][]stateSymbolPair,
	minNonterminalIndex, maxNonterminalIdex, numberOfStates int) map[stateSymbolPair][]int {

	var S Stack[stateSymbolPair]
	var N map[stateSymbolPair]int = make(map[stateSymbolPair]int, 0)
	var F map[stateSymbolPair][]int = make(map[stateSymbolPair][]int, 0)

	// Iterujemy po wszystkich możliwych kodach nieterminali oraz stanach automatu i
	// ustawiamy wartości początkowe - zera w mapie N oraz zbiory nil dla wartości obliczanej funkcji F
	for index := minNonterminalIndex; index <= maxNonterminalIdex; index++ {
		for state := 0; state < numberOfStates; state++ {
			pair := stateSymbolPair{state, index}

			N[pair] = 0
			F[pair] = nil
		}
	}

	var traverse func(stateSymbolPair)

	traverse = func(x stateSymbolPair) {
		S.Push(x)
		d := S.Size()
		N[x] = d

		F[x] = predefinedSets[x]

		pairsInRelation, _ := relation[x]

		for _, y := range pairsInRelation {
			if N[y] == 0 {
				traverse(y)
			}
			N[x] = minimum(N[x], N[y])
			F[x] = setUnion(F[x], F[y])
		}

		if N[x] == d {
			for true {

				topValue, _ := S.Top()
				N[topValue] = math.MaxInt
				F[topValue] = F[x]

				S.Pop()

				if topValue == x {
					break
				}
			}
		}

	}

	for key, value := range N {
		if value == 0 {
			traverse(key)
		}
	}

	return F

}
