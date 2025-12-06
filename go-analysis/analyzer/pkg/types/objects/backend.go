package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
)

type BackendObject struct {
	Object
	Origin     *Object
	ObjectInfo *ObjectInfo
}

func (v *BackendObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo *ObjectInfo `json:"backend"`
	}{
		ObjectInfo: v.ObjectInfo,
	})
}

func (v *BackendObject) String() string {
	return v.ObjectInfo.String()
}

func (v *BackendObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *BackendObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for backend variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *BackendObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *BackendObject) GetDependencies() []Object {
	return v.GetVariableInfo().GetDependencies()
}

func (v *BackendObject) GetNestedDependencies(includeRefBy bool) []Object {
	return nil
}

func (v *BackendObject) NewObject() Object {
	logger.Logger.Fatalf("[OBJ BACKEND] attempted to create new backend: %s", v.String())
	return nil
}

func (v *BackendObject) Copy(force bool) Object {
	return &BackendObject{
		ObjectInfo: v.ObjectInfo.Copy(force),
	}
}

func (v *BackendObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS BACKEND - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	return &BackendObject{
		ObjectInfo: v.ObjectInfo.DeepCopy(),
	}
}

func (v *BackendObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
}

func (v *BackendObject) GetUnassaignedVariables() []Object {
	return nil
}

func (t *BackendObject) GetNestedFieldVariables(prefix string, noSQL bool) ([]Object, []string) {
	return nil, nil
}

func (t *BackendObject) CopyFrom(target *BackendObject) {
	// nothing to be done here
}
