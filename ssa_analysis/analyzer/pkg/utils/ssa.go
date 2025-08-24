package utils

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/ssa"
)

func ExtractStringFromValue(val ssa.Value) (string, bool) {
	if c, ok := val.(*ssa.Const); ok {
		return strings.Trim(c.Value.ExactString(), "\""), true
	}
	fmt.Printf("[UTILS] could not extract string from non-constant: [%T] %v\n", val, val)
	return "", false
}

func SSAValueIsBuiltinFuncCall(val ssa.Value) (bool, *ssa.Builtin) {
	if call, ok := val.(*ssa.Call); ok && !call.Call.IsInvoke() {
		if builtin, ok := call.Call.Value.(*ssa.Builtin); ok {
			return true, builtin
		}
	}
	return false, nil
}

// direct => can be tainted
func SSABuiltinFuncIsDirect(builtin *ssa.Builtin) (bool, bool, string) {
	// append(slice []Type, elems ...Type) []Type
	// -----------------------------------
	// copy(dst, src []Type) int
	// delete(m map[Type]Type1, key Type)
	// -----------------------------------
	// len(v Type) int
	// cap(v Type) int
	// make(t Type, size ...IntegerType) Type
	// max[T cmp.Ordered](x T, y ...T) T
	// new(Type) *Type
	// complex(r, i FloatType) ComplexType
	// real(c ComplexType) FloatType
	// imag(c ComplexType) FloatType
	// clear[T ~[]Type | ~map[Type]Type1](t T)
	// close(c chan<- Type)
	// panic(v any)
	// recover() any
	// print(args ...Type)
	// println(args ...Type)
	// error
	switch builtin.Name() {
	case "append":
		return true, true, builtin.Name()
	case "copy", "delete":
		return true, false, builtin.Name()
	}
	return false, false, ""
}
