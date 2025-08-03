package foreignkeycoordination

import "fmt"

func (detector *ForeignKeyCoordinationDetector) GetResults() string {
	return detector.results
}

func (detector *ForeignKeyCoordinationDetector) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "---------------- ANALYSIS: FOREIGN KEY COORDINATION -----------------\n"
	header += "---------------------------------------------------------------------\n"

	var results string
	for request, foreignreads := range detector.inconsistencies {
		results += fmt.Sprintf("request #%d @ %s:\n", request.idx, request.entry.String())
		for i, foreignread := range foreignreads {
			results += fmt.Sprintf("\tforeign read #%d:\n", i)
			results += fmt.Sprintf("\t\tread_1 = %s\n", foreignread.read1.call.String())
			results += fmt.Sprintf("\t\tread_2 = %s\n", foreignread.read2.call.String())
		}
	}
	detector.results = header + results
}
