package parser

import (
	"errors"
	"fmt"
	"goparser/lexer"
	"os"
	"strconv"
	"strings"
	"unicode"
)

//Terminale oraz nietermiale będą reprezentowane liczbami naturnalymi
//(np. 0-10 terminale (te same co w lekserze) 11-14 nieterminale)

// Inna nazwa: Stack item
type Object struct {
	id    int
	Value any
}

// func (o Object) GetIntegerValue() int {
// 	return o.IntegerValue
// }

// func (o Object) GetStringValue() string {
// 	return o.StringValue
// }

// func (o *Object) SetIntegerValue(i int) {
// 	o.IntegerValue = i
// }

// func (o *Object) SetStringValue(s string) {
// 	o.StringValue = s
// }

func (o *Object) setValue(s any) {
	o.Value = s
}

type ParserRule struct {
	leftHandSide  int
	rightHandSide []int
	action        func([]any)
}

// Funkcja tylko do tesów to do celów debugowania
func GetSymbolName(id int) string {

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

func CreateParserRule(leftHandSide int, rightHandSide []int, action func([]any)) ParserRule {
	return ParserRule{leftHandSide, rightHandSide, action}
}

func (p ParserRule) GetRightHandSideLength() int {
	return len(p.rightHandSide)
}

func (p ParserRule) GetRightHandSideSymbol(index int) int {
	return p.rightHandSide[index]
}

func (p ParserRule) GetRightHandSide() []int {
	return p.rightHandSide
}

func (p ParserRule) GetLeftHandSideSymbol() int {
	return p.leftHandSide
}

var S Stack[int]
var ActionS Stack[Object]

func checkNonterminalName(s string) bool {

	for _, c := range s {
		if !(c == '_' || (unicode.IsLetter(c) && unicode.IsUpper(c))) {
			return false
		}
	}

	return true

}

var nonTerminalNames map[string]int = make(map[string]int)

func GetNumberOfGrammarSymbols() int {
	return len(lexer.GetTokenNames()) + len(nonTerminalNames) + 1
}

func toParserRule(s string, tokenNames map[string]int, action func([]any)) (ParserRule, error) {

	splitStrings := strings.Split(s, " ")
	splitStringsClear := make([]string, 0, 5)

	leftHandSide := 0
	rightHandSide := make([]int, 0)

	nextFreeId := len(tokenNames) + len(nonTerminalNames) + 1
	//nonTerminalNames := make(map[string]int)

	for _, split := range splitStrings {
		if split != "" {
			splitStringsClear = append(splitStringsClear, split)
			//fmt.Printf("%#v\n", split)
		}
	}

	if len(splitStringsClear) < 3 {
		return ParserRule{}, errors.New("This is not a valid parser rule.")
	}

	if splitStringsClear[1] != "->" {
		return ParserRule{}, errors.New("This is not a valid parser rule.")
	}

	//Rozpatrujemy najpierw oddzielnie symbol z lewej strony produkcji

	if !checkNonterminalName(splitStringsClear[0]) {
		return ParserRule{}, errors.New(fmt.Sprintf("Wrong nonterminal symbol name : %q. Names of nonterminals can contain only capital letters and underscores!", splitStringsClear[0]))
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
			return ParserRule{}, errors.New(fmt.Sprintf("Wrong nonterminal symbol name : %q. Names of nonterminals can contain only capital letters and underscores!", str))
		}

	}

	return ParserRule{leftHandSide, rightHandSide, action}, nil
}

var rules []ParserRule

func GetParserRules() []ParserRule {
	return rules
}

// Zwraca pierwszy indeks (liczbę), jaki został nadany symbolowi nieterminalnemu.
// Jest to również (zgodnie z konwencją przyjętą w tym programie) indeks symbolu
// startowego wprowadzonej przez użytkownika gramatyki.
func GetMinimalNonTerminalIndex() int {
	return len(lexer.GetTokenNames()) + 1
}

func GetEndOfInputSymbolId() int {
	return len(lexer.GetTokenNames())
}

func AddParserRule(s string, action func([]any)) error {

	result, err := toParserRule(s, lexer.GetTokenNames(), action)

	if err == nil {
		rules = append(rules, result)
	}

	return err
}

