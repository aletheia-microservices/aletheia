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
		var printedEntry = false
		for _, fread := range foreignreads {
			//detector.restrictive = false

			// this information is only accurate after the entire schema is built at the end of the iteration
			detector.updateForeignReadConstraints(fread)
			if !detector.isValidForeignRead(fread) {
				continue
			}

			if !printedEntry {
				results += fmt.Sprintf("entry request: %s()\n", request.entry.String())
				printedEntry = true
			}

			total++
			if detector.isTypePrimaryKey() {
				results += fmt.Sprintf("\tPRIMARY KEY READS #%d:\n", total)
			} else if detector.isTypeForeignKey() {
				results += fmt.Sprintf("\tFOREIGN KEY READS #%d:\n", total)
			} else {
				log.Fatalf("unexpected")
			}
			results += fmt.Sprintf("\t\tREAD 1: %s\n", fread.op1.call.String())
			results += fmt.Sprintf("\t\t\t- read field: %s\n", fread.field1.GetPath())
			if fread.constraint1 != nil {
				results += fmt.Sprintf("\t\t\t\t- constraint: %s\n", fread.constraint1.String())
			}
			results += fmt.Sprintf("\t\tREAD 2: %s\n", fread.op2.call.String())
			results += fmt.Sprintf("\t\t\t- read field: %s\n", fread.field2.GetPath())
			if fread.constraint2 != nil {
				results += fmt.Sprintf("\t\t\t\t- constraint: %s\n", fread.constraint2.String())
			}
		}
	}
	detector.results = header + results
}
