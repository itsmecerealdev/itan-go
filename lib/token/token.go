package token

import (
	"errors"
	// "fmt"
	"unicode"

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
	Line int
	Col int
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

func Tokenize(buf string, lineNumber int) ([]Token, error) {
	if len(buf) == 0 {
		return nil, errors.New("buffer is empty")
	}

	buffer := []rune(buf)
	var returnTokens []Token
	bufferIndex := 0

	for bufferIndex < len(buffer)  {
		c := buffer[bufferIndex] 
		// fmt.Printf("%c\n", c)
		if unicode.IsSpace(c) { 
			bufferIndex++
			continue
		}
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && ! isSymbol(c) {
			return nil, errors.New("invalid character: " + string(c))		
		}
		if isSymbol(c) {
			bufferIndex++
			returnTokens = append(returnTokens, Token{
				Type : symbolToTokenType[c], 
				Name : "",
				Line: lineNumber,
				Col: bufferIndex,
			})
		}
		if unicode.IsLetter(c) {
			var str []rune 
			for !unicode.IsSpace(c) && !isSymbol(c) {
				str = append(str, c)
				bufferIndex++
				if bufferIndex >= len(buf) {
					break
				}
				c = buffer[bufferIndex]
			}
			returnTokens = append(returnTokens, Token{
				Name : string(str),
				Type : Identifier,
				Line: lineNumber,
				Col: bufferIndex,
			})
		}
		if unicode.IsNumber(c) {
			var str []rune 
			for !unicode.IsSpace(c) && !isSymbol(c) && !unicode.IsLetter(c) {
				str = append(str, c)
				bufferIndex++
				if bufferIndex >= len(buf) {
					break
				}
				c = buffer[bufferIndex]
			}
			returnTokens = append(returnTokens, Token{
				Name : string(str),
				Type : Number,
				Line: lineNumber,
				Col: bufferIndex,
			})
		}
	}
	returnTokens = append(returnTokens, Token{
		Name : "End",
		Type : End,
		Line: lineNumber,
		Col: bufferIndex,
	})

	length := resolveMultiSymbol(returnTokens)
	return returnTokens[0:length], nil
}

func isSymbol(c rune) bool {
	validSymbols := []rune("()+-*/=;^,{}<>!")
	return slices.Contains(validSymbols, c)
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
			if haskey {
				tokens[i + 1].Type = newType
				tokens = slices.Delete(tokens, i, i+1)
			}
		}
	}
	return len(tokens)
}


