package xcy

import "slices"

type DetectionMode int

const (
	FOREIGN_KEYS_DEFAULT DetectionMode = iota
	FOREIGN_KEYS_LINEAGES
	XCY_ALL_DATASTORES
	XCY_EQUAL_DATASTORES
	DEBUG_LINEAGES
	DEBUG_XCY_MISSING_DEPENDENCIES
	DEBUG_XCY_MINIMIZE_DEPENDENCIES
)

func GetActiveDetectionModes() []DetectionMode {
	return []DetectionMode{
		//XCY_ALL_DATASTORES,
		//XCY_EQUAL_DATASTORES,
		//FOREIGN_KEYS_DEFAULT,
		DEBUG_LINEAGES,
		DEBUG_XCY_MISSING_DEPENDENCIES,
		DEBUG_XCY_MINIMIZE_DEPENDENCIES,
		FOREIGN_KEYS_LINEAGES,
	}
}

func DetectionModeName(detector *XCYDetector) string {
	var detectionMap = map[DetectionMode]string{
		FOREIGN_KEYS_DEFAULT:            "FOREIGN_KEYS_DEFAULT",
		FOREIGN_KEYS_LINEAGES:           "FOREIGN_KEYS_LINEAGES",
		XCY_ALL_DATASTORES:              "XCY_ALL_DATASTORES",
		XCY_EQUAL_DATASTORES:            "XCY_EQUAL_DATASTORES",
		DEBUG_LINEAGES:                  "DEBUG_LINEAGES",
		DEBUG_XCY_MISSING_DEPENDENCIES:  "DEBUG_XCY_MISSING_DEPENDENCIES",
		DEBUG_XCY_MINIMIZE_DEPENDENCIES: "DEBUG_XCY_MINIMIZE_DEPENDENCIES",
	}
	return detectionMap[detector.detectionMode]
}

func DetectionModeUsesLineages(detector *XCYDetector) bool {
	var modesWithLineages = []DetectionMode{
		FOREIGN_KEYS_LINEAGES,
		XCY_ALL_DATASTORES,
		XCY_EQUAL_DATASTORES,
		DEBUG_LINEAGES,
		DEBUG_XCY_MISSING_DEPENDENCIES,
		DEBUG_XCY_MINIMIZE_DEPENDENCIES,
	}
	return slices.Contains(modesWithLineages, detector.detectionMode)
}
