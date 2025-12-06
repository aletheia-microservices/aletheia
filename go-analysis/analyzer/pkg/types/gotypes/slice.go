package gotypes

import "analyzer/pkg/logger"

type SliceType struct {
	Type  `json:"-"`
	UnderlyingType Type
}

// ------------
// Type Methods
// ------------

func (t *SliceType) DeepCopy() Type {
	return &SliceType{
		UnderlyingType: t.UnderlyingType.DeepCopy(),
	}
}

func (t *SliceType) IsSameType(other Type) bool {
	_, ok := other.(*SliceType)
	return ok
}
func (t *SliceType) String() string {
	if t.UnderlyingType == nil {
		logger.Logger.Fatal("[TYPES SLICE] unexpected nil underlying type for Slice Type")
	}
	return "[]" + t.UnderlyingType.String()
}
func (t *SliceType) LongString() string {
	return t.String()
}
func (t *SliceType) GetName() string {
	return t.String()
}
func (t *SliceType) GetLongName() string {
	return t.String()
}
func (t *SliceType) GetUnderlyingType() Type {
	if t.UnderlyingType == nil {
		logger.Logger.Fatal("[TYPES SLICE] unexpected nil underlying type for Slice Type")
	}
	return t.UnderlyingType
}
func (t *SliceType) GetBasicValue() string {
	logger.Logger.Fatalf("unable to get value for tuple type %s", t.String())
	return ""
}
func (t *SliceType) AddValue(value string) {
	logger.Logger.Fatalf("unable to add value for tuple type %s", t.String())
}
func (t *SliceType) GetNestedFieldTypes(prefix string, noSQL bool) ([]Type, []string) {
	logger.Logger.Warnf("GETTING NESTED FIELD TYPES FOR SLICE %v", t.String())
	return nil, nil
}

// -------------
// Tuple Methods
// -------------
