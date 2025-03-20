package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
)

type TupleObject struct {
	Object
	ObjectInfo *ObjectInfo
	Objects    []Object
}

func NewTupleObject(info *ObjectInfo) *TupleObject {
	return &TupleObject{
		ObjectInfo: info,
	}
}

func (v *TupleObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo *ObjectInfo `json:"tuple"`
		Variables  []Object    `json:"tuple_variables,omitempty"`
	}{
		ObjectInfo: v.ObjectInfo,
		Variables:  v.Objects,
	})
}

func (v *TupleObject) GetVariableInfo() *ObjectInfo {
	return nil
}

func (v *TupleObject) AddObjectAndType(variable Object) {
	logger.Logger.Warnf("[VARS TUPLE] adding variable (%s) to tuple variable (%s)", variable.String(), v.String())
	v.Objects = append(v.Objects, variable)
	v.GetTupleType().AddType(variable.GetType())
}

func (v *TupleObject) GetVariables() []Object {
	return v.Objects
}

func (v *TupleObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
	for i := 0; i < len(v.Objects); i++ {
		v.AddReferenceWithID(target.GetElementAt(i), creator)
	}
}

func (v *TupleObject) MergeTupleVariable(tupleVariable *TupleObject) {
	logger.Logger.Warnf("[VARS TUPLE] merging tuple variable (%s) to (%s)", tupleVariable.String(), v.String())
	v.Objects = append(v.Objects, tupleVariable.GetVariables()...)
}

func (v *TupleObject) GetVariableAt(index int) Object {
	if index < len(v.Objects) {
		return v.Objects[index]
	}
	logger.Logger.Fatalf("index (%d) out of bounds for tuple variable (%s) -- variables: %v", index, v.LongString(), v.Objects)
	return nil
}

func (v *TupleObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *TupleObject) NumVariables() int {
	return len(v.Objects)
}

func (v *TupleObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for tuple variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *TupleObject) GetTupleType() *gotypes.TupleType {
	return v.ObjectInfo.GetType().(*gotypes.TupleType)
}

func (v *TupleObject) GetDependencies() []Object {
	return append(v.GetVariableInfo().GetDependencies(), v.Objects...)
}

func (v *TupleObject) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	for _, elem := range v.GetDependencies() { // to include underlying dependencies from variable info
		logger.Logger.Debugf("[TUPLE OBJECT] GOT NESTED DEP (INC/ VINFO) FOR ELEM %s (%s)", elem.String(), VariableTypeName(elem))
		deps = append(deps, elem.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}

func (v *TupleObject) String() string {
	s := "("
	for i, elem := range v.Objects {
		s += elem.String()
		if i < len(v.Objects)-1 {
			s += ", "
		}
	}
	return s + ")"
}

func (v *TupleObject) LongString() string {
	// USE ONLY STRING HERE TO SIMPLIFY DEBUG
	s := v.ObjectInfo.String() + " = ("
	for i, elem := range v.Objects {
		s += elem.String()
		if i < len(v.Objects)-1 {
			s += ", "
		}
	}
	return s + ")"
}

func (v *TupleObject) NewObject() Object {
	newObject := NewTupleObject(NewObjectInfo(v.GetType()))
	// (1) new array points to the same elements as the original one
	// (2) we don't add new dependency edge from new to old object
	// because we keep the same elements
	for _, elem := range v.Objects {
		newObject.AddObjectAndType(elem)
		elem.GetVariableInfo().AddParent(elem, newObject)
	}
	return newObject
}

func (v *TupleObject) Copy(force bool) Object {
	copy := &TupleObject{ObjectInfo: v.ObjectInfo.Copy(force)}
	for _, v := range v.Objects {
		newElem := v.Copy(force)
		copy.Objects = append(copy.Objects, newElem)
		newElem.GetVariableInfo().AddParent(newElem, copy)
	}
	return copy
}

func (v *TupleObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS TUPLE - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &TupleObject{ObjectInfo: v.ObjectInfo.DeepCopy()}
	for _, v := range v.Objects {
		newElem := v.DeepCopy()
		copy.Objects = append(copy.Objects, newElem)
		newElem.GetVariableInfo().AddParent(newElem, copy)
	}
	return copy
}

func (v *TupleObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
		for _, v := range v.Objects {
			variables = append(variables, v.GetUnassaignedVariables()...)
		}
	}
	return variables
}
