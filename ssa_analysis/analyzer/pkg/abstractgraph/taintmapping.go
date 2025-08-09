package abstractgraph

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"analyzer/pkg/common"
)

// TaintMapping is a mapping between primary taint (key) that already existed
// and the list of secondary taints (value) that were recently propagated
//
// in practice, this means that all objects with primary taint will also inherit
// the secondary taint and possibly originate a foreign key when written to the database
type TaintMapping struct {
	mapping map[AbstractTaint][]AbstractTaint // do not use pointers so that we can compare easily
}

func NewTaintMapping() *TaintMapping {
	return &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)}
}

func (tm *TaintMapping) GetMapping() map[AbstractTaint][]AbstractTaint {
	return tm.mapping
}

func (tm *TaintMapping) AddIfNotExists(key AbstractTaint, valElem AbstractTaint) {
	if !slices.Contains(tm.mapping[key], valElem) {
		tm.mapping[key] = append(tm.mapping[key], valElem)
	}
	/* if existingElems, exists := tm.mapping[key]; exists {
		if !slices.Contains(existingElems, valElem) {
			tm.mapping[key] = append(existingElems, valElem)
		}
	} else {
		tm.mapping[key] = []AbstractTaint{valElem}
	} */
}

func (tm *TaintMapping) Merge(other *TaintMapping) {
	for otherKey, otherValLst := range other.mapping {
		for _, otherValElem := range otherValLst {
			tm.AddIfNotExists(otherKey, otherValElem)
		}
	}
}

func (tm *TaintMapping) String() string {
	if len(tm.mapping) == 0 {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{\n")

	keys := make([]AbstractTaint, 0, len(tm.mapping))
	for k := range tm.mapping {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].GetDatabasePath() < keys[j].GetDatabasePath()
	})

	for _, key := range keys {
		var valsStr []string
		for _, val := range tm.mapping[key] {
			valsStr = append(valsStr, val.GetDatabasePath())
		}

		sort.Strings(valsStr)

		builder.WriteString(fmt.Sprintf("  %s: [%s]\n", key.GetDatabasePath(), strings.Join(valsStr, ", ")))
	}

	builder.WriteString("}")
	return builder.String()
}

func MergeTaints(obj *AbstractObject, otherTaintsMap map[string][]*AbstractTaint, primary bool, traced bool) *TaintMapping {
	fmt.Printf("[TAINTMAPPING] merging taints (primary = %t): %v\n", primary, otherTaintsMap)
	var taintMapping *TaintMapping

	if !primary {
		taintMapping = &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)}
	}

	//TODO: deal with upper/lower paths
	for objpath, otherTaints := range otherTaintsMap {
		existingTaints := obj.taints[objpath]

		exists := func(t *AbstractTaint) (string, bool) {
			for _, e := range existingTaints {
				if e.Equals(t) {
					return objpath, true
				}
				if ok, subpath := e.IsUpperPath(t); ok {
					return objpath + subpath, false // must be false so that we create a new abstract taint
				}
			}
			return objpath, false
		}

		for _, otherTaint := range otherTaints {
			if objpath, ok := exists(otherTaint); !ok {
				// need to create new AbstractTaint to avoid just storing the pointer and modifying its fields
				newTaint := NewAbstractTaint(
					otherTaint.dbpath, 
					otherTaint.dbcallID, 
					otherTaint.dbOpType,
					primary, traced,
				)

				// trace info for arguments and (especially) returns
				// can still lead to secondary taints that we still want to track
				if !primary {
					for _, existingTaint := range obj.taints[objpath] {
						// filter by writes to reduce number of foreign keys for now
						if existingTaint.IsPrimary() || traced /* && existingTaint.IsWrite() */ {
							taintMapping.AddIfNotExists(*existingTaint, *newTaint)
						}
					}
				}
				
				fmt.Printf("\t\t[TAINTMAPPING] [DATABASE] adding new taint (%s, traced=%t) on obj path (%s): %v\n", common.OperationTypeToString(newTaint.dbOpType), newTaint.traced, objpath, newTaint)
				obj.taints[objpath] = append(obj.taints[objpath], newTaint)
			}
		}
	}
	return taintMapping
}
