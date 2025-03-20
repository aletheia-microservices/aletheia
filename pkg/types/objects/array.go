package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/utils"
)

type ArrayObject struct {
	Object
	ObjectInfo      *ObjectInfo
	Elements        []Object
	DynamicElements []Object
}

func NewArrayObject(info *ObjectInfo) *ArrayObject {
	return &ArrayObject{
		ObjectInfo: info,
	}
}

func (v *ArrayObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo *ObjectInfo `json:"array"`
		Elements   []Object    `json:"array_elems,omitempty"`
	}{
		ObjectInfo: v.ObjectInfo,
		Elements:   v.Elements,
	})
}

func (v *ArrayObject) String() string {
	return v.ObjectInfo.String()
}

func (v *ArrayObject) AddDynamicElement(obj Object) {
	v.DynamicElements = append(v.DynamicElements, obj)
}

func (v *ArrayObject) UpgradeToSlice() *SliceObject {
	if v.ObjectInfo.Name == "usermentions" {
		logger.Logger.Infof("[ARRAY] upgrading array object to slice: %s", v.String())
	}
	v.ObjectInfo.Type = &gotypes.SliceType{
		UnderlyingType: v.GetArrayType().ElementsType,
	}
	sliceVariable := &SliceObject{
		ObjectInfo: v.ObjectInfo,
		Elements:   v.Elements,
	}
	for _, elem := range sliceVariable.Elements {
		elem.GetVariableInfo().AddParent(elem, sliceVariable)
	}
	return sliceVariable
}

func (v *ArrayObject) LongString() string {
	s := v.ObjectInfo.String() + " = ("
	for i, elem := range v.Elements {
		s += elem.String()
		if i < len(v.Elements)-1 {
			s += ", "
		}
	}
	return s + ")"
}

func (v *ArrayObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
	for i := 0; i < len(v.Elements); i++ {
		v.AddReferenceWithID(target.GetElementAt(i), creator)
	}
}

func (v *ArrayObject) GetElements() []Object {
	return v.Elements
}

func (v *ArrayObject) SetElementAt(idx int, elem Object) {
	if idx > len(v.Elements)-1 { // this can actually happen if length is a dynamic value
		logger.Logger.Warnf("[ARRAY OBJECT] setting new element (%s) at index (%d) out of bounds for array with elements: %v", elem.String(), idx, v.Elements)
	}
	for i := len(v.Elements); i <= idx; i++ {
		v.Elements = append(v.Elements, nil)
	}
	v.Elements[idx] = elem
}

func (v *ArrayObject) AppendElement(varElements Object) {
	if varElementsSlice, ok := varElements.(*ArrayObject); ok {
		v.Elements = append(v.Elements, varElementsSlice.GetElements()...)
	} else if v.GetElementsType().IsSameType(varElements.GetType()) {
		v.Elements = append(v.Elements, varElements)
	} else if v.GetArrayType().ElementsTypeIsInterface() {
		v.Elements = append(v.Elements, varElements)
	} else {
		//logger.Logger.Fatalf("[VARS SLICE] slice variable underlying type (%s) does not match elements type (%s)", utils.GetType(v.GetArrayType().ElementsType), utils.GetType(varElements.GetType()))
		v.Elements = append(v.Elements, varElements)
	}
}

func (v *ArrayObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *ArrayObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for array variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *ArrayObject) GetArrayType() *gotypes.ArrayType {
	if userType, ok := v.ObjectInfo.GetType().(*gotypes.UserType); ok {
		return userType.UserType.(*gotypes.ArrayType)
	}
	if sliceType, ok := v.ObjectInfo.GetType().(*gotypes.SliceType); ok {
		logger.Logger.Fatalf("GOT SLICE TYPE (%s) FOR ARRAY OBJECT: %s", sliceType.String(), v.String())
	}
	return v.ObjectInfo.GetType().(*gotypes.ArrayType)
}

func (v *ArrayObject) GetArrayOrSliceType() gotypes.Type {
	if userType, ok := v.ObjectInfo.GetType().(*gotypes.UserType); ok {
		return userType.UserType.(*gotypes.ArrayType)
	}
	if sliceType, ok := v.ObjectInfo.GetType().(*gotypes.SliceType); ok {
		logger.Logger.Warnf("GOT SLICE TYPE (%s) FOR ARRAY OBJECT: %s", sliceType.String(), v.String())
		return sliceType
	}
	return v.ObjectInfo.GetType().(*gotypes.ArrayType)
}

