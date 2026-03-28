package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/itsmecerealdev/itan-go/lib/node"
	"github.com/itsmecerealdev/itan-go/lib/token"
)

type tokenBuffer struct {
	tokens []token.Token	
	tokenIndex int
}

func (tb *tokenBuffer)getCurrentToken() token.Token {
	return tb.tokens[tb.tokenIndex]
}

//Helper funcs
func peekAhead(buf *tokenBuffer) (token.TokenType, error) {
	if buf.tokenIndex + 1 >= len(buf.tokens) {
		err := errors.New("attempting to peek at an out of bounds index")
		return -1, err
	}
	return buf.tokens[buf.tokenIndex + 1].Type, nil
}

func peek(buf *tokenBuffer) token.TokenType {
	return buf.getCurrentToken().Type
}

func expect(buf *tokenBuffer, ttype token.TokenType) (token.Token, error) {
	t := buf.getCurrentToken()
	if t.Type != ttype {
		err := fmt.Errorf("expected: %s Got: %s at index -> %d", ttype, t.Type, buf.tokenIndex)
		return token.Token{}, err
	}
	return consume(buf)
}

func consume(buf *tokenBuffer) (token.Token, error) {
	if buf.tokenIndex >= len(buf.tokens) {
		err := errors.New("buffer index is outside bounds of buffer")
		return token.Token{}, err
	}
	t := buf.getCurrentToken()
	buf.tokenIndex++
	return t, nil
}

func ParseProgram(tokens []token.Token) node.ProgramNode {
	fmt.Print("Program\n")
	buf := tokenBuffer{tokens, 0}
	return node.ProgramNode{
		Scope : parseScope(&buf),
	}
}

func parseScope(buf *tokenBuffer) node.ScopeNode {
	fmt.Print("Scope\n")
	n := node.ScopeNode{}
	//TODO add conditions for parsing func decls and stuff later. current implementation is just arithmetic for now
	for peek(buf) != token.End {
		n.Statements = append(n.Statements, parseStatement(buf))	
	}
	return n
}

func parseStatement(buf *tokenBuffer) node.Node {
	fmt.Print("Stmt.\n")
	tt := peek(buf)
	for tt != token.StatementEnd {
		if tt == token.Number {
			return parseExpression(buf)
		}
		tt = peek(buf)
	}
	return node.NumberNode{}
}

func parseExpression(buf *tokenBuffer) node.Node {
	fmt.Print("Expr.\n")
	left := parseTerm(buf)	
	fmt.Printf("%s\n", left)
	next := peek(buf)	
	for next == token.Plus || next == token.Minus {
		oper, _ := consume(buf)	
		right := parseTerm(buf)
		temp := node.OperandNode {
			Left : left,
			Right : right,
			Type : oper.Type,
		}
		next = peek(buf)
		left = temp
	}
	fmt.Printf("%s\n", left)
	return left
}

func parseTerm(buf *tokenBuffer) node.Node {
	fmt.Print("Term\n")
	left := parseExponent(buf)
	next := peek(buf)	
	for next == token.Star || next == token.Slash {
		oper, _ := consume(buf)	
		right := parseFactor(buf)
		temp := node.OperandNode {
			Left : left,
			Right : right,
			Type : oper.Type,
		}
		next = peek(buf)
		left = temp
	}
	return left
}

func parseExponent(buf *tokenBuffer) node.Node {
	fmt.Print("Exponent\n")
	left := parseFactor(buf)	
	if peek(buf) == token.Exponent {
		oper, _ := consume(buf)
		right := parseExponent(buf)
		temp := node.OperandNode {
			Left : left,
			Right : right,
			Type : oper.Type,
		}
		return temp 
	}
	return left
}

func parseFactor(buf *tokenBuffer) node.Node {
	fmt.Printf("Factor\n%s\n", peek(buf))
	if peek(buf) == token.Number {
		n, _ := expect(buf, token.Number)
		val, err := strconv.ParseInt(n.Name, 10, 64)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("%s : %d", n.Name, val)		
		return node.NumberNode{
			Value : val,
		}
	}
	return node.NumberNode{}
}
