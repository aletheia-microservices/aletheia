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

// flags potential inconsistency if there is a write on a numerical constraint
func (info *RequestInfo) hasPotentialInconsistencies() bool {
	return info.writeOnConstraint // && len(info.operations) != 1
}

func (info *RequestInfo) getOperations() []*Operation {
	return info.operations
}
