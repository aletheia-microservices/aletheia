package abstractgraph

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"analyzer/pkg/common"
	"analyzer/pkg/utils"
)

type AbstractObject struct {
	name     string   // ssa name
	allNames []string // all ssa value names for debugging purposes
	taints   map[string][]*AbstractTaint
	traces   map[string][]*AbstractTrace
}

func NewAbstractObject(ssaStr string, directTaints map[string][]*AbstractTaint, traces map[string][]*AbstractTrace) *AbstractObject {
	obj := &AbstractObject{
		name:   ssaStr,
		taints: directTaints,
		traces: traces,
	}
	return obj
}

func (obj *AbstractObject) addToAllNames(name string) {
	obj.allNames = append(obj.allNames, name)
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

func (obj *AbstractObject) GetTaintsForCurrentAndLowerPaths(currPath string) (map[string][]*AbstractTaint, map[string]string) {
	var taints map[string][]*AbstractTaint = make(map[string][]*AbstractTaint)
	var pathsDiffs map[string]string = make(map[string]string)
	for existingPath, taintsLst := range obj.taints {
		if existingPath == currPath {
			taints[existingPath] = append(taints[existingPath], taintsLst...)
			pathsDiffs[existingPath] = ""
		} else {
			if ok, diff := utils.IsUpperPath(currPath, existingPath); ok {
				taints[existingPath] = append(taints[existingPath], taintsLst...)
				pathsDiffs[existingPath] = diff
			}
		}
	}
	return taints, pathsDiffs
}

func (obj *AbstractObject) SetTaintsForObjectPath(objpath string, taints []*AbstractTaint) {
	obj.taints[objpath] = taints
}

func (obj *AbstractObject) GetTraces() map[string][]*AbstractTrace {
	return obj.traces
}

func (obj *AbstractObject) SetTracesForObjectPath(objpath string, traces []*AbstractTrace) {
	obj.traces[objpath] = traces
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
func (obj *AbstractObject) Annotations() string {
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
		builder.WriteString(objpath)
		builder.WriteByte('\n')

		sortedTaints := obj.GetTaintsForObjectPath(objpath)
		sort.Slice(sortedTaints, func(i, j int) bool {
			return utils.LessT(sortedTaints[i].GetT(), sortedTaints[j].GetT())
		})

		for _, taint := range sortedTaints {
			builder.WriteString("[")
			builder.WriteString(common.OperationTypeToString(taint.dbOpType))

			if taint.IsTraced() {
				builder.WriteString(", traced]")
			} else if taint.IsPrimary() {
				builder.WriteString(", primary]")
			} else {
				builder.WriteString(", secondary]")
			}

			if taint.IsReadKey() {
				builder.WriteString(" [K]")
			} else if taint.IsReadValue() {
				builder.WriteString(" [V]")
			}

			builder.WriteString(fmt.Sprintf(" [%s]", taint.GetT()))

			builder.WriteString(" @ ")
			builder.WriteString(taint.String())
			builder.WriteByte('\n')
		}

		sortedTraces := obj.GetTracesForObjectPath(objpath)
		sort.Slice(sortedTraces, func(i, j int) bool {
			return utils.LessT(sortedTraces[i].GetT(), sortedTraces[j].GetT())
		})

		for _, trace := range sortedTraces {
			builder.WriteString("[rpc]")
			builder.WriteString(fmt.Sprintf(" [%s]", trace.GetT()))
			builder.WriteString(" @ ")
			builder.WriteString(trace.String())
			builder.WriteByte('\n')
		}
	}
	return builder.String()
}

func (obj *AbstractObject) GetAllTaints() map[string][]*AbstractTaint {
	return obj.taints
}

func (obj *AbstractObject) GetAllTaintsBeforeT(otherT string) map[string][]*AbstractTaint {
	var taints = make(map[string][]*AbstractTaint)
	for objpath, taintLst := range obj.taints {
		for _, taint := range taintLst {
			if utils.LessT(taint.GetT(), otherT) || utils.EqualT(taint.GetT(), otherT) {
				taints[objpath] = append(taints[objpath], taint)

			}
		}
	}
	return taints
}

func (obj *AbstractObject) GetAllTaintsAfterT(otherT string) map[string][]*AbstractTaint {
	var taints = make(map[string][]*AbstractTaint)
	for objpath, taintLst := range obj.taints {
		for _, taint := range taintLst {
			if utils.GreaterT(taint.GetT(), otherT) || utils.EqualT(taint.GetT(), otherT) {
				taints[objpath] = append(taints[objpath], taint)
			}
		}
	}
	return taints
}

func (obj *AbstractObject) GetAllTraces() map[string][]*AbstractTrace {
	return obj.traces
}

func (obj *AbstractObject) GetAllAbstractLocationsWithTraces() []string {
	var locations []string
	for objpath := range obj.traces {
		locations = append(locations, objpath)
	}
	// sort in reverse from lower locations (more specific) to upper locations (more general)
	sort.Sort(sort.Reverse(sort.StringSlice(locations)))
	return locations
}

func (obj *AbstractObject) GetAllAbstractLocationsWithTaints() []string {
	var locations []string
	for objpath := range obj.taints {
		locations = append(locations, objpath)
	}
	// sort in reverse from lower locations (more specific) to upper locations (more general)
	sort.Sort(sort.Reverse(sort.StringSlice(locations)))
	return locations
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

func (obj *AbstractObject) GetAffectedDatabaseFieldsForCall(callID string) []string {
	var fieldpaths []string
	for _, taintLst := range obj.GetPrimaryTaints() {
		for _, taint := range taintLst {
			if taint.GetDatabaseCallID() == callID {
				if !slices.Contains(fieldpaths, taint.fieldpath) {
					fieldpaths = append(fieldpaths, taint.fieldpath)
				}
			}
		}
	}
	return fieldpaths
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

func (obj *AbstractObject) GetAllTaintsFlatList() []*AbstractTaint {
	var lst []*AbstractTaint
	for _, taints := range obj.taints {
		lst = append(lst, taints...)
	}
	return lst
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
		// write index for in-place compaction
		w := 0

		for _, taint := range taints {
			if taint.IsPrimary() {
				taints[w] = taint
				w++
			}
		}

		if w == 0 {
			// no primary taints left
			delete(obj.taints, objpath)
		} else if w < len(taints) {
			// some removed: shrink slice
			obj.taints[objpath] = taints[:w]
		}
		// else: all were primary, slice unchanged
	}
}

// argument 'other' must not be a pointer because the objective is to compare taints with same content
func (obj *AbstractObject) FindObjectPathWithEqualOrUpperTaint(other AbstractTaint) (string, bool) {
	for objpath, taintLst := range obj.GetAllTaints() {
		for _, taint := range taintLst {
			if taint.Similar(&other) {
				return objpath, true
			}
			// taint.dbfield: notification
			// other.dbfield: notification.PostID
			if ok, subpath := taint.IsUpperTaint(&other); ok {
				return objpath + subpath, true
			}
		}
	}
	return "", false
}

// argument 'other' must not be a pointer because the objective is to compare taints with same content
func (obj *AbstractObject) HasSimilarTaint(other AbstractTaint) bool {
	for _, taintLst := range obj.GetAllTaints() {
		for _, taint := range taintLst {
			if taint.Similar(&other) {
				return true
			}
		}
	}
	return false
}

// argument 'other' must not be a pointer because the objective is to compare taints with same content
func (obj *AbstractObject) HasSimilarTaintOnObjectPath(objpath string, other AbstractTaint) bool {
	for _, taint := range obj.GetTaintsForObjectPath(objpath) {
		if taint.Similar(&other) {
			return true
		}
	}
	return false
}

// argument 'other' must not be a pointer because the objective is to compare taints with same content
func (obj *AbstractObject) HasEqualTaint(objpath string, other AbstractTaint) bool {
	for _, taint := range obj.GetTaintsForObjectPath(objpath) {
		if taint.EqualExceptReadKeyAndReadVal(&other) {
			return true
		}
	}
	return false
}

// argument 'newtaint' must not be a pointer because the objective is is to compare taints with the same content
func (obj *AbstractObject) AddTaintIfSimilarNotExists(objpath string, newtaint AbstractTaint) {
	exists := obj.HasSimilarTaint(newtaint)
	if !exists {
		taint := newtaint.Copy()
		obj.taints[objpath] = append(obj.taints[objpath], taint)
	}
}

func (obj *AbstractObject) AddTaintIfNotExists(objpath string, newtaint *AbstractTaint) bool {
	exists := obj.HasEqualTaint(objpath, *newtaint)
	if !exists {
		obj.taints[objpath] = append(obj.taints[objpath], newtaint)
		return false
	}
	return true
}
