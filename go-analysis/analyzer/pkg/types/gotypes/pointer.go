package gotypes

import (
	"fmt"

	"analyzer/pkg/logger"
)

type PointerType struct {
	Type      `json:"-"`
	PointerTo Type
}

// ------------
// Type Methods
// ------------

func (t *PointerType) DeepCopy() Type {
	return &PointerType{
		PointerTo: t.PointerTo.DeepCopy(),
	}
}

func (t *PointerType) IsSameType(other Type) bool {
	_, ok := other.(*PointerType)
	return ok
}
func (t *PointerType) String() string {
	return fmt.Sprintf("(*%s)", t.PointerTo.String())
}
func (t *PointerType) LongString() string {
	return fmt.Sprintf("(*%s)", t.PointerTo.LongString())
}
func (t *PointerType) GetName() string {
	return "(*" + t.PointerTo.GetName() + ")"
}
func (t *PointerType) GetLongName() string {
	return "(*" + t.PointerTo.GetLongName() + ")"
}
func (t *PointerType) GetBasicValue() string {
	return t.PointerTo.GetBasicValue()
}
func (t *PointerType) AddValue(value string) {
	t.PointerTo.AddValue(value)
}
func (t *PointerType) GetPointerTo() Type {
	logger.Logger.Debugf("[TYPES POINTER] getting pointerTo type for pointer: %s", t.String())
	if t.PointerTo == nil {
		logger.Logger.Fatalf("[TYPES POINTER] unexpected nil pointer to type in pointer type (%s)", t.String())
	}
	return t.PointerTo
}
func (t *PointerType) GetNestedFieldTypes(prefix string, noSQL bool) ([]Type, []string) {
	return t.PointerTo.GetNestedFieldTypes(prefix, noSQL)
}

// ---------------
// Pointer Methods
// ---------------
