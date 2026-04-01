package foreignkeycascade

import (
	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app/backends"
)

type DeleteOperation struct {
	call      *abstractgraph.AbstractEdge
	arguments []*abstractgraph.AbstractObject
	database  string
	schema    string
}

type WriteOperation struct {
	call      *abstractgraph.AbstractEdge
	arguments []*abstractgraph.AbstractObject
	database  string
	schema    string
	fields    []*backends.Field
}

func NewDeleteOperation(call *abstractgraph.AbstractEdge, arguments []*abstractgraph.AbstractObject, database string, schema string) *DeleteOperation {
	return &DeleteOperation{
		call:      call,
		arguments: arguments,
		database:  database,
		schema:    schema,
	}
}

func NewWriteOperation(call *abstractgraph.AbstractEdge, arguments []*abstractgraph.AbstractObject, database string, schema string) *WriteOperation {
	return &WriteOperation{
		call:      call,
		arguments: arguments,
		database:  database,
		schema:    schema,
	}
}

func (op *DeleteOperation) GetCallID() string {
	return op.call.GetID()
}

func (op *WriteOperation) GetCallID() string {
	return op.call.GetID()
}
