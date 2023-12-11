package parser

import (
	"sort"
)

// Zgodnie z sugestiami z "Kompilatory. Reguły..." sytuację LR(0) przedstawiamy jako parę liczb,
// z których pierwsza odnosi się do numeru produkcji, której dotyczy sytuacja, zaś druga
// oznacza pozycję znacznika (kropki) w tej produkcji. Na przykład pozycja kropki 0 oznacza, że
// znajduje się ona przed elementem tablcy 0 (na samym początku produkcji).

func isNonTerminal(index int, minimalNonTerminalIndex int) bool {
	if index >= minimalNonTerminalIndex {
		return true
	} else {
		return false
	}
}

type lr0Item struct {
	ruleNumber     int
	markerLocation int
}

type lr0ItemSet = []lr0Item

type automatonTransition struct {
	sourceState int
	destState   int
	symbol      int
}

func createAutomatonTransition(source int, dest int, symbol int) automatonTransition {
	return automatonTransition{source, dest, symbol}
}

func (at automatonTransition) GetSourceState() int {
	return at.sourceState
}

func (at automatonTransition) GetDestState() int {
	return at.destState
}

func (at automatonTransition) GetSymbol() int {
	return at.symbol
}

func (I lr0Item) isComplete(rules []parserRule) bool {
	if rules[I.ruleNumber].getRightHandSideLength() == I.markerLocation {
		return true
	}
	return false
}

func (p *Parser) closure(I lr0ItemSet) lr0ItemSet {

	var J lr0ItemSet = make([]lr0Item, len(I))
	copy(J, I)

	var usedProductions []bool = make([]bool, len(p.rules))

	for i := 0; i < len(J); i++ {

		currentItem := J[i]

		if !currentItem.isComplete(p.rules) && p.rules[currentItem.ruleNumber].getRightHandSideSymbol(currentItem.markerLocation) >= p.getMinimalNonTerminalIndex() {

			nonterminal := p.rules[currentItem.ruleNumber].getRightHandSideSymbol(currentItem.markerLocation)

			for j, rule := range p.rules {
				if rule.getLeftHandSideSymbol() == nonterminal && !usedProductions[j] {
					J = append(J, lr0Item{j, 0})
					usedProductions[j] = true
				}
			}

		}

	}

	return J
}

func (p *Parser) gotoFunction(I lr0ItemSet, symbol int) lr0ItemSet {

	var J lr0ItemSet = make([]lr0Item, 0)

	for i := 0; i < len(I); i++ {

		currentItem := I[i]

		if !currentItem.isComplete(p.rules) && p.rules[currentItem.ruleNumber].getRightHandSideSymbol(currentItem.markerLocation) == symbol {
			J = append(J, lr0Item{currentItem.ruleNumber, currentItem.markerLocation + 1})
		}

	}

	return p.closure(J)

}

func isElement(I lr0ItemSet, C []lr0ItemSet) (bool, int) {

	for index, element := range C {

		if checkEqualLr0ItemSets(element, I) {
			return true, index
		}
	}

	return false, -1

}

func checkEqualLr0ItemSets(I1 lr0ItemSet, I2 lr0ItemSet) bool {

	if len(I1) != len(I2) {
		return false
	}

	sort.Slice(I1, func(i, j int) bool {
		if I1[i].ruleNumber != I1[j].ruleNumber {
			return I1[i].ruleNumber < I1[j].ruleNumber
		} else {
			return I1[i].markerLocation < I1[j].markerLocation
		}
	})

	sort.Slice(I2, func(i, j int) bool {
		if I2[i].ruleNumber != I2[j].ruleNumber {
			return I2[i].ruleNumber < I2[j].ruleNumber
		} else {
			return I2[i].markerLocation < I2[j].markerLocation
		}
	})

	for i := 0; i < len(I1); i++ {
		if I1[i] != I2[i] {
			return false
		}
	}

	return true
}

func (p *Parser) createLr0ItemSets() {

	// Uzupełniamy gramatykę o nowy symbol startowy (dodajemy regułę S' -> .S)

	p.rules = append(p.rules, createParserRule(-1, []int{p.getMinimalNonTerminalIndex(), p.getEndOfInputSymbolId()}, nil))

	var C []lr0ItemSet = make([][]lr0Item, 0)
	var firstItem lr0Item = lr0Item{len(p.rules) - 1, 0}

	// Inicjalizujemy zmienną do przechowywania przejść automatu LR(0) (krawędzi grafu automatu LR(0))

	transitions := make([][]automatonTransition, 1)

	// C - kolekcja zbiorów sytuacji LR(0)

	C = append(C, p.closure([]lr0Item{firstItem}))

	for i := 0; i < len(C); i++ {
		for j := 0; j < p.getNumberOfGrammarSymbols(); j++ {

			gotoResult := p.gotoFunction(C[i], j)

			isElem, index := isElement(gotoResult, C)

			if len(gotoResult) != 0 && !isElem {
				C = append(C, gotoResult)
				transitions = append(transitions, make([]automatonTransition, 0))
				transitions[i] = append(transitions[i], automatonTransition{i, len(C) - 1, j})
			} else if isElem {
				transitions[i] = append(transitions[i], automatonTransition{i, index, j})
			}

		}
	}

	p.lr0Sets = C
	p.transitions = transitions

}
