package keycoordination

import (
	"fmt"
	"log"

	"analyzer/pkg/app"
)

func (detector *KeyCoordinationDetector) GetResults() string {
	return detector.results
}

func (detector *KeyCoordinationDetector) ComputeResults(app *app.App) {
	header := "---------------------------------------------------------------------\n"
	if detector.isTypePrimaryKey() {
		header += "---------------- ANALYSIS: PRIMARY KEY COORDINATION -----------------\n"
	} else if detector.isTypeForeignKey() {
		header += "---------------- ANALYSIS: FOREIGN KEY COORDINATION -----------------\n"
	}
	header += "---------------------------------------------------------------------\n"

	var results string
	total := 0
	for request, foreignreads := range detector.foreignReads {
		results += fmt.Sprintf("entry request: %s()\n", request.entry.String())
		for _, foreignread := range foreignreads {
			total++
			if detector.isTypePrimaryKey() {
				results += fmt.Sprintf("\tPRIMARY KEY READS #%d:\n", total)
			} else if detector.isTypeForeignKey() {
				results += fmt.Sprintf("\tFOREIGN KEY READS #%d:\n", total)
			} else {
				log.Fatalf("unexpected")
			}
			results += fmt.Sprintf("\t\tREAD 1: %s\n", foreignread.op1.call.String())
			var shift = 0
			if detector.isTypePrimaryKey() {
				shift = 1
			}
			if foreignread.field1 != nil {
				results += fmt.Sprintf("\t\t\t- read field: %s\n", foreignread.field1.GetPath())
				if detector.isTypePrimaryKey() {
					results += fmt.Sprintf("\t\t\t\t- constraint #0: PRIMARY KEY %s\n", foreignread.field1.GetName())
				}
				for i, constraint := range foreignread.field1.GetConstraintForeignKey() {
					results += fmt.Sprintf("\t\t\t\t- constraint #%d: %s\n", i+shift, constraint.String())
				}
			}
			results += fmt.Sprintf("\t\tREAD 2: %s\n", foreignread.op2.call.String())
			results += fmt.Sprintf("\t\t\t- read field: %s\n", foreignread.field2.GetPath())
			if detector.isTypePrimaryKey() {
				results += fmt.Sprintf("\t\t\t\t- constraint #0: PRIMARY KEY %s\n", foreignread.field2.GetName())
			}
			for i, constraint := range foreignread.field2.GetConstraintForeignKey() {
				results += fmt.Sprintf("\t\t\t\t- constraint #%d: %s\n", i+shift, constraint.String())
			}
		}
	}
	detector.results = header + results
}
