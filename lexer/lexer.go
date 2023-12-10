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

type Lexer struct {
	tokenNames            map[string]int
	nextTokenId           int
	tokenDefinitions      []TokenDefinition
	compiledRegexes       []*regexp.Regexp
	ignoreRegexes         []string
	ignoreCompiledRegexes []*regexp.Regexp
	fileBuffer            []byte
	isInitalized          bool
}

func NewLexer() *Lexer {
	return &Lexer{
		tokenDefinitions:      make([]TokenDefinition, 0),
		compiledRegexes:       make([]*regexp.Regexp, 0),
		ignoreRegexes:         make([]string, 0),
		ignoreCompiledRegexes: make([]*regexp.Regexp, 0),
		tokenNames:            make(map[string]int),
	}
}

func (l *Lexer) AddTokenDefinition(name, regex string) error {
	if !l.isInitalized {
		l.tokenDefinitions = append(l.tokenDefinitions, TokenDefinition{name, regex})
		l.tokenNames[name] = l.nextTokenId
		l.nextTokenId++
		return nil
	} else {
		return errors.New("Lexer has already been initialized!")
	}
}

func (l *Lexer) Ignore(regex string) error {
	if !l.isInitalized {
		l.ignoreRegexes = append(l.ignoreRegexes, regex)
		return nil
	} else {
		return errors.New("Lexer has already been initialized!")
	}
}

func (l *Lexer) PrintTokens() {
	for i := range l.tokenDefinitions {
		fmt.Println(l.tokenDefinitions[i].name, l.tokenDefinitions[i].regex)
	}
}

func (l *Lexer) GetTokenNames() map[string]int {
	return l.tokenNames
}

func (l *Lexer) OpenFile(fileName string) error {
	var err error = nil
	l.fileBuffer, err = os.ReadFile(fileName)

	if err != nil {
		err = errors.New(fmt.Sprintf("Can't open file: %v!", fileName))
	}

	return err
}

func (l *Lexer) SetInputString(input string) {
	l.fileBuffer = []byte(input)
}

// Prints file to test whether is has been properly open
func (l *Lexer) TestPrintFile() {
	fmt.Print(string(l.fileBuffer))
}

func (l *Lexer) Init() error {

	if l.isInitalized {
		return nil
	}

	if len(l.tokenDefinitions) == 0 {
		return errors.New("The set of tokens cannot be empty!")
	}

	// Sprawdzamy definicje tokenów

	for i := range l.tokenDefinitions {

		if len(l.tokenDefinitions[i].name) == 0 {
			return errors.New("A name of a token cannot be an empty string!")
		}

		for _, c := range l.tokenDefinitions[i].name {
			if !(c == '_' || (unicode.IsLetter(c))) {
				return errors.New(fmt.Sprintf("Wrong character : %q. Names of tokens can contain only letters and underscores!", c))
			}
		}

		compiledRegex, err := regexp.Compile(l.tokenDefinitions[i].regex)

		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't compile regular expression for token %v. \"%v\" is not a valid regular expression!",
				l.tokenDefinitions[i].name, l.tokenDefinitions[i].regex))
		}

		l.compiledRegexes = append(l.compiledRegexes, compiledRegex)

	}

	for _, re := range l.compiledRegexes {
		re.Longest()
	}

	// Sprawdzamy wyrażenia regularne dla ignorowanych ciągów znaków

	for i := range l.ignoreRegexes {

		compiledRegex, err := regexp.Compile(l.ignoreRegexes[i])

		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't compile regular expression for ignored token. \"%v\" is not a valid regular expression!",
				l.tokenDefinitions[i].regex))
		}

		l.ignoreCompiledRegexes = append(l.ignoreCompiledRegexes, compiledRegex)

	}

	for _, re := range l.ignoreCompiledRegexes {
		re.Longest()
	}

	l.isInitalized = true
	return nil
}

func PrintToken(tok Token) {
	fmt.Printf("{%v, %q}\n", tok.name, tok.matchedText)
}

func (l *Lexer) NextToken() (Token, int, error) {

	if l.isInitalized == false {
		return Token{}, 0, errors.New("Lexer hasn't been initialized!")
	}

	var matchedText string
	var matchedLoc []int

	if len(l.fileBuffer) == 0 {
		return Token{"", ""}, len(l.tokenDefinitions), nil
	}

	ignored := true

	for ignored {
		for i, re := range l.compiledRegexes {
			matchedLoc = re.FindIndex(l.fileBuffer)

			if matchedLoc != nil && matchedLoc[0] == 0 {
				matchedText = string(l.fileBuffer[matchedLoc[0]:matchedLoc[1]])
				l.fileBuffer = l.fileBuffer[matchedLoc[1]:]
				return Token{l.tokenDefinitions[i].name, matchedText}, i, nil
			}

		}

		ignored = false

		for _, re := range l.ignoreCompiledRegexes {
			matchedLoc = re.FindIndex(l.fileBuffer)
			if matchedLoc != nil && matchedLoc[0] == 0 {
				l.fileBuffer = l.fileBuffer[matchedLoc[1]:]
				ignored = true
			}
		}
	}

	return Token{"", ""}, 0, errors.New("The lexer was not able to match given input!")
}
