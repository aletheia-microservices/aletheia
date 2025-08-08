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
	fmt.Printf("[UTILS] could not extract string from non-constant: [%T] %v", val, val)
	return "", false
}
