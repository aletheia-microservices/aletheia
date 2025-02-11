package unicity

import (
	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
)

type Operation struct {
	call                *abstractgraph.AbstractDatabaseCall
	datastore           *datastores.Datastore
	onUnicityConstraint bool
}

func NewOperation(call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore) *Operation {
	return &Operation{
		call:      call,
		datastore: datastore,
	}
}

func NewOperationOnUnicityConstraint(call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore) *Operation {
	return &Operation{
		call:                call,
		datastore:           datastore,
		onUnicityConstraint: true,
	}
}
