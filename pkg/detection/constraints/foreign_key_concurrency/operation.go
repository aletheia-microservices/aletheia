package foreign_key_concurrency

import (
	"fmt"
	"slices"

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

func (del *delete) getDatastore() *datastores.Datastore {
	return del.datastore
}

type write struct {
	idx           int
	call          *abstractgraph.AbstractDatabaseCall
	datastore     *datastores.Datastore
	writtenFields []*datastores.Field
	writtenObjs   []objects.Object
}

func (w *write) addWrittenObjects(obj ...objects.Object) {
	w.writtenObjs = append(w.writtenObjs, obj...)
}

func (w *write) getWrittenFields() []*datastores.Field {
	return w.writtenFields
}

func (w *write) writesToField(field *datastores.Field) bool {
	return slices.Contains(w.writtenFields, field)
}

func (w *write) getDatastore() *datastores.Datastore {
	return w.datastore
}

func (w *write) getDbCall() *abstractgraph.AbstractDatabaseCall {
	return w.call
}

func NewOperation(idx int, call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore, writtenFields []*datastores.Field) *write {
	return &write{
		idx:           idx,
		call:          call,
		datastore:     datastore,
		writtenFields: writtenFields,
	}
}
