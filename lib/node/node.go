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
	// fmt.Println("Program:")
	visitor.Action(node)
	node.Scope.Accept(visitor)
}

func (node FuncDeclNode)Accept(visitor Visitor) {
	// fmt.Printf("Func: %s Type: %s\n", node.Name, node.Type.Type)
	visitor.Action(node)
	for _, p := range node.Params {
		p.Accept(visitor)
	}
	visitor.MiddleAction(node)
	node.Scope.Accept(visitor)
	visitor.ExitAction(node)
}

func (node CastOrCallNode)Accept(visitor Visitor) {
	// fmt.Println(node.Name)
	visitor.Action(node)
	for _, leaf := range node.Arguments {
		leaf.Accept(visitor)
	}
	visitor.ExitAction(node)
}

func (node ParamNode)Accept(visitor Visitor) {
	// fmt.Printf("Param: %s Type: %s ", node.Name, node.Type.Type)
	visitor.Action(node)
	if node.HasDefault {
		node.Default.Accept(visitor)
	}
}

func (node ReturnNode)Accept(visitor Visitor) {
	// fmt.Printf("Return: ")
	visitor.Action(node)
	node.Expression.Accept(visitor)
	visitor.ExitAction(node)
}

func (node ScopeNode)Accept(visitor Visitor) { 
	// fmt.Println("Scope:")
	visitor.Action(node)
	for _, statement := range node.Statements {
		statement.Accept(visitor)
		// fmt.Print("\n")
	}
	visitor.ExitAction(node)
}

func (node OperandNode)Accept(visitor Visitor) { 
	visitor.Action(node)
	node.Left.Accept(visitor)
	visitor.MiddleAction(node)
	// fmt.Println(node.Type)
	node.Right.Accept(visitor)
	visitor.ExitAction(node)
}

func (node DeclarationNode)Accept(visitor Visitor) { 
	// fmt.Println(node.Type)
	// fmt.Println(node.Name)
	visitor.Action(node)
	node.Expression.Accept(visitor)
	visitor.ExitAction(node)
}

func (node AssignmentNode)Accept(visitor Visitor) {
	// fmt.Println(node.Name)
	visitor.Action(node)
	node.Expression.Accept(visitor)
	visitor.ExitAction(node)
}

func (node VariableNode)Accept(visitor Visitor) {
	// fmt.Print(node.Name)
	visitor.Action(node)
}

func (node NumberNode)Accept(visitor Visitor) {
	// fmt.Print(node.Value)
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
