package node

import (
	"github.com/itsmecerealdev/itan-go/lib/types"
)

type Node interface { //This gives us the ability to do []*Node, equivalent to cpp vector<Node*> abstract,
					  //but we retain the concrete unlike c++
	isNode()
}

type ExpressionNode interface {
	isExpression()
}

type ProgramNode struct {
	Scope *ScopeNode
}

type FuncDeclNode {
	Type TypeStruct
	name string
	params []*ParamNode
	scope *ScopeNode
}

type ParamNode {
	Type TypeStruct
	name string
	value *ExpressionNode
}

type ScopeNode struct {
	Statements []*Node 
}

type Scope struct {
	Parent *Scope
	Variables map[string]DeclarationNode
	Params map[string]ParamNode
	Funcs map[string]FuncDeclNode
}

type DeclarationNode struct {
	Type TypeStruct
	name string
	expression *ExpressionNode
}

type AssignmentNode struct {
	name string expression *ExpressionNode
}

type NumberNode struct {
	value int64
}

//The node interface thing needs these though, even though they are stubs that do nada, as they restrict what types a []*Node fits
//[]interface{} works too but it's unsafe, as ALL types in go implement the empty interface interface{}
func (pn ProgramNode) isNode() {}
func (sn ScopeNode) isNode() {}
func (dn DeclarationNode) isNode() {}
func (an AssignmentNode) isNode() {}
func (nn NumberNode) isNode() {}

//These are for *ExpressionNode capability to clamp what can fit in what
//These will be things like OperandNode, NumberNode, BoolNode, etc
func (nn NumberNode) isExpression(){}
