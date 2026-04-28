package uniquenessconcurrency

import (
	"analyzer/pkg/abstractgraph"
)

type WriteOperation struct {
	call      *abstractgraph.AbstractEdge
	arguments []*abstractgraph.AbstractObject
}

func NewWriteOperation(call *abstractgraph.AbstractEdge, arguments []*abstractgraph.AbstractObject) *WriteOperation {
	return &WriteOperation{
		call:      call,
		arguments: arguments,
	}
}

func (op *WriteOperation) GetCallID() string {
	return op.call.GetID()
}
