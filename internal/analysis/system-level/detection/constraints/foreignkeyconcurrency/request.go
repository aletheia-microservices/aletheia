package foreignkeyconcurrency

import "github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"

type Request struct {
	idx     int
	deletes []*DeleteOperation
	writes  []*WriteOperation
	entry   *abstractgraph.AbstractNode
}

func NewRequest(idx int, entry *abstractgraph.AbstractNode) *Request {
	return &Request{
		idx:   idx,
		entry: entry,
	}
}

func (req *Request) addDeleteOperation(delete *DeleteOperation) {
	req.deletes = append(req.deletes, delete)
}

func (req *Request) addWriteOperation(write *WriteOperation) {
	req.writes = append(req.writes, write)
}

func (req *Request) getAllDeleteOperations() []*DeleteOperation {
	return req.deletes
}

func (req *Request) getAllWriteOperations() []*WriteOperation {
	return req.writes
}
