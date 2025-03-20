package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
)

type FuncObject struct {
	Object
	ObjectInfo *ObjectInfo
	Name       string
}

func NewFuncObject(info *ObjectInfo, name string) *FuncObject {
	return &FuncObject{
		ObjectInfo: info,
		Name:       name,
	}
}

func (v *FuncObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo *ObjectInfo `json:"func"`
		Name       string      `json:"name"`
	}{
		ObjectInfo: v.ObjectInfo,
		Name:       v.Name,
	})
}

func (v *FuncObject) String() string {
	return v.ObjectInfo.String()
}

func (v *FuncObject) LongString() string {
	return v.ObjectInfo.LongString()
}

func (v *FuncObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *FuncObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS FUNC] unexpected nil type for interface variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *FuncObject) GetFuncType() *gotypes.FuncTypeType {
	return v.ObjectInfo.GetType().(*gotypes.FuncTypeType)
}

func (v *FuncObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *FuncObject) GetDependencies() []Object {
	return v.GetVariableInfo().GetDependencies()
}

func (v *FuncObject) GetNestedDependencies(includeRefBy bool) []Object {
	return nil
}

func (v *FuncObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
}

func (v *FuncObject) Copy(force bool) Object {
	copy := &FuncObject{
		ObjectInfo: v.ObjectInfo.Copy(force),
	}
	return copy
}

func (v *FuncObject) NewObject() Object {
	newObject := NewFuncObject(NewObjectInfo(v.GetType()), v.Name)
	// add dependency edge since object is completely new
	newObject.ObjectInfo.AddDependency(v)
	return newObject
}

func (v *FuncObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS INTERFACE - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &FuncObject{
		ObjectInfo: v.ObjectInfo.DeepCopy(),
	}
	return copy
}

func (v *FuncObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
	}
	return variables
}

func (v *FuncObject) UpgradeToAssertedType(newType gotypes.Type) {
	v.GetVariableInfo().Type = newType
}
