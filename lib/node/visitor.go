package node

import (
	"fmt"
	"strings"
)

type Printer struct {
	TabDepth int
}

func (printer Printer)TabHelper() {
	fmt.Print(strings.Repeat("   ", printer.TabDepth))
}

type Declarator struct {
	scopes []Scope
}

type Visitor interface {
	Action(currNode Node)
	MiddleAction(currNode Node)
	ExitAction(currNode Node)
}

