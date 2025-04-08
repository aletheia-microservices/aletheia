package foreign_key_cascade

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
)

type deleteOperation struct {
	call         *abstractgraph.AbstractDatabaseCall
	datastore    *datastores.Datastore
	dependencies []*deleteDependency
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

func (op *deleteOperation) getDependencies() []*deleteDependency {
	return op.dependencies
}

func (op *deleteOperation) getDependenciesWithMissingCascade() []*deleteDependency {
	var lst []*deleteDependency
	for _, dep := range op.dependencies {
		if !dep.hasCascadingAction() {
			lst = append(lst, dep)
		}
	}
	return lst
}

func (op *deleteOperation) addDependency(dep *deleteDependency) {
	op.dependencies = append(op.dependencies, dep)
}

func (op *deleteOperation) addDependencyIfNotExists(dep *deleteDependency) {
	if !op.hasDependency(dep) {
		op.addDependency(dep)
		//logger.Logger.Debugf("[CASCADE DETECTOR] added dependency %s to %s", dep.String(), op.String())
	}
}

func (op *deleteOperation) hasDependency(other *deleteDependency) bool {
	for _, dep := range op.getDependencies() {
		if dep.hasDatastore(other.datastore) {
			return true
		}
	}
	return false
}

func (op *deleteOperation) getDependency(datastore *datastores.Datastore) *deleteDependency {
	for _, dep := range op.getDependencies() {
		if dep.hasDatastore(datastore) {
			return dep
		}
	}
	logger.Logger.Warnf("[CASCADE DETECTOR] could not find dependency for datastore (%s) for origin delete operation %s", datastore.Name, op.String())
	return nil
}

func (op *deleteOperation) String() string {
	return fmt.Sprintf("datastore = %s", op.datastore.Name)
}
