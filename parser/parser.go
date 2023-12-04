package parser

import (
	"errors"
	"fmt"
	"goparser/lexer"
	"strings"
	"unicode"
)

//Terminale oraz nietermiale będą reprezentowane liczbami naturnalymi
//(np. 0-10 terminale (te same co w lekserze) 11-14 nieterminale)

// Inna nazwa: Stack item
type object struct {
	id    int
	Value any
}

func (o *object) setValue(s any) {
	o.Value = s
}

type parserRule struct {
	leftHandSide  int
	rightHandSide []int
	action        func([]any)
}

// Funkcja tylko do tesów to do celów debugowania
func getSymbolName(id int) string {

	for name, index := range lexer.GetTokenNames() {
		if index == id {
			return name
		}
	}

	for name, index := range nonTerminalNames {
		if index == id {
			return name
		}
	}

	if id == -1 {
		return "S'"
	}

	if id == len(lexer.GetTokenNames()) {
		return "$"
	}

	return "Unknown symbol!"
}

func createParserRule(leftHandSide int, rightHandSide []int, action func([]any)) parserRule {
	return parserRule{leftHandSide, rightHandSide, action}
}

func (p parserRule) getRightHandSideLength() int {
	return len(p.rightHandSide)
}

func (p parserRule) getRightHandSideSymbol(index int) int {
	return p.rightHandSide[index]
}

func (p parserRule) getRightHandSide() []int {
	return p.rightHandSide
}

func (p parserRule) getLeftHandSideSymbol() int {
	return p.leftHandSide
}

var actionStack Stack[object]

func checkNonterminalName(s string) bool {

	for _, c := range s {
		if !(c == '_' || (unicode.IsLetter(c) && unicode.IsUpper(c))) {
			return false
		}
	}

	return true

}

var nonTerminalNames map[string]int = make(map[string]int)

func getNumberOfGrammarSymbols() int {
	return len(lexer.GetTokenNames()) + len(nonTerminalNames) + 1
}

func toParserRule(s string, tokenNames map[string]int, action func([]any)) (parserRule, error) {

	splitStrings := strings.Split(s, " ")
	splitStringsClear := make([]string, 0, 5)

	leftHandSide := 0
	rightHandSide := make([]int, 0)

	nextFreeId := len(tokenNames) + len(nonTerminalNames) + 1

	for _, split := range splitStrings {
		if split != "" {
			splitStringsClear = append(splitStringsClear, split)
			//fmt.Printf("%#v\n", split)
		}
	}

	if len(splitStringsClear) < 3 {
		return parserRule{}, errors.New("This is not a valid parser rule.")
	}

	if splitStringsClear[1] != "->" {
		return parserRule{}, errors.New("This is not a valid parser rule.")
	}

	//Rozpatrujemy najpierw oddzielnie symbol z lewej strony produkcji

	if !checkNonterminalName(splitStringsClear[0]) {
		return parserRule{}, errors.New(fmt.Sprintf("Wrong nonterminal symbol name : %q. Names of nonterminals can contain only capital letters and underscores!", splitStringsClear[0]))
	} else {
		id, foundNonTerminal := nonTerminalNames[splitStringsClear[0]]

		if foundNonTerminal {
			leftHandSide = id
		} else {
			nonTerminalNames[splitStringsClear[0]] = nextFreeId
			leftHandSide = nextFreeId
			nextFreeId++
		}

	}

	for index := 2; index < len(splitStringsClear); index++ {

		str := splitStringsClear[index]
		id, foundTerminal := tokenNames[str]

		// Przypadek 0. - Dany string jest równy "epsilon" (napis pusty)

		if str == "epsilon" {
			continue
		}

		//Przypadek 1. - Dany string został odnaleziony w tablicy z nazwami tokenów
		//(czyli jest terminalem)

		if foundTerminal {
			rightHandSide = append(rightHandSide, id)
			continue
		}

		id, foundNonTerminal := nonTerminalNames[str]

		//Przypadek 2. - Dany string został odnaleziony w tablicy z nazwami symboli nieterminalnych
		//(jest on nieterminalem, który został już wcześniej napotkany)

		if foundNonTerminal {
			rightHandSide = append(rightHandSide, id)
			continue
		}

		//Przypadek 3. - Dany string nie został odnaleziony w żadnej z tablic
		//(musi być to nieterminal, którego jeszcze nie napotkaliśmy)

		if checkNonterminalName(str) {
			nonTerminalNames[str] = nextFreeId
			rightHandSide = append(rightHandSide, nextFreeId)
			nextFreeId++
		} else {
			return parserRule{}, errors.New(fmt.Sprintf("Wrong nonterminal symbol name : %q. Names of nonterminals can contain only capital letters and underscores!", str))
		}

	}

	return parserRule{leftHandSide, rightHandSide, action}, nil
}

var rules []parserRule

func getParserRules() []parserRule {
	return rules
}

// Zwraca pierwszy indeks (liczbę), jaki został nadany symbolowi nieterminalnemu.
// Jest to również (zgodnie z konwencją przyjętą w tym programie) indeks symbolu
// startowego wprowadzonej przez użytkownika gramatyki.
func getMinimalNonTerminalIndex() int {
	return len(lexer.GetTokenNames()) + 1
}

func getEndOfInputSymbolId() int {
	return len(lexer.GetTokenNames())
}

func AddParserRule(s string, action func([]any)) error {

	result, err := toParserRule(s, lexer.GetTokenNames(), action)

	if err == nil {
		rules = append(rules, result)
	}

	return err
}

func setParseTable(pt [][]string) {
	parsingTable = pt
}

var parsingTable [][]string
var parsingError error = nil

func RaiseError(err error) {
	parsingError = err
}

var tablesGenerated bool = false
