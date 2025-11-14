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

		results += fmt.Sprintf("entry request: %s()\n", request.entry.String())
		for _, dangerousDelete := range sortedDangerousDeleteLst {
			results += fmt.Sprintf("\tDELETE: %s\n", dangerousDelete.CallString())
			var sortedConcurrentWrites []*ConcurrentWrite = dangerousDelete.concurrentWrites
			sort.Slice(sortedConcurrentWrites, func(i, j int) bool {
				if sortedConcurrentWrites[i].CallString() != sortedConcurrentWrites[j].CallString() {
					return sortedConcurrentWrites[i].CallString() < sortedConcurrentWrites[j].CallString()
				}
				return sortedConcurrentWrites[i].EntryString() < sortedConcurrentWrites[j].EntryString()
			})

			for _, concurrentWrite := range dangerousDelete.concurrentWrites {
				numWarnings++
				results += fmt.Sprintf("\t\tCONCURRENT WRITE #%d: %s\n", numWarnings, concurrentWrite.CallString())
				var orderedFieldNames []string

				for _, field := range concurrentWrite.affectedFields {
					orderedFieldNames = append(orderedFieldNames, field.GetName())
				}
				sort.Strings(orderedFieldNames)
				results += fmt.Sprintf("\t\t- entry={%s}, database={%s}, written fields={%s}\n", 
					concurrentWrite.EntryString(), 
					concurrentWrite.database.GetName(), 
					strings.Join(orderedFieldNames, ", "),
				)
			}
		}
		results += "\n"
	}
	detector.results = header + fmt.Sprintf("[NUM_WARNINGS = %d]\n", numWarnings) + results
}
