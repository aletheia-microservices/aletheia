package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
)

type GenericObject struct {
	Object
	ObjectInfo *ObjectInfo
	Params     []Object
}

func NewGenericObject(info *ObjectInfo) *GenericObject {
	return &GenericObject{
		ObjectInfo: info,
	}
}

func (v *GenericObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo *ObjectInfo `json:"generic"`
		Params     []Object    `json:"generic_params,omitempty"`
	}{
		ObjectInfo: v.ObjectInfo,
		Params:     v.Params,
	})
}

func (v *GenericObject) String() string {
	return v.ObjectInfo.String()
}

func (v *GenericObject) LongString() string {
	return v.ObjectInfo.LongString()
}

func (v *GenericObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *GenericObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for generic variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *GenericObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *GenericObject) GetDependencies() []Object {
	return append(v.GetVariableInfo().GetDependencies(), v.Params...)
}

func (v *GenericObject) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	for _, elem := range v.GetDependencies() { // to include underlying dependencies from variable info
		logger.Logger.Debugf("[GENERIC OBJECT] GOT NESTED DEP (INC/ VINFO) FOR ELEM %s (%s)", elem.String(), VariableTypeName(elem))
		deps = append(deps, elem.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}

func (v *GenericObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
	if targetGeneric, ok := target.(*GenericObject); ok {
		for i, p := range v.Params {
			if i > len(targetGeneric.Params) {
				logger.Logger.Fatalf("[VARS GENERIC - REF] invalid argument index %d for length %d in params: %v", i, (targetGeneric.Params), targetGeneric.Params)
			} else {
				p.AddReferenceWithID(targetGeneric.Params[i], creator)
			}
		}
	} else {
		logger.Logger.Fatalf("[VARS GENERIC - REF] skipping referenced variables with different types (%s vs %s) (%s vs %s)", v.String(), target.String(), VariableTypeName(v), VariableTypeName(target))
	}
	logger.Logger.Tracef("[VARS GENERIC - REF] added reference (%s) -> (%s) with id = %d (creator: %s)", v.ObjectInfo.Name, target.GetVariableInfo().GetName(), v.ObjectInfo.Id, creator)
}

func (v *GenericObject) NewObject() Object {
	newObject := NewGenericObject(NewObjectInfo(v.GetType()))
	// add dependency edge since object is completely new
	newObject.ObjectInfo.AddDependency(v)
	return newObject
}

func (v *GenericObject) Copy(force bool) Object {
	copy := &GenericObject{
		ObjectInfo: v.ObjectInfo.Copy(force),
	}
	for _, p := range v.Params {
		copy.Params = append(copy.Params, p.Copy(force))
	}
	return copy
}

func (v *GenericObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS GENERIC - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &GenericObject{
		ObjectInfo: v.ObjectInfo.DeepCopy(),
	}
	for _, p := range v.Params {
		copy.Params = append(copy.Params, p.DeepCopy())
	}
	return copy
}

func (v *GenericObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
		for _, p := range v.Params {
			variables = append(variables, p.GetUnassaignedVariables()...)
		}
	}
	return variables
}
