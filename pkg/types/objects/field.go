package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/utils"
)

type FieldObject struct {
	Object
	ObjectInfo      *ObjectInfo
	WrappedVariable Object
}

func NewFieldObject(info *ObjectInfo, wrappedObj Object) *FieldObject {
	return &FieldObject{
		ObjectInfo:      info,
		WrappedVariable: wrappedObj,
	}
}

func (v *FieldObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo      *ObjectInfo `json:"field"`
		WrappedVariable Object      `json:"wrapped_variable,omitempty"`
	}{
		ObjectInfo:      v.ObjectInfo,
		WrappedVariable: v.WrappedVariable,
	})
}

func (v *FieldObject) String() string {
	return v.ObjectInfo.String()
}

func (v *FieldObject) LongString() string {
	return v.ObjectInfo.LongString()
}

func (v *FieldObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *FieldObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for field variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *FieldObject) GetFieldType() *gotypes.FieldType {
	return v.ObjectInfo.GetType().(*gotypes.FieldType)
}

func (v *FieldObject) GetWrappedVariable() Object {
	return v.WrappedVariable
}

func (v *FieldObject) FieldTypeIsDirectStructType() bool {
	// if it is directly a struct type (and not user type)
	_, ok := v.ObjectInfo.GetType().(*gotypes.FieldType).GetWrappedType().(*gotypes.StructType)
	return ok
}

func (v *FieldObject) AssignVariable(rvariable Object) {
	// e.g. post.Text = post2.Text2
	if v.GetType().IsSameType(rvariable.GetType()) {
		logger.Logger.Infof("[VAR FIELD] (a) assigning variables with types (%s --> %s) for lvariable (%s) and rvariable (%s)", utils.GetType(v.WrappedVariable), utils.GetType(rvariable), v.WrappedVariable.String(), rvariable.String())
		v.WrappedVariable = rvariable.(*FieldObject).WrappedVariable.Copy(false)
		v.WrappedVariable.GetVariableInfo().SetParent(v.WrappedVariable, v)
		return
	}
	// e.g. post.Text = text
	if v.WrappedVariable.GetType().IsSameType(rvariable.GetType()) {
		logger.Logger.Infof("[VAR FIELD] (b) assigning variables with types (%s --> %s) for lvariable (%s) and rvariable (%s)", utils.GetType(v.WrappedVariable), utils.GetType(rvariable), v.WrappedVariable.String(), rvariable.String())
		v.WrappedVariable = rvariable.Copy(false)
		v.WrappedVariable.GetVariableInfo().SetParent(v.WrappedVariable, v)
		return
	}

	// e.g. some_slice = some_array
	_, leftSliceOk := v.WrappedVariable.GetType().(*gotypes.SliceType)
	_, rightArrayOk := rvariable.GetType().(*gotypes.ArrayType)
	if leftSliceOk && rightArrayOk {
		// maintain left slice and add copy right array
		v.WrappedVariable = rvariable.Copy(false)
		v.WrappedVariable.(*ArrayObject).UpgradeToSlice()
		v.WrappedVariable.GetVariableInfo().SetParent(v.WrappedVariable, v)
		return
	}

	// e.g. some_slice = some_array
	_, leftArrayOk := v.WrappedVariable.GetType().(*gotypes.ArrayType)
	_, rightSliceOk := rvariable.GetType().(*gotypes.SliceType)
	if leftArrayOk && rightSliceOk {
		// maintain left slice and add copy right array
		v.WrappedVariable = rvariable.Copy(false)
		v.WrappedVariable.(*SliceObject).UpgradeToArray()
		v.WrappedVariable.GetVariableInfo().SetParent(v.WrappedVariable, v)
		return
	}

	logger.Logger.Fatalf("[VAR FIELD] cannot assign variable with unmatched types (%s --> %s) for lvariable (%s) and rvariable (%s)", utils.GetType(v.WrappedVariable), utils.GetType(rvariable), v.WrappedVariable.String(), rvariable.String())
}

func (v *FieldObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *FieldObject) GetDependencies() []Object {
	return append(v.GetVariableInfo().GetDependencies(), v.WrappedVariable)
}

