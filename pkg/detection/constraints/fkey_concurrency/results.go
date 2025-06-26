package fkey_concurrency

import (
	"fmt"

	"analyzer/pkg/utils"
)

func (detector *ForeignKeyConcurrencyDetector) GetResults() string {
	return detector.results
}

func (detector *ForeignKeyConcurrencyDetector) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "----------------- FOREIGN KEY CONCURRENCY ANALYSIS ------------------\n"
	header += "---------------------------------------------------------------------\n"

	detector.loadInconsistencies()
	header += fmt.Sprintf(">> (# DELETES; # WRITTEN CONSTRAINTS AFFECTED BY DELETES):\n>> (%d;%d)\n", detector.numDeletes, detector.numAffectedWrittenFields)
	detector.results = header + "---------------------------------------------------------------------\n" + utils.TEXT_RESET_COLOR + detector.results
}

func (detector *ForeignKeyConcurrencyDetector) loadInconsistencies() {
	for _, dels := range detector.deletes {
		for _, del := range dels {
			if len(del.affectedWrittenFields) > 0 {
				detector.numDeletes++
				detector.results += fmt.Sprintf("[%d] delete affecting %d written fields:\n", detector.numDeletes, len(del.affectedWrittenFields))
				detector.results += fmt.Sprintf("%s: %s\n", del.call.GetCallerStr(), del.call.ShortString())
				for _, writtenField := range del.affectedWrittenFields {
					detector.results += writtenField.String() + "\n"
					detector.numAffectedWrittenFields++
				}
				detector.results += "\n"
			}
		}
	}
}
