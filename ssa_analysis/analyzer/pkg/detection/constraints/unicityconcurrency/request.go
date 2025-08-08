package unicityconcurrency

import "analyzer/pkg/abstractgraph"

type Request struct {
	idx   int
	ops   []*WriteOperation
	entry *abstractgraph.AbstractNode
}

func NewRequest(idx int, entry *abstractgraph.AbstractNode) *Request {
	return &Request{
		idx:   idx,
		entry: entry,
	}
}

func (req *Request) AddOperation(op *WriteOperation) {
	req.ops = append(req.ops, op)
}

func (req *Request) GetAllOperations() []*WriteOperation {
	return req.ops
}

func (req *Request) FindOperationByCallID(callID string) *WriteOperation {
	for _, op := range req.ops {
		if op.GetCallID() == callID {
			return op
		}
	}
	return nil
}
