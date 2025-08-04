package foreignkeycascade

import (
	"analyzer/pkg/abstractgraph"
)

type DeleteOperation struct {
	call      *abstractgraph.AbstractEdge
	arguments []*abstractgraph.AbstractObject
}

func NewDeleteOperation(call *abstractgraph.AbstractEdge, arguments []*abstractgraph.AbstractObject) *DeleteOperation {
	return &DeleteOperation{
		call:      call,
		arguments: arguments,
	}
}

func (op *DeleteOperation) GetCallID() string {
	return op.call.GetID()
}
