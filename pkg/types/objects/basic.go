package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
)

type BasicObject struct {
	Object
	ObjectInfo        *ObjectInfo
	UnderlyingObjects []Object // variables that influence the value of basic variables - e.g. "some_variable" in len(some_variable)
}

func NewBasicObject(info *ObjectInfo) *BasicObject {
	return &BasicObject{ObjectInfo: info}
}

func (v *BasicObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo          *ObjectInfo `json:"basic"`
		UnderlyingVariables []Object    `json:"-"`
	}{
		ObjectInfo:          v.ObjectInfo,
		UnderlyingVariables: v.UnderlyingObjects,
	})
}

func (v *BasicObject) String() string {
	return v.ObjectInfo.String()
}

func (v *BasicObject) LongString() string {
	return v.ObjectInfo.String()
}

func (v *BasicObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *BasicObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for basic variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *BasicObject) GetBasicType() *gotypes.BasicType {
	if userType, ok := v.ObjectInfo.GetType().(*gotypes.UserType); ok {
		return userType.UserType.(*gotypes.BasicType)
	}
	return v.ObjectInfo.GetType().(*gotypes.BasicType)
}

func (v *BasicObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *BasicObject) GetDependencies() []Object {
	return append(v.GetVariableInfo().GetDependencies(), v.UnderlyingObjects...)
}

func (v *BasicObject) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	for _, elem := range v.GetDependencies() { // to include underlying dependencies from variable info
		logger.Logger.Debugf("[BASIC OBJECT] GOT NESTED DEP (INC/ VINFO) FOR ELEM %s (%s)", elem.String(), VariableTypeName(elem))
		deps = append(deps, elem.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}

func (v *BasicObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
	for _, v := range v.UnderlyingObjects {
		v.AddReferenceWithID(target, creator)
	}
}

func (v *BasicObject) NewObject() Object {
	newObject := NewBasicObject(NewObjectInfo(v.GetType()))
	// we attach the dependency because we created a fully new object
	newObject.GetVariableInfo().AddDependency(v)
	return newObject
}

func (v *BasicObject) Copy(force bool) Object {
	var underlyingVariablesCopy []Object
	for _, v := range v.UnderlyingObjects {
		underlyingVariablesCopy = append(underlyingVariablesCopy, v.Copy(force))
	}
	copy := &BasicObject{
		ObjectInfo:        v.ObjectInfo.Copy(force),
		UnderlyingObjects: underlyingVariablesCopy,
	}
	return copy
}

func (v *BasicObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS BASIC - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	var underlyingVariablesCopy []Object
	for _, v := range v.UnderlyingObjects {
		underlyingVariablesCopy = append(underlyingVariablesCopy, v.DeepCopy())
	}
	copy := &BasicObject{
		ObjectInfo:        v.ObjectInfo.DeepCopy(),
		UnderlyingObjects: underlyingVariablesCopy,
	}
	return copy
}

func (v *BasicObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
	}
	for _, v := range v.UnderlyingObjects {
		variables = append(variables, v.GetUnassaignedVariables()...)
	}
	return variables
}

func (v *BasicObject) UpgradeFromPreviousInterface(interfaceVariable *InterfaceObject) {
	/* v.GetVariableInfo().Reference = interfaceVariable.GetVariableInfo().Reference
	v.GetVariableInfo().Dataflows = interfaceVariable.GetVariableInfo().Dataflows
	v.GetVariableInfo().IndirectDataflows = interfaceVariable.GetVariableInfo().IndirectDataflows */
	v.UnderlyingObjects = append(v.UnderlyingObjects, interfaceVariable)
	interfaceVariable.GetVariableInfo().AddParent(interfaceVariable, v)
}
