package parser

import (
	"fmt"

	"github.com/itsmecerealdev/itan-go/lib/node"
	"github.com/itsmecerealdev/itan-go/lib/token"
)

type tokenBuffer struct {
	tokens []token.Token	
	tokenIndex int
}

//Helper funcs
func peekAhead(buf tokenBuffer) token.TokenType {

}

func peek(buf tokenBuffer) token.TokenType {

}

func expect(buf tokenBuffer, ttype token.TokenType) (token.Token, error) {
	t := buf.tokens[buf.tokenIndex]
	if t.Type != ttype {
		err := fmt.Errorf("Expected: %i Got: %i at index -> %i",ttype, t.Type, buf.tokenIndex)
		return token.Token{}, err
	}
	return buf.tokens[buf.tokenIndex], nil
}

func consume(buf tokenBuffer) token.Token {
	t := buf.tokens[buf.tokenIndex]
	buf.tokenIndex++
	return t
}

func ParseProgram(tokens []token.Token) node.ProgramNode {
	buf := tokenBuffer{tokens, 0}
	parseScope(buf)
}

func parseScope(buf tokenBuffer) node.ScopeNode {
	parseStatement(buf)
}

func parseStatement(buf tokenBuffer) node.Node {

}
