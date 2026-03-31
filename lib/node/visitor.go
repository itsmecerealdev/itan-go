package node

import (
	"fmt"
	"strings"
)

type Visitor interface {
	Action(currNode Node)
	MiddleAction(currNode Node)
	ExitAction(currNode Node)
}

func (printer *Printer)Action(currNode Node) {
	switch concrete := currNode.(type) {
	case ProgramNode:
		fmt.Println("Program:")
	case FuncDeclNode:
		printer.TabHelper()
		fmt.Printf("%s %s(", concrete.Type.Type, concrete.Name)
	case CastOrCallNode:
		printer.TabHelper()
		fmt.Printf("%s(", concrete.Name)
	case ParamNode:	
		fmt.Printf("%s %s ", concrete.Type.Type, concrete.Name)
		if concrete.HasDefault {
			fmt.Printf("= ")
		}
	case ReturnNode:
		printer.TabHelper()
		fmt.Print("=> ")
	case ScopeNode:
		printer.TabHelper()
		fmt.Println("{")
		printer.TabDepth++
	case OperandNode:
		printer.TabHelper()
	case DeclarationNode:
		printer.TabHelper()
		fmt.Printf("%s %s = ", concrete.Type.Type, concrete.Name)
	case AssignmentNode:
		printer.TabHelper()
		fmt.Printf("%s = ", concrete.Name)
	case VariableNode:
		fmt.Print(concrete.Name)
	case NumberNode:
		fmt.Print(concrete.Value)
	}
}

func (printer *Printer)MiddleAction(currNode Node) {
	switch concrete := currNode.(type) {
		case FuncDeclNode:
			fmt.Print(") ")
		case OperandNode:
			fmt.Printf(" %s ", concrete.Type.String())
	}
}

func (printer *Printer)ExitAction(currNode Node) {
	switch currNode.(type) {
	case ReturnNode:
		fmt.Println()
	case CastOrCallNode:
		fmt.Println(")")
	case ScopeNode:
		printer.TabDepth--;
		printer.TabHelper()
		fmt.Println("}")
	case OperandNode:
		fmt.Println()
	case DeclarationNode:
		fmt.Println()
	case AssignmentNode:
		fmt.Println()
	}
}

func (printer Printer)TabHelper() {
	fmt.Print(strings.Repeat("   ", printer.TabDepth))
}

type Printer struct {
	TabDepth int
}
