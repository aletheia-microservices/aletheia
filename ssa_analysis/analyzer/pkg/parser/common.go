package parser

import "strings"

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
