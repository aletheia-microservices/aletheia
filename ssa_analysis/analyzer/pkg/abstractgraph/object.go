package abstractgraph

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"analyzer/pkg/common"
)

type AbstractObject struct {
	name   string // ssa name
	taints map[string][]*AbstractTaint
	traces map[string][]*AbstractTrace
}

func NewAbstractObject(ssaStr string, directTaints map[string][]*AbstractTaint, traces map[string][]*AbstractTrace) *AbstractObject {
	obj := &AbstractObject{
		name:   ssaStr,
		taints: directTaints,
		traces: traces,
	}
	return obj
}

func (obj *AbstractObject) GetName() string {
	return obj.name
}

func (obj *AbstractObject) GetTaints() map[string][]*AbstractTaint {
	return obj.taints
}

func (obj *AbstractObject) GetTaintsForObjectPath(objpath string) []*AbstractTaint {
	return obj.taints[objpath]
}

func (obj *AbstractObject) GetTraces() map[string][]*AbstractTrace {
	return obj.traces
}

func (obj *AbstractObject) GetTracesForObjectPath(objpath string) []*AbstractTrace {
	return obj.traces[objpath]
}

func (obj *AbstractObject) String() string {
	return obj.name
}

func (obj *AbstractObject) IsTainted() bool {
	return len(obj.taints) > 0
}

func (obj *AbstractObject) IsTraced() bool {
	return len(obj.traces) > 0
}

func (obj *AbstractObject) TaintLongString() string {
	if len(obj.taints) == 0 {
		return ""
	}
	var objpaths []string
	for objpath := range obj.taints {
		objpaths = append(objpaths, objpath)
	}
	for objpath := range obj.GetTraces() {
		if !slices.Contains(objpaths, objpath) {
			objpaths = append(objpaths, objpath)
		}
	}
	sort.Strings(objpaths)

	var builder strings.Builder
	for _, objpath := range objpaths {
		for _, taint := range obj.taints[objpath] {
			builder.WriteString("\t" + objpath + " @ " + taint.LongString() + "\n")
		}
		for _, trace := range obj.traces[objpath] {
			builder.WriteString("\t" + objpath + " @ " + trace.LongString() + "\n")
		}
	}
	return builder.String()
}

// same logic as SSAGraph Node
func (obj *AbstractObject) TaintString() string {
	if len(obj.taints) == 0 && len(obj.traces) == 0 {
		return ""
	}

	var objpaths []string
	for objpath := range obj.GetTaints() {
		objpaths = append(objpaths, objpath)
	}
	for objpath := range obj.GetTraces() {
		if !slices.Contains(objpaths, objpath) {
			objpaths = append(objpaths, objpath)
		}
	}
	sort.Strings(objpaths)

	var builder strings.Builder
	for _, objpath := range objpaths {
		taints := obj.GetTaintsForObjectPath(objpath)
		builder.WriteString(objpath)
		builder.WriteByte('\n')
		for _, taint := range taints {
			builder.WriteString("[")
			builder.WriteString(common.OperationTypeToString(taint.dbOpType))

			if taint.IsTraced() {
				builder.WriteString(", traced]")
			} else if taint.IsPrimary() {
				builder.WriteString(", primary]")
			} else {
				builder.WriteString(", secondary]")
			}

			builder.WriteString(" @ ")
			builder.WriteString(taint.String())
			builder.WriteByte('\n')
		}

		for _, taint := range obj.GetTracesForObjectPath(objpath) {
			builder.WriteString("[rpc] @ ")
			builder.WriteString(taint.String())
			builder.WriteByte('\n')
		}
	}
	return builder.String()
}

func (obj *AbstractObject) GetAllTaints() map[string][]*AbstractTaint {
	return obj.taints
}

func (obj *AbstractObject) GetWriteTaints() map[string][]*AbstractTaint {
	writeTaints := make(map[string][]*AbstractTaint, 0)
	for objpath, taints := range obj.taints {
		for _, taint := range taints {
			if taint.IsWrite() {
				writeTaints[objpath] = append(writeTaints[objpath], taint)
			}
		}
	}
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

func (obj *AbstractObject) GetPrimaryTaintsFlatList() []*AbstractTaint {
	var lst []*AbstractTaint
	for _, taints := range obj.taints {
		for _, taint := range taints {
			if taint.IsPrimary() {
				lst = append(lst, taint)
			}
		}
	}
	return lst
}

func (obj *AbstractObject) GetSecondaryTaintsFlatList() []*AbstractTaint {
	var lst []*AbstractTaint
	for _, taints := range obj.taints {
		for _, taint := range taints {
			if !taint.IsPrimary() {
				lst = append(lst, taint)
			}
		}
	}
	return lst
}

func (obj *AbstractObject) CleanSecondaryTaints() {
	for objpath, taints := range obj.taints {
		var filtered []*AbstractTaint
		for _, taint := range taints {
			if taint.IsPrimary() {
				filtered = append(filtered, taint)
			}
		}
		if len(filtered) > 0 {
			obj.taints[objpath] = filtered
		} else {
			delete(obj.taints, objpath)
		}
	}
}

// argument 'other' must not be a pointer because the objective is to compare taints with same content
func (obj *AbstractObject) FindObjectPathWithEqualOrUpperTaint(other AbstractTaint) (string, bool) {
	fmt.Printf("[ABSTRACT OBJECT] finding object path with equal taint\n")
	for objpath, taintLst := range obj.GetAllTaints() {
		for _, taint := range taintLst {
			if taint.Equals(&other) {
				return objpath, true
			}
			// taint.dbfield: notification
			// other.dbfield: notification.PostID
			if ok, subpath := taint.IsUpperPath(&other); ok {
				return objpath + subpath, true
			}
		}
	}
	return "", false
}

// argument 'other' must not be a pointer because the objective is to compare taints with same content
func (obj *AbstractObject) HasEqualTaint(other AbstractTaint) bool {
	fmt.Printf("[ABSTRACT OBJECT] finding object path with equal taint\n")
	for _, taintLst := range obj.GetAllTaints() {
		for _, taint := range taintLst {
			if taint.Equals(&other) {
				return true
			}
		}
	}
	return false
}

// argument 'newtaint' must not be a pointer because the objective is is to compare taints with the same content
func (obj *AbstractObject) AddTaintIfNotExists(objpath string, newtaint AbstractTaint) {
	fmt.Printf("[ABSTRACT OBJECT] propagate new taint for obj path (%s): %s\n", objpath, newtaint.LongString())
	exists := obj.HasEqualTaint(newtaint)
	if !exists {
		taint := &AbstractTaint{
			dbpath:   newtaint.dbpath,
			dbcallID: newtaint.dbcallID,
			primary:  newtaint.primary,
			dbOpType: newtaint.dbOpType,
		}
		obj.taints[objpath] = append(obj.taints[objpath], taint)
		fmt.Printf("[ABSTRACT OBJECT] added new taint to obj path (%s): %s\n", objpath, taint)
	}
}
