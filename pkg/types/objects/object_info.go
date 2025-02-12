package objects

import (
	"encoding/json"
	"fmt"
	"runtime"
	"slices"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
)

type ObjectInfo struct {
	Name string
	Type gotypes.Type
	Id   int64

	// dependencies are added to variables computed from returns
	// from function calls that cannot be statically analyzed
	// in practice, all returned variables will have all the parameters passed to the function call has dependencies
	Dependencies []Object

	ReferencedBy []Object
	References   []*Reference
	Parents      []Object

	IsBlockParam  bool
	BlockParamIdx int

	Dataflows         []*ObjectDataflow
	IndirectDataflows []*ObjectDataflow

	Dynamic bool
}

func NewObjectInfoInline(t gotypes.Type) *ObjectInfo {
	return &ObjectInfo{
		Type: t,
		Id:   VARIABLE_INLINE_ID,
	}
}

func (vinfo *ObjectInfo) SetDynamic() {
	vinfo.Dynamic = true
}

func (vinfo *ObjectInfo) IsDynamic() bool {
	return vinfo.Dynamic
}

func (vinfo *ObjectInfo) ResetAllDataflows() {
	vinfo.Dataflows = nil
	vinfo.IndirectDataflows = nil
}

func (vinfo *ObjectInfo) GetDependencies() []Object {
	return vinfo.Dependencies
}

func (vinfo *ObjectInfo) SetType(t gotypes.Type) {
	vinfo.Type = t
}

func (vinfo *ObjectInfo) GetDataflows() []*ObjectDataflow {
	return vinfo.Dataflows
}

func (vinfo *ObjectInfo) AddDependency(dep Object) {
	vinfo.Dependencies = append(vinfo.Dependencies, dep)
}

func (vinfo *ObjectInfo) AddDependencies(deps []Object) {
	vinfo.Dependencies = append(vinfo.Dependencies, deps...)
}

func (vinfo *ObjectInfo) GetIndirectDataflows() []*ObjectDataflow {
	return vinfo.IndirectDataflows
}

func (vinfo *ObjectInfo) GetAllDataflows() []*ObjectDataflow {
	return append(vinfo.Dataflows, vinfo.IndirectDataflows...)
}

func (vinfo *ObjectInfo) GetAllDataflowsForDatastore(datastore string) []*ObjectDataflow {
	var dataflows []*ObjectDataflow
	for _, df := range vinfo.Dataflows {
		if df.IsOpInDatastore(datastore) {
			dataflows = append(dataflows, df)
		}
	}
	for _, df := range vinfo.IndirectDataflows {
		if df.IsOpInDatastore(datastore) {
			dataflows = append(dataflows, df)
		}
	}
	return dataflows
}

func (vinfo *ObjectInfo) GetAllReadDataflowsForDatastore(datastore string) []*ObjectDataflow {
	var dataflows []*ObjectDataflow
	for _, df := range vinfo.Dataflows {
		if df.IsOpInDatastore(datastore) && !df.IsWriteOp() {
			dataflows = append(dataflows, df)
		}
	}
	for _, df := range vinfo.IndirectDataflows {
		if df.IsOpInDatastore(datastore) && !df.IsWriteOp() {
			dataflows = append(dataflows, df)
		}
	}
	return dataflows
}

func (vinfo *ObjectInfo) GetAllReadDataflowsExceptDatastore(datastore string) []*ObjectDataflow {
	var dataflows []*ObjectDataflow
	for _, df := range vinfo.Dataflows {
		if !df.IsOpInDatastore(datastore) && !df.IsWriteOp() {
			dataflows = append(dataflows, df)
		}
	}
	for _, df := range vinfo.IndirectDataflows {
		if !df.IsOpInDatastore(datastore) && !df.IsWriteOp() {
			dataflows = append(dataflows, df)
		}
	}
	return dataflows
}

func (vinfo *ObjectInfo) GetAllWriteDataflowsForDatastore(datastore string) []*ObjectDataflow {
	var dataflows []*ObjectDataflow
	for _, df := range vinfo.Dataflows {
		if df.IsOpInDatastore(datastore) && df.IsWriteOp() {
			dataflows = append(dataflows, df)
		}
	}
	for _, df := range vinfo.IndirectDataflows {
		if df.IsOpInDatastore(datastore) && df.IsWriteOp() {
			dataflows = append(dataflows, df)
		}
	}
	return dataflows
}

