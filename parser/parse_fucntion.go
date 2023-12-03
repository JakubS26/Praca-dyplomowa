package parser

import (
	"errors"
	"goparser/lexer"
	"strconv"
)

func Parse() error {

	if !tablesGenerated {
		generateParser()
	}

	//Pobieramy pierwszy token
	tok, a, _ := lexer.NextToken()

	//Na stosie stan początkowy
	actionStack.Push(object{0, nil})

	for true {

		s, _ := actionStack.Peek()

		if parsingTable[s.id][a] == "" {
			return errors.New("Syntax error!")
		} else if string(parsingTable[s.id][a][0]) == "s" {

			t, _ := strconv.Atoi(parsingTable[s.id][a][1:])
			actionStack.Push(object{t, tok.GetMatchedText()})
			tok, a, _ = lexer.NextToken()

		} else if string(parsingTable[s.id][a][0]) == "r" {

			parsingError = nil

			n, _ := strconv.Atoi(parsingTable[s.id][a][1:])
			symbolsToPop := len(rules[n].rightHandSide)

			semanticValues := make([]any, symbolsToPop+1)
			valuesFromStack := actionStack.TopSubStack(symbolsToPop)

			for i := 1; i <= symbolsToPop; i++ {
				semanticValues[i] = valuesFromStack[i-1].Value
			}

			if rules[n].action != nil {
				rules[n].action(semanticValues)
			}

			for i := 1; i <= symbolsToPop; i++ {
				actionStack.Pop()
			}

			t, _ := actionStack.Peek()
			A := rules[n].leftHandSide

			if parsingTable[t.id][A] == "" {
				return errors.New("Syntax error!")
			}

			gotoSymbol, _ := strconv.Atoi(parsingTable[t.id][A])
			actionStack.Push(object{gotoSymbol, semanticValues[0]})

			if parsingError != nil {
				return parsingError
			}

		} else if string(parsingTable[s.id][a][0]) == "a" {
			break
		}

	}

	return nil
}
