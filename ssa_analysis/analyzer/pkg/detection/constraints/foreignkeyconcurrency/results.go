package foreignkeyconcurrency

import (
	"fmt"
	"sort"

	"analyzer/pkg/app"
)

func (detector *ForeignKeyConcurrencyDetector) GetResults() string {
	return detector.results
}

func (detector *ForeignKeyConcurrencyDetector) ComputeResults(app *app.App) {
	header := "---------------------------------------------------------------------\n"
	header += "----------------- ANALYSIS: FOREIGN KEY CONCURRENCY -----------------\n"
	header += "---------------------------------------------------------------------\n"

	var results string
	var numWarnings int

	var sortedRequests []*Request
	for request := range detector.dangerousDeletes {
		sortedRequests = append(sortedRequests, request)
	}
	sort.Slice(sortedRequests, func(i, j int) bool {
		return sortedRequests[i].entry.String() < sortedRequests[j].entry.String()
	})

	for _, request := range sortedRequests {
		dangerousDeleteLst := detector.dangerousDeletes[request]
		results += fmt.Sprintf("entry request: %s()\n", request.entry.String())
		for _, dangerousDelete := range dangerousDeleteLst {
			results += fmt.Sprintf("\tDELETE: %s\n", dangerousDelete.delete.call.String())
			/* for _, field := range dangerousDelete.delete.schema.GetAllFieldsLst() {
				results += fmt.Sprintf("\t- deleted field: %s\n", field.GetPath())
			} */
			for _, concWrite := range dangerousDelete.concurrentWrites {
				numWarnings++
				results += fmt.Sprintf("\t\tCONCURRENT WRITE #%d: %s\n", numWarnings, concWrite.write.call.String())
				var fieldspaths string
				for i, field := range concWrite.affectedFields {
					fieldspaths += field.GetPath()
					if i < len(concWrite.affectedFields)-1 {
						fieldspaths += ", "
					}
				}
				fieldspaths = "{" + fieldspaths + "}"
				results += fmt.Sprintf("\t\t- written fields: %s\n", fieldspaths)
				results += fmt.Sprintf("\t\t- entry request: %s\n", concWrite.write.request.entry.String())

				/* for _, field := range write.affectedFields {
					results += fmt.Sprintf("\t\t- written fields: %s\n", strings.field.GetPath())
					for i, constraint := range field.GetConstraints() {
						if !constraint.IsMandatory() {
							results += fmt.Sprintf("\t\t\t- affected constraint #%d: %s\n", i, constraint.String())
						}
						//results += fmt.Sprintf("\t\t\t- affected constraint #%d: %s\n", i, constraint.String())
					}
				} */
			}
		}
	}
	detector.results = header + results
}