func SetParseTable(pt [][]string) {
	parsingTable = pt
}

// Tabela parsowania z książki (do testów)
// int, przesunięcia bitowe
//
//	var parsingTable [][]string = [][]string{
//		{"s5", "", "", "s4", "", "", "1", "2", "3"},
//		{"", "s6", "", "", "", "a", "", "", ""},
//		{"", "r2", "s7", "", "r2", "r2", "", "", ""},
//		{"", "r4", "r4", "", "r4", "r4", "", "", ""},
//		{"s5", "", "", "s4", "", "", "8", "2", "3"},
//		{"", "r6", "r6", "", "r6", "r6", "", "", ""},
//		{"s5", "", "", "s4", "", "", "", "9", "3"},
//		{"s5", "", "", "s4", "", "", "", "", "10"},
//		{"", "s6", "", "", "s11", "", "", "", ""},
//		{"", "r1", "s7", "", "r1", "r1", "", "", ""},
//		{"", "r3", "r3", "", "r3", "r3", "", "", ""},
//		{"", "r5", "r5", "", "r5", "r5", "", "", ""}}

var parsingTable [][]string

func Parse() {
	//Pobieramy pierwszy token
	a, _ := lexer.NextTokenId()

	//Na stosie stan początkowy
	S.Push(0)

	for true {

		fmt.Println("STOS:", S)

		s, _ := S.Peek()

		if parsingTable[s][a] == "" {
			fmt.Println("Syntax error!")
			os.Exit(1)
		} else if string(parsingTable[s][a][0]) == "s" {

			t, _ := strconv.Atoi(parsingTable[s][a][1:])
			S.Push(t)
			fmt.Println("Wykonano przesunięcie: ", parsingTable[s][a])
			a, _ = lexer.NextTokenId()

		} else if string(parsingTable[s][a][0]) == "r" {

			n, _ := strconv.Atoi(parsingTable[s][a][1:])
			symbolsToPop := len(rules[n].rightHandSide)
			for i := 1; i <= symbolsToPop; i++ {
				S.Pop()
			}

			t, _ := S.Peek()
			A := rules[n].leftHandSide
			gotoSymbol, _ := strconv.Atoi(parsingTable[t][A])
			S.Push(gotoSymbol)
			fmt.Println("Wykonano redukcję: ", parsingTable[s][a])

		} else if string(parsingTable[s][a][0]) == "a" {
			fmt.Println("Parsowanie zakończone")
			break
		}

	}
}

func ParseWithSemanticActions() {
	//Pobieramy pierwszy token
	tok, a, _ := lexer.NextTokenWithId()

	//Na stosie stan początkowy
	ActionS.Push(Object{0, nil})

	for true {

		s, _ := ActionS.Peek()

		if parsingTable[s.id][a] == "" {
			fmt.Printf("Syntax error!\n\n")
			break
		} else if string(parsingTable[s.id][a][0]) == "s" {

			t, _ := strconv.Atoi(parsingTable[s.id][a][1:])
			ActionS.Push(Object{t, tok.GetMatchedText()})
			tok, a, _ = lexer.NextTokenWithId()

		} else if string(parsingTable[s.id][a][0]) == "r" {

			n, _ := strconv.Atoi(parsingTable[s.id][a][1:])
			symbolsToPop := len(rules[n].rightHandSide)

			semanticValues := make([]any, symbolsToPop+1)
			valuesFromStack := ActionS.TopSubStack(symbolsToPop)

			for i := 1; i <= symbolsToPop; i++ {
				semanticValues[i] = valuesFromStack[i-1].Value
			}

			if rules[n].action != nil {
				rules[n].action(semanticValues)
			}

			for i := 1; i <= symbolsToPop; i++ {
				ActionS.Pop()
			}

			t, _ := ActionS.Peek()
			A := rules[n].leftHandSide
			gotoSymbol, _ := strconv.Atoi(parsingTable[t.id][A])
			ActionS.Push(Object{gotoSymbol, semanticValues[0]})

		} else if string(parsingTable[s.id][a][0]) == "a" {
			break
		}

	}
}
