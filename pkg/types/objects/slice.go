package objects

import (
	"encoding/json"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/utils"
)

type SliceObject struct {
	Object
	ObjectInfo      *ObjectInfo
	Elements        []Object
	DynamicElements []Object
}

func NewSliceObject(info *ObjectInfo) *SliceObject {
	return &SliceObject{
		ObjectInfo: info,
	}
}

func (v *SliceObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ObjectInfo *ObjectInfo `json:"slice"`
		Elements   []Object    `json:"slice_variables,omitempty"`
	}{
		ObjectInfo: v.ObjectInfo,
		Elements:   v.Elements,
	})
}

func (v *SliceObject) UpgradeToArray() *ArrayObject {
	if v.ObjectInfo.Name == "usermentions" {
		logger.Logger.Infof("[ARRAY] upgrading array object to slice: %s", v.String())
	}
	v.ObjectInfo.Type = &gotypes.SliceType{
		UnderlyingType: v.GetSliceType().UnderlyingType,
	}
	sliceVariable := &ArrayObject{
		ObjectInfo: v.ObjectInfo,
		Elements:   v.Elements,
	}
	for _, elem := range sliceVariable.Elements {
		elem.GetVariableInfo().AddParent(elem, sliceVariable)
	}
	return sliceVariable
}

func (v *SliceObject) SetElementAt(idx int, elem Object) {
	if idx > len(v.Elements)-1 {
		logger.Logger.Fatalf("[ARRAY OBJECT] attempted to set new element (%s) at index (%d) out of bounds for array with elements: %v", elem.String(), idx, v.Elements)
	}
	v.Elements[idx] = elem
}

func (v *SliceObject) AddDynamicElement(obj Object) {
	v.DynamicElements = append(v.DynamicElements, obj)
}

func (v *SliceObject) GetVariableInfo() *ObjectInfo {
	return v.ObjectInfo
}

func (v *SliceObject) GetElements() []Object {
	return v.Elements
}

func (v *SliceObject) AddReferenceWithID(target Object, creator string) {
	v.ObjectInfo.AddReferenceWithID(v, target, creator)
	for i := 0; i < len(v.Elements); i++ {
		v.AddReferenceWithID(target.GetElementAt(i), creator)
	}
}

func (v *SliceObject) AppendElements(varElements Object) {
	if varElementsSlice, ok := varElements.(*SliceObject); ok {
		v.Elements = append(v.Elements, varElementsSlice.GetElements()...)
	} else {
		if v.GetSliceType().UnderlyingType.IsSameType(varElements.GetType()) {
			v.Elements = append(v.Elements, varElements)
		} else {
			logger.Logger.Fatalf("[VARS SLICE] slice variable underlying type (%s) does not match elements type (%s)", utils.GetType(v.GetSliceType().UnderlyingType), utils.GetType(varElements.GetType()))
		}
	}
}

func (v *SliceObject) AddElement(element Object) {
	v.Elements = append(v.Elements, element)
	element.GetVariableInfo().AddParent(element, v)
}

func (v *SliceObject) AddElements(elements []Object) {
	v.Elements = append(v.Elements, elements...)
	for _, elem := range elements {
		elem.GetVariableInfo().AddParent(elem, v)
	}
}

func (v *SliceObject) GetVariableAt(index int) Object {
	if index < len(v.Elements) {
		return v.Elements[index]
	}
	logger.Logger.Fatalf("index (%d) out of bounds for tuple variable (%s) -- variables: %v", index, v.LongString(), v.Elements)
	return nil
}

func (v *SliceObject) GetId() int64 {
	return v.ObjectInfo.GetId()
}

func (v *SliceObject) NumVariables() int {
	return len(v.Elements)
}

func (v *SliceObject) GetType() gotypes.Type {
	if v.ObjectInfo.GetType() == nil {
		logger.Logger.Fatalf("[VARS ADDRESS] unexpected nil type for slice variable: %s", v.String())
	}
	return v.ObjectInfo.GetType()
}

