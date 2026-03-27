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
		// fmt.Println(buffer_index)
		var c = buffer[buffer_index]
		if unicode.IsSpace(c) {
			buffer_index++
			continue
		}
		if (!unicode.IsLetter(c) && !isSymbol(c) && !unicode.IsNumber(c)) {
			return nil, errors.New("Invalid character: " + string(c) )
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
				c = buffer[buffer_index]
				buffer_index ++;
			}
			var str = buffer[start_index:buffer_index]
			return_tokens = append(return_tokens, Token{Identifier, string(str)})
		}
		if(unicode.IsNumber(c)) {
			start_index := buffer_index
			for (buffer_index < len(buffer)) && !unicode.IsSpace(c) && !isSymbol(c) {
				c = buffer[buffer_index]
				buffer_index++;
			}
			str := buffer[start_index:buffer_index]
			return_tokens = append(return_tokens, Token{Number, string(str)})
		}
	}

	return_tokens = append(return_tokens, Token{End, "end"})
	length := resolveMultiSymbol(return_tokens)
	return return_tokens[0:length], nil
}

func isSymbol(c rune) bool {
	var valid_symbols []rune = []rune("()+-*/=;^,{}<>!")
	return slices.Contains(valid_symbols, c)
}

var multiSymMap = map[[2]TokenType]TokenType {
	{Assignment, Assignment} : Equal,
	{Not, Assignment} : NotEqual,
	{LessThan, Assignment} : LessEqual,
	{GreaterThan, Assignment} : GreaterEqual,
	{Assignment, GreaterThan} : Return,
}

func resolveMultiSymbol(tokens []Token) int {
	for i, token := range tokens {
		if(i < (len(tokens) - 1)) {
			tt1 := token.Type
			tt2 := tokens[i + 1].Type
			newType, haskey := multiSymMap[[2]TokenType{tt1, tt2}]
			if !haskey {
				fmt.Printf("Not In multiSymMap table: %i : %i\n", tt1, tt2)
			} else {
				tokens[i + 1].Type = newType
				tokens = slices.Delete(tokens, i, i+1)
			}
		}
	}
	return len(tokens)
}


