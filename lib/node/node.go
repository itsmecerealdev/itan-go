package node

import (
	"fmt"

	"github.com/itsmecerealdev/itan-go/lib/token"
	"github.com/itsmecerealdev/itan-go/lib/types"
) 

type Node interface {
	Visit()
}

func (pn ProgramNode)Visit() {
	fmt.Println("Program:")
	pn.Scope.Visit()
}

func (fdl FuncDeclNode)Visit() {
	fmt.Printf("Func: %s Type: %s\n", fdl.Name, fdl.Type.Type)
	for _, p := range fdl.Params {
		p.Visit()
	}
	fdl.Scope.Visit()
}

func (pan ParamNode)Visit() {
	fmt.Printf("Param: %s Type: %s\n", pan.Name, pan.Type.Type)
	if pan.HasDefault {
		pan.Default.Visit()
	}
}

func (sn ScopeNode)Visit() { 
	fmt.Println("Scope:")
	for _, statement := range sn.Statements {
		statement.Visit()
		fmt.Print("\n")
	}
}

func (on OperandNode)Visit() { 
	on.Left.Visit()
	fmt.Println(on.Type)
	on.Right.Visit()
}

func (dn DeclarationNode)Visit() { 
	fmt.Println(dn.Type)
	fmt.Println(dn.Name)
	dn.Expression.Visit()
}

func (an AssignmentNode)Visit() {
	fmt.Println(an.Name)
	an.Expression.Visit()
}

func (vn VariableNode)Visit() {
	fmt.Println(vn.Name)
}

func (nn NumberNode)Visit() {
	fmt.Println(nn.Value)
}

type ProgramNode struct {
	Scope ScopeNode
}

type LeafNode interface {
	GetVal()
}

func (vn VariableNode)GetVal() string {
	return vn.Name
}

func (nn NumberNode)GetVal() int64 {
	return nn.Value
}

type FuncDeclNode struct {
	Type types.TypeStruct
	Name string
	Params []ParamNode
	Scope ScopeNode
}

type ParamNode struct {
	Type types.TypeStruct
	Name string
	HasDefault bool
	Default Node
}

type ScopeNode struct {
	Statements []Node 
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
