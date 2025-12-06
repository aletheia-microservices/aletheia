package gotypes

import (
	"analyzer/pkg/logger"
)

type PackageType struct {
	Type  `json:"-"`
	Alias string // import alias
	Name  string
	Path  string
}

// ------------
// Type Methods
// ------------

func (t *PackageType) DeepCopy() Type {
	return &PackageType{
		Alias: t.Alias,
		Name: t.Name,
		Path: t.Path,
	}
}

func (t *PackageType) IsSameType(other Type) bool {
	_, ok := other.(*PackageType)
	return ok
}
func (t *PackageType) String() string {
	return t.Name
}
func (t *PackageType) LongString() string {
	return t.Name + " " + t.Path
}
func (t *PackageType) GetName() string {
	return t.Name
}
func (t *PackageType) GetLongName() string {
	return t.LongString()
}
func (t *PackageType) GetBasicValue() string {
	logger.Logger.Fatalf("[TYPES PACKAGE] unable to get value for package type %s", t.String())
	return ""
}
func (t *PackageType) AddValue(value string) {
	logger.Logger.Fatalf("[TYPES PACKAGE] unable to add value for package type %s", t.String())
}
func (t *PackageType) GetNestedFieldTypes(prefix string, noSQL bool) ([]Type, []string) {
	return nil, nil
}
