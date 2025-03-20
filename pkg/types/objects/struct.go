package objects

import (
	"encoding/json"
	"fmt"
	"slices"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/utils"
)

type StructObject struct {
	Object
	Origin     *Object
	ObjectInfo *ObjectInfo
	fields     map[string]Object
	fieldslst  []Object
}

func (v *StructObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Origin     *Object           `json:"-"`
		ObjectInfo *ObjectInfo       `json:"struct"`
		Fields     map[string]Object `json:"struct_fields,omitempty"`
		FieldsLst  []Object          `json:"-"`
	}{
		Origin:     v.Origin,
		ObjectInfo: v.ObjectInfo,
		Fields:     v.fields,
		FieldsLst:  v.fieldslst,
	})
}

func NewStructObject(info *ObjectInfo) *StructObject {
	return &StructObject{ObjectInfo: info, fields: make(map[string]Object)}
}

func (v *StructObject) String() string {
	return v.ObjectInfo.String()
}

func (v *StructObject) LongString() string {
	s := v.ObjectInfo.LongString() + " \n\t\t\t\t\t\t = {"
	i := 0
	for name, f := range v.fields {
		s += fmt.Sprintf("%s: %s", name, f.LongString())
		if i < len(v.fields)-1 {
			s += ", "
		}
		i++
	}
	for i, f := range v.fieldslst {
		s += f.LongString()
		if i < len(v.fieldslst)-1 {
			s += ", "
		}
	}
	return s + "}"
}

func (v *StructObject) GetElementAt(index int) Object {
	if index > v.NumFieldsList()-1 {
		logger.Logger.Fatalf("[VARS STRUCT] element at index (%d) does not exist in array variable with len (%d): %s", index, v.NumFieldsList(), v.String())
	}
	return v.fieldslst[index]
}

func (v *StructObject) NumFieldsList() int {
	return len(v.fieldslst)
}

func (v *StructObject) GetFieldsList() []Object {
	return v.fieldslst
}

func (v *StructObject) GetFieldsMap() map[string]Object {
	return v.fields
}

func (v *StructObject) GetFieldAt(index int) *FieldObject {
	if index < v.NumFieldsList() {
		return v.fieldslst[index].(*FieldObject)
	}
	logger.Logger.Fatalf("[STRUCT] index (%d) out of bounds for list with %d elements: %v", index, v.NumFieldsList(), v.String())
	return nil
}

func (v *StructObject) GetFieldByKey(key string) *FieldObject {
	if field, ok := v.fields[key]; ok {
		return field.(*FieldObject)
	}
	logger.Logger.Fatalf("[STRUCT] field with key (%s) does not exist: %s", key, v.LongString())
	return nil
}

func (v *StructObject) GetFieldByKeyIfExists(key string) *FieldObject {
	if field, ok := v.fields[key]; ok {
		return field.(*FieldObject)
	}
	return nil
}

func (v *StructObject) SetFieldByKey(key string, field Object) {
	if _, exists := v.fields[key]; !exists {
		v.fields[key] = field
		return
	}

	// FIXME: should we allow this?
	logger.Logger.Fatalf("[STRUCT] field with key (%s) already exists: %s", key, v.LongString())
}

func (v *StructObject) AddFieldToList(field Object) {
	v.fieldslst = append(v.fieldslst, field)
}

func (v *StructObject) GetFieldVariablesLst() []*FieldObject {
	var fields []*FieldObject
	for _, f := range v.fields {
		fields = append(fields, f.(*FieldObject))
	}
	return fields
}

func (v *StructObject) GetFieldVariables() map[string]*FieldObject {
	fields := make(map[string]*FieldObject)
	for n, f := range v.fields {
		fields[n] = f.(*FieldObject)
	}
	return fields
}

func (v *StructObject) attachFieldVariable(fieldVariable *FieldObject) {
	fieldType := fieldVariable.GetFieldType()
	fieldType.SetOrigin(v.GetType())
	fieldType.EnableStructField()
	if fieldVariable.FieldTypeIsDirectStructType() {
		fieldType.EnableEmbedded()
	}
}

func (v *StructObject) AddOrGetFieldVariable(fieldVariable *FieldObject) {
	key := fieldVariable.GetFieldType().GetFieldName()

	logger.Logger.Warnf("[VARS STRUCT] adding field named (%s) to struct variable (%s): (%s) %s", key, v.String(), VariableTypeName(fieldVariable.WrappedVariable), fieldVariable.WrappedVariable.String())
	if key != "" {
		v.SetFieldByKey(key, fieldVariable)
	} else {
		logger.Logger.Warnf("[VARS STRUCT] FIXME! SOMETIMES WE CAN HAVE UNKEYED FIELDS ")
	}

	//v.SetFieldByKey(key, fieldVariable)
	v.AddFieldToList(fieldVariable)
	v.attachFieldVariable(fieldVariable)
	fieldVariable.GetVariableInfo().AddParent(fieldVariable, v)

	v.GetStructType().UpdateFieldAtIfExists(fieldVariable.GetFieldType())
}

