package unicityconcurrency

import (
	"fmt"

	"analyzer/pkg/app"
)

func (detector *UnicityConcurrencyDetector) GetResults() string {
	return detector.results
}

func (detector *UnicityConcurrencyDetector) ComputeResults(app *app.App) {
	header := "---------------------------------------------------------------------\n"
	header += "------------------- ANALYSIS: UNICITY CONCURRENCY -------------------\n"
	header += "---------------------------------------------------------------------\n"

	var numWarnings int
	var results string
	for request, vulnerableWriteSets := range detector.vulnerableWriteSets {
		var ok bool
		for _, writeSet := range vulnerableWriteSets {
			if len(writeSet.otherOps) > 0 {
				ok = true
			}
		}
		if !ok {
			continue
		}
		results += fmt.Sprintf("entry request: %s()\n", request.entry.String())
		for _, writeSet := range vulnerableWriteSets {
			if request.entry.String() != writeSet.constrainedOp.call.GetFromNode().String() {
				results += fmt.Sprintf("write (origin): %s() ... %s\n", request.entry.String(), writeSet.constrainedOp.call.String())
			} else {
				results += fmt.Sprintf("write (origin): %s\n", writeSet.constrainedOp.call.String())
			}

			for _, field := range writeSet.constrainedFields {
				results += fmt.Sprintf("\t\t- field (constrained): %s", field.GetPath())
				if field.IsPrimaryKey() {
					results += " (PRIMARY KEY)"
				} else if field.IsUnique() {
					results += " (UNIQUE)"
				}
				results += "\n"
			}
			for _, op := range writeSet.otherOps {
				numWarnings++
				results += fmt.Sprintf("\t- affected write #%d: %s\n", numWarnings, op.call.String())
			}
		}
		results += "\n"
	}
	detector.results = header + fmt.Sprintf("[NUM_WARNINGS = %d]\n", numWarnings) + results

}
