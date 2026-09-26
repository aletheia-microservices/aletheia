package uniquenessconcurrency

import (
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
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
