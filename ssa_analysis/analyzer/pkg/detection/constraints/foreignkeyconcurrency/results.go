package foreignkeyconcurrency

import (
	"fmt"

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
	for request, dangerousDeleteLst := range detector.dangerousDeletes {
		results += fmt.Sprintf("entry request: %s()\n", request.entry.String())
		for _, dangerousDelete := range dangerousDeleteLst {
			results += fmt.Sprintf("\tDELETE: %s\n", dangerousDelete.delete.call.String())
			/* for _, field := range dangerousDelete.delete.schema.GetAllFieldsLst() {
				results += fmt.Sprintf("\t- deleted field: %s\n", field.GetPath())
			} */
			for _, write := range dangerousDelete.concurrentWrites {
				results += fmt.Sprintf("\t\tCONCURRENT WRITE: %s\n", write.write.call.String())
				for _, field := range write.affectedFields {
					results += fmt.Sprintf("\t\t- written field: %s\n", field.GetPath())
					for i, constraint := range field.GetConstraints() {
						results += fmt.Sprintf("\t\t\t- affected constraint #%d: %s\n", i, constraint.String())
					}
				}
			}
		}
	}
	detector.results = header + results
}
