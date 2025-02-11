package cascade

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

func (op *deleteOperation) getCall() *abstractgraph.AbstractDatabaseCall {
	return op.call
}

func (op *deleteOperation) getDependencies() []*deleteDependency {
	return op.dependencies
}

func (op *deleteOperation) addDependency(dep *deleteDependency) {
	op.dependencies = append(op.dependencies, dep)
}

func (op *deleteOperation) hasDependency(other *deleteDependency) bool {
	for _, dep := range op.getDependencies() {
		if dep.isEqual(other) {
			return true
		}
	}
	return false
}

func (op *deleteOperation) getDependency(serviceName string, datastore *datastores.Datastore) *deleteDependency {
	for _, dep := range op.getDependencies() {
		if dep.hasServiceAndDatastore(serviceName, datastore) {
			return dep
		}
	}
	logger.Logger.Warnf("[CASCADE DETECTOR] could not find dependency for service (%s) and datastore (%s) for origin delete operation %s", serviceName, datastore.Name, op.String())
	return nil
}

func (op *deleteOperation) String() string {
	return fmt.Sprintf("datastore = %s", op.datastore.Name)
}
