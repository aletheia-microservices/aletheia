package unicity

import (
	"analyzer/pkg/abstractgraph"
)

type RequestInfo struct {
	entry         *abstractgraph.AbstractServiceCall
	operations    []*Operation
	inconsistency bool
}

func (info *RequestInfo) addOperation(operation *Operation) {
	info.operations = append(info.operations, operation)
}

// flags potential inconsistency if there is a write on a numerical constraint
func (info *RequestInfo) hasPotentialInconsistencies() bool {
	return info.inconsistency // && len(info.operations) != 1
}

func (info *RequestInfo) flagInconsistency() {
	info.inconsistency = true
}

func (info *RequestInfo) getOperations() []*Operation {
	return info.operations
}

func (info *RequestInfo) numOperations() int {
	return len(info.operations)
}
