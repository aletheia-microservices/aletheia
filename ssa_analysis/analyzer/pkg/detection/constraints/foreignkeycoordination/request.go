package foreignkeycoordination

import "analyzer/pkg/abstractgraph"

type Request struct {
	idx   int
	reads []*ReadOperation
	entry *abstractgraph.AbstractNode
}

func NewRequest(idx int, entry *abstractgraph.AbstractNode) *Request {
	return &Request{
		idx:   idx,
		entry: entry,
	}
}

func (req *Request) AddRead(read *ReadOperation) {
	req.reads = append(req.reads, read)
}

func (req *Request) GetAllReads() []*ReadOperation {
	return req.reads
}

func (req *Request) FindOperationByCallID(callID string) *ReadOperation {
	for _, op := range req.reads {
		if op.GetCallID() == callID {
			return op
		}
	}
	return nil
}
