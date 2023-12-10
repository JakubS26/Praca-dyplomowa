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

type Parser struct {
	nonTerminalNames        map[string]int
	rules                   []parserRule
	parsingTable            [][]string
	parsingError            error
	tablesGenerated         bool
	transitions             [][]automatonTransition
	lr0Sets                 []lr0ItemSet
	lexer                   *lexer.Lexer
	nullableSymbols         map[int]struct{}
	endOfInputSymbolId      int
	minimalNonTerminalIndex int
	numberOfGrammarSymbols  int
	isInitialized           bool
}

func NewParser(lexer *lexer.Lexer) *Parser {
	return &Parser{
		nonTerminalNames:        make(map[string]int),
		rules:                   make([]parserRule, 0),
		parsingTable:            nil,
		tablesGenerated:         false,
		lexer:                   lexer,
		endOfInputSymbolId:      len(lexer.GetTokenNames()),
		minimalNonTerminalIndex: len(lexer.GetTokenNames()) + 1,
		numberOfGrammarSymbols:  len(lexer.GetTokenNames()) + 1,
	}
}

func (p *Parser) Init() error {

	if p.isInitialized {
		return nil
	}

	if len(p.rules) == 0 {
		return errors.New("The set of grammar rules cannot be empty!")
	}

	p.generateParser()

	p.isInitialized = true
	p.tablesGenerated = true

	return nil
}

// Funkcja tylko do tesów do celów debugowania
func (p *Parser) getSymbolName(id int) string {

	for name, index := range p.lexer.GetTokenNames() {
		if index == id {
			return name
		}
	}

	for name, index := range p.nonTerminalNames {
		if index == id {
			return name
		}
	}

	if id == -1 {
		return "S'"
	}

	if id == len(p.lexer.GetTokenNames()) {
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

//var actionStack Stack[object]

func checkNonterminalName(s string) bool {

	for _, c := range s {
		if !(c == '_' || (unicode.IsLetter(c))) {
			return false
		}
	}

	return true

}

func (p *Parser) toParserRule(s string, tokenNames map[string]int, action func([]any), nonTerminalNames map[string]int) (parserRule, error) {

	splitStrings := strings.Split(s, " ")
	splitStringsClear := make([]string, 0, 5)

	leftHandSide := 0
	rightHandSide := make([]int, 0)

	nextFreeId := len(tokenNames) + len(nonTerminalNames) + 1

	for _, split := range splitStrings {
		if split != "" {
			splitStringsClear = append(splitStringsClear, split)
		}
	}

	if len(splitStringsClear) < 3 {
		return parserRule{}, errors.New("This is not a valid parser rule.")
	}

	if splitStringsClear[1] != "->" {
		return parserRule{}, errors.New("This is not a valid parser rule.")
	}

	_, alreadyIsTerminal := tokenNames[splitStringsClear[0]]

	if alreadyIsTerminal {
		return parserRule{}, errors.New(fmt.Sprintf("Wromg symbol name : %q. This symbol is already a terminal!", splitStringsClear[0]))
	}

	//Rozpatrujemy najpierw oddzielnie symbol z lewej strony produkcji

	if !checkNonterminalName(splitStringsClear[0]) {
		return parserRule{}, errors.New(fmt.Sprintf("Wrong nonterminal symbol name : %q. Names of nonterminals can contain only letters and underscores!", splitStringsClear[0]))
	} else {
		id, foundNonTerminal := nonTerminalNames[splitStringsClear[0]]

		if foundNonTerminal {
			leftHandSide = id
		} else {
			nonTerminalNames[splitStringsClear[0]] = nextFreeId
			p.numberOfGrammarSymbols++
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
			p.numberOfGrammarSymbols++
			rightHandSide = append(rightHandSide, nextFreeId)
			nextFreeId++
		} else {
			return parserRule{}, errors.New(fmt.Sprintf("Wrong nonterminal symbol name : %q. Names of nonterminals can contain only letters and underscores!", str))
		}

	}

	return parserRule{leftHandSide, rightHandSide, action}, nil
}

//var rules []parserRule

func (p *Parser) getParserRules() []parserRule {
	return p.rules
}

// Zwraca pierwszy indeks (liczbę), jaki został nadany symbolowi nieterminalnemu.
// Jest to również (zgodnie z konwencją przyjętą w tym programie) indeks symbolu
// startowego wprowadzonej przez użytkownika gramatyki.
func (p *Parser) getMinimalNonTerminalIndex() int {
	return p.minimalNonTerminalIndex
}

func (p *Parser) getEndOfInputSymbolId() int {
	return p.endOfInputSymbolId
}

func (p *Parser) getNumberOfGrammarSymbols() int {
	return p.numberOfGrammarSymbols
}

func (p *Parser) AddParserRule(s string, action func([]any)) error {

	if p.isInitialized {
		return errors.New("Parser has already been initialized!")
	}

	result, err := p.toParserRule(s, p.lexer.GetTokenNames(), action, p.nonTerminalNames)

	if err == nil {
		p.rules = append(p.rules, result)
	}

	return err
}

func (p *Parser) setParseTable(pt [][]string) {
	p.parsingTable = pt
}

// var parsingTable [][]string
//var parsingError error = nil

func (p *Parser) RaiseError(err error) {
	p.parsingError = err
}

//var tablesGenerated bool = false
