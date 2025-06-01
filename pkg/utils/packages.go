package utils

import (
	"strings"
)

func RemoveQuotesFromPathImport(path string) string {
	return strings.Trim(path, "\"")
}

func IsAppPackage(appPath string, packagePath string) bool {
	return strings.HasPrefix(packagePath, appPath)
}

func IsBlueprintBackendPath(path string) bool {
	return path == BLUEPRINT_PATH_CORE_BACKEND
}
