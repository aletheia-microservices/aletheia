package numerical

import (
	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
)

type Operation struct {
	call                  *abstractgraph.AbstractDatabaseCall
	datastore             *datastores.Datastore
	onNumericalConstraint bool
}

func NewOperation(call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore) *Operation {
	return &Operation{
		call:      call,
		datastore: datastore,
	}
}

func NewOperationOnNumericalConstraint(call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore) *Operation {
	return &Operation{
		call:                  call,
		datastore:             datastore,
		onNumericalConstraint: true,
	}
}
