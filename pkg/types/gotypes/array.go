package gotypes

import (
	"runtime"

	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

type ArrayType struct {
	Type         `json:"-"`
	ElementsType Type
	NumElems     int // value set to -1 means that array has no length specified
}

// ------------
// Type Methods
// ------------

func (t *ArrayType) DeepCopy() Type {
	return &ArrayType{
		ElementsType: t.ElementsType,
	}
}

func (t *ArrayType) ElementsTypeIsInterface() bool {
	_, ok := t.ElementsType.(*InterfaceType)
	return ok
}

func (t *ArrayType) IsSameType(other Type) bool {
	_, ok := other.(*ArrayType)
	return ok
}
func (t *ArrayType) String() string {
	if t.ElementsType == nil {
		logger.Logger.Fatal("[TYPES ARRAY] unexpected nil underlying type for Array Type")
	}
	return "[]" + t.ElementsType.String()
}
func (t *ArrayType) LongString() string {
	if t.ElementsType == nil {
		logger.Logger.Fatal("[TYPES ARRAY] unexpected nil underlying type for Array Type")
	}
	return "[]" + t.ElementsType.LongString()
}
func (t *ArrayType) GetName() string {
	return t.String()
}
func (t *ArrayType) GetLongName() string {
	return t.String()
}
func (t *ArrayType) GetBasicValue() string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		logger.Logger.Fatalf("unable to get value for array type %s", t.String())
	}
	callerFunc := runtime.FuncForPC(pc).Name()
	logger.Logger.Fatalf("unable to get value for array type (%s)... caller: %s %s:%d", t.String(), callerFunc, file, line)
	return ""
}
func (t *ArrayType) AddValue(value string) {
	logger.Logger.Fatalf("unable to add value for array type %s", t.String())
}
func (t *ArrayType) GetNestedFieldTypes(prefix string, noSQL bool) ([]Type, []string) {
	logger.Logger.Warnf("[TYPES ARRAY] attempted to get nested fields and variables: %s", t.String())
	var nestedTypes = []Type{t.ElementsType}
	var nestedIDs = []string{prefix + utils.DYNAMIC_FIELD}
	return nestedTypes, nestedIDs
}

// -------------
// Array Methods
// -------------
