package foreign_key_concurrency

import (
	"analyzer/pkg/abstractgraph"
)

type RequestInfo struct {
	index int
	entry *abstractgraph.AbstractServiceCall
}
