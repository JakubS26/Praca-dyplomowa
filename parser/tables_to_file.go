package parser

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func addPadding(s string, finalLength int) string {
	return strings.Repeat(" ", finalLength-len(s)) + s
}

func (p *Parser) ExportParseTablesToFile(filename string) error {

	if p.isInitialized == false {
		return errors.New("Parser has not been initialized!")
	}

	os.Remove(filename)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	maxLength := 0

	for i := 0; i < p.numberOfGrammarSymbols; i++ {
		maxLength = max(maxLength, len(p.getSymbolName(i)))
	}

	padding := maxLength + 1

	maxNumberLength := len(strconv.Itoa(len(p.transitions) - 1))

	file.WriteString(addPadding("", maxNumberLength+1))
	for i := 0; i < p.getNumberOfGrammarSymbols(); i++ {
		file.WriteString(addPadding(p.getSymbolName(i), padding))
	}
	file.WriteString("\n")

	for i := 0; i < len(p.parsingTable); i++ {
		file.WriteString(addPadding(strconv.Itoa(i), maxNumberLength))
		file.WriteString(" ")
		for _, action := range p.parsingTable[i] {
			file.WriteString(addPadding(action, padding))
		}
		file.WriteString("\n")
	}

	file.Close()
	return nil
}
