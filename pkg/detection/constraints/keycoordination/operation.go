package keycoordination

import (
	"analyzer/pkg/abstractgraph"
)

type ReadOperation struct {
	call      *abstractgraph.AbstractEdge
	arguments []*abstractgraph.AbstractObject
	reqIdx    int
}

func NewReadOperation(call *abstractgraph.AbstractEdge, arguments []*abstractgraph.AbstractObject, reqIdx int) *ReadOperation {
	return &ReadOperation{
		call:      call,
		arguments: arguments,
		reqIdx:    reqIdx,
	}
}

func (op *ReadOperation) GetCallID() string {
	return op.call.GetID()
}
