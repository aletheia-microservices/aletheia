package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/tools/go/ssa"
)

func ExtractStringFromValue(val ssa.Value) (string, bool) {
	if c, ok := val.(*ssa.Const); ok {
		return strings.Trim(c.Value.ExactString(), "\""), true
	}
	// EVAL: logrus.Tracef("[UTILS] could not extract string from non-constant: [%T] %v\n", val, val)
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
