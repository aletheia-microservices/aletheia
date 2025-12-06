package gotypes

import (
	"fmt"
	"strings"

	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

type BasicType struct {
	Type     `json:"-"`
	Name     string
	Value    string
	NotValue bool
}

// ------------
// Type Methods
// ------------

func NewBasicType(name string, value string) *BasicType {
	return &BasicType{
		Name:  name,
		Value: value,
	}
}

func (t *BasicType) DeepCopy() Type {
	return &BasicType{
		Name:  t.Name,
		Value: t.Value,
	}
}

func (t *BasicType) IsNil() bool {
	return t.Name == "nil"
}
func (t *BasicType) IsSameType(other Type) bool {
	_, ok := other.(*BasicType)
	return ok
}
func (t *BasicType) String() string {
	if t.Value != "" {
		return fmt.Sprintf("%s %s", t.Value, t.Name)
	}
	return t.Name
}
func (t *BasicType) LongString() string {
	return t.String()
}
func (t *BasicType) GetName() string {
	return t.Name
}
func (t *BasicType) GetBasicValue() string {
	// remove double quotes " " if they exist
	value := t.Value
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
	}
	return value
}

func (t *BasicType) AddValue(value string) {
	if t.Value == "" {
		t.Value = value
		return
	}
	t.Value = utils.AddValue(t.Name, t.Value, value)
}

func (t *BasicType) InvertBool() {
	if t.Value == "true" {
		t.Value = "false"
	} else if t.Value == "false" {
		t.Value = "true"
	} else {
		logger.Logger.Warnf("[BASIC TYPE] found dynamic boolean (value='%s') for basic type: %s", t.Value, t.String())
	}

	t.NotValue = !t.NotValue
}

func (t *BasicType) GetNestedFieldTypes(prefix string, noSQL bool) ([]Type, []string) {
	return nil, nil
}

// -------------
// Basic Methods
// -------------