func (v *SliceObject) GetSliceType() *gotypes.SliceType {
	if userType, ok := v.ObjectInfo.GetType().(*gotypes.UserType); ok {
		return userType.UserType.(*gotypes.SliceType)
	}
	return v.ObjectInfo.GetType().(*gotypes.SliceType)
}

func (v *SliceObject) GetDependencies() []Object {
	return append(v.GetVariableInfo().GetDependencies(), v.Elements...)
}

func (v *SliceObject) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{v}
	if v.GetVariableInfo().HasReferences() {
		deps = append(deps, v.GetVariableInfo().GetReferencesNestedDependencies(includeRefBy, v)...)
	}
	if includeRefBy && v.GetVariableInfo().IsReferencedBy() {
		deps = append(deps, v.GetVariableInfo().GetNestedRefByDependencies(nil)...)
	}
	for _, elem := range v.GetDependencies() { // to include underlying dependencies from variable info
		logger.Logger.Debugf("[SLICE OBJECT] GOT NESTED DEP (INC/ VINFO) FOR ELEM %s (%s)", elem.String(), VariableTypeName(elem))
		deps = append(deps, elem.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}

func (v *SliceObject) GetElementAt(index int) Object {
	if index > len(v.Elements)-1 {
		logger.Logger.Fatalf("[VARS SLICE] element at index (%d) does not exist in array variable with len (%d): %s", index, len(v.Elements), v.String())
	}
	return v.Elements[index]
}

func (v *SliceObject) GetElementAtIfExists(index int) Object {
	if index > len(v.Elements)-1 {
		logger.Logger.Warnf("[VARS SLICE] element at index (%d) does not exist in array variable with len (%d): %s", index, len(v.Elements), v.String())
		return nil
	}
	return v.Elements[index]
}

func (v *SliceObject) String() string {
	return v.ObjectInfo.String()
}

func (v *SliceObject) LongString() string {
	// USE ONLY STRING HERE TO SIMPLIFY DEBUG
	s := v.ObjectInfo.LongString() + " = ("
	for i, elem := range v.Elements {
		s += elem.String()
		if i < len(v.Elements)-1 {
			s += ", "
		}
	}
	return s + ")"
}

func (v *SliceObject) NewVersion() Object {
	copy := &SliceObject{
		ObjectInfo: v.ObjectInfo.Copy(true),
	}
	for k, v := range v.Elements { // FIXME: this is not 100% correct
		copy.Elements[k] = v
		copy.Elements[k].GetVariableInfo().AddParent(copy.Elements[k], copy)
	}
	return copy
}

func (v *SliceObject) NewObject() Object {
	newObject := NewSliceObject(NewObjectInfo(v.GetType()))
	// (1) new slice points to the same elements as the original one
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

func (v *SliceObject) Copy(force bool) Object {
	copy := &SliceObject{ObjectInfo: v.ObjectInfo.Copy(force)}
	for _, v := range v.Elements {
		newElem := v.Copy(force)
		copy.Elements = append(copy.Elements, newElem)
		newElem.GetVariableInfo().AddParent(newElem, copy)
	}
	return copy
}

func (v *SliceObject) DeepCopy() Object {
	logger.Logger.Warnf("[VARS SLICE - DEEP COPY] (%s) %s", VariableTypeName(v), v.String())
	copy := &SliceObject{ObjectInfo: v.ObjectInfo.DeepCopy()}
	for _, v := range v.Elements {
		newElem := v.DeepCopy()
		copy.Elements = append(copy.Elements, newElem)
		newElem.GetVariableInfo().AddParent(newElem, copy)
	}
	return copy
}

func (v *SliceObject) GetUnassaignedVariables() []Object {
	var variables []Object
	if v.GetVariableInfo().IsUnassigned() {
		variables = append(variables, v)
		for _, v := range v.Elements {
			variables = append(variables, v.GetUnassaignedVariables()...)
		}
	}
	return variables
}
