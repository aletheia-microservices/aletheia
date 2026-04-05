package foreignkeycascade

import "analyzer/pkg/abstractgraph"

type Request struct {
	idx      int
	ops      []*DeleteOperation
	writeOps []*WriteOperation
	entry    *abstractgraph.AbstractNode
}

func NewRequest(idx int, entry *abstractgraph.AbstractNode) *Request {
	return &Request{
		idx:   idx,
		entry: entry,
	}
}

func (req *Request) AddOperation(op *DeleteOperation) {
	req.ops = append(req.ops, op)
}

func (req *Request) AddWriteOperation(op *WriteOperation) {
	req.writeOps = append(req.writeOps, op)
}

func (req *Request) GetAllOperations() []*DeleteOperation {
	return req.ops
}

func (req *Request) GetAllWriteOperations() []*WriteOperation {
	return req.writeOps
}

func (req *Request) FindOperationByCallID(callID string) *DeleteOperation {
	for _, op := range req.ops {
		if op.GetCallID() == callID {
			return op
		}
	}
	return nil
}
