package unicity

import (
	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
)

type Operation struct {
	call                *abstractgraph.AbstractDatabaseCall
	datastore           *datastores.Datastore
	constraints         []*datastores.Constraint
	onUnicityConstraint bool
}

func NewOperation(call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore) *Operation {
	return &Operation{
		call:      call,
		datastore: datastore,
	}
}

func NewOperationOnUnicityConstraint(call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore, constraints []*datastores.Constraint) *Operation {
	return &Operation{
		call:                call,
		datastore:           datastore,
		constraints:         constraints,
		onUnicityConstraint: true,
	}
}
