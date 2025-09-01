package foreignkeyconcurrency

import (
	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app/backends"
)

type DeleteOperation struct {
	call     *abstractgraph.AbstractEdge
	database *backends.Database
	schema   *backends.Schema
}

type WriteOperation struct {
	call     *abstractgraph.AbstractEdge
	database *backends.Database

	// fields in current database with constraint foreign key + mandatory
	fields  []*backends.Field
	request *Request
}

func NewDeleteOperation(call *abstractgraph.AbstractEdge, database *backends.Database) *DeleteOperation {
	return &DeleteOperation{
		call:     call,
		database: database,
	}
}

func (delete *DeleteOperation) setSchema(schema *backends.Schema) {
	delete.schema = schema
}

func NewWriteOperation(call *abstractgraph.AbstractEdge, database *backends.Database, entry *Request) *WriteOperation {
	return &WriteOperation{
		call:     call,
		database: database,
		request:  entry,
	}
}

func (write *WriteOperation) SetFields(fields []*backends.Field) {
	write.fields = fields
}
