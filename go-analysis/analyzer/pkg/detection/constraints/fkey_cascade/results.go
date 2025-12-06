package fkey_cascade

import (
	"fmt"

	"analyzer/pkg/utils"
)

func (detector *CascadeDetector) GetResults() string {
	return detector.results
}

func (detector *CascadeDetector) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "-------------------- FOREIGN KEY CASCADE ANALYSIS -------------------\n"
	header += "---------------------------------------------------------------------\n"

	detector.checkInconsistencies()

	header += fmt.Sprintf(">> (# DELETES ON REFERENCED OBJECT; # ABSENCE OF CASCADING DELETES):\n>> (%d;%d)\n", detector.numDeletes, detector.numMissingCascadingDeletes)
	detector.results = header + "---------------------------------------------------------------------\n" + utils.TEXT_RESET_COLOR + detector.results
}

func (detector *CascadeDetector) checkInconsistencies() {
	for detector.requestInfoStack.Len() > 0 {
		requestInfo := detector.requestInfoStack.Pop().(*RequestInfo)
		for i, op := range requestInfo.getDeleteOperations() {
			depsWithMissingCascading := op.getDependenciesWithMissingCascade()
			detector.results += fmt.Sprintf("[%d] delete with %d missing cascades:\n", i+1, len(depsWithMissingCascading))
			detector.results += fmt.Sprintf("%s: %s\n", op.getCall().GetCallerStr(), op.call.ShortString())
			detector.numDeletes++
			for _, dep := range depsWithMissingCascading {
				if !dep.cascading {
					detector.results += fmt.Sprintf("\t- %s\n", dep.LongString())
					detector.numMissingCascadingDeletes++
				}
			}
		}
	}
}
