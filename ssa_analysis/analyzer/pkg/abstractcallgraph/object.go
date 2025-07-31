package abstractcallgraph

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type AbstractObject struct {
	name   string // ssa name
	taints map[string][]*AbstractTaint
}

func NewAbstractObject(ssaStr string, directTaints map[string][]*AbstractTaint) *AbstractObject {
	obj := &AbstractObject{
		name:   ssaStr,
		taints: directTaints,
	}
	return obj
}

func (obj *AbstractObject) String() string {
	return obj.name
}

func (obj *AbstractObject) IsTainted() bool {
	return len(obj.taints) > 0
}

func (obj *AbstractObject) taintString() string {
	if len(obj.taints) == 0 {
		return ""
	}

	var taintStr string
	for objpath, taints := range obj.taints {
		taintStr += fmt.Sprintf("%s\n", objpath)
		for _, taint := range taints {
			if taint.IsPrimary() {
				taintStr += fmt.Sprintf("[P] @ %s\n", taint.String())
			} else {
				taintStr += fmt.Sprintf("[S] @ %s\n", taint.String())
			}
		}
	}
	return taintStr
}

// TaintMapping is a mapping between primary taint (key) that already existed
// and the list of secondary taints (value) that were recently propagated
//
// in practice, this means that all objects with primary taint will also inherit
// the secondary taint and possibly originate a foreign key when written to the database
type TaintMapping struct {
	mapping map[string][]string
}

func (tm *TaintMapping) AddIfNotExists(key string, valElem string) {
	if !slices.Contains(tm.mapping[key], valElem) {
		tm.mapping[key] = append(tm.mapping[key], valElem)
	}
}

func (tm *TaintMapping) Merge(other *TaintMapping) {
	for otherKey, otherValLst := range other.mapping {
		for _, otherValElem := range otherValLst {
			if !slices.Contains(tm.mapping[otherKey], otherValElem) {
				tm.mapping[otherKey] = append(tm.mapping[otherKey], otherValElem)
			}
		}
	}
}

func (tm *TaintMapping) String() string {
	if len(tm.mapping) == 0 {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{\n")

	keys := make([]string, 0, len(tm.mapping))
	for k := range tm.mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		vals := tm.mapping[key]
		sort.Strings(vals)
		builder.WriteString(fmt.Sprintf("  %s: [%s]\n", key, strings.Join(vals, ", ")))
	}

	builder.WriteString("}")
	return builder.String()
}

func (obj *AbstractObject) MergeTaints(otherTaintsMap map[string][]*AbstractTaint, primary bool) *TaintMapping {
	fmt.Printf("[ABSTARCTOBJECT] merging taints (primary = %t): %v\n", primary, otherTaintsMap)
	var taintMapping *TaintMapping

	if !primary {
		taintMapping = &TaintMapping{mapping: make(map[string][]string)}
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
				newTaint := &AbstractTaint{
					dbfield:  otherTaint.dbfield,
					dbcallID: otherTaint.dbcallID,
					primary:  primary,
				}
				fmt.Printf("\t [ABSTRACTOBJECT] adding new taint: %v\n", newTaint)
				obj.taints[objpath] = append(obj.taints[objpath], newTaint)

				if !primary {
					for _, existingTaint := range obj.taints[objpath] {
						if existingTaint.IsPrimary() {
							taintMapping.AddIfNotExists(existingTaint.dbfield, newTaint.dbfield)
						}
					}
				}
			}
		}
	}
	return taintMapping
}

func (obj *AbstractObject) GetAllTaints() map[string][]*AbstractTaint {
	return obj.taints
}

func (obj *AbstractObject) GetPrimaryTaints() map[string][]*AbstractTaint {
	primaryTaints := make(map[string][]*AbstractTaint, 0)
	for objpath, taints := range obj.taints {
		for _, taint := range taints {
			if taint.IsPrimary() {
				primaryTaints[objpath] = append(primaryTaints[objpath], taint)
			}
		}
	}
	return obj.taints
}

func (obj *AbstractObject) GetSecondaryTaints() map[string][]*AbstractTaint {
	secondaryTaints := make(map[string][]*AbstractTaint, 0)
	for objpath, taints := range obj.taints {
		for _, taint := range taints {
			if !taint.IsPrimary() {
				secondaryTaints[objpath] = append(secondaryTaints[objpath], taint)
			}
		}
	}
	return obj.taints
}
