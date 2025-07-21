package utils

import (
	"go/types"

	"golang.org/x/tools/go/ssa"
)

func FieldIndexToName(t *ssa.FieldAddr) string {
	return t.X.Type().Underlying().(*types.Pointer).Elem().(*types.Named).Underlying().(*types.Struct).Field(t.Field).Name()
}