func (v *StructObject) AddOrGetFieldVariableAndType(fieldVariable *FieldObject) {
	v.AddOrGetFieldVariable(fieldVariable)
	v.GetStructType().AddOrGetFieldType(fieldVariable.GetFieldType())
	fieldVariable.GetVariableInfo().AddParent(fieldVariable, v)
}

func (v *StructObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *StructObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for struct variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *StructObject) GetStructType() *gotypes.StructType {
	if userType, ok := v.GetType().(*gotypes.UserType); ok {
		return userType.UserType.(*gotypes.StructType)
	}
	return v.GetType().(*gotypes.StructType)
}

func (v *StructObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *StructObject) GetOrderedFields() []Object {
	var orderedKeys []string
	for k := range v.GetFieldsMap() {
		orderedKeys = append(orderedKeys, k)
	}
	slices.Sort(orderedKeys)

	var orderedFields []Object
	for _, k := range orderedKeys {
		orderedFields = append(orderedFields, v.GetFieldByKey(k))
	}
	return orderedFields
}

func (v *StructObject) GetDependencies() []Object {
	return append(v.GetVariableInfo().GetDependencies(), v.GetOrderedFields()...)
}

func (v *StructObject) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	for _, elem := range v.GetDependencies() { // to include underlying dependencies from variable info
		logger.Logger.Debugf("[STRUCT OBJECT] GOT NESTED DEP (INC/ VINFO) FOR ELEM %s (%s)", elem.String(), VariableTypeName(elem))
		deps = append(deps, elem.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}

func (v *StructObject) NewVersion() Object {
	copy := &StructObject{
		ObjectInfo: v.ObjectInfo.Copy(true),
		fields:     make(map[string]Object),
	}
	for n, p := range v.GetFieldsMap() { // FIXME: this is not 100% correct
		copy.SetFieldByKey(n, p)
		copy.GetFieldByKey(n).GetVariableInfo().AddParent(copy.GetFieldByKey(n), copy)
	}
	for _, f := range v.GetFieldsList() {
		copy.AddFieldToList(f)
	}
	return copy
}

func (v *StructObject) Copy(force bool) Object {
	logger.Logger.Debugf("[VARS STRUCT - COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &StructObject{
		ObjectInfo: v.ObjectInfo.Copy(force),
		fields:     make(map[string]Object),
	}
	for n, p := range v.GetFieldsMap() {
		copy.SetFieldByKey(n, p.Copy(force))
		copy.GetFieldByKey(n).GetVariableInfo().AddParent(copy.GetFieldByKey(n), copy)
	}
	for _, f := range v.GetFieldsList() {
		copy.AddFieldToList(f)
	}
	return copy
}

func (v *StructObject) NewObject() Object {
	logger.Logger.Debugf("[OBJ STRUCT - NEW COPY] (%s) %s", VariableTypeName(v), v.String())
	newObject := NewStructObject(NewObjectInfo(v.GetType()))
	for k, p := range v.GetFieldsMap() {
		newField := p.NewObject()
		newObject.SetFieldByKey(k, newField)
		newFieldInfo := newObject.GetFieldByKey(k).GetVariableInfo()
		newFieldInfo.AddParent(newField, newObject)
	}
	for _, f := range v.GetFieldsList() {
		newObject.AddFieldToList(f)
	}
	newObject.GetVariableInfo().AddDependency(v)
	return newObject
}

func (v *StructObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS STRUCT - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &StructObject{
		ObjectInfo: v.ObjectInfo.DeepCopy(),
		fields:     make(map[string]Object),
	}
	for n, p := range v.GetFieldsMap() {
		copy.SetFieldByKey(n, p.DeepCopy())
		copy.GetFieldByKey(n).GetVariableInfo().AddParent(copy.GetFieldByKey(n), copy)
	}
	for _, f := range v.GetFieldsList() {
		copy.AddFieldToList(f)
	}
	return copy
}

func (v *StructObject) GetFieldVariableIfExists(name string) Object {
	if f, ok := v.fields[name]; ok {
		return f
	}
	//logger.Logger.Warnf("[VARS STRUCT] unknown field (%s) for struct variable (%s)", name, v.String())
	return nil
}

func (v *StructObject) AddOrGetFieldKeyVariable(name string, field Object) {
	v.fields[name] = field
	v.AddFieldToList(field)
	field.GetVariableInfo().AddParent(field, v)
}

func (v *StructObject) AddOrGetFieldKeyVariableIfNotExists(name string, field Object) bool {
	if _, exists := v.fields[name]; !exists {
		v.fields[name] = field
		v.AddFieldToList(field)
		field.GetVariableInfo().AddParent(field, v)
		return true
	} else {
		logger.Logger.Warnf("[VARS STRUCT] field (%s) already exists in structure (%s)", name, v.String())
	}
	return false
}

func (v *StructObject) attachReferenceToFields(target Object, creator string) {
	if targetField, ok := target.(*FieldObject); ok {
		logger.Logger.Warnf("------- FOUND FIELD THAT IS ALSO A STRUCT!!!!")
		target = targetField.WrappedVariable
	}

	if referenceStruct, ok := target.(*StructObject); ok {
		logger.Logger.Infof("[VARS STRUCT - REF] current variable has %d fields & reference variable has (%d)", len(v.fields), len(referenceStruct.fields))
		for name, field := range referenceStruct.fields {
			if _, ok := v.fields[name]; ok {
				v.fields[name].AddReferenceWithID(field, creator)
			} else {
				logger.Logger.Warnf("[VARS STRUCT - REF] skipping reference for unknown field (%s) in struct variable %s", name, v.String())
			}
		}
	} else {
		logger.Logger.Fatalf("[VARS STRUCT - REF] referenced variables with different types (%s vs %s) (%s vs %s)", v.String(), target.String(), utils.GetType(v), utils.GetType(target))
	}
}

func (v *StructObject) AddReferenceWithID(target Object, creator string) {
	logger.Logger.Infof("[VARS STRUCT - REF] trying to reference struct variable type (%s) for target type (%s) for variables: \n\t\t\t\t\t - current [%s]: %s (INFO = %v) \n\t\t\t\t\t - target [%s]: %s (INFO = %v)", utils.GetType(v), utils.GetType(target), utils.GetType(v), v.String(), v.GetVariableInfo(), utils.GetType(target), target.String(), target.GetVariableInfo())
	v.ObjectInfo.AddReferenceWithID(v, target, creator)

	// target can also have its own references e.g.
	// (i) 1st reference: StorageService pushes to queue and Analytics service pops from queue using the workerThread() method
	// (ii) 2nd reference: AnalyticsService queue handler calls another function to process the message with the handleMessage() method
	// the current variable (v) is e.g. the "message" in the handleMessage() method
	// which references the message pop by AnalyticsService and passed in (ii)
	// which in turn references the message queued by StorageService in (i)
	v.attachReferenceToFields(target, creator)
	for _, ref := range target.GetVariableInfo().GetReferences() {
		v.AddReferenceWithID(ref.Object, ref.Creator)
	}

	logger.Logger.Debugf("[VARS STRUCT - REF] added reference (%s) -> (%s) with id = %d (creator: %s)", v.ObjectInfo.Name, target.GetVariableInfo().GetName(), v.ObjectInfo.Id, creator)
}

func (v *StructObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
		for _, p := range v.fields {
			variables = append(variables, p.GetUnassaignedVariables()...)
		}
	}
	return variables
}

