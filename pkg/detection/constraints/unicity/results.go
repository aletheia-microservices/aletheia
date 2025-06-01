package unicity

import "fmt"

func (detector *UnicityDetector) GetResults() string {
	return detector.results
}

func (detector *UnicityDetector) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "--------------------------- UNICITY ANALYSIS ------------------------\n"
	header += "---------------------------------------------------------------------\n"

	var numRequests, numOps int

	for detector.requestInfoStack.Len() > 0 {
		requestInfo := detector.requestInfoStack.Pop().(*RequestInfo)
		if requestInfo.hasPotentialInconsistencies() {
			detector.results += fmt.Sprintf("\n[ENTRY] %s\n", requestInfo.entry.GetMethodStr())
			numRequests++
			for _, op := range requestInfo.getOperations() {
				if op.AffectsOps() {
					detector.results += fmt.Sprintf("[%d] (%s, %s)\n", op.idx, op.call.Service, op.datastore.GetName())
					detector.results += "-> " + op.call.String() + "\n"

					detector.results += "\t -------------------- UNICITY CONSTRAINTS --------------------\n"
					for _, constraint := range op.constraints {
						detector.results += "\t @ " + constraint.String() + "\n"
					}

					detector.results += "\t -------------------------------------------------------------\n"
					detector.results += "\t -------------------- AFFECTED OPERATIONS --------------------\n"
					for affectedOp, refFields := range op.GetAffectedOps() {
						detector.results += fmt.Sprintf("\t [%d] (%s, %s)\n", affectedOp.idx, affectedOp.call.Service, affectedOp.datastore.GetName())
						detector.results += "\t -> " + affectedOp.call.String() + "\n"
						detector.results += "\t\t referencing fields from prev. operation:\n"
						for _, field := range refFields {
							detector.results += "\t\t - " + field.GetFullName() + "\n"
						}
						numOps++
					}
					detector.results += "\t -------------------------------------------------------------\n"
				}
				detector.results += "\n"
			}
		}
	}

	header += fmt.Sprintf(">> (# END-TO-END REQUESTS; # AFFECTED OPERATIONS):\n>> (%d;%d)\n", numRequests, numOps)
	detector.results = header + "---------------------------------------------------------------------\n" + detector.results
}
