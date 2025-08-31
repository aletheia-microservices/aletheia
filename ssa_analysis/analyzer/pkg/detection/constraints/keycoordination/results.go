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
	numWarnings := 0
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

			numWarnings++
			if detector.isTypePrimaryKey() {
				results += fmt.Sprintf("\tPRIMARY KEY READS #%d:\n", numWarnings)
			} else if detector.isTypeForeignKey() {
				results += fmt.Sprintf("\tFOREIGN KEY READS #%d:\n", numWarnings)
			} else {
				log.Fatalf("unexpected")
			}

			if fread.constraint1 == nil {
				results += fmt.Sprintf("\t\tREAD (ORIGIN): %s\n", fread.op1.call.String())
				results += fmt.Sprintf("\t\t\t- field: %s\n", fread.field1.GetPath())
			}
			
			if fread.constraint2 == nil {
				results += fmt.Sprintf("\t\tREAD (ORIGIN): %s\n", fread.op2.call.String())
				results += fmt.Sprintf("\t\t\t- field: %s\n", fread.field2.GetPath())
			}

			if fread.constraint1 != nil {
				if detector.isTypePrimaryKey() {
					results += fmt.Sprintf("\t\tREAD: %s\n", fread.op1.call.String())
				} else if detector.isTypeForeignKey() {
					results += fmt.Sprintf("\t\tREAD (FOREIGN KEY): %s\n", fread.op1.call.String())
				}
				results += fmt.Sprintf("\t\t\t- field: %s\n", fread.field1.GetPath())
				results += fmt.Sprintf("\t\t\t- constraint: %s\n", fread.constraint1.String())
			}

			if fread.constraint2 != nil {
				if detector.isTypePrimaryKey() {
					results += fmt.Sprintf("\t\tREAD: %s\n", fread.op2.call.String())
				} else if detector.isTypeForeignKey() {
					results += fmt.Sprintf("\t\tREAD (FOREIGN KEY): %s\n", fread.op2.call.String())
				}
				results += fmt.Sprintf("\t\t\t- field: %s\n", fread.field2.GetPath())
				results += fmt.Sprintf("\t\t\t- constraint: %s\n", fread.constraint2.String())
			}
		}
	}
	detector.results = header + results
}
