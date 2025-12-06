package gotypes

import "analyzer/pkg/logger"

type GenericType struct {
	Type `json:"-"`
	Name string
}

// ------------
// Type Methods
// ------------

func (t *GenericType) DeepCopy() Type {
	return &GenericType{
		Name: t.Name,
	}
}

func (t *GenericType) IsSameType(other Type) bool {
	_, ok := other.(*GenericType)
	return ok
}
func (t *GenericType) String() string {
	/* if t.Name != "" {
		return t.Name
	} */
	return "undefined"
}
func (t *GenericType) LongString() string {
	return t.String()
}
func (t *GenericType) GetName() string {
	return t.String()
}
func (t *GenericType) GetLongName() string {
	return t.String()
}
func (t *GenericType) GetBasicValue() string {
	logger.Logger.Fatalf("unable to get value for generic type %s", t.String())
	return ""
}
func (t *GenericType) AddValue(value string) {
	logger.Logger.Fatalf("unable to add value for generic type %s", t.String())
}

// ---------------
// Generic Methods
// ---------------
