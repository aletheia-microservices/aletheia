package foreign_key_concurrency

import (
	"analyzer/pkg/abstractgraph"
)

type request struct {
	index      int
	entry      *abstractgraph.AbstractServiceCall
	operations []*write
}

func (info *request) addOperation(operation *write) {
	info.operations = append(info.operations, operation)
}

func (info *request) getOperations() []*write {
	return info.operations
}

func (info *request) numOperations() int {
	return len(info.operations)
}
