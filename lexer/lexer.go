package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"unicode"
)

type TokenDefinition struct {
	name  string
	regex string
}

var tokens []TokenDefinition

func AddTokenDefinition(name, regex string) {
	tokens = append(tokens, TokenDefinition{name, regex})
}

func PrintTokens() {
	for i := range tokens {
		fmt.Println(tokens[i].name, tokens[i].regex)
	}
}

var file *os.File

func OpenFile(fileName string) error {
	var err error = nil
	file, err = os.Open(fileName)

	if err != nil {
		err = errors.New(fmt.Sprintf("Błąd otwierania pliku o nazwie: %v!", fileName))
	}

	return err
}

// Prints file to test whether is has been properly read
func TestPrintFile() {
	reader := bufio.NewReader(file)

	var err error = nil
	var b byte = 0

	for err == nil {
		b, err = reader.ReadByte()
		fmt.Printf("%c", b)
	}

}

func Init() error {

	if len(tokens) == 0 {
		return errors.New("The set of tokens cannot be empty!")
	}

	for i := range tokens {

		if len(tokens[i].name) == 0 {
			return errors.New("A name of a token cannot be an empty string!")
		}

		for _, c := range tokens[i].name {
			if !(c == '_' || (unicode.IsLetter(c) && unicode.IsUpper(c))) {
				return errors.New(fmt.Sprintf("Wrong character : %q. Names of tokens can contain only capital letters and underscores!", c))
			}
		}

	}

	return nil
}
