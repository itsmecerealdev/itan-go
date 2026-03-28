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

	buffer := []rune(buf)
	var return_tokens []Token
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
			return_tokens = append(return_tokens, Token{
				Type : symbolToTokenType[c], 
				Name : "",
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
			return_tokens = append(return_tokens, Token{
				Name : string(str),
				Type : Identifier,
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
			return_tokens = append(return_tokens, Token{
				Name : string(str),
				Type : Number,
			})
		}
	}
	return_tokens = append(return_tokens, Token{
		Name : "End",
		Type : End,
	})

	length := resolveMultiSymbol(return_tokens)
	return return_tokens[0:length], nil

	// for buffer_index < len(buffer) {
		// fmt.Println(buffer_index)
		// var c = buffer[buffer_index]
		// if unicode.IsSpace(c) {
			// buffer_index++
			// continue
		// }
		// if (!unicode.IsLetter(c) && !isSymbol(c) && !unicode.IsNumber(c)) {
			// return nil, errors.New("Invalid character: " + string(c) )
		// }

		// if isSymbol(c) {
			// return_tokens = append(return_tokens, Token{symbolToTokenType[c], ""})
			// buffer_index++
			// continue
		// }

		// if (unicode.IsLetter(c)) {
			// var start_index = buffer_index
			// var c = buffer[buffer_index];
			// for (buffer_index < len(buffer)) && !unicode.IsSpace(c) && !isSymbol(c) {
				// fmt.Printf("%b isSpace, rune: %c\n", unicode.IsSpace(c), c)
				// c = buffer[buffer_index]
				// buffer_index ++;
			// }
			// var str = buffer[start_index:buffer_index]
			// return_tokens = append(return_tokens, Token{Identifier, string(str)})
		// }
		// if(unicode.IsNumber(c)) {
			// start_index := buffer_index
			// for _, c = range buffer {
				// fmt.Printf("%c : %b\n",c ,unicode.IsDigit(c))
				// if !unicode.IsDigit(c) {
					// break
				// }
				// buffer_index++
			// }
			// str := buffer[start_index:buffer_index]
			// return_tokens = append(return_tokens, Token{Number, string(str)})
		// }
	// }
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
			if haskey {
				tokens[i + 1].Type = newType
				tokens = slices.Delete(tokens, i, i+1)
			}
		}
	}
	return len(tokens)
}