func (vinfo *ObjectInfo) GetAllWriteDataflows() []*ObjectDataflow {
	var dataflows []*ObjectDataflow
	for _, df := range vinfo.Dataflows {
		if df.IsWriteOp() {
			dataflows = append(dataflows, df)
		}
	}
	for _, df := range vinfo.IndirectDataflows {
		if df.IsWriteOp() {
			dataflows = append(dataflows, df)
		}
	}
	return dataflows
}

func (vinfo *ObjectInfo) GetAllReadDataflows() []*ObjectDataflow {
	var dataflows []*ObjectDataflow
	for _, df := range vinfo.Dataflows {
		if !df.IsWriteOp() {
			dataflows = append(dataflows, df)
		}
	}
	for _, df := range vinfo.IndirectDataflows {
		if !df.IsWriteOp() {
			dataflows = append(dataflows, df)
		}
	}
	return dataflows
}

func (vinfo *ObjectInfo) Copy(force bool) *ObjectInfo {
	if vinfo == nil {
		logger.Logger.Fatal("?????")
	}
	// skip deep copy, especially since we can have references and dataflows
	if !force {
		return vinfo
	}

	var dataflows []*ObjectDataflow
	for _, df := range vinfo.Dataflows {
		dataflows = append(dataflows, df.Copy(force))
	}
	var indirectDataflows []*ObjectDataflow
	for _, df := range vinfo.IndirectDataflows {
		indirectDataflows = append(indirectDataflows, df.Copy(force))
	}
	var refsCopy []*Reference
	for _, r := range vinfo.GetReferences() {
		refsCopy = append(refsCopy, r.Copy(force).(*Reference))
	}
	return &ObjectInfo{
		Name:              vinfo.Name,
		Type:              vinfo.Type,
		Id:                vinfo.Id,
		References:        refsCopy,
		IsBlockParam:      vinfo.IsBlockParam,
		BlockParamIdx:     vinfo.BlockParamIdx,
		Dataflows:         dataflows,
		IndirectDataflows: indirectDataflows,
	}
}

func (vinfo *ObjectInfo) DeepCopy() *ObjectInfo {
	if vinfo == nil {
		logger.Logger.Fatal("?????")
	}

	var dataflows []*ObjectDataflow
	for _, df := range vinfo.Dataflows {
		dataflows = append(dataflows, df.DeepCopy())
	}
	var indirectDataflows []*ObjectDataflow
	for _, df := range vinfo.IndirectDataflows {
		indirectDataflows = append(indirectDataflows, df.DeepCopy())
	}
	var refsCopy []*Reference
	for _, r := range vinfo.GetReferences() {
		logger.Logger.Infof("[VARS INFO - DEEP COPY] deep copy reference (%s) for creator = %s", vinfo.String(), r.Creator)
		refsCopy = append(refsCopy, r.DeepCopy().(*Reference))
	}
	return &ObjectInfo{
		Name:              vinfo.Name,
		Type:              vinfo.Type,
		Id:                vinfo.Id,
		References:        refsCopy,
		IsBlockParam:      vinfo.IsBlockParam,
		BlockParamIdx:     vinfo.BlockParamIdx,
		Dataflows:         dataflows,
		IndirectDataflows: indirectDataflows,
	}
}

func (vinfo *ObjectInfo) MarshalJSON() ([]byte, error) {
	/* refStr := "<nil>"
	if v.Reference != nil {
		refStr = v.Reference.Variable.GetVariableInfo().GetName() + fmt.Sprintf(" (%d) @ ", v.Reference.Variable.GetId()) + v.Reference.Creator
	} */
	return json.Marshal(&struct {
		Name string `json:"name"`
		Type string `json:"type"`
		//Id   int64  `json:"id,omitempty"`
		//Reference bool   `json:"ref,omitempty"`
		References []*Reference `json:"refs,omitempty"`
	}{
		Name: vinfo.Name,
		Type: vinfo.Type.String(),
		//Id:   vinfo.Id,
		//Reference: v.Reference != nil,
		References: vinfo.GetReferences(),
	})
}

func (vinfo *ObjectInfo) String() string {
	if vinfo.Type != nil {
		if vinfo.Name != "" {
			return fmt.Sprintf("%s %s", vinfo.Name, vinfo.Type.String())
		}
		return vinfo.Type.String()
	}
	return fmt.Sprintf("%s (unknown)", vinfo.Name)
}

