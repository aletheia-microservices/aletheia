package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
)

type ServiceObject struct {
	Object
	ObjectInfo         *ObjectInfo
	UnderlyingVariable Object
}

func (v *ServiceObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo         *ObjectInfo `json:"service"`
		UnderlyingVariable Object      `json:"underlying_variable"`
	}{
		ObjectInfo:         v.ObjectInfo,
		UnderlyingVariable: v.UnderlyingVariable,
	})
}

func (v *ServiceObject) String() string {
	return v.ObjectInfo.String()
}

func (v *ServiceObject) LongString() string {
	return v.ObjectInfo.String()
}

func (v *ServiceObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *ServiceObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for service variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *ServiceObject) GetServiceName() string {
	return v.ObjectInfo.GetType().GetName()
}

func (v *ServiceObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *ServiceObject) GetDependencies() []Object {
	return v.GetVariableInfo().GetDependencies()
}

func (v *ServiceObject) GetNestedDependencies(includeRefBy bool) []Object {
	return nil
}

func (v *ServiceObject) NewObject() Object {
	logger.Logger.Fatalf("[OBJ SERVICE] attempted to create new object: %s", v.String())
	return nil
}

func (v *ServiceObject) Copy(force bool) Object {
	return &ServiceObject{
		ObjectInfo: v.ObjectInfo.Copy(force),
	}
}

func (v *ServiceObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS SERVICE - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	return &ServiceObject{
		ObjectInfo: v.ObjectInfo.DeepCopy(),
	}
}

func (v *ServiceObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
}

func (v *ServiceObject) GetUnassaignedVariables() []Object {
	return nil
}

func (t *ServiceObject) GetNestedFieldVariables(prefix string, noSQL bool) ([]Object, []string) {
	return nil, nil
}

func (t *ServiceObject) CopyFrom(target *ServiceObject) {
	// nothing to be done here
}
