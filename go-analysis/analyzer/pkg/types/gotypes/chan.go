package gotypes

import (
	"fmt"

	"analyzer/pkg/logger"
)

type ChanType struct {
	Type     `json:"-"`
	ChanType Type
}

// ------------
// Type Methods
// ------------

func (t *ChanType) DeepCopy() Type {
	return &ChanType{
		ChanType: t.ChanType,
	}
}

func (t *ChanType) IsSameType(other Type) bool {
	_, ok := other.(*ChanType)
	return ok
}
func (t *ChanType) String() string {
	return fmt.Sprintf("chan %s", t.ChanType.String())
}
func (t *ChanType) LongString() string {
	return fmt.Sprintf("chan %s", t.ChanType.LongString())
}
func (t *ChanType) GetName() string {
	return t.String()
}
func (t *ChanType) GetLongName() string {
	return t.String()
}
func (t *ChanType) GetBasicValue() string {
	logger.Logger.Fatalf("unable to get value for chan type %s", t.String())
	return ""
}
func (t *ChanType) AddValue(value string) {
	logger.Logger.Fatalf("unable to add value for chan type %s", t.String())
}

// ------------
// Chan Methods
// ------------