func (v *ArrayObject) GetElementsType() gotypes.Type {
	var t gotypes.Type
	if userType, ok := v.ObjectInfo.GetType().(*gotypes.UserType); ok {
		t = userType.UserType
		if t == nil {
			logger.Logger.Fatalf("[ARRAY OBJECT] got unexpected nil underlying user type: %s", userType.String())
		}
	} else {
		t = v.GetType()
	}
	if arrayType, ok := t.(*gotypes.ArrayType); ok {
		return arrayType.ElementsType
	} else if sliceType, ok := t.(*gotypes.SliceType); ok {
		logger.Logger.Warnf("GOT SLICE TYPE (%s) FOR ARRAY OBJECT: %s", sliceType.String(), v.String())
		return sliceType.UnderlyingType
	}
	logger.Logger.Fatalf("[ARRAY OBJECT] attempted to get elements type but got unexpected obj type [%s] for object: %s", utils.GetType(t), v.String())
	return nil
}

func (v *ArrayObject) GetSliceType() *gotypes.SliceType {
	if userType, ok := v.ObjectInfo.GetType().(*gotypes.UserType); ok {
		return userType.UserType.(*gotypes.SliceType)
	}
	return v.ObjectInfo.GetType().(*gotypes.SliceType)
}

func (v *ArrayObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *ArrayObject) AddElement(element Object) {
	v.Elements = append(v.Elements, element)
}

func (v *ArrayObject) AddElements(element []Object) {
	v.Elements = append(v.Elements, element...)
}

func (v *ArrayObject) GetElementAt(index int) Object {
	if index > len(v.Elements)-1 {
		logger.Logger.Fatalf("[VARS ARRAY] element at index (%d) does not exist in array variable with len (%d): %s", index, len(v.Elements), v.String())
	}
	return v.Elements[index]
}

func (v *ArrayObject) GetElementAtIfExists(index int) Object {
	if index > len(v.Elements)-1 {
		logger.Logger.Warnf("[VARS ARRAY] element at index (%d) does not exist in array variable with len (%d): %s", index, len(v.Elements), v.String())
		return nil
	}
	return v.Elements[index]
}

func (v *ArrayObject) GetDependencies() []Object {
	return append(v.GetVariableInfo().GetDependencies(), v.Elements...)
}

func (v *ArrayObject) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	for _, elem := range v.GetDependencies() { // to include underlying dependencies from variable info
		logger.Logger.Debugf("[ARRAY OBJECT] GOT NESTED DEP (INC/ VINFO) FOR ELEM %s (%s)", elem.String(), VariableTypeName(elem))
		deps = append(deps, elem.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}

func (v *ArrayObject) NewVersion() Object {
	copy := &ArrayObject{
		ObjectInfo: v.ObjectInfo.Copy(true),
	}
	for k, v := range v.Elements { // FIXME: this is not 100% correct
		copy.Elements[k] = v
		copy.Elements[k].GetVariableInfo().AddParent(copy.Elements[k], copy)
	}
	return copy
}

func (v *ArrayObject) NewObject() Object {
	newObject := NewArrayObject(NewObjectInfo(v.GetType()))
	// (1) new array points to the same elements as the original one
	// (2) we don't add new dependency edge from new to old object
	// because we keep the same elements
	for _, elem := range v.Elements {
		newObject.AddElement(elem)
		elem.GetVariableInfo().AddParent(elem, newObject)
	}
	for _, elem := range v.DynamicElements {
		newObject.AddDynamicElement(elem)
		elem.GetVariableInfo().AddParent(elem, newObject)
	}
	return newObject
}

func (v *ArrayObject) Copy(force bool) Object {
	copy := &ArrayObject{
		ObjectInfo: v.ObjectInfo.Copy(force),
	}
	for _, elem := range v.Elements {
		newElem := elem.Copy(force)
		copy.Elements = append(copy.Elements, newElem)
		newElem.GetVariableInfo().AddParent(newElem, copy)
	}
	return copy
}

func (v *ArrayObject) DeepCopy() Object {
	logger.Logger.Debugf("[VARS ARRAY - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &ArrayObject{
		ObjectInfo: v.ObjectInfo.DeepCopy(),
	}
	for _, elem := range v.Elements {
		newElem := elem.DeepCopy()
		copy.Elements = append(copy.Elements, newElem)
		newElem.GetVariableInfo().AddParent(newElem, copy)
	}
	return copy
}

func (v *ArrayObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
		for _, v := range v.Elements {
			variables = append(variables, v.GetUnassaignedVariables()...)
		}
	}
	return variables
}
