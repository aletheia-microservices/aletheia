package abstractcallgraph

import "fmt"

type AbstractObject struct {
	name           string // ssa name
	taints         map[string][]*AbstractTaint
}

func NewAbstractObject(ssaStr string, directTaints map[string][]*AbstractTaint) *AbstractObject {
	obj := &AbstractObject{
		name:           ssaStr,
		taints:         directTaints,
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

func (obj *AbstractObject) AddSecondaryTaints(taintsMap map[string][]*AbstractTaint) {
	for objpath, taints := range taintsMap {
		// need to create new AbstractTaint to avoid just storing the pointer and modifying its fields
		for _, taint := range taints {
			obj.taints[objpath] = append(obj.taints[objpath], &AbstractTaint{
				dbfield:  taint.dbfield,
				dbcallID: taint.dbcallID,
				primary:  false,
			})
		}
	}
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
