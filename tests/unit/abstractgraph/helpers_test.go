package abstractgraph_test

import (
	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
)

func primaryWrite(t string, dbpath string, callID string) *abstractgraph.AbstractTaint {
	return abstractgraph.NewAbstractTaint(t, dbpath, callID, common.OP_WRITE, true, false, false, false)
}

func secondaryWrite(t string, dbpath string, callID string) *abstractgraph.AbstractTaint {
	return abstractgraph.NewAbstractTaint(t, dbpath, callID, common.OP_WRITE, false, false, false, false)
}

func newObject(taints map[string][]*abstractgraph.AbstractTaint) *abstractgraph.AbstractObject {
	return newTracedObject("t0", taints, nil)
}

func newTracedObject(name string, taints map[string][]*abstractgraph.AbstractTaint, traces map[string][]*abstractgraph.AbstractTrace) *abstractgraph.AbstractObject {
	if taints == nil {
		taints = make(map[string][]*abstractgraph.AbstractTaint)
	}
	if traces == nil {
		traces = make(map[string][]*abstractgraph.AbstractTrace)
	}
	return abstractgraph.NewAbstractObject(name, taints, traces)
}

// mappingPairs returns the taint mapping as "<key db path> -> <value db path>" pairs, in mapping order
func mappingPairs(tm *abstractgraph.TaintMapping) []string {
	var pairs []string
	for _, key := range tm.GetMappingKeys() {
		for _, val := range tm.GetMappingForKey(key) {
			pairs = append(pairs, key.GetDatabasePath()+" -> "+val.GetDatabasePath())
		}
	}
	return pairs
}
