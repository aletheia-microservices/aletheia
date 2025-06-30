package gotypes

import (
	"fmt"
	"go/ast"

	"analyzer/pkg/logger"
)

type FuncTypeType struct {
	Type   `json:"-"`
	Body   *ast.BlockStmt
	Name   string
	Params *ast.FieldList
}

// ------------
// Type Methods
// ------------

func (t *FuncTypeType) DeepCopy() Type {
	return &FuncTypeType{
		Body:   t.Body,
		Name:   t.Name,
		Params: t.Params,
	}
}

func (t *FuncTypeType) GetBody() *ast.BlockStmt {
	return t.Body
}

func (t *FuncTypeType) GetParams() *ast.FieldList {
	return t.Params
}

func (t *FuncTypeType) IsSameType(other Type) bool {
	_, ok := other.(*FuncTypeType)
	return ok
}
func (t *FuncTypeType) String() string {
	return t.Name
}
func (t *FuncTypeType) LongString() string {
	s := fmt.Sprintf("%s(", t.Name)
	i := 0
	for _, m := range t.Params.List {
		s += m.Names[0].Name
		if i < len(t.Params.List) {
			s += ", "
		}
		i++
	}
	return s + "}"
}
func (t *FuncTypeType) GetName() string {
	return t.String()
}
func (t *FuncTypeType) SetName(name string) {
	t.Name = name
}
func (t *FuncTypeType) GetLongName() string {
	return t.String()
}
func (t *FuncTypeType) GetBasicValue() string {
	logger.Logger.Fatalf("unable to get value for FuncType type %s", t.String())
	return ""
}
func (t *FuncTypeType) AddValue(value string) {
	logger.Logger.Fatalf("unable to add value for FuncType type %s", t.String())
}

// -----------------
// FuncType Methods
// -----------------

func (t *FuncTypeType) SetParams(params *ast.FieldList) {
	t.Params = params
}
