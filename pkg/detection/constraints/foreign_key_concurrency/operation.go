package foreign_key_concurrency

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
	"analyzer/pkg/types/objects"
)

type delete struct {
	call                  *abstractgraph.AbstractDatabaseCall
	datastore             *datastores.Datastore
	affectedWrittenFields []*writtenField
}

type writtenField struct {
	call       *abstractgraph.AbstractDatabaseCall
	field      *datastores.Field
	constraint *datastores.Constraint
}

// e.g.,
// - notifications_queue.Push() NOTIFICATIONS_QUEUE.Message.ReqID             --> POSTS_DB.Post.ReqID
// - notifications_queue.Push() NOTIFICATIONS_QUEUE.Message.PostID_MESSAGE    --> POSTS_DB.Post.PostID
func (w *writtenField) String() string {
	return fmt.Sprintf("\t- %-45s ----> foreign key %s\n\t  @ %s: %s",
		w.field.GetFullName(),
		w.constraint.GetReferencedByField().GetFullName(),
		w.call.GetCallerStr(),
		w.call.ShortString(),
	)
}

func (del *delete) flagAffectedWriteOnField(call *abstractgraph.AbstractDatabaseCall, field *datastores.Field, constraint *datastores.Constraint) {
	del.affectedWrittenFields = append(del.affectedWrittenFields, &writtenField{
		call:       call,
		field:      field,
		constraint: constraint,
	})
}

type write struct {
	idx           int
	call          *abstractgraph.AbstractDatabaseCall
	datastore     *datastores.Datastore
	writtenFields []*datastores.Field
	writtenObjs   []objects.Object
}

func (op *write) addWrittenObjects(obj ...objects.Object) {
	op.writtenObjs = append(op.writtenObjs, obj...)
}

func (op *write) getWrittenFields() []*datastores.Field {
	return op.writtenFields
}

func (op *write) getDbCall() *abstractgraph.AbstractDatabaseCall {
	return op.call
}

func NewOperation(idx int, call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore, writtenFields []*datastores.Field) *write {
	return &write{
		idx:           idx,
		call:          call,
		datastore:     datastore,
		writtenFields: writtenFields,
	}
}
