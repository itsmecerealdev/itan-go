package node

import (
	"github.com/itsmecerealdev/itan-go/lib/types"
	"github.com/itsmecerealdev/itan-go/lib/token"
) 

type Node interface {
	Visit()
}

func (sn ScopeNode)Visit() { 

}

func (on OperandNode)Visit() { 

}

func (dn DeclarationNode)Visit() { 

}

func (an AssignmentNode)Visit() {

}

func (nn NumberNode)Visit() {

}

type ProgramNode struct {
	Scope ScopeNode
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
	Value Node
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

type NumberNode struct {
	Value int64
}
