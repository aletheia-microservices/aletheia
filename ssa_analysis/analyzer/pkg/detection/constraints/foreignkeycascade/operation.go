package foreignkeycascade

import (
	"analyzer/pkg/abstractgraph"
)

type DeleteOperation struct {
	call      *abstractgraph.AbstractEdge
	arguments []*abstractgraph.AbstractObject
	database string
	schema string
}

func NewDeleteOperation(call *abstractgraph.AbstractEdge, arguments []*abstractgraph.AbstractObject, database string, schema string) *DeleteOperation {
	return &DeleteOperation{
		call:      call,
		arguments: arguments,
		database:  database,
		schema:    schema,
	}
}

func (op *DeleteOperation) GetCallID() string {
	return op.call.GetID()
}
