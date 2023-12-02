package lexer

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"unicode"
)

type TokenDefinition struct {
	name  string
	regex string
}

type Token struct {
	name        string
	matchedText string
}

func (t Token) GetMatchedText() string {
	return t.matchedText
}

var tokenNames map[string]int = make(map[string]int)
var nextTokenId int = 0

var tokenDefinitions []TokenDefinition
var compiledRegexes []*regexp.Regexp

var ignoreRegexes []string
var ignoreCompiledRegexes []*regexp.Regexp

var fileBuffer []byte

func AddTokenDefinition(name, regex string) {
	tokenDefinitions = append(tokenDefinitions, TokenDefinition{name, regex})
	tokenNames[name] = nextTokenId
	nextTokenId++
}

func Ignore(regex string) {
	ignoreRegexes = append(ignoreRegexes, regex)
}

func PrintTokens() {
	for i := range tokenDefinitions {
		fmt.Println(tokenDefinitions[i].name, tokenDefinitions[i].regex)
	}
}

func GetTokenNames() map[string]int {
	return tokenNames
}

func OpenFile(fileName string) error {
	var err error = nil
	fileBuffer, err = os.ReadFile(fileName)

	if err != nil {
		err = errors.New(fmt.Sprintf("Błąd otwierania pliku o nazwie: %v!", fileName))
	}

	return err
}

func SetInputString(input string) {
	fileBuffer = []byte(input)
}

// Prints file to test whether is has been properly open
func TestPrintFile() {
	fmt.Print(string(fileBuffer))
}

func Init() error {

	if len(tokenDefinitions) == 0 {
		return errors.New("The set of tokens cannot be empty!")
	}

	// Sprawdzamy definicje tokenów

	for i := range tokenDefinitions {

		if len(tokenDefinitions[i].name) == 0 {
			return errors.New("A name of a token cannot be an empty string!")
		}

		for _, c := range tokenDefinitions[i].name {
			if !(c == '_' || (unicode.IsLetter(c) /*&& unicode.IsUpper(c)*/)) {
				return errors.New(fmt.Sprintf("Wrong character : %q. Names of tokens can contain only capital letters and underscores!", c))
			}
		}

		compiledRegex, err := regexp.Compile(tokenDefinitions[i].regex)

		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't compile regular expression for token %v. \"%v\" is not a valid regular expression!",
				tokenDefinitions[i].name, tokenDefinitions[i].regex))
		}

		compiledRegexes = append(compiledRegexes, compiledRegex)

	}

	for _, re := range compiledRegexes {
		re.Longest()
	}

	// Sprawdzamy wyrażenia regularne dla ignorowanych ciągów znaków

	for i := range ignoreRegexes {

		compiledRegex, err := regexp.Compile(ignoreRegexes[i])

		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't compile regular expression for ignored token. \"%v\" is not a valid regular expression!",
				tokenDefinitions[i].regex))
		}

		ignoreCompiledRegexes = append(ignoreCompiledRegexes, compiledRegex)

	}

	for _, re := range ignoreCompiledRegexes {
		re.Longest()
	}

	return nil
}

func PrintToken(tok Token) {
	fmt.Printf("{%v, %q}\n", tok.name, tok.matchedText)
}

// func NextToken() (Token, error) {
// 	var matchedText string
// 	var matchedLoc []int

// 	if len(fileBuffer) == 0 {
// 		return Token{"", ""}, errors.New("End of input.")
// 	}

// 	for i, re := range compiledRegexes {
// 		matchedLoc = re.FindIndex(fileBuffer)

// 		if matchedLoc != nil && matchedLoc[0] == 0 {
// 			matchedText = string(fileBuffer[matchedLoc[0]:matchedLoc[1]])
// 			fileBuffer = fileBuffer[matchedLoc[1]:]
// 			return Token{tokenDefinitions[i].name, matchedText}, nil
// 		}

// 	}

// 	fmt.Println("The lexer was not able to match given input!")
// 	return Token{"", ""}, errors.New("The lexer was not able to match given input!")
// }

func NextTokenWithId() (Token, int, error) {
	var matchedText string
	var matchedLoc []int

	if len(fileBuffer) == 0 {
		return Token{"", ""}, len(tokenDefinitions), errors.New("End of input.")
	}

	ignored := true

	for ignored {
		for i, re := range compiledRegexes {
			matchedLoc = re.FindIndex(fileBuffer)

			if matchedLoc != nil && matchedLoc[0] == 0 {
				matchedText = string(fileBuffer[matchedLoc[0]:matchedLoc[1]])
				fileBuffer = fileBuffer[matchedLoc[1]:]
				return Token{tokenDefinitions[i].name, matchedText}, i, nil
			}

		}

		ignored = false

		for _, re := range ignoreCompiledRegexes {
			matchedLoc = re.FindIndex(fileBuffer)
			if matchedLoc != nil && matchedLoc[0] == 0 {
				fileBuffer = fileBuffer[matchedLoc[1]:]
				ignored = true
			}
		}
	}

	//fmt.Println("BUFFER:", fileBuffer)
	fmt.Println("The lexer was not able to match given input!")
	return Token{"", ""}, 0, errors.New("The lexer was not able to match given input!")
}

func NextTokenId() (int, error) {
	var matchedLoc []int

	if len(fileBuffer) == 0 {
		return len(tokenDefinitions), errors.New("End of input.")
	}

	for i, re := range compiledRegexes {
		matchedLoc = re.FindIndex(fileBuffer)

		if matchedLoc != nil && matchedLoc[0] == 0 {
			fileBuffer = fileBuffer[matchedLoc[1]:]
			//fmt.Println(i)
			return i, nil
		}

	}

	fmt.Println("The lexer was not able to match given input!")
	return 0, errors.New("The lexer was not able to match given input!")
}
