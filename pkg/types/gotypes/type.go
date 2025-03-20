package gotypes

import "strings"

type Type interface {
	String() string
	LongString() string
	GetName() string
	GetLongName() string
	GetPackage() string // package is always the real package path
	GetBasicValue() string
	AddValue(value string)
	GetNestedFieldTypes(prefix string, noSQL bool) ([]Type, []string)
	GetLongFieldName() string
	GetParentUserType() *UserType
	AddParentUserType(userType *UserType)
	IsSameType(t Type) bool
	DeepCopy() Type
}

func packageAlias(pkg string) string {
	splits := strings.Split(pkg, "/")
	return splits[len(splits)-1]
}

func UnwrapUserAndFieldTypes(t Type) (Type, Type) {
	unwrappedType, _ := UnwrapUserType(t)
	return UnwrapFieldType(unwrappedType)

}

// UnwrapUserTypeUnderlying returns the most inner unwrapped type
func UnwrapUserTypeUnderlying(t Type) Type {
	if userType, ok := t.(*UserType); ok {
		return UnwrapUserTypeUnderlying(userType.UserType)
	}
	return t
}

// UnwrapUserType returns the most inner unwrapped type along with the most inner user type
// if no user type is found, then it returns the current type twice
func UnwrapUserType(t Type) (Type, Type) {
	if userType, ok := t.(*UserType); ok {
		unwrappedType, t := UnwrapUserType(userType.UserType)
		// always return the inned user type
		// if we found a wrapped value, we return the current user type (which is the most inner one)
		if _, ok := t.(*UserType); !ok {
			t = userType
		}
		return unwrappedType, t
	}
	return t, t
}

// UnwrapFieldType returns the most inner unwrapped type along with the most inner field type
// if no field type is found, then it returns the current type twice
func UnwrapFieldType(t Type) (Type, Type) {
	if fieldType, ok := t.(*FieldType); ok {
		unwrappedType, t := UnwrapFieldType(fieldType.WrappedType)
		// always return the inned user type
		// if we found a wrapped value, we return the current user type (which is the most inner one)
		if _, ok := t.(*FieldType); !ok {
			t = unwrappedType
		}
		return unwrappedType, t
	}
	return t, t
}
