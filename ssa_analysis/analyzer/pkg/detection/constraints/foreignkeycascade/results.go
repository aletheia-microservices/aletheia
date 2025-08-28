package foreignkeycascade

import (
	"fmt"

	"analyzer/pkg/app"
)

func (detector *ForeignKeyCascadeDetector) GetResults() string {
	return detector.summary
}

func (detector *ForeignKeyCascadeDetector) ComputeResults(app *app.App) {
	header := "---------------------------------------------------------------------\n"
	header += "------------------- ANALYSIS: FOREIGN KEY CASCADE -------------------\n"
	header += "---------------------------------------------------------------------\n"

	var results string
	for request, cascadeDeletes := range detector.cascadeDeletes {
		var found bool
		for _, cascadeDelete := range cascadeDeletes {
			if len(cascadeDelete.pendingFields) != 0 {
				found = true
			}
		}
		if !found {
			continue
		}
		results += fmt.Sprintf("entry request: %s()\n", request.entry.String())
		for _, cascadeDelete := range cascadeDeletes {
			if len(cascadeDelete.pendingFields) == 0 {
				continue
			}
			results += fmt.Sprintf("\tDELETE: %s\n", cascadeDelete.op.call.String())
			results += fmt.Sprintln("\t\tMISSING CASCADE DELETE on:")
			for _, pendingField := range cascadeDelete.pendingFields {
				results += fmt.Sprintf("\t\t- %s\n", pendingField.GetPath())
				/* for i, constraint := range pendingField.GetConstraints() {
					results += fmt.Sprintf("\t\t\t- affected constraint #%d: %s\n", i, constraint.String())
				} */
			}
		}
		results += "\n"
	}
	detector.summary = header + results
}
