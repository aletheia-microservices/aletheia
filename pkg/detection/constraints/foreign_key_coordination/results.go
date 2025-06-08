package foreign_key_coordination

import (
	"fmt"

	"analyzer/pkg/utils"
)

func (detector *ForeignKeyDetector) GetResults() string {
	return detector.results
}

func (detector *ForeignKeyDetector) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "--------------------- FOREIGN KEY READ ANALYSIS ---------------------\n"
	header += "---------------------------------------------------------------------\n"

	for i, read := range detector.reads {
		detector.results += fmt.Sprintf("[%d] foreign key read:\n%s\n", i+1, read.String())
		if i < len(detector.reads)-1 {
			detector.results += "\n" // enforce empty line between each foreign key read result
		}
	}

	header += fmt.Sprintf(">> (# READS USING FOREIGN REFERENCES):\n>> (%d)\n", len(detector.reads))
	detector.results = header + "---------------------------------------------------------------------\n" + utils.TEXT_RESET_COLOR + detector.results
}
