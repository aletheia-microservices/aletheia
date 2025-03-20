package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/utils"
)

type AddressObject struct {
	Object
	ObjectInfo *ObjectInfo
	AddressOf  Object
}

func (v *AddressObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo *ObjectInfo `json:"address"`
		AddressOf  Object      `json:"addr_of,omitempty"`
	}{
		ObjectInfo: v.ObjectInfo,
		AddressOf:  v.AddressOf,
	})
}

func (v *AddressObject) String() string {
	return v.ObjectInfo.String()
}

func (v *AddressObject) LongString() string {
	return v.ObjectInfo.String() + " = " + "*" + v.GetAddressOf().LongString()
}

func (v *AddressObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *AddressObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for address variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *AddressObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *AddressObject) GetAddressOf() Object {
	logger.Logger.Debugf("[VARS ADDRESS] getting addressOf variable for address variable: %s", v.String())
	if v.AddressOf == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil address to variable in address variable (%s)", v.String())
	}
	return v.AddressOf
}

func (v *AddressObject) GetDependencies() []Object {
	return append(v.GetVariableInfo().GetDependencies(), v.AddressOf)
}

func (v *AddressObject) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	for _, elem := range v.GetDependencies() { // to include underlying dependencies from variable info
		logger.Logger.Debugf("[ADDRESS OBJECT] GOT NESTED DEP (INC/ VINFO) FOR ELEM %s (%s)", elem.String(), VariableTypeName(elem))
		deps = append(deps, elem.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}

func (v *AddressObject) AddReferenceWithID(target Object, creator string) {
	if targetAddressOf, ok := target.(*AddressObject); ok {
		v.AddressOf.AddReferenceWithID(targetAddressOf.AddressOf, creator)
	} else {
		logger.Logger.Fatalf("referenced variables with different types (%s vs %s) (%s vs %s)", v.String(), target.String(), utils.GetType(v), utils.GetType(target))
	}
}

func (v *AddressObject) Copy(force bool) Object {
	copy := &AddressObject{
		ObjectInfo: v.ObjectInfo.Copy(force),
		AddressOf:  v.AddressOf, // underlying values of addresses are never deep copied
	}
	copy.AddressOf.GetVariableInfo().AddParent(copy.AddressOf, copy)
	return copy
}

func (v *AddressObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS ADDRESS - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &AddressObject{
		ObjectInfo: v.ObjectInfo.DeepCopy(),
		AddressOf:  v.AddressOf.DeepCopy(), // underlying values of addresses are never deep copied
	}
	copy.AddressOf.GetVariableInfo().AddParent(copy.AddressOf, copy)
	return copy
}

func (v *AddressObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
		variables = append(variables, v.AddressOf.GetUnassaignedVariables()...)
	}
	return variables
}
