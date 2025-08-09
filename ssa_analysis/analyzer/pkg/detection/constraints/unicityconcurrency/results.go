package unicityconcurrency

import "fmt"

func (detector *UnicityConcurrencyDetector) GetResults() string {
	return detector.summary
}

func (detector *UnicityConcurrencyDetector) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "------------------- ANALYSIS: UNICITY CONCURRENCY -------------------\n"
	header += "---------------------------------------------------------------------\n"

	var results string
	for request, vulnerableWriteSets := range detector.vulnerableWriteSets {
		results += fmt.Sprintf("entry request: %s\n", request.entry.String())
		for i, writeSet := range vulnerableWriteSets {
			results += fmt.Sprintf("\twrite (constrained) #%d: %s\n", i, writeSet.constrainedOp.call.String())
			for _, field := range writeSet.constrainedFields {
				results += fmt.Sprintf("\t\tfield (constrained): %s\n", field.GetPath())
			}
			for _, op := range writeSet.otherOps {
				results += fmt.Sprintf("\t\tAFFECTED WRITE: %s\n", op.call.String())
			}
		}
		results += "\n"
	}
	detector.summary = header + results
}
