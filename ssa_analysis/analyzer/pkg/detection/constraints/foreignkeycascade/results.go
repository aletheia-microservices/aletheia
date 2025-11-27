package foreignkeycascade

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
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
		for _, cascadeDelete := range cascadeDeletes {
			if len(cascadeDelete.pendingFields) == 0 {
				continue
			}
			if request.entry.String() != cascadeDelete.op.call.GetFromNode().String() {
				results += fmt.Sprintf("delete: %s() ... %s\n", request.entry.String(), cascadeDelete.op.call.String())				
			} else {
				results += fmt.Sprintf("delete: %s\n", cascadeDelete.op.call.String())
			}

			var schemaToPendingField = make(map[*backends.Schema][]string)
			var sortedSchemas []*backends.Schema
			for _, pendingField := range cascadeDelete.pendingFields {
				schema := pendingField.GetSchema()
				fieldname := pendingField.GetName()
				schemaToPendingField[schema] = append(schemaToPendingField[schema], fieldname)
				if !slices.Contains(sortedSchemas, schema) {
					sortedSchemas = append(sortedSchemas, schema)
				}
			}

			sort.Slice(sortedSchemas, func(i, j int) bool {
				return sortedSchemas[i].GetName() < sortedSchemas[j].GetName()
			})

			for _, schema := range sortedSchemas {
				database := schema.GetDatabase()
				fieldsLst := schemaToPendingField[schema]
				sort.Strings(fieldsLst)
				numWarnings++
				results += fmt.Sprintf("\tmissing cascade #%d: database={%s}, entity={%s}, pending_fields={%s}\n", numWarnings, database.GetName(), schema.GetName(), strings.Join(fieldsLst, ", "))
			}
		}
		results += "\n"
	}
	detector.results = header + fmt.Sprintf("[NUM_WARNINGS = %d]\n", numWarnings) + results
}
