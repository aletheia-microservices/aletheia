package foreign_key_cascade

import (
	"analyzer/pkg/abstractgraph"
)

type RequestInfo struct {
	index            int
	entry            *abstractgraph.AbstractServiceCall
	deleteOperations []*deleteOperation
}

func (info *RequestInfo) addDeleteOperation(op *deleteOperation) {
	info.deleteOperations = append(info.deleteOperations, op)
}

func (info *RequestInfo) getDeleteOperations() []*deleteOperation{
	return info.deleteOperations
}
