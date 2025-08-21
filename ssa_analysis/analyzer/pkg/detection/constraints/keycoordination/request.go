package keycoordination

import "analyzer/pkg/abstractgraph"

type Request struct {
	idx   int
	ops   []*ReadOperation
	entry *abstractgraph.AbstractNode
}

func NewRequest(idx int, entry *abstractgraph.AbstractNode) *Request {
	return &Request{
		idx:   idx,
		entry: entry,
	}
}

func (req *Request) AddOperation(read *ReadOperation) {
	req.ops = append(req.ops, read)
}

func (req *Request) GetAllOperations() []*ReadOperation {
	return req.ops
}

func (req *Request) FindOperationByCallID(callID string) *ReadOperation {
	for _, op := range req.ops {
		if op.GetCallID() == callID {
			return op
		}
	}
	return nil
}
