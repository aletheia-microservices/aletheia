package numerical

import "fmt"

func (detector *NumericalDetector) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "-------------------------- NUMERICAL ANALYSIS -----------------------\n"
	header += "---------------------------------------------------------------------\n"

	var numRequests, numOps int

	for detector.requestInfoStack.Len() > 0 {
		requestInfo := detector.requestInfoStack.Pop().(*RequestInfo)
		if requestInfo.hasPotentialInconsistencies() {
			detector.results += fmt.Sprintf("\n[ENTRY] %s\n", requestInfo.entry.GetMethodStr())
			numRequests++
			for _, op := range requestInfo.getOperations() {
				if op.onNumericalConstraint {
					detector.results += "\t* "
				} else {
					detector.results += "\t- "
				}
				detector.results += fmt.Sprintf("(%s, %s)\n", op.call.Service, op.datastore.GetName())
				detector.results += "\t -> " + op.call.LongString() + "\n"

				// workaround because currently the FullString() method of objects does not show the full content
				// so we are not able to get the insides of mongodb operations such as: $inc(NumClaims, 1)
				if op.repr != "" {
					detector.results += "\t\t " + op.repr + "\n"
				}
				for _, constraint := range op.constraints {
					detector.results += "\t\t @ " + constraint.String() + "\n"
				}
				detector.results += "\n"
				numOps++
			}
		}
	}

	header += fmt.Sprintf(">> (# END-TO-END REQUESTS; # AFFECTED OPERATIONS):\n>> (%d;%d)\n", numRequests, numOps)
	detector.results = header + "---------------------------------------------------------------------\n" + detector.results
}

func (detector *NumericalDetector) GetResults() string {
	return detector.results
}
