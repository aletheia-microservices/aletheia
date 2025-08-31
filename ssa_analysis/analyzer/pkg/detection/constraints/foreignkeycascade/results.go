package foreignkeycascade

import (
	"fmt"
	"strings"

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
	var numWarnings int
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

			var dbToPendingField = make(map[string][]string)
			for _, pendingField := range cascadeDelete.pendingFields {
				dbname := pendingField.GetDatabase().GetName()
				fieldname := pendingField.GetName()
				dbToPendingField[dbname] = append(dbToPendingField[dbname], fieldname)
			}

			for db, fieldsLst := range dbToPendingField {
				numWarnings++
				results += fmt.Sprintf("\t\tMISSING CASCADE DELETE #%d: database={%s}, pending fields={%s}\n", numWarnings, db, strings.Join(fieldsLst, ", "))
			}
		}
		results += "\n"
	}
	detector.summary = header + results
}