func (v *StructObject) GetNestedFieldVariables(prefix string, noSQL bool) ([]Object, []string) {
	var nestedVariables []Object
	var nestedIDs []string

	logger.Logger.Debugf("[VARS STRUCT] found (%d) field VARIABLES for (%s); current prefix = %s", len(v.fields), v.String(), prefix)
	for _, f := range v.fields {
		if fieldVariable, ok := f.(*FieldObject); ok {
			nestedFieldObjects, nestedFieldIDs := fieldVariable.GetNestedFieldVariables(prefix, noSQL)
			nestedVariables = append(nestedVariables, nestedFieldObjects...)
			nestedIDs = append(nestedIDs, nestedFieldIDs...)
		} else {
			logger.Logger.Fatalf("[VARS STRUCT] ignoring field variable (%s) for (%s)", f.String(), VariableTypeName(f))
		}
	}
	return nestedVariables, nestedIDs
}

func (v *StructObject) GetNestedFieldVariablesWithReferences(prefix string, noSQL bool) ([]Object, []string) {
	nestedVariables, nestedIDs := v.GetNestedFieldVariables(prefix, noSQL)
	for _, reference := range v.GetVariableInfo().GetReferences() {
		nestedVariablesRef, nestedIDsRef := reference.Object.(*StructObject).GetNestedFieldVariablesWithReferences(prefix, noSQL)
		nestedVariables = append(nestedVariables, nestedVariablesRef...)
		nestedIDs = append(nestedIDs, nestedIDsRef...)
	}
	return nestedVariables, nestedIDs
}

func (v *StructObject) CopyFrom(target *StructObject) {
	logger.Logger.Debugf("[VARS STRUCT COPY] copying from %s", target.LongString())
	for name, targetField := range target.fields {
		if _, ok := v.fields[name]; !ok {
			v.fields[name] = targetField
		}
	}
}
