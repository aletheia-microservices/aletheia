package utils

import "regexp"

// Add space before "Service"
// e.g. StorageService → Storage Service
func GetSpacedName(name string) string {
	re := regexp.MustCompile(`(?i)(Service)$`) // case-insensitive match for "Service"
    if re != nil {
        return re.ReplaceAllString(name, " Service")
    }
    return name
}
