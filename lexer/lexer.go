package lexer

import (
	"errors"
	"fmt"
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
