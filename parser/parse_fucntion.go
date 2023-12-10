package parser

import (
	"errors"
	"strconv"
)

func (p *Parser) Parse() error {

	if p.isInitialized == false {
		return errors.New("The parser hasn't been initialized!")
	}

	var actionStack Stack[object]

	//Pobieramy pierwszy token
	tok, a, err := p.lexer.NextToken()

	if err != nil {
		return err
	}

	//Na stosie stan początkowy
	actionStack.Push(object{0, nil})

	for true {

		s, _ := actionStack.Peek()

		if p.parsingTable[s.id][a] == "" {
			return errors.New("Syntax error!")
		} else if string(p.parsingTable[s.id][a][0]) == "s" {

			t, _ := strconv.Atoi(p.parsingTable[s.id][a][1:])
			actionStack.Push(object{t, tok.GetMatchedText()})
			tok, a, err = p.lexer.NextToken()

			if err != nil {
				return err
			}

		} else if string(p.parsingTable[s.id][a][0]) == "r" {

			p.parsingError = nil

			n, _ := strconv.Atoi(p.parsingTable[s.id][a][1:])
			symbolsToPop := len(p.rules[n].rightHandSide)

			semanticValues := make([]any, symbolsToPop+1)
			valuesFromStack := actionStack.TopSubStack(symbolsToPop)

			for i := 1; i <= symbolsToPop; i++ {
				semanticValues[i] = valuesFromStack[i-1].Value
			}

			if p.rules[n].action != nil {
				p.rules[n].action(semanticValues)
			}

			for i := 1; i <= symbolsToPop; i++ {
				actionStack.Pop()
			}

			t, _ := actionStack.Peek()
			A := p.rules[n].leftHandSide

			if p.parsingTable[t.id][A] == "" {
				return errors.New("Syntax error!")
			}

			gotoSymbol, _ := strconv.Atoi(p.parsingTable[t.id][A])
			actionStack.Push(object{gotoSymbol, semanticValues[0]})

			if p.parsingError != nil {
				return p.parsingError
			}

		} else if string(p.parsingTable[s.id][a][0]) == "a" {
			break
		}

	}

	return nil
}
