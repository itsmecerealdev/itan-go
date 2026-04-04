package parser

import (
	"errors"
	"fmt"
	"log"
	//"strconv"

	"github.com/itsmecerealdev/itan-go/lib/node"
	"github.com/itsmecerealdev/itan-go/lib/token"
	"github.com/itsmecerealdev/itan-go/lib/types"
)

// Lower number means higher prescedence
var tokenPrescedence = map[token.TokenType]int {
	token.Number: 1,
	token.Identifier: 1,
	token.Plus: 2,
	token.Minus: 2,
	token.Star: 3,
	token.Slash: 3,
	token.Exponent: 4,
	token.LessThan: 5,
	token.GreaterThan: 5,
	token.LessEqual: 5,
	token.GreaterEqual: 5,
	token.Equal: 6,
	token.NotEqual: 6,
	token.Or: 8,
	token.And: 9,
	//token.Not: 10
	//token.LParen: "LParen",
	//token.RParen: "RParen",
}


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
		err := fmt.Errorf("expected: %s got: %s \n\tline, col: %d, %d", ttype, t.Type, t.Line, t.Col)
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
	global := node.ScopeNode{}
	global.Symbols = &node.Scope{}
	tt := peek(&buf)
	for tt != token.End {
		if tt == token.LBrace {
			global.Statements = append(global.Statements, parseScope(&buf))
		} else if tt == token.Identifier && peekAhead(&buf) == token.LParen {
			global.Statements = append(global.Statements, parseFuncDecl(&buf))
		} else {
			global.Statements = append(global.Statements, parseStatement(&buf))
		}
		tt = peek(&buf)
	}
	return node.ProgramNode{
		Scope: global,
	}
}

func parseFuncDecl(buf *tokenBuffer) node.FuncDeclNode {
	idTok, err := expect(buf, token.Identifier)
	if err != nil {
		log.Fatal(err)
	}
	funcNode := node.FuncDeclNode{
		Name: idTok.Name,
	}
	_, err = expect(buf, token.LParen)
	if err != nil {
		log.Fatal(err)
	}
	funcNode.Params = parseParams(buf)
	if peek(buf) == token.Return {
		consume(buf)
		typeTok, err := expect(buf, token.Identifier)
		if err != nil {
			log.Fatal(err)
		}
		_, isType := types.TypeKeywords[typeTok.Name]
		if !isType {
			err := fmt.Errorf("identifier: %s is not a type\n\tline, col %d, %d", typeTok.Name, typeTok.Line, typeTok.Col)
			log.Fatal(err)
		}
		funcNode.Type = types.TypeStruct{ Type : typeTok.Name }
	} else {
		funcNode.Type = types.TypeStruct{ Type : "void" }
	}
	funcNode.Scope = parseScope(buf)
	return funcNode
}

func parseParams(buf *tokenBuffer) []node.ParamNode {
	params := []node.ParamNode{}
	for {
		typeTok, err := expect(buf, token.Identifier)
		if err != nil {
			log.Fatal(err)
		}
		_, isType := types.TypeKeywords[typeTok.Name]
		if !isType {
			err := fmt.Errorf("identifier: %s is not a type\n\tline, col %d, %d", typeTok.Name, typeTok.Line, typeTok.Col)
			log.Fatal(err)
		}
		nameTok, err := expect(buf, token.Identifier)
		if err != nil {
			log.Fatal(err)
		}
		currentParam := node.ParamNode{
			Type: types.TypeStruct{ Type : typeTok.Name},
			Name: nameTok.Name,
		}
		if peek(buf) == token.Assignment {
			consume(buf);
			currentParam.HasDefault = true;
			currentParam.Default = parseExpression(buf)
		}
		params = append(params, currentParam)
		if peek(buf) == token.Comma {
			consume(buf)
			continue
		} else { break }
	}
	_, err := expect(buf, token.RParen)
	if err != nil {
		log.Fatal(err)
	}
	return params
}

func parseScope(buf *tokenBuffer) node.ScopeNode {
	n := node.ScopeNode{}
	_, err := expect(buf, token.LBrace)
	if err != nil {
		log.Fatal(err)
	}
	for peek(buf) != token.RBrace {
		n.Statements = append(n.Statements, parseStatement(buf))
	}
	_, err = expect(buf, token.RBrace)
	if err != nil {
		log.Fatal(err)
	}
	n.Symbols = &node.Scope{}
	return n
}

func parseStatement(buf *tokenBuffer) node.Node {
	tt := peek(buf)
	tn := peekAhead(buf)
	var n node.Node
	if tt == token.Identifier && tn == token.Identifier {
		n = parseDeclaration(buf)
	} else if tt == token.Identifier && tn == token.LParen {
		n = parseFuncCall(buf)
	} else if tt == token.Return {
		n = parseReturn(buf)
	} else { 
		n = parseExprStatement(buf) 
	}
	_, err := expect(buf, token.StatementEnd)
	if err != nil {
		log.Fatal(err)
	}
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

func parseDeclaration(buf *tokenBuffer) node.Node {
	typetok, err := expect(buf, token.Identifier)
	if err != nil {
		log.Fatal(err)
	}
	idtok, err := expect(buf, token.Identifier)
	if err != nil {
		log.Fatal(err)
	}
	_, isType := types.TypeKeywords[typetok.Name]
	if !isType {
		err := fmt.Errorf("\nidentifier: %s is not a type, unrecoverable state in declaration\n\t line, col %d, %d", typetok.Name, typetok.Line, typetok.Col )
		log.Fatal(err)
	}
	_, err = expect(buf, token.Assignment)
	if err != nil {
		log.Fatal(err)
	}
	expr := parseExpression(buf)
	return node.DeclarationNode{
		Type: types.TypeStruct{
			Type: typetok.Name,
		},
		Name: idtok.Name,
		Expression: expr,
	}
}

func parseAssignment(buf *tokenBuffer) node.Node {
	tok, _ := consume(buf)
	name := tok.Name
	_, isType := types.TypeKeywords[name]
	if isType {
		err := fmt.Errorf("\nidentifier: %s is a type and cannot be assigned \n\t line, col: %d, %d", name, tok.Line, tok.Col)
		log.Fatal(err)
	}
	_, err := expect(buf, token.Assignment)
	if err != nil {
		log.Fatal(err)
	}
	assignedVal := parseExpression(buf)
	return node.AssignmentNode{
		Name: name,
		Expression: assignedVal,
	}
}

func parseFuncCall(buf *tokenBuffer) node.Node {
	nameTok, _ := consume(buf)
	node := node.CastOrCallNode {
		Name: nameTok.Name,
	}
	consume(buf)
	for peek(buf) != token.RParen{
		node.Arguments = append(node.Arguments, parseExpression(buf))
		if peek(buf) == token.Comma {
			consume(buf)
			continue
		} else { break }
	}
	_, err := expect(buf, token.RParen)
	if err != nil {
		log.Fatal(err)
	}
	return node
}

func parseReturn(buf *tokenBuffer) node.ReturnNode {
	consume(buf)
	node := node.ReturnNode {
		Expression: parseExpression(buf),
	}
	return node
}

func parseExpression(buf *tokenBuffer) node.Node {
	var curNode node.OperandNode
	for {
		next := peek(buf)

		prec, isInPrec := tokenPrescedence[next]
		if !isInPrec {
			return curNode
		}
		consume(buf)
		// Traverse down the tree
		for {
			nodePrec := tokenPrescedence[curNode.Type]
			if nodePrec > prec {
				var newNode node.OperandNode

				curNode = newNode
			}
		}
	}
}
