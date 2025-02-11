package xcy

import (
	"encoding/json"
	"fmt"
	"slices"

	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
	"analyzer/pkg/types/objects"
)

type Operation struct {
	ID                     int
	LineageID              int
	Service                string
	Method                 string
	Key                    objects.Object
	Object                 objects.Object
	Datastore              *datastores.Datastore
	Write                  bool
	VisibleDependencies    []*Operation
	MinimizedDependencySet []*Operation
}

func (op *Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Service  string `json:"service"`
		Method   string `json:"method"`
		Key      string `json:"key"`
		Object   string `json:"object"`
		Database string `json:"database"`
	}{
		Service:  op.Service,
		Method:   op.Method,
		Key:      fmt.Sprintf("%s %s", op.Key.GetVariableInfo().GetName(), op.Key.GetVariableInfo().GetType().GetLongName()),
		Object:   fmt.Sprintf("%s %s", op.Object.GetVariableInfo().GetName(), op.Object.GetVariableInfo().GetType().GetLongName()),
		Database: fmt.Sprintf("%s %s", op.Datastore.GetName(), op.Datastore.GetTypeLongName()),
	})
}

func (op *Operation) IsAfter(other *Operation) bool {
	if op.LineageID > other.LineageID {
		return true
	}
	if op.LineageID < other.LineageID {
		return false
	}
	return op.ID > other.ID
}

func (op *Operation) AddMinimizedDependencyIfNotExists(dep *Operation) {
	if !slices.Contains(op.MinimizedDependencySet, dep) {
		op.MinimizedDependencySet = append(op.MinimizedDependencySet, dep)
	}
}

func (op *Operation) GetVisibleDependenciesForDatastore(datastore *datastores.Datastore) []*Operation {
	var dependencies []*Operation
	for _, dep := range op.VisibleDependencies {
		if dep.HasDatastore(datastore) {
			dependencies = append(dependencies, dep)
		}
	}
	return dependencies
}

func (op *Operation) GetMinimizedDependencySet() []*Operation {
	return op.MinimizedDependencySet
}

func (op *Operation) HasDatastore(datastore *datastores.Datastore) bool {
	return op.Datastore == datastore
}

func (op *Operation) GetAllUnderlyingVariables(includeRefBy bool) []objects.Object {
	return append(op.Key.GetNestedDependencies(false), op.Object.GetNestedDependencies(false)...)
}

func (op *Operation) GetKeyUnderlyingVariables(includeRefBy bool) []objects.Object {
	return op.Key.GetNestedDependencies(includeRefBy)
}

func (op *Operation) GetObjectUnderlyingVariables(includeRefBy bool) []objects.Object {
	return op.Object.GetNestedDependencies(includeRefBy)
}

func (op *Operation) GetVisibleDependencies() []*Operation {
	return op.VisibleDependencies
}

func (op *Operation) AppendVisibleDependencies(dependencies []*Operation) {
	op.VisibleDependencies = append(op.VisibleDependencies, dependencies...)
}

func (op *Operation) GetDatastore() *datastores.Datastore {
	return op.Datastore
}

func (op *Operation) String() string {
	if op.Key == nil {
		logger.Logger.Fatalf("nil key in %s", op.Method)
	}
	if op.Object == nil {
		logger.Logger.Fatalf("nil object in %s", op.Method)
	}
	opType := "write"
	if !op.Write {
		opType = "read"
	}
	return fmt.Sprintf("(%d, %d) %s(%s, %s, %s) @ %s",
		op.LineageID,
		op.ID,
		opType,
		op.Datastore.GetName(),
		op.Key.GetVariableInfo().GetName(),
		op.Object.GetVariableInfo().GetName(),
		op.Service,
	)
}

func (operation *Operation) GetVisibleDependenciesString() []string {
	var visibleDependencies []string
	for _, op := range operation.GetVisibleDependencies() {
		visibleDependencies = append(visibleDependencies, op.String())
	}
	return visibleDependencies
}

func (operation *Operation) GetMinimizedDependencySetString() []string {
	var visibleDependencies []string
	for _, op := range operation.GetMinimizedDependencySet() {
		visibleDependencies = append(visibleDependencies, op.String())
	}
	return visibleDependencies
}
