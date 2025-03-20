package gotypes

import (
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

type StructType struct {
	Type           `json:"-"`
	ParentUserType *UserType
	FieldTypes     []*FieldType
	//FieldTypesMap  map[string]*FieldType
	Methods map[string]string // maps method name to package path
}

func NewStructType() *StructType {
	return &StructType{
		//FieldTypesMap: make(map[string]*FieldType),
		Methods: make(map[string]string),
	}
}

// ------------
// Type Methods
// ------------

func (t *StructType) DeepCopy() Type {
	var fieldTypesCopy []*FieldType
	//fieldTypesMapCopy := make(map[string]*FieldType)
	for _, fieldType := range t.FieldTypes {
		fieldTypeCopy := fieldType.DeepCopy().(*FieldType)
		fieldTypeCopy.Origin = t
		fieldTypesCopy = append(fieldTypesCopy, fieldTypeCopy)
		//fieldTypesMapCopy[fieldType.GetName()] = fieldTypeCopy
	}

	var methodsCopy map[string]string = make(map[string]string)
	for k, v := range t.Methods {
		methodsCopy[k] = v
	}

	var parentUserType *UserType
	if t.ParentUserType != nil {
		parentUserType = t.ParentUserType.DeepCopy().(*UserType)
	}

	return &StructType{
		ParentUserType: parentUserType,
		FieldTypes:     fieldTypesCopy,
		Methods:        t.Methods,
	}
}

func (t *StructType) AddNewFieldType(newType *FieldType) {
	t.FieldTypes = append(t.FieldTypes, newType)
}

func (t *StructType) GetFieldTypes() []*FieldType {
	return t.FieldTypes
}

func (t *StructType) IsSameType(other Type) bool {
	_, ok := other.(*StructType)
	return ok
}
func (t *StructType) String() string {
	if len(t.FieldTypes) == 0 {
		return "struct{}"
	}
	s := "struct{"
	i := 0
	for _, f := range t.FieldTypes {
		s += f.FieldName + " " + f.String()
		if i < len(t.FieldTypes)-1 {
			s += ", "
		}
		i++
	}
	s += "}"
	return s
}
func (t *StructType) MethodsString() string {
	s := "{"
	i := 0
	for m := range t.Methods {
		s += m
		if i < len(t.Methods)-1 {
			s += ", "
		}
		i++
	}
	return s + "}"
}
func (t *StructType) LongString() string {
	if len(t.FieldTypes) == 0 {
		return "struct{}"
	}
	s := "struct{"
	i := 0
	for _, f := range t.FieldTypes {
		s += f.FieldName + " " + f.LongString()
		if i < len(t.FieldTypes)-1 {
			s += ", "
		}
		i++
	}
	s += "}"
	return s
}
func (t *StructType) GetName() string {
	return t.String()
}
func (t *StructType) GetLongName() string {
	return t.String()
}
func (t *StructType) GetBasicValue() string {
	logger.Logger.Fatalf("[TYPES STRUCT] unable to get value for struct type %s", t.String())
	return ""
}
func (t *StructType) AddValue(value string) {
	logger.Logger.Fatalf("[TYPES STRUCT] unable to add value for struct type %s", t.String())
}
func (t *StructType) GetNestedFieldTypes(prefix string, noSQL bool) ([]Type, []string) {
	var nestedTypes []Type
	var nestedIDs []string

	for _, f := range t.FieldTypes {
		nestedFieldTypes, nestedFieldIDs := f.GetNestedFieldTypes(prefix, noSQL)
		logger.Logger.Debugf("getting nested field types for field: %s", f.String())
		nestedTypes = append(nestedTypes, nestedFieldTypes...)
		nestedIDs = append(nestedIDs, nestedFieldIDs...)
	}
	return nestedTypes, nestedIDs
}

func (t *StructType) GetParentUserType() *UserType {
	return t.ParentUserType
}

func (t *StructType) AddParentUserType(userType *UserType) {
	t.ParentUserType = userType
}

// --------------
// Struct Methods
// --------------

func (v *StructType) GetMethodPackagePath(method string) string {
	if pkgPath, ok := v.Methods[method]; ok {
		return pkgPath
	}
	logger.Logger.Fatalf("unknown method (%s) for struct type (%s): methods = %s", method, v.String(), v.MethodsString())
	return ""
}

func (t *StructType) GetFieldTypeByName(name string) *FieldType {
	for _, f := range t.FieldTypes {
		if f.FieldName == name {
			return f
		}
	}
	logger.Logger.Fatalf("[TYPES STRUCT] unknown field (%s) for struct type (%s)", name, t.String())
	return nil
}

func (t *StructType) UpdateFieldSubTypeByName(name string, newType Type) {
	for i, f := range t.FieldTypes {
		if f.FieldName == name {
			t.FieldTypes[i].WrappedType = newType
			return
		}
	}
	logger.Logger.Fatalf("[TYPES STRUCT] unknown field (%s) for struct type (%s)", name, t.String())
}

func (t *StructType) GetFieldTypeByNameIfExists(name string) *FieldType {
	for _, f := range t.FieldTypes {
		if f.FieldName == name {
			return f
		}
	}
	//logger.Logger.Warnf("[TYPES STRUCT] unknown field type (%s) for struct type (%s)", name, t.String())
	return nil
}

func (t *StructType) GetFieldTypeAt(index int) *FieldType {
	if index > len(t.FieldTypes)-1 {
		logger.Logger.Fatalf("[TYPES STRUCT] unknown field at index (%d) for struct type (%s)", index, t.String())
	}
	return t.FieldTypes[index]
}

func (t *StructType) UpdateFieldAtIfExists(fieldType *FieldType) {
	logger.Logger.Warnf("[TYPES STRUCT] updating field type at for struct type (%s): NEW [%s] = %s", t.String(), fieldType.LongString(), utils.GetType(fieldType))

	index := -1
	// search for the field (created when struct type was created for the first time e.g. for a user type) whose
	// field name matches the current field type for the field object that is being used
	for i, ft := range t.FieldTypes {
		if ft.GetName() == fieldType.GetName() {
			index = i
			break
		}
	}
	if index != -1 {
		// update field type with current one
		t.FieldTypes[index] = fieldType
	} else {
		// may happen when we create an object with a nil type (e.g. in ast.CompositeLit == nil for bson.D{{myobjid, 0}})
		logger.Logger.Warnf("[STRUCT TYPE] could not find index for matching field type: %s", fieldType.String())
		t.FieldTypes = append(t.FieldTypes, fieldType)
	}
}

func (t *StructType) AddOrGetFieldType(field *FieldType) {
	t.FieldTypes = append(t.FieldTypes, field)
}

func (t *StructType) AddFieldType(field *FieldType) {
	t.FieldTypes = append(t.FieldTypes, field)
}

func (t *StructType) AddMethod(name string, pkgPath string) {
	t.Methods[name] = pkgPath
}
