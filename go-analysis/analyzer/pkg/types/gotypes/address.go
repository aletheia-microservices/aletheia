package gotypes

import (
	"fmt"

	"analyzer/pkg/logger"
)

type AddressType struct {
	Type      `json:"-"`
	AddressOf Type
}

// ------------
// Type Methods
// ------------

func (t *AddressType) DeepCopy() Type {
	return &AddressType{
		AddressOf: t.AddressOf.DeepCopy(),
	}
}

func (t *AddressType) IsSameType(other Type) bool {
	_, ok := other.(*AddressType)
	return ok
}
func (t *AddressType) String() string {
	return fmt.Sprintf("(&%s)", t.AddressOf.String())
}
func (t *AddressType) LongString() string {
	return fmt.Sprintf("(&%s)", t.AddressOf.LongString())
}
func (t *AddressType) GetName() string {
	return "(&" + t.AddressOf.GetName() + ")"
}
func (t *AddressType) GetLongName() string {
	return "(&" + t.AddressOf.GetLongName() + ")"
}
func (t *AddressType) GetBasicValue() string {
	return t.AddressOf.GetBasicValue()
}
func (t *AddressType) AddValue(value string) {
	t.AddressOf.AddValue(value)
}
func (t *AddressType) GetAddressOf() Type {
	return t.AddressOf
}
func (t *AddressType) GetPointerTo() Type {
	logger.Logger.Debugf("[TYPES ADDRESS] getting addressOf type for address: %s", t.String())
	if t.AddressOf == nil {
		logger.Logger.Fatalf("[TYPES ADDRESS] unexpected nil address to type in address type (%s)", t.String())
	}
	return t.AddressOf
}
func (t *AddressType) GetNestedFieldTypes(prefix string, noSQL bool) ([]Type, []string) {
	return t.AddressOf.GetNestedFieldTypes(prefix, noSQL)
}

// ---------------
// Address Methods
// ---------------
