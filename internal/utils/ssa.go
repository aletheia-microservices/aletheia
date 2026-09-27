package utils

import (
	"crypto/rand"
	"fmt"
	"go/constant"
	"go/token"
	"math/big"
	"strings"

	"golang.org/x/tools/go/ssa"
)

func ExtractStringFromValue(val ssa.Value) (string, bool) {
	if c, ok := val.(*ssa.Const); ok {
		return strings.Trim(c.Value.ExactString(), "\""), true
	}
	return "", false
}

// ExtractBaseStringFromValue returns the constant string that val starts with, and whether it is the complete
// string (complete = false when parts are added at runtime), e.g., for a SQL query built as below it returns
// the value of baseQuery and complete = false:
//
//	var baseQuery = "SELECT ... FROM sock"
//	query := baseQuery
//	if len(tags) > 0 {
//		query += " WHERE tag.Name=?"
//	}
//	query += ";"
//
// it follows the left operand of each concatenation (e.g., query + ";"), the incoming values of phis
// (e.g., query after the if), and package-level variables set to a constant in the package init
func ExtractBaseStringFromValue(val ssa.Value) (str string, complete bool, ok bool) {
	return extractBaseStringFromValue(val, make(map[ssa.Value]bool))
}

func extractBaseStringFromValue(val ssa.Value, visited map[ssa.Value]bool) (string, bool, bool) {
	if visited[val] {
		return "", false, false
	}
	visited[val] = true

	switch v := val.(type) {
	case *ssa.Const:
		// e.g., "SELECT ... FROM sock"
		if v.Value != nil && v.Value.Kind() == constant.String {
			return constant.StringVal(v.Value), true, true
		}
	case *ssa.UnOp:
		// e.g., t0 = *baseQuery
		if global, ok := v.X.(*ssa.Global); ok && v.Op == token.MUL {
			str, ok := extractStringFromGlobal(global)
			return str, ok, ok
		}
	case *ssa.BinOp:
		// e.g., t10 = t9 + ";"
		if v.Op == token.ADD {
			str, _, ok := extractBaseStringFromValue(v.X, visited)
			return str, false, ok
		}
	case *ssa.Phi:
		// e.g., t9 = phi [t0, t7] (query)
		for _, edge := range v.Edges {
			if str, _, ok := extractBaseStringFromValue(edge, visited); ok {
				return str, false, true
			}
		}
	}
	return "", false, false
}

// extractStringFromGlobal returns the constant string stored in a package-level variable by the package init,
// e.g., *baseQuery = "SELECT ... FROM sock"
func extractStringFromGlobal(global *ssa.Global) (string, bool) {
	init := global.Pkg.Func("init")
	if init == nil {
		return "", false
	}
	for _, block := range init.Blocks {
		for _, instr := range block.Instrs {
			if store, ok := instr.(*ssa.Store); ok && store.Addr == global {
				if c, ok := store.Val.(*ssa.Const); ok && c.Value != nil && c.Value.Kind() == constant.String {
					return constant.StringVal(c.Value), true
				}
			}
		}
	}
	return "", false
}

type FUNC_TYPE int

const (
	FUNC_TYPE_IGNORE FUNC_TYPE = iota
	FUNC_TYPE_APPEND
	FUNC_TYPE_TRANSFER
	FUNC_TYPE_MAP_ELEMS
)

// direct => can be tainted
func SSABuiltinFuncIsDirect(builtin *ssa.Builtin) (bool, FUNC_TYPE, string) {
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
		return true, FUNC_TYPE_APPEND, builtin.Name()
	case "copy":
		return true, FUNC_TYPE_TRANSFER, builtin.Name()
	case "delete":
		return true, FUNC_TYPE_MAP_ELEMS, builtin.Name()
	}
	return false, FUNC_TYPE_IGNORE, ""
}

func ComputeInstructionID(instr ssa.Instruction) string {
	if !instr.Pos().IsValid() { // meaning there is no position
		n, err := rand.Int(rand.Reader, big.NewInt(1<<31))
		if err != nil {
			return ""
		}
		return "instr_" + instrString(instr) + "_" + fmt.Sprintf("%d", n)
	}
	return "instr_" + instrString(instr) + "_" + fmt.Sprintf("%d", instr.Pos())
}

func instrString(instr ssa.Instruction) string {
	switch instr.(type) {
	case *ssa.Store:
		return "store"
	case *ssa.Return:
		return "ret"
	}
	return ""
}
