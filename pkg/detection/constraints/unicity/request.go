package unicity

import (
	"analyzer/pkg/abstractgraph"
)

type RequestInfo struct {
	entry             *abstractgraph.AbstractServiceCall
	operations        []*Operation
	writeOnConstraint bool
}

func (info *RequestInfo) addOperation(operation *Operation) {
	info.operations = append(info.operations, operation)
}

func (info *RequestInfo) hasOperations() bool {
	return len(info.operations) > 0
}

func (info *RequestInfo) hasPotentialInconsistencies() bool {
	return len(info.operations) > 1 // only if we have more than 2 ops
}

func (info *RequestInfo) getOperations() []*Operation {
	return info.operations
}
