package foreign_key_concurrency

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
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
	return fmt.Sprintf("\t - %s: %s \t %-45s --> %s\n",
		w.call.GetCallerStr(),
		w.call.ShortString(),
		w.field.GetFullName(),
		w.constraint.GetReferencedByField().GetFullName(),
	)
}

func (del *delete) addAffectedWrittenField(call *abstractgraph.AbstractDatabaseCall, field *datastores.Field, constraint *datastores.Constraint) {
	del.affectedWrittenFields = append(del.affectedWrittenFields, &writtenField{
		call:       call,
		field:      field,
		constraint: constraint,
	})
}
