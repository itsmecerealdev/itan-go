package node

import (
	"fmt"
	"log"
)

func (declarator *Declarator)Action(node Node) {
	switch concrete := node.(type) {
	case FuncDeclNode:
		scope := declarator.scopes[len(declarator.scopes)-1]
		if _, exists := scope.Funcs[concrete.Name]; !exists {
			scope.Funcs[concrete.Name] = concrete
		} else {
			err := fmt.Errorf("double declaration of func %s", concrete.Name)
			log.Fatal(err)
		}
	case ParamNode:
		scope := declarator.scopes[len(declarator.scopes)-1]
		if _, exists := scope.Params[concrete.Name]; !exists {
			scope.Params[concrete.Name] = concrete
		} else {
			err := fmt.Errorf("double declaration of param %s", concrete.Name)
			log.Fatal(err)
		}
	case ScopeNode:
		if len(declarator.scopes) > 0 {
			concrete.Symbols.Parent = &declarator.scopes[len(declarator.scopes)-1]
		}
		newScope := Scope {
			Variables: make(map[string]DeclarationNode),
			Params: make(map[string]ParamNode),
			Funcs: make(map[string]FuncDeclNode),
		}
		declarator.scopes = append(declarator.scopes, newScope)
	case DeclarationNode:
		scope := declarator.scopes[len(declarator.scopes)-1]
		if _, exists := scope.Variables[concrete.Name]; !exists {
			scope.Variables[concrete.Name] = concrete
		} else {
			err := fmt.Errorf("double declaration of var %s", concrete.Name)
			log.Fatal(err)
		}
	}
}

func (declarator *Declarator)MiddleAction(node Node) {

}

func (declarator *Declarator)ExitAction(node Node) {
	switch node.(type) {
	case ScopeNode:
		declarator.scopes = declarator.scopes[:len(declarator.scopes)-1]
	}
}
