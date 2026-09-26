package foreignkeyconcurrency

import (
	"fmt"
	"sort"
	"strings"

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
		sortedDangerousDeleteLst := detector.dangerousDeletes[request]
		sort.Slice(sortedDangerousDeleteLst, func(i, j int) bool {
			return sortedDangerousDeleteLst[i].CallString() < sortedDangerousDeleteLst[j].CallString()
		})

		for _, dangerousDelete := range sortedDangerousDeleteLst {
			if request.entry.String() != dangerousDelete.delete.call.GetFromNode().String() {
				results += fmt.Sprintf("delete: %s() ... %s\n", request.entry.String(), dangerousDelete.CallString())
			} else {
				results += fmt.Sprintf("delete: %s\n", dangerousDelete.CallString())
			}

			var sortedConcurrentWrites []*ConcurrentWrite = dangerousDelete.concurrentWrites
			sort.Slice(sortedConcurrentWrites, func(i, j int) bool {
				if sortedConcurrentWrites[i].CallString() != sortedConcurrentWrites[j].CallString() {
					return sortedConcurrentWrites[i].CallString() < sortedConcurrentWrites[j].CallString()
				}
				return sortedConcurrentWrites[i].EntryString() < sortedConcurrentWrites[j].EntryString()
			})

			for _, concurrentWrite := range dangerousDelete.concurrentWrites {
				numWarnings++
				results += fmt.Sprintf("\twrite #%d: %s() ... %s\n", numWarnings, concurrentWrite.EntryString(), concurrentWrite.CallString())
				var orderedFieldNames []string
				var seenOrderedFieldName = make(map[string]bool)

				for _, field := range concurrentWrite.affectedFields {
					if _, exists := seenOrderedFieldName[field.GetName()]; !exists {
						orderedFieldNames = append(orderedFieldNames, field.GetName())
						seenOrderedFieldName[field.GetName()] = true
					}
				}
				sort.Strings(orderedFieldNames)
				results += fmt.Sprintf("\t\t- database={%s}, entity={%s}, written_fields={%s}\n",
					concurrentWrite.database.GetName(),
					concurrentWrite.schema.GetName(),
					strings.Join(orderedFieldNames, ", "),
				)
			}
		}
		results += "\n"
	}
	detector.results = header + fmt.Sprintf("[NUM_WARNINGS = %d]\n", numWarnings) + results
}
