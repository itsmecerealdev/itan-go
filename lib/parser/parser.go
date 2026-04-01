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
	global := node.ScopeNode {}
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
		Scope : global,
	}
}

func parseFuncDecl(buf *tokenBuffer) node.FuncDeclNode {
	idTok, err := expect(buf, token.Identifier)
	if err != nil {
		log.Fatal(err)
	}
	funcNode := node.FuncDeclNode{
		Name : idTok.Name,
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
			Type : types.TypeStruct{ Type : typeTok.Name},
			Name : nameTok.Name,
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
	//TODO add conditions for parsing func decls and stuff later. current implementation is just arithmetic for now
	for peek(buf) != token.RBrace {
		n.Statements = append(n.Statements, parseStatement(buf))
	}
	_, err = expect(buf, token.RBrace)
	if err != nil {
		log.Fatal(err)
	}
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
	err := fmt.Errorf("reached end of expression parse, unexpected token encountered: " + tok.Type.String() + "\n\tline, col: %d, %d", tok.Line, tok.Col)
	log.Fatal(err)
	return node.NumberNode{}
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
		Name : name,
		Expression : assignedVal,
	}
}

func parseFuncCall(buf *tokenBuffer) node.Node {
	nameTok, _ := consume(buf)
	node := node.CastOrCallNode {
		Name : nameTok.Name,
	}
	consume(buf)
	for peek(buf) != token.RParen{
		node.Arguments = append(node.Arguments, parseFactor(buf))
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
		Expression : parseExpression(buf),
	}
	return node
}
