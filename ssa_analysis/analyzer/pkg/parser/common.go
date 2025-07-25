package parser

import "strings"

// getShortFunctionPath returns string of with possible formats:
// - <pkg name>.<member type>.<func type>
// - <pkg name>.<func type>
func getShortFunctionPath(s string) string {
	// remove leading (* if present
	if strings.HasPrefix(s, "(*") {
		s = s[2:]
	}

	// extract everything after "workflow/"
	if idx := strings.Index(s, "workflow/"); idx != -1 {
		s = s[idx+len("workflow/"):]
	}

	s = strings.ReplaceAll(s, ")", "")

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


func extractServiceNameFromShortFunctionPath(s string) string {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return ""
	}
	memberImpl := parts[1]
	member, _ := strings.CutSuffix(memberImpl, "Impl")
	return member
}
