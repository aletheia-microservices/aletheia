package fkey_concurrency

import (
	"fmt"
	"slices"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
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
		w.constraint.GetReferenceToField().GetFullName(),
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

func (w *write) getWrittenObjectAt(idx int) objects.Object {
	if idx >= len(w.writtenObjs) {
		// force to only object that exists
		// FIXME: this could be cleaner
		if w.datastore.IsNoSQLDatabase() || w.datastore.IsQueue() {
			return w.writtenObjs[0]
		}
		logger.Logger.Fatalf("[FK CONCURRENCY] len (%d) out of bounds for written objects: %v\n // written fields: %v", idx, w.writtenObjs, w.writtenFields)
	}
	return w.writtenObjs[idx]
}

func (w *write) writesToField(field *datastores.Field) (bool, int) {
	if slices.Contains(w.writtenFields, field) {
		for idx, f := range w.writtenFields {
			if f == field {
				return true, idx
			}
		}
	}
	return false, -1
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