func (v *FieldObject) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	for _, elem := range v.GetDependencies() { // to include underlying dependencies from variable info
		logger.Logger.Debugf("[FIELD OBJECT] GOT NESTED DEP (INC/ VINFO) FOR ELEM %s (%s)", elem.String(), VariableTypeName(elem))
		deps = append(deps, elem.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}

func (v *FieldObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
	v.WrappedVariable.AddReferenceWithID(target, creator)
}

/* func (v *FieldObject) NewVersion() Object {
	copy := &FieldObject{
		ObjectInfo:      v.ObjectInfo.Copy(true),
		WrappedVariable: v.WrappedVariable.Copy(true),
	}
	copy.WrappedVariable.GetVariableInfo().SetParent(v.WrappedVariable, copy)
	return copy
} */

func (v *FieldObject) NewObject() Object {
	newWrappedObject := v.GetWrappedVariable().NewObject()
	newObject := NewFieldObject(NewObjectInfo(v.GetType()), newWrappedObject)
	newWrappedObject.GetVariableInfo().SetParent(v.WrappedVariable, newObject)
	newObject.GetVariableInfo().AddDependency(v)
	return newObject
}

func (v *FieldObject) Copy(force bool) Object {
	copy := &FieldObject{
		ObjectInfo:      v.ObjectInfo.Copy(force),
		WrappedVariable: v.WrappedVariable.Copy(force),
	}
	copy.WrappedVariable.GetVariableInfo().SetParent(v.WrappedVariable, copy)
	return copy
}

func (v *FieldObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS FIELD - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &FieldObject{
		ObjectInfo:      v.ObjectInfo.DeepCopy(),
		WrappedVariable: v.WrappedVariable.DeepCopy(),
	}
	copy.WrappedVariable.GetVariableInfo().SetParent(v.WrappedVariable, copy)
	return copy
}

func (v *FieldObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
		variables = append(variables, v.WrappedVariable.GetUnassaignedVariables()...)
	}
	return variables
}

func (t *FieldObject) GetNestedFieldVariables(prefix string, noSQL bool) ([]Object, []string) {
	variableType := t.ObjectInfo.Type
	if userType, ok := t.GetType().(*gotypes.UserType); ok {
		variableType = userType.UserType
	}

	logger.Logger.Warnf("PREFIX BEFORE = %s", prefix)
	prefix = variableType.(*gotypes.FieldType).GetNestedFieldName(prefix, noSQL)
	logger.Logger.Warnf("PREFIX AFTER = %s", prefix)

	nestedVariables := []Object{t}
	nestedNames := []string{prefix}
	if nestedField, ok := t.WrappedVariable.(*FieldObject); ok {
		r1, r2 := nestedField.GetNestedFieldVariables(prefix, noSQL)
		nestedVariables = append(nestedVariables, r1...)
		nestedNames = append(nestedNames, r2...)
	} else if nestedStruct, ok := t.WrappedVariable.(*StructObject); ok {
		r1, r2 := nestedStruct.GetNestedFieldVariables(prefix, noSQL)
		nestedVariables = append(nestedVariables, r1...)
		nestedNames = append(nestedNames, r2...)
	} else {
		nestedVariables = append(nestedVariables, t.WrappedVariable)
		nestedNames = append(nestedNames, prefix)
	}
	return nestedVariables, nestedNames
}

func (t *FieldObject) GetNestedFieldVariablesWithReferences(prefix string, noSQL bool) ([]Object, []string) {
	nestedVariables, nestedIDs := t.GetNestedFieldVariables(prefix, noSQL)
	for _, reference := range t.GetVariableInfo().GetReferences() {
		logger.Logger.Debugf("HEREEEEE FOR REFERENCE %s", reference.String())
		nestedVariablesRef, nestedIDsRef := reference.Object.(*FieldObject).GetNestedFieldVariablesWithReferences(prefix, noSQL)
		nestedVariables = append(nestedVariables, nestedVariablesRef...)
		nestedIDs = append(nestedIDs, nestedIDsRef...)
	}
	return nestedVariables, nestedIDs
}