func (vinfo *ObjectInfo) LongString() string {
	if vinfo.Type != nil {
		if vinfo.Name != "" {
			return fmt.Sprintf("%s %s", vinfo.Name, vinfo.Type.LongString())
		}
		return vinfo.Type.LongString()
	}
	return fmt.Sprintf("%s (unknown)", vinfo.Name)
}

func (vinfo *ObjectInfo) IsBlockParameter() bool {
	return vinfo.IsBlockParam
}

func (vinfo *ObjectInfo) EqualBlockParamIndex(idx int) bool {
	return vinfo.BlockParamIdx == idx
}

func (vinfo *ObjectInfo) SetName(name string) {
	vinfo.Name = name
}

func (vinfo *ObjectInfo) SetUnassigned() {
	vinfo.Id = VARIABLE_UNASSIGNED_ID
}

func (vinfo *ObjectInfo) IsUnassigned() bool {
	return vinfo.Id == VARIABLE_UNASSIGNED_ID || vinfo.Id == VARIABLE_INLINE_ID
}

func (vinfo *ObjectInfo) HasReferences() bool {
	return vinfo.References != nil
}

func (vinfo *ObjectInfo) IsReferencedBy() bool {
	if vinfo.Name == "postID_STORAGE_SVC" {
		logger.Logger.Warnf("HERE: %s", vinfo.String())
	}
	return vinfo.ReferencedBy != nil
}

func (vinfo *ObjectInfo) AssignID(id int64) {
	vinfo.Id = id
}

func (vinfo *ObjectInfo) AddReferenceWithID(source Object, target Object, creator string) {
	targetInfo := target.GetVariableInfo()
	targetInfo.ReferencedBy = append(targetInfo.ReferencedBy, source)
	vinfo.Id = target.GetId()
	vinfo.References = append(vinfo.References, &Reference{
		Creator: creator,
		Object:  target,
		Id:      target.GetId(),
	})
	logger.Logger.Debugf("[VARS INFO] added new reference (%s) @ (%s) for variable (%s) with references list (len=%d): %v", target.String(), creator, vinfo.String(), len(vinfo.GetReferences()), vinfo.GetReferences())
}

func (vinfo *ObjectInfo) GetName() string {
	return vinfo.Name
}

func (vinfo *ObjectInfo) GetId() int64 {
	return vinfo.Id
}

func (vinfo *ObjectInfo) GetType() gotypes.Type {
	return vinfo.Type
}

func (vinfo *ObjectInfo) GetReferences() []*Reference {
	return vinfo.References
}

func (vinfo *ObjectInfo) SetParent(current Object, parent Object) {
	vinfo.AddParent(current, parent)
}

func (vinfo *ObjectInfo) AddParent(current Object, parent Object) {
	if current == parent || vinfo == parent.GetVariableInfo() {
		pc, file, line, ok := runtime.Caller(1)
		if !ok {
			logger.Logger.Fatalf("RECURSION!! (%s) %s", VariableTypeName(parent), parent.String())
		}
		callerFunc := runtime.FuncForPC(pc).Name()
		logger.Logger.Fatalf("RECURSION!! (%s) %s \n\t\t\t\t\t(caller: %s) \n\t\t\t\t\t %s:%d", VariableTypeName(parent), parent.String(), callerFunc, file, line)
	}
	vinfo.Parents = append(vinfo.Parents, parent)
}
func (vinfo *ObjectInfo) GetReferencesNestedDependencies(includeRefBy bool, v Object) []Object {
	var deps []Object
	for _, ref := range vinfo.References {
		deps = append(deps, ref.GetNestedDependencies(false)...)
	}
	var visited []*ObjectInfo
	deps = append(deps, vinfo.GetNestedRefByDependencies(visited)...)
	return deps
}

func (vinfo *ObjectInfo) GetNestedRefByDependencies(visited []*ObjectInfo) []Object {
	var deps []Object
	if !slices.Contains(visited, vinfo) {
		visited = append(visited, vinfo)
		for _, refBy := range vinfo.ReferencedBy {
			deps = append(deps, refBy)
			deps = append(deps, refBy.GetVariableInfo().GetNestedRefByDependencies(visited)...)
		}
	}
	return deps
}
