package blueprint

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
)

type BlueprintBackendObject struct {
	objects.Object
	ObjectInfo     *objects.ObjectInfo
	TargetVariable objects.Object // variables that read cursors with cursor.One() or cursor.All() - applies only if Type is NoSQLCursorType
}

func (v *BlueprintBackendObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo     *objects.ObjectInfo `json:"blueprint_backend"`
		TargetVariable objects.Object      `json:"target_variable,omitempty"`
	}{
		ObjectInfo:     v.ObjectInfo,
		TargetVariable: v.TargetVariable,
	})
}

func (v *BlueprintBackendObject) IsNoSQLCollection() bool {
	return v.GetBlueprintBackendType().IsNoSQLComponent() && v.GetBlueprintBackendType().IsNoSQLCollection()
}

func (v *BlueprintBackendObject) DeepCopy() objects.Object {
	return &BlueprintBackendObject{
		ObjectInfo: v.ObjectInfo.DeepCopy(),
	}
}

func (v *BlueprintBackendObject) IsNoSQLCursor() bool {
	return v.GetBlueprintBackendType().IsNoSQLComponent() && v.GetBlueprintBackendType().IsNoSQLCursor()
}

func (v *BlueprintBackendObject) IsNoSQLComponent() bool {
	return v.GetBlueprintBackendType().IsNoSQLComponent()
}

func (v *BlueprintBackendObject) String() string {
	return v.ObjectInfo.String()
}

func (v *BlueprintBackendObject) LongString() string {
	return v.ObjectInfo.LongString()
}

func (v *BlueprintBackendObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *BlueprintBackendObject) GetType() gotypes.Type {
	return v.ObjectInfo.GetType()
}

func (v *BlueprintBackendObject) GetBlueprintBackendType() *BlueprintBackendType {
	return v.ObjectInfo.GetType().(*BlueprintBackendType)
}

func (v *BlueprintBackendObject) GetVariableInfo() *objects.ObjectInfo {
	return v.ObjectInfo
}

func (v *BlueprintBackendObject) GetDependencies() []objects.Object {
	deps := v.GetVariableInfo().GetDependencies()
	if v.TargetVariable != nil {
		deps = append(deps, v.TargetVariable.GetDependencies()...)
	}
	return deps
}

func (v *BlueprintBackendObject) GetNestedDependencies(includeRefBy bool) []objects.Object {
	var deps = []objects.Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	deps = append(deps, v.TargetVariable.GetNestedDependencies(includeRefBy)...)
	return deps
}

func (v *BlueprintBackendObject) AddReferenceWithID(target objects.Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
}

func (v *BlueprintBackendObject) Copy(force bool) objects.Object {
	copy := &BlueprintBackendObject{
		ObjectInfo: v.ObjectInfo,
	}
	return copy
}

func (v *BlueprintBackendObject) GetUnassaignedVariables() []objects.Object {
	var variables []objects.Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
	}
	return variables
}

func (v *BlueprintBackendObject) SetTargetVariable(target objects.Object) {
	if v.TargetVariable != nil {
		logger.Logger.Fatalf("[BLUEPRINT BACKEND] target (%s) already exists for cursor (%s)", target.String(), v.String())
	}
	v.TargetVariable = target
}

func (v *BlueprintBackendObject) GetTargetObject() objects.Object {
	return v.TargetVariable
}
