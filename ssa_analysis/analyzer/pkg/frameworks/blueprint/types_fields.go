package blueprint

import (
	"fmt"
	"go/types"

	"analyzer/pkg/frameworks/components"
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
	GetTypeString() string
	GetTypeName() string
	GetTypeLongName() string
	GetIndex() int
}

type FieldInfo struct {
	Name string
	Type types.Type
	Idx  int
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
	FieldInfo FieldInfo
	Queue     bool
	Datastore *components.DatastoreInfo // instance is the same as getting FieldInfo.Type
}

// -------------
// GENERIC FIELD
// -------------
func (f *GenericField) String() string {
	return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.String())
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
func (f *GenericField) GetType() types.Type {
	return f.FieldInfo.Type
}
func (f *GenericField) GetTypeName() string {
	return f.FieldInfo.Type.String()
}
func (f *GenericField) SetType(t types.Type) {
	f.FieldInfo.Type = t
}

// -------------
// SERVICE FIELD
// -------------
func (f *ServiceField) String() string {
	return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.String())
}
func (f *ServiceField) GetTypeString() string {
	return f.FieldInfo.Type.String()
}
func (f *ServiceField) GetTypeName() string {
	return f.FieldInfo.Type.String()
}
func (f *ServiceField) GetIndex() int {
	return f.FieldInfo.Idx
}
func (f *ServiceField) GetName() string {
	return f.FieldInfo.Name
}
func (f *ServiceField) GetType() types.Type {
	return f.FieldInfo.Type
}
func (f *ServiceField) SetType(t types.Type) {
	f.FieldInfo.Type = t
}

// --------------
// DATABASE FIELD
// --------------
func (f *DatabaseField) String() string {
	return fmt.Sprintf("%s %s", f.FieldInfo.Name, f.FieldInfo.Type.String())
}
func (f *DatabaseField) GetTypeString() string {
	return f.FieldInfo.Type.String()
}
func (f *DatabaseField) GetTypeName() string {
	return f.FieldInfo.Type.String()
}
func (f *DatabaseField) GetIndex() int {
	return f.FieldInfo.Idx
}
func (f *DatabaseField) GetName() string {
	return f.FieldInfo.Name
}
func (f *DatabaseField) GetType() types.Type {
	return f.FieldInfo.Type
}
func (f *DatabaseField) SetType(t types.Type) {
	f.FieldInfo.Type = t
}
func (f *DatabaseField) IsQueue() bool {
	return f.Queue
}
func (f *DatabaseField) GetDatastore() *components.DatastoreInfo {
	return f.Datastore
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
func (f *MethodField) GetTypeString() string {
	return f.FieldInfo.Type.String()
}
func (f *MethodField) GetTypeName() string {
	return f.FieldInfo.Type.String()
}
func (f *MethodField) GetIndex() int {
	return -1
}
func (f *MethodField) GetName() string {
	return f.FieldInfo.Name
}
func (f *MethodField) GetType() types.Type {
	return f.FieldInfo.Type
}
func (f *MethodField) SetType(t types.Type) {
	f.FieldInfo.Type = t
}
