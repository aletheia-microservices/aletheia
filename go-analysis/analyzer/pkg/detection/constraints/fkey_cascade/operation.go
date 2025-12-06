package fkey_cascade

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
)

type deleteOperation struct {
	call           *abstractgraph.AbstractDatabaseCall
	datastore      *datastores.Datastore
	pendingDeletes []*pendingDelete
}

func newDeleteOperation(call *abstractgraph.AbstractDatabaseCall, datastore *datastores.Datastore) *deleteOperation {
	return &deleteOperation{
		call:      call,
		datastore: datastore,
	}

}

func (op *deleteOperation) getCall() *abstractgraph.AbstractDatabaseCall {
	return op.call
}

func (op *deleteOperation) getPendingDeletes() []*pendingDelete {
	return op.pendingDeletes
}

func (op *deleteOperation) getDependenciesWithMissingCascade() []*pendingDelete {
	var lst []*pendingDelete
	for _, dep := range op.pendingDeletes {
		if !dep.hasCascadingAction() {
			lst = append(lst, dep)
		}
	}
	return lst
}

func (op *deleteOperation) addDependency(dep *pendingDelete) {
	op.pendingDeletes = append(op.pendingDeletes, dep)
}

func (op *deleteOperation) addPendingDeleteIfNotExists(dep *pendingDelete) {
	if !op.hasDependency(dep) {
		op.addDependency(dep)
		//logger.Logger.Debugf("[CASCADE DETECTOR] added dependency %s to %s", dep.String(), op.String())
	}
}

func (op *deleteOperation) hasDependency(other *pendingDelete) bool {
	for _, dep := range op.getPendingDeletes() {
		if dep.isOnDatastore(other.datastore) && dep.isOnConstraint(other.constraint) {
			return true
		}
	}
	return false
}

// DEPRECATED
func (op *deleteOperation) getDependency(datastore *datastores.Datastore) *pendingDelete {
	for _, dep := range op.getPendingDeletes() {
		if dep.isOnDatastore(datastore) {
			return dep
		}
	}
	logger.Logger.Warnf("[CASCADE DETECTOR] could not find dependency for datastore (%s) for origin delete operation %s", datastore.Name, op.String())
	return nil
}

func (op *deleteOperation) String() string {
	return fmt.Sprintf("datastore = %s", op.datastore.Name)
}
