package parsergen

import (
	"fmt"
	"goparser/parser"
	"reflect"
)

// Zgodnie z sugestiami z "Kompilatory. Reguły..." sytuację LR(0) przedstawiamy jako parę liczb,
// z których pierwsza odnosi się do numeru produkcji, której dotyczy sytuacja, zaś druga
// oznacza pozycję znacznika (kropki) w tej produkcji. Na przykład pozycja kropki 0 oznacza, że
// znajduje się ona przed elementem tablcy 0 (na samym początku produkcji).

var rules []parser.ParserRule
var minimalNonTerminalIndex = parser.GetMinimalNonTerminalIndex()

func isNonTerminal(index int) bool {
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

var itemSets []lr0ItemSet
var numberOfSymbols int = parser.GetNumberOfGrammarSymbols()

func (I lr0Item) isComplete() bool {
	if rules[I.ruleNumber].GetRightHandSideLength() == I.markerLocation {
		return true
	}
	return false
}

func (I lr0Item) print() {

	name := parser.GetSymbolName(rules[I.ruleNumber].GetLeftHandSideSymbol())
	fmt.Print(name)

	fmt.Print(" -> ")

	for i := 0; i < rules[I.ruleNumber].GetRightHandSideLength(); i++ {

		if I.markerLocation == i {
			fmt.Print(" . ")
		}

		name = parser.GetSymbolName(rules[I.ruleNumber].GetRightHandSideSymbol(i))
		fmt.Print(name, " ")

	}

	if I.isComplete() {
		fmt.Print(" .")
	}

}

func Print(I lr0ItemSet) {
	for _, item := range I {
		item.print()
		fmt.Println("  ")
	}
}

func closure(I lr0ItemSet) lr0ItemSet {

	var J lr0ItemSet = make([]lr0Item, len(I))
	copy(J, I)

	var usedProductions []bool = make([]bool, len(rules))

	for i := 0; i < len(J); i++ {

		currentItem := J[i]

		if !currentItem.isComplete() && isNonTerminal(rules[currentItem.ruleNumber].GetRightHandSideSymbol(currentItem.markerLocation)) {

			nonterminal := rules[currentItem.ruleNumber].GetRightHandSideSymbol(currentItem.markerLocation)

			for j, rule := range rules {
				if rule.GetLeftHandSideSymbol() == nonterminal && !usedProductions[j] {
					J = append(J, lr0Item{j, 0})
					usedProductions[j] = true
				}
			}

		}

	}

	return J
}

func gotoFunction(I lr0ItemSet, symbol int) lr0ItemSet {

	var J lr0ItemSet = make([]lr0Item, 0)

	for i := 0; i < len(I); i++ {

		currentItem := I[i]

		if !currentItem.isComplete() && rules[currentItem.ruleNumber].GetRightHandSideSymbol(currentItem.markerLocation) == symbol {
			J = append(J, lr0Item{currentItem.ruleNumber, currentItem.markerLocation + 1})
		}

	}

	return closure(J)

}

func isElement(I lr0ItemSet, C []lr0ItemSet) bool {

	for _, element := range C {
		if reflect.DeepEqual(element, I) {
			return true
		}
	}

	return false

}

func CreateLr0ItemSets() []lr0ItemSet {

	// Uzupełniamy gramatykę o nowy symbol startowy (dodajemy regułę S' -> .S)

	rules = parser.GetParserRules()
	rules = append(rules, parser.CreateParserRule(-1, []int{parser.GetMinimalNonTerminalIndex()}, nil))

	var C []lr0ItemSet = make([][]lr0Item, 0)
	var firstItem lr0Item = lr0Item{len(rules) - 1, 0}

	C = append(C, closure([]lr0Item{firstItem}))

	for i := 0; i < len(C); i++ {
		for j := 0; j < parser.GetNumberOfGrammarSymbols(); j++ {

			gotoResult := gotoFunction(C[i], j)

			//fmt.Println("i = ", i, "; j = ", j, "; GOTO = ", gotoResult)

			if len(gotoResult) != 0 && !isElement(gotoResult, C) {
				C = append(C, gotoResult)
			}

		}
	}

	return C

}
