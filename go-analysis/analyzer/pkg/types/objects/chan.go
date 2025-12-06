package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
)

type ChanObject struct {
	Object
	ObjectInfo *ObjectInfo
}

func NewChanObject(info *ObjectInfo) *ChanObject {
	return &ChanObject{
		ObjectInfo: info,
	}
}

func (v *ChanObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo *ObjectInfo `json:"chan"`
	}{
		ObjectInfo: v.ObjectInfo,
	})
}

func (v *ChanObject) String() string {
	return v.ObjectInfo.String()
}

func (v *ChanObject) LongString() string {
	return v.ObjectInfo.LongString()
}

func (v *ChanObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *ChanObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for chan variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *ChanObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *ChanObject) GetDependencies() []Object {
	return v.GetVariableInfo().GetDependencies()
}

func (v *ChanObject) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	for _, elem := range v.GetDependencies() { // to include underlying dependencies from variable info
		logger.Logger.Debugf("[CHAN OBJECT] GOT NESTED DEP (INC/ VINFO) FOR ELEM %s (%s)", elem.String(), VariableTypeName(elem))
		deps = append(deps, elem.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}

func (v *ChanObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
}

func (v *ChanObject) NewObject() Object {
	newObject := NewChanObject(NewObjectInfo(v.GetType()))
	// we attach the dependency because we created a fully new object
	newObject.ObjectInfo.AddDependency(v)
	return newObject
}

func (v *ChanObject) Copy(force bool) Object {
	copy := &ChanObject{
		ObjectInfo: v.ObjectInfo.Copy(force),
	}
	return copy
}

func (v *ChanObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS CHAN - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &ChanObject{
		ObjectInfo: v.ObjectInfo.DeepCopy(),
	}
	return copy
}

func (v *ChanObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
	}
	return variables
}
