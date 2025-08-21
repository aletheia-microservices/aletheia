package keycoordination

import (
	"analyzer/pkg/abstractgraph"
)

type ReadOperation struct {
	call      *abstractgraph.AbstractEdge
	arguments []*abstractgraph.AbstractObject
}

func NewReadOperation(call *abstractgraph.AbstractEdge, arguments []*abstractgraph.AbstractObject) *ReadOperation {
	return &ReadOperation{
		call:      call,
		arguments: arguments,
	}
}

func (op *ReadOperation) GetCallID() string {
	return op.call.GetID()
}
