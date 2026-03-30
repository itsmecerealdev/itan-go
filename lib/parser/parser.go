package parser

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/itsmecerealdev/itan-go/lib/node"
	"github.com/itsmecerealdev/itan-go/lib/token"
	"github.com/itsmecerealdev/itan-go/lib/types"
)

type tokenBuffer struct {
	tokens []token.Token	
	tokenIndex int
	//eventually, there will be an []error field
}

func (tb *tokenBuffer)getCurrentToken() token.Token {
	return tb.tokens[tb.tokenIndex]
}

//Helper funcs
func peekAhead(buf *tokenBuffer) token.TokenType {
	if buf.tokenIndex + 1 >= len(buf.tokens) {
		err := errors.New("attempting to peek at an out of bounds index, unrecoverable state")
		log.Fatal(err)
	}
	return buf.tokens[buf.tokenIndex + 1].Type
}

func peek(buf *tokenBuffer) token.TokenType {
	return buf.getCurrentToken().Type
}

func expect(buf *tokenBuffer, ttype token.TokenType) (token.Token, error) {
	t := buf.getCurrentToken()
	if t.Type != ttype {
		err := fmt.Errorf("expected: %s got: %s at line, col: %d, %d", ttype, t.Type, t.Line, t.Col)
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
	buf := tokenBuffer{tokens, 0}
	return node.ProgramNode{
		Scope : parseScope(&buf),
	}
}

func parseScope(buf *tokenBuffer) node.ScopeNode {
	n := node.ScopeNode{}
	//TODO add conditions for parsing func decls and stuff later. current implementation is just arithmetic for now
	for peek(buf) != token.End {
		n.Statements = append(n.Statements, parseStatement(buf))	
	}
	return n
}

func parseStatement(buf *tokenBuffer) node.Node {
	tt := peek(buf)
	tn := peekAhead(buf)
	if tt == token.Identifier && tn == token.Identifier {
		n := parseDeclaration(buf)
		expect(buf, token.StatementEnd)
		return n 
	}
	n := parseExprStatement(buf)
	expect(buf, token.StatementEnd)
	return n 
}

func parseExprStatement(buf *tokenBuffer) node.Node {
	tt := peek(buf)
	tn := peekAhead(buf)
	if tt == token.Identifier && tn == token.Assignment  {
		assign := parseAssignment(buf)
		return assign
	} else {
		expr := parseExpression(buf)
		return expr
	}
}

func parseExpression(buf *tokenBuffer) node.Node {
	left := parseTerm(buf)	
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
	return left
}

func parseTerm(buf *tokenBuffer) node.Node {
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
	if peek(buf) == token.Identifier {
		tok, _ := consume(buf)
		return node.VariableNode{
			Name : tok.Name,
		}
	}
	if peek(buf) == token.Number {
		n, _ := expect(buf, token.Number)
		val, err := strconv.ParseInt(n.Name, 10, 64)
		if err != nil {
			fmt.Print(err)
		}
		return node.NumberNode{
			Value : val,
		}
	}
	tok, _ := consume(buf)
	err := fmt.Errorf("reached end of expression parse, unexpected token encountered: " + tok.Name + " line, col: %d, %d", tok.Line, tok.Col)
	log.Fatal(err)
	return node.NumberNode{}
}

func parseDeclaration(buf *tokenBuffer) node.Node {
	typetok, _ := expect(buf, token.Identifier)
	idtok, _ := expect(buf, token.Identifier)
	_, isType := types.TypeKeywords[typetok.Name]
	if !isType {
		err := fmt.Errorf("identifier: %s at line, col: %d, %d  is not a type, unrecoverable state", typetok.Name, typetok.Line, typetok.Col )
		log.Fatal(err)
	}
	expect(buf, token.Assignment)
	expr := parseExpression(buf)
	return node.DeclarationNode{
		Type : types.TypeStruct{
			Type : typetok.Name,
		},
		Name : idtok.Name,
		Expression : expr,
	}
}

func parseAssignment(buf *tokenBuffer) node.Node {
	tok, _ := consume(buf)
	name := tok.Name
	expect(buf, token.Assignment)
	assignedVal := parseExpression(buf)
	expect(buf, token.StatementEnd)
	return node.AssignmentNode{
		Name : name,
		Expression : assignedVal,
	}
}
