package foreign_key_concurrency

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
)

type delete struct {
	call      *abstractgraph.AbstractDatabaseCall
	datastore *datastores.Datastore
}

type fieldWithReference struct {
	field       *datastores.Field
	constraints []*datastores.Constraint // foreign key constraints only
}

type writeWithReference struct {
	call      *abstractgraph.AbstractDatabaseCall
	datastore *datastores.Datastore
	fields    []*fieldWithReference
}

func (write *writeWithReference) String() string {
	var fieldsWithReferenceStr string
	for _, field := range write.fields {
		fieldsWithReferenceStr += fmt.Sprintf("- field = %s; constraints = %s\n", field.field.GetName(), field.constraints)
	}
	return fmt.Sprintf("write with reference to datastore (%s); fields with reference: \n%v", write.datastore.Name, fieldsWithReferenceStr)
}
