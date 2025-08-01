package abstractgraph

import (
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

	var builder strings.Builder
	for objpath, taints := range obj.taints {
		builder.WriteString(objpath)
		builder.WriteByte('\n')
		for _, taint := range taints {
			builder.WriteString("[")
			if taint.IsPrimary() {
				builder.WriteString("P,")
			} else {
				builder.WriteString("S,")
			}
			if taint.IsWrite() {
				builder.WriteString("W]")
			} else {
				builder.WriteString("R]")
			}
			builder.WriteString(" @ ")
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
