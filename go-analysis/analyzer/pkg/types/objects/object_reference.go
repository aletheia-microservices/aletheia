package objects

import (
	"encoding/json"
	"fmt"
)

type Reference struct {
	Object
	Creator string
	Id      int64
}

func (ref *Reference) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name    string `json:"name"`
		Creator string `json:"creator"`
		//Id      int64  `json:"id"`
		//Variable Variable `json:"variable"`
		/* IndirectDeps string `json:"indirect_deps"` */
		/* Params 			[]Variable 		`json:"ref_params"` */
	}{
		Name:    ref.Object.GetVariableInfo().GetName(),
		Creator: ref.Creator,
		//Id:      ref.Variable.GetId(),
		//Variable: ref.Variable,
		/* IndirectDeps: getDependenciesString(GetIndirectDependencies(ref)...), */
		/* Params:  		ref.Variable.GetDependencies(), */
	})
}

func (ref *Reference) Copy(force bool) Object {
	return &Reference{
		Object:  ref.Object.Copy(force),
		Creator: ref.Creator,
	}
}

func (ref *Reference) DeepCopy() Object {
	return &Reference{
		Object:  ref.Object.DeepCopy(),
		Creator: ref.Creator,
	}
}

func (ref *Reference) String() string {
	if ref.Object != nil {
		return fmt.Sprintf("ref <%s> @ %s", ref.Object.String(), ref.Creator)
	}
	return "nil-reference"
}

func (ref *Reference) GetVariableInfo() *ObjectInfo {
	if ref.Object != nil {
		return ref.Object.GetVariableInfo()
	}
	return nil
}

func (ref *Reference) GetDependencies() []Object {
	if ref.Object != nil {
		return ref.Object.GetDependencies()
	}
	return nil
}

func (ref *Reference) GetNestedDependencies(includeRefBy bool) []Object {
	var deps = []Object{ref}
	if ref.Object != nil {
		deps = append(deps, ref.Object.GetNestedDependencies(includeRefBy)...)
	}
	return deps
}
