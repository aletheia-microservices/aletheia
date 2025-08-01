package abstractgraph

import (
	"fmt"
	"slices"
	"sort"
	"strings"
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

func (tm *TaintMapping) AddIfNotExists(key AbstractTaint, valElem AbstractTaint) {
	if !slices.Contains(tm.mapping[key], valElem) {
		tm.mapping[key] = append(tm.mapping[key], valElem)
	}
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
		return keys[i].GetField() < keys[j].GetField()
	})

	for _, key := range keys {
		var valsStr []string
		for _, val := range tm.mapping[key] {
			valsStr = append(valsStr, val.GetField())
		}

		sort.Strings(valsStr)

		builder.WriteString(fmt.Sprintf("  %s: [%s]\n", key.GetField(), strings.Join(valsStr, ", ")))
	}

	builder.WriteString("}")
	return builder.String()
}

func MergeTaints(obj *AbstractObject, otherTaintsMap map[string][]*AbstractTaint, primary bool) *TaintMapping {
	fmt.Printf("[TAINTMAPPING] merging taints (primary = %t): %v\n", primary, otherTaintsMap)
	var taintMapping *TaintMapping

	if !primary {
		taintMapping = &TaintMapping{mapping: make(map[AbstractTaint][]AbstractTaint)}
	}

	for objpath, otherTaints := range otherTaintsMap {
		existingTaints := obj.taints[objpath]

		exists := func(t *AbstractTaint) bool {
			for _, e := range existingTaints {
				if e.Equals(t) {
					return true
				}
			}
			return false
		}

		for _, otherTaint := range otherTaints {
			if !exists(otherTaint) {
				// need to create new AbstractTaint to avoid just storing the pointer and modifying its fields
				newTaint := NewAbstractTaint(otherTaint.dbfield, otherTaint.dbcallID, primary, otherTaint.write)
				fmt.Printf("\t [TAINTMAPPING] adding new taint (write = %t): %v\n", newTaint.write, newTaint)
				obj.taints[objpath] = append(obj.taints[objpath], newTaint)

				if !primary {
					for _, existingTaint := range obj.taints[objpath] {
						// filter by writes to reduce number of foreign keys for now
						if existingTaint.IsPrimary() /* && existingTaint.IsWrite() */ {
							taintMapping.AddIfNotExists(*existingTaint, *newTaint)
						}
					}
				}
			}
		}
	}
	return taintMapping
}
