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

// e.g.,
// objpath1 (upper): notification
// objpath2 (lower): notification.PostID
func IsUpperOrEqualPath(objpath1 string, objpath2 string) (bool, string) {
	if objpath1 == objpath2 {
		return true, ""
	}
	if strings.HasPrefix(objpath2, objpath1) {
		var diff string
		_, diff, _ = strings.Cut(objpath2, objpath1)
		return true, diff
	}
	return false, ""
}

// InsertAfterFieldName inserts prefix right after the field name of path, where path has the form
// <database>.<table>.<fieldname><subpath> and <subpath> is either empty or starts with '.' or '['
// e.g., InsertAfterFieldName("db.users.address.city", ".home") returns "db.users.address.home.city"
func InsertAfterFieldName(path, prefix string) string {
	// find first '.'
	first := strings.IndexByte(path, '.')
	// find second '.' (after first)
	second := strings.IndexByte(path[first+1:], '.') + first + 1

	// fieldname starts at second+1
	// ends at next '.' or '[' or end
	i := second + 1
	for i < len(path) && path[i] != '.' && path[i] != '[' {
		i++
	}
	// insert prefix right after fieldname
	return path[:i] + prefix + path[i:]
}
