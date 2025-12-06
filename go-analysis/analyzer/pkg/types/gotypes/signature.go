package gotypes

import "analyzer/pkg/logger"

type SignatureType struct {
	Type        `json:"-"`
	ReturnTypes []Type
}

// ------------
// Type Methods
// ------------

func (t *SignatureType) DeepCopy() Type {
	var returnTypesCopy []Type
	for _, returnType := range t.ReturnTypes {
		returnTypesCopy = append(returnTypesCopy, returnType.DeepCopy())
	}
	return &SignatureType{
		ReturnTypes: returnTypesCopy,
	}
}

func (t *SignatureType) IsSameType(other Type) bool {
	_, ok := other.(*SignatureType)
	return ok
}
func (t *SignatureType) String() string {
	if len(t.ReturnTypes) == 0 {
		return ""
	}
	if len(t.ReturnTypes) == 1 {
		return t.ReturnTypes[0].LongString()
	}
	s := "( "
	for i, ft := range t.ReturnTypes {
		s += ft.LongString()
		if i < len(t.ReturnTypes)-1 {
			s += ", "
		}
	}
	s += " )"
	return s
}
func (t *SignatureType) LongString() string {
	return t.String()
}
func (t *SignatureType) GetName() string {
	return ""
}
func (t *SignatureType) GetLongName() string {
	return ""
}
func (t *SignatureType) GetPackage() string {
	return ""
}
func (t *SignatureType) GetReturns() []Type {
	return t.ReturnTypes
}
func (t *SignatureType) GetBasicValue() string {
	logger.Logger.Fatalf("unable to get value for signature type %s", t.String())
	return ""
}
func (t *SignatureType) AddValue(value string) {
	logger.Logger.Fatalf("unable to add value for signature type %s", t.String())
}

// -----------------
// Signature Methods
// -----------------
