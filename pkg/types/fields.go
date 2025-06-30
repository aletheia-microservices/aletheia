package types

import (
	"fmt"

	"analyzer/pkg/datastores"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
)

// Field represents all the fields of a ServiceImpl (Blueprint) or a Method (parsed from a FuncDecl)
// Fields can be (1) MethodField, (2) GenericField, (3) ServiceField, or (4) DatabaseField
//
// e.g.
//
//	type StorageServiceImpl struct {
//		analyticsService AnalyticsService
//		mediaService     MediaService
//		posts_cache      backend.Cache
//		postsDb          backend.NoSQLDatabase
//		analyticsQueue   backend.Queue
//	}
type Field interface {
	String() string
	LongString() string
	GetName() string
	GetType() gotypes.Type
	SetType(t gotypes.Type)
	GetTypeString() string
	GetTypeName() string
	GetTypeLongName() string
	GetIndex() int
}

type FieldInfo struct {
	Name   string
	Type   gotypes.Type
	Object objects.Object
	Idx    int
}

type MethodField struct {
	Field
	FieldInfo FieldInfo
}
type GenericField struct {
	Field
	FieldInfo FieldInfo
}
type ServiceField struct {
	Field
	FieldInfo FieldInfo
}
type DatabaseField struct {
	Field
	FieldInfo         FieldInfo
	Queue             bool
	DatastoreInstance datastores.DatabaseInstance // instance is the same as getting FieldInfo.Type
}

// -------------
// GENERIC FIELD
// -------------
func (f *GenericField) String() string {
	return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.String())
}
func (f *GenericField) LongString() string {
	return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.LongString())
}
func (f *GenericField) GetTypeString() string {
	return f.FieldInfo.Type.String()
}
func (f *GenericField) GetIndex() int {
	return f.FieldInfo.Idx
}
func (f *GenericField) GetName() string {
	return f.FieldInfo.Name
}
func (f *GenericField) GetType() gotypes.Type {
	return f.FieldInfo.Type
}
func (f *GenericField) GetTypeName() string {
	return f.FieldInfo.Type.GetName()
}
func (f *GenericField) GetTypeLongName() string {
	return f.FieldInfo.Type.GetName()
}
func (f *GenericField) SetType(t gotypes.Type) {
	f.FieldInfo.Type = t
}

// -------------
// SERVICE FIELD
// -------------
func (f *ServiceField) String() string {
	return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.String())
}
func (f *ServiceField) LongString() string {
	return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.LongString())
}
func (f *ServiceField) GetTypeString() string {
	return f.FieldInfo.Type.String()
}
func (f *ServiceField) GetTypeName() string {
	return f.FieldInfo.Type.GetName()
}
func (f *ServiceField) GetTypeLongName() string {
	return f.FieldInfo.Type.GetLongName()
}
func (f *ServiceField) GetIndex() int {
	return f.FieldInfo.Idx
}
func (f *ServiceField) GetName() string {
	return f.FieldInfo.Name
}
func (f *ServiceField) GetType() gotypes.Type {
	return f.FieldInfo.Type
}
func (f *ServiceField) SetType(t gotypes.Type) {
	f.FieldInfo.Type = t
}

// --------------
// DATABASE FIELD
// --------------
func (f *DatabaseField) String() string {
	return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.String())
}
func (f *DatabaseField) LongString() string {
	return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.LongString())
}
func (f *DatabaseField) GetTypeString() string {
	return f.FieldInfo.Type.String()
}
func (f *DatabaseField) GetTypeName() string {
	return f.FieldInfo.Type.GetName()
}
func (f *DatabaseField) GetTypeLongName() string {
	return f.FieldInfo.Type.GetLongName()
}
func (f *DatabaseField) GetIndex() int {
	return f.FieldInfo.Idx
}
func (f *DatabaseField) GetName() string {
	return f.FieldInfo.Name
}
func (f *DatabaseField) GetType() gotypes.Type {
	return f.FieldInfo.Type
}
func (f *DatabaseField) SetType(t gotypes.Type) {
	f.FieldInfo.Type = t
}
func (f *DatabaseField) IsQueue() bool {
	return f.Queue
}
func (f *DatabaseField) GetDatastoreInstance() datastores.DatabaseInstance {
	return f.DatastoreInstance
}

// ------------------
// FUNCTION PARAMETER
// ------------------
func (f *MethodField) String() string {
	if f.FieldInfo.Name != "" {
		return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.String())
	}
	return f.FieldInfo.Type.String()
}
func (f *MethodField) LongString() string {
	if f.FieldInfo.Name != "" {
		return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.LongString())
	}
	return f.FieldInfo.Type.String()
}
func (f *MethodField) GetTypeString() string {
	return f.FieldInfo.Type.String()
}
func (f *MethodField) GetTypeName() string {
	return f.FieldInfo.Type.GetName()
}
func (f *MethodField) GetTypeLongName() string {
	return f.FieldInfo.Type.GetLongName()
}
func (f *MethodField) GetIndex() int {
	return -1
}
func (f *MethodField) GetName() string {
	return f.FieldInfo.Name
}
func (f *MethodField) GetType() gotypes.Type {
	return f.FieldInfo.Type
}
func (f *MethodField) SetType(t gotypes.Type) {
	f.FieldInfo.Type = t
}
