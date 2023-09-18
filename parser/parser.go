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

var S Stack[int]

//Terminale oraz nietermiale będą reprezentowane licznbami naturnalymi
//(np. 0-10 terminale (te same co w lekserze) 11-14 nieterminale)

type ParserRule struct {
	leftHandSide  int
	rightHandSide []int
}

// Liczba terminali i nieterminali gramatyki
var teminals int = 5
var nonteminals int = 3

//Produkcje gramatyki z książki (do testów)
/*Numeracja symboli:
0 id
1 +
2 *
3 (
4 )
5 $
6 E
7 T
8 F
*/

func checkNonterminalName(s string) bool {

	for _, c := range s {
		if !(c == '_' || (unicode.IsLetter(c) && unicode.IsUpper(c))) {
			return false
		}
	}

	return true

}

func ToParserRule(s string, tokenNames map[string]int, nonTerminalNames map[string]int) (ParserRule, error) {

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

	return ParserRule{leftHandSide, rightHandSide}, nil
}

var rules []ParserRule = []ParserRule{
	//E -> E + T
	{6, []int{6, 1, 7}},
	//E -> T
	{6, []int{7}},
	//T -> T * F
	{7, []int{7, 2, 8}},
	//T -> F
	{7, []int{8}},
	//F -> (E)
	{8, []int{3, 6, 4}},
	//F -> id
	{8, []int{0}},
}

// Tabela parsowania z książki (do testów)
// int, przesunięcia bitowe
var parsingTable [][]string = [][]string{
	{"s5", "", "", "s4", "", "", "1", "2", "3"},
	{"", "s6", "", "", "", "a", "", "", ""},
	{"", "r2", "s7", "", "r2", "r2", "", "", ""},
	{"", "r4", "r4", "", "r4", "r4", "", "", ""},
	{"s5", "", "", "s4", "", "", "8", "2", "3"},
	{"", "r6", "r6", "", "r6", "r6", "", "", ""},
	{"s5", "", "", "s4", "", "", "", "9", "3"},
	{"s5", "", "", "s4", "", "", "", "", "10"},
	{"", "s6", "", "", "s11", "", "", "", ""},
	{"", "r1", "s7", "", "r1", "r1", "", "", ""},
	{"", "r3", "r3", "", "r3", "r3", "", "", ""},
	{"", "r5", "r5", "", "r5", "r5", "", "", ""}}

func Parse() {
	//Pobieramy pierwszy token
	a, _ := lexer.NextTokenId()

	//Na stosie stan początkowy
	S.Push(0)

	for true {

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
			symbolsToPop := len(rules[n-1].rightHandSide)
			for i := 1; i <= symbolsToPop; i++ {
				S.Pop()
			}

			t, _ := S.Peek()
			A := rules[n-1].leftHandSide
			gotoSymbol, _ := strconv.Atoi(parsingTable[t][A])
			S.Push(gotoSymbol)
			fmt.Println("Wykonano redukcję: ", parsingTable[s][a])

		} else if string(parsingTable[s][a][0]) == "a" {
			fmt.Println("Parsowanie zakończone")
			break
		}

	}
}
