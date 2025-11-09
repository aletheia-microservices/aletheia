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

// e.g.,
// objpath1 (upper): notification
// objpath2 (lower): notification.PostID
func IsUpperPath(objpath1 string, objpath2 string) (bool, string) {
	if objpath1 != objpath2 && strings.HasPrefix(objpath2, objpath1) {
		var diff string
		_, diff, _ = strings.Cut(objpath2, objpath1)
		return true, diff
	}
	return false, ""
}
