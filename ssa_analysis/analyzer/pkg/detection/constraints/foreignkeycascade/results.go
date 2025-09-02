package foreignkeycascade

import (
	"fmt"
	"sort"
	"strings"

	"analyzer/pkg/app"
)

func (detector *ForeignKeyCascadeDetector) GetResults() string {
	return detector.results
}

func (detector *ForeignKeyCascadeDetector) ComputeResults(app *app.App) {
	header := "---------------------------------------------------------------------\n"
	header += "------------------- ANALYSIS: FOREIGN KEY CASCADE -------------------\n"
	header += "---------------------------------------------------------------------\n"

	var results string
	var numWarnings int

	var sortedRequests []*Request
	for request := range detector.cascadeDeletes {
		sortedRequests = append(sortedRequests, request)
	}
	sort.Slice(sortedRequests, func(i, j int) bool {
		return sortedRequests[i].entry.String() < sortedRequests[j].entry.String()
	})

	for _, request := range sortedRequests {
		cascadeDeletes := detector.cascadeDeletes[request]
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

			var sortedDbs []string
			for db := range dbToPendingField {
				sortedDbs = append(sortedDbs, db)
			}
			sort.Strings(sortedDbs)

			for _, db := range sortedDbs {
				fieldsLst := dbToPendingField[db]
				sort.Strings(fieldsLst)
				numWarnings++
				results += fmt.Sprintf("\t\tMISSING CASCADE DELETE #%d: database={%s}, pending fields={%s}\n", numWarnings, db, strings.Join(fieldsLst, ", "))
			}
		}
		results += "\n"
	}
	detector.results = header + fmt.Sprintf("[NUM_WARNINGS = %d]\n", numWarnings) + results
}
