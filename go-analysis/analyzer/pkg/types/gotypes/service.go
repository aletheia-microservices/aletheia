package gotypes

import (
	"fmt"

	"analyzer/pkg/logger"
)

type ServiceType struct {
	Type          `json:"-"`
	ImplName      string   // Service Info
	Name          string   // UserType
	PackagePath   string   // UserType
	Methods       []string // StructType or InterfaceType
	StructType    *StructType
	InterfaceType *InterfaceType
}

// ------------
// Type Methods
// ------------

func (t *ServiceType) DeepCopy() Type {
	var methodsCopy []string
	methodsCopy = append(methodsCopy, t.Methods...)
	var structTypeCopy *StructType
	if t.StructType != nil {
		structTypeCopy = t.StructType.DeepCopy().(*StructType)
	}
	var interfaceTypeCopy *InterfaceType
	if t.InterfaceType != nil {
		interfaceTypeCopy = t.InterfaceType.DeepCopy().(*InterfaceType)
	}
	return &ServiceType{
		ImplName:      t.ImplName,
		Name:          t.Name,
		PackagePath:   t.PackagePath,
		Methods:       methodsCopy,
		StructType:    structTypeCopy,
		InterfaceType: interfaceTypeCopy,
	}
}

func (t *ServiceType) IsSameType(other Type) bool {
	_, ok := other.(*ServiceType)
	return ok
}
func (t *ServiceType) HasPackagePath(path string) bool {
	return t.PackagePath == path
}
func (t *ServiceType) IsCurrentServiceType(other Type) bool {
	ptrType, ok := other.(*PointerType)
	if ok {
		other = ptrType.PointerTo
	}
	otherUserType, ok := other.(*UserType)
	if !ok {
		return false
	}
	_, ok = otherUserType.UserType.(*StructType)
	if !ok {
		return false
	}
	return t.PackagePath == otherUserType.PackagePath && (otherUserType.Name == t.Name || otherUserType.Name == t.ImplName)
}

func (t *ServiceType) String() string {
	return fmt.Sprintf("%s.%s", packageAlias(t.PackagePath), t.Name)
}
func (t *ServiceType) LongString() string {
	return t.String()
}
func (t *ServiceType) GetName() string {
	return t.Name
}
func (t *ServiceType) GetLongName() string {
	return t.Name
}
func (t *ServiceType) GetPackage() string {
	return t.PackagePath
}
func (t *ServiceType) GetBasicValue() string {
	logger.Logger.Fatalf("unable to get value for service type %s", t.String())
	return ""
}
func (t *ServiceType) AddValue(value string) {
	logger.Logger.Fatalf("unable to add value for service type %s", t.String())
}
