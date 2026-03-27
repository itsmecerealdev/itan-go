package token

import (
	"errors"
	"unicode"
	"fmt"
	"slices"
)


type TokenType int

const (
	Number TokenType = iota
	Plus
	Minus
	Star
	Slash
	Exponent
	LParen
	RParen
	Comma
	Assignment
	Not
	Equal
	NotEqual
	LessThan
	GreaterThan
	LessEqual
	GreaterEqual
	And
	Or
	Type
	StatementEnd
	Identifier
	FuncIdentifier
	Return
	LBrace
	RBrace
	End
)

var tokenName = map[TokenType]string {
	Number: "Number",
	Plus: "Plus",
	Minus: "Minus",
	Star: "Star",
	Slash: "Slash",
	Exponent: "Exponent",
	LParen: "LParen",
	RParen: "RParen",
	Comma: "Comma",
	Assignment: "Assignment",
	Not: "Not",
	Equal: "Equal",
	NotEqual: "NotEqual",
	LessThan: "LessThan",
	GreaterThan: "GreaterThan",
	LessEqual: "LessEqual",
	GreaterEqual: "GreaterEqual",
	And: "And",
	Or: "Or",
	Type: "Type",
	StatementEnd: "StatementEnd",
	Identifier: "Identifier",
	FuncIdentifier: "FuncIdentifier",
	Return: "Return",
	LBrace: "LBrace",
	RBrace: "RBrace",
	End: "End",
}

func (tt TokenType)String() string {
	return tokenName[tt]
}


type Token struct {
	Type TokenType
	Name string
}

var symbolToTokenType = map[rune]TokenType {
	'+': Plus,
	'-': Minus,
	'*': Star,
	'/': Slash,
	'^': Exponent,
	'(': LParen,
	')': RParen,
	',': Comma,
	'=': Assignment,
	'!': Not,
	'<': LessThan,
	'>': GreaterThan,
	'&': And,
	'|': Or,
	';': StatementEnd,
	'{': LBrace,
	'}': RBrace,
}

func Tokenize(buf string) ([]Token, error) {
	if len(buf) == 0 {
		return nil, errors.New("Buffer is empty!")
	}

	var buffer []rune = []rune(buf)
	var return_tokens []Token
	var buffer_index int = 0
	for buffer_index < len(buffer) {
		fmt.Println(buffer_index)
		var c = buffer[buffer_index]
		if unicode.IsSpace(c) {
			buffer_index++
			continue
		}
		if (!unicode.IsLetter(c) && !isSymbol(c)) {
			return nil, errors.New("Invalid character")
		}

		if isSymbol(c) {
			return_tokens = append(return_tokens, Token{symbolToTokenType[c], ""})
			buffer_index++
			continue
		}

		if (unicode.IsLetter(c)) {
			var start_index = buffer_index
			var c = buffer[buffer_index];
			for (buffer_index < len(buffer)) && !unicode.IsSpace(c) && !isSymbol(c) {
				// fmt.Printf("%b isSpace, rune: %c\n", unicode.IsSpace(c), c)
				buffer_index ++;
				c = buffer[buffer_index]
			}
			var str = buffer[start_index:buffer_index]
			return_tokens = append(return_tokens, Token{Identifier, string(str)})
		}
	}

	return_tokens = append(return_tokens, Token{End, "end"})
	return return_tokens, nil
}

func isSymbol(c rune) bool {
	var valid_symbols []rune = []rune("()+-*/=;^,{}<>!")
	return slices.Contains(valid_symbols, c)
}
