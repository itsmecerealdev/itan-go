package node

import (
	// "fmt"

	"github.com/itsmecerealdev/itan-go/lib/token"
	"github.com/itsmecerealdev/itan-go/lib/types"
) 

type Node interface {
	Accept(visitor Visitor)
}

func (node ProgramNode)Accept(visitor Visitor) {
	visitor.Action(node)
	node.Scope.Accept(visitor)
}

func (node FuncDeclNode)Accept(visitor Visitor) {
	visitor.Action(node)
	for _, p := range node.Params {
		p.Accept(visitor)
	}
	visitor.MiddleAction(node)
	node.Scope.Accept(visitor)
	visitor.ExitAction(node)
}

func (node CastOrCallNode)Accept(visitor Visitor) {
	visitor.Action(node)
	for _, leaf := range node.Arguments {
		leaf.Accept(visitor)
	}
	visitor.ExitAction(node)
}

func (node ParamNode)Accept(visitor Visitor) {
	visitor.Action(node)
	if node.HasDefault {
		node.Default.Accept(visitor)
	}
}

func (node ReturnNode)Accept(visitor Visitor) {
	visitor.Action(node)
	node.Expression.Accept(visitor)
	visitor.ExitAction(node)
}

func (node ScopeNode)Accept(visitor Visitor) { 
	visitor.Action(node)
	for _, statement := range node.Statements {
		statement.Accept(visitor)
	}
	visitor.ExitAction(node)
}

func (node OperandNode)Accept(visitor Visitor) { 
	visitor.Action(node)
	node.Left.Accept(visitor)
	visitor.MiddleAction(node)
	node.Right.Accept(visitor)
	visitor.ExitAction(node)
}

func (node DeclarationNode)Accept(visitor Visitor) { 
	visitor.Action(node)
	node.Expression.Accept(visitor)
	visitor.ExitAction(node)
}

func (node AssignmentNode)Accept(visitor Visitor) {
	visitor.Action(node)
	node.Expression.Accept(visitor)
	visitor.ExitAction(node)
}

func (node VariableNode)Accept(visitor Visitor) {
	visitor.Action(node)
}

func (node NumberNode)Accept(visitor Visitor) {
	visitor.Action(node)
}

type ProgramNode struct {
	Scope ScopeNode
}

type LeafNode interface {
	GetVal()
}

func (node VariableNode)GetVal() string {
	return node.Name
}

func (node NumberNode)GetVal() int64 {
	return node.Value
}

type FuncDeclNode struct {
	Type types.TypeStruct
	Name string
	Params []ParamNode
	Scope ScopeNode
}

type CastOrCallNode struct {
	Name string
	Arguments []Node
}

type ParamNode struct {
	Type types.TypeStruct
	Name string
	HasDefault bool
	Default Node
}

type ReturnNode struct {
	Type types.TypeStruct
	Expression Node
}

type ScopeNode struct {
	Statements []Node 
	Symbols Scope
}

type Scope struct {
	Parent *Scope
	Variables map[string]DeclarationNode
	Params map[string]ParamNode
	Funcs map[string]FuncDeclNode
}

type OperandNode struct {
	Left Node
	Right Node
	Type token.TokenType
}

type DeclarationNode struct {
	Type types.TypeStruct
	Name string
	Expression Node
}

type AssignmentNode struct {
	Name string 
	Expression Node
}

type VariableNode struct {
	Name string
}

type NumberNode struct {
	Value int64
}
