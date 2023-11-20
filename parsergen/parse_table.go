package parsergen

import "goparser/parser"

func GenerateLalrParseTables(automatonTransitions [][]automatonTransition,
	lookaheadSets map[stateProductionPair][]int, rules []parser.ParserRule) [][]string {
	return nil
}
