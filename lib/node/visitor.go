package node

type Visitor interface {
	Action(currNode Node)
	MiddleAction(currNode Node)
	ExitAction(currNode Node)
}

