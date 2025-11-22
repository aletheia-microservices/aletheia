package keycoordination

import (
	"fmt"
	"sort"

	"analyzer/pkg/app"

	"github.com/sirupsen/logrus"
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

	var sortedRequests []*Request
	for request := range detector.foreignReads {
		sortedRequests = append(sortedRequests, request)
	}
	sort.Slice(sortedRequests, func(i, j int) bool {
		return sortedRequests[i].entry.String() < sortedRequests[j].entry.String()
	})

	for _, request := range sortedRequests {
		sortedForeignReads := detector.foreignReads[request]
		sort.Slice(sortedForeignReads, func(i, j int) bool {
			if sortedForeignReads[i].FirstCallString() != sortedForeignReads[j].FirstCallString() {
				return sortedForeignReads[i].FirstCallString() < sortedForeignReads[j].FirstCallString()
			}
			return sortedForeignReads[i].SecondCallString() < sortedForeignReads[j].SecondCallString()
		})

		var printedEntry = false
		for _, fread := range sortedForeignReads {
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
				logrus.Fatalf("unexpected")
			}

			if fread.constraint1 != nil {
				if detector.isTypePrimaryKey() {
					results += fmt.Sprintf("\t\tREAD: %s\n", fread.FirstCallString())
				} else if detector.isTypeForeignKey() {
					results += fmt.Sprintf("\t\tREAD (FOREIGN KEY): %s\n", fread.FirstCallString())
				}
				results += fmt.Sprintf("\t\t\t- field: %s\n", fread.field1.GetPath())
				results += fmt.Sprintf("\t\t\t- constraint: %s\n", fread.constraint1.String())
			}

			if fread.constraint2 != nil {
				if detector.isTypePrimaryKey() {
					results += fmt.Sprintf("\t\tREAD: %s\n", fread.SecondCallString())
				} else if detector.isTypeForeignKey() {
					results += fmt.Sprintf("\t\tREAD (FOREIGN KEY): %s\n", fread.SecondCallString())
				}
				results += fmt.Sprintf("\t\t\t- field: %s\n", fread.field2.GetPath())
				results += fmt.Sprintf("\t\t\t- constraint: %s\n", fread.constraint2.String())
			}

			if fread.constraint1 == nil {
				results += fmt.Sprintf("\t\tREAD (ORIGIN): %s\n", fread.FirstCallString())
				results += fmt.Sprintf("\t\t\t- field: %s\n", fread.field1.GetPath())
			}

			if fread.constraint2 == nil {
				results += fmt.Sprintf("\t\tREAD (ORIGIN): %s\n", fread.SecondCallString())
				results += fmt.Sprintf("\t\t\t- field: %s\n", fread.field2.GetPath())
			}
		}
	}
	detector.results = header + fmt.Sprintf("[NUM_WARNINGS = %d]\n", numWarnings) + results
}
