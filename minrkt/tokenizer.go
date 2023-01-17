package minrkt

import (
	"fmt"
	"regexp"
	"strings"
)

type Token struct {
	tokType TokenType
	val     string
}

type TokenType int

const (
	TOK_INVALID TokenType = iota - 1
	TOK_LPAREN
	TOK_RPAREN
	TOK_NUM
	TOK_ADD
	TOK_SUB
	TOK_MUL
	TOK_DIV
	TOK_EQ
	TOK_GTEQ
	TOK_LTEQ
	TOK_GT
	TOK_LT
	TOK_AND
	TOK_OR
	TOK_NOT
	TOK_IF
	TOK_TRUE
	TOK_FALSE
	TOK_DEFINE
	TOK_VAR
)

var tokenRegexList = []string{
	`^(\()`,
	`^(\))`,
	`^([1-9][0-9]*(?:\.[0-9]*)?)`,
	`^(\+)`,
	`^(\-)`,
	`^(\*)`,
	`^(/)`,
	`^(=)`,
	`^(>=)`,
	`^(<=)`,
	`^(>)`,
	`^(<)`,
	`^(and)`,
	`^(or)`,
	`^(not)`,
	`^(if)`,
	`^(true|#t)`,
	`^(false|#f)`,
	`^(define)`,
	`^([a-zA-Z][a-zA-z0-9_]*)`,
}

type InvalidCharError struct {
	c string
}

func (e *InvalidCharError) Error() string {
	return fmt.Sprintf("Invalid Character %s", e.c)
}

// returns -1 if no matching token
func getTokenIndx(tokenList []string) int {
	indx := -1
	for i := 1; i < len(tokenList); i++ {
		if tokenList[i] != "" {
			indx = i - 1
		}
	}
	return indx
}

// Maps first character of input to a Token
// currently also responsible for validating characters
// does not accept leading zeroes
func NextToken(remainder string) (Token, string, error) {
	re := regexp.MustCompile(strings.Join(tokenRegexList, "|"))
	wsRe := regexp.MustCompile(`^\s+`)
	var token Token
	var newRemainder string
	var err error

	if len(remainder) != 0 {
		if len(wsRe.FindStringSubmatch(remainder)) != 0 {
			remainder = remainder[1:]
		}

		tokenList := re.FindStringSubmatch(remainder)
		indx := getTokenIndx(tokenList)

		if indx < 0 {
			err = &InvalidCharError{remainder[0:1]}
			newRemainder = ""
		} else {
			token.tokType, token.val = TokenType(indx), tokenList[0]

			tokenSize := len(tokenList[0])
			newRemainder = remainder[tokenSize:]
		}
	}

	return token, newRemainder, err
}

// returns a list of tokens
// assumes no empty input will be passed
func Tokenizer(line string) ([]Token, error) {

	remainder := line
	var tokens []Token
	for {
		token, newRemainder, err := NextToken(remainder)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
		if len(newRemainder) == 0 {
			return tokens, nil
		}
		remainder = newRemainder
	}
}
