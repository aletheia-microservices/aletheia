package utils

import "strings"

func ExtractUpperPath(objpath string) (string, string, bool) {
	idx := strings.LastIndex(objpath, ".")
	if idx == -1 {
		return "", "", false
	}
	// containts . before e.g. (.ID)
	subpath := objpath[idx:]
	return objpath[:idx], subpath, true
}
