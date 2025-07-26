package utils

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/ssa"
)

func FieldIndexToName(t *ssa.FieldAddr) string {
	return t.X.Type().Underlying().(*types.Pointer).Elem().(*types.Named).Underlying().(*types.Struct).Field(t.Field).Name()
}

// getShortFunctionPath returns string of with possible formats:
// - <pkg name>.<member type>.<func name>
// - <pkg name>.<func name>
func GetShortFunctionPath(s string) string {
	// remove leading (* if present
	if strings.HasPrefix(s, "(*") {
		s = s[2:]
	}

	// extract everything after "workflow/"
	if idx := strings.Index(s, "workflow/"); idx != -1 {
		s = s[idx+len("workflow/"):]
	}

	s = strings.ReplaceAll(s, ")", "")

	// split into parts
	parts := strings.Split(s, ".")
	if len(parts) == 3 {
		// remove "Impl" from member type
		parts[1] = strings.ReplaceAll(parts[1], "Impl", "")
		s = strings.Join(parts, ".")
	}

	// remove $... suffix if present
	/* if idx := strings.IndexByte(s, '$'); idx != -1 {
		s = s[:idx]
	} */

	return s
}

/* func getPathWithoutFunctionName(s string) string {
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return s // nothing to remove
	}
	return strings.Join(parts[:len(parts)-1], ".")
} */


func ExtractServiceNameFromShortFunctionPath(s string) string {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return ""
	}
	memberImpl := parts[1]
	member, _ := strings.CutSuffix(memberImpl, "Impl")
	return member
}

func ExtractMethodNameFromShortFunctionPath(s string) string {
	parts := strings.Split(s, ".")
	return parts[len(parts)-1]
}
