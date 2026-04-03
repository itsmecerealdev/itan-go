package node

import (
	"fmt"
	"log"
	"github.com/itsmecerealdev/itan-go/lib/types"
)

type Semantic struct {
	InFunc bool
	Returned bool
	CurrScope *Scope
	CurrInit string
}

func (semantic *Semantic)Action(node Node) {
	switch concrete := node.(type) {
	case FuncDeclNode:
		semantic.InFunc = true;
	case CastOrCallNode:
		_, isType := types.TypeKeywords[concrete.Name]
		if isType && len(concrete.Arguments) == 1 {
			concrete.Res = Cast	
		} else if isType {
			err := fmt.Errorf("cast %s has %d parameters, expecting one", concrete.Name, len(concrete.Arguments))
			log.Fatal(err)
		} else {
			concrete.Res = Call
		}
	case ParamNode:
		semantic.CurrInit = concrete.Name
	case ReturnNode:
		semantic.Returned = true
	case ScopeNode:
		if concrete.Symbols == nil { return }
		semantic.CurrScope = concrete.Symbols
	case DeclarationNode:
		semantic.CurrInit = concrete.Name
		if !semantic.doesScopeContain() {
			err := fmt.Errorf("assignment to undeclared variable/parameter %s", concrete.Name)
			log.Fatal(err)
		}
	case AssignmentNode:
		semantic.CurrInit = concrete.Name
		if !semantic.doesScopeContain() {
			err := fmt.Errorf("assignment to undeclared variable %s", concrete.Name)
			log.Fatal(err)
		}
		semantic.CurrInit = ""
	case VariableNode:
		if semantic.CurrInit == concrete.Name {
			err := fmt.Errorf("self assignment to %s is undefined behavior and not allowed", concrete.Name)
			log.Fatal(err)
		}
		semantic.CurrInit = concrete.Name
		if !semantic.doesScopeContain() {
			err := fmt.Errorf("use of undeclared variable %s", concrete.Name)
			log.Fatal(err)
		}
		semantic.CurrInit = ""
	}
}

func (semantic *Semantic)MiddleAction(node Node) {
	// switch concrete := node.(type) {

	// }
}

func (semantic *Semantic)ExitAction(node Node) {
	switch concrete := node.(type) {
		case FuncDeclNode:
			if semantic.Returned == false && concrete.Type.Type != "void" {
				err := fmt.Errorf("function %s is typed and expecting a return. violates paired => rule", concrete.Name)
				log.Fatal(err)
			}
			semantic.Returned = false
			semantic.InFunc = false
		case ParamNode:
			semantic.CurrInit = ""
		case ScopeNode:
			semantic.CurrScope = semantic.CurrScope.Parent
		case DeclarationNode:
			semantic.CurrInit = ""
	}
}

func (semantic *Semantic)doesScopeContain() bool {
	tempScope := semantic.CurrScope
	for tempScope != nil {
		_, isParam := tempScope.Params[semantic.CurrInit]
		_, isVar := tempScope.Variables[semantic.CurrInit]
		if isVar || isParam { return true }
		tempScope = tempScope.Parent
	}
	return false
}
