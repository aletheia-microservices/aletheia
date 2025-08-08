package foreignkeycascade

import "fmt"

func (detector *ForeignKeyCascadeDetector) GetResults() string {
	return detector.summary
}

func (detector *ForeignKeyCascadeDetector) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "------------------- ANALYSIS: FOREIGN KEY CASCADE -------------------\n"
	header += "---------------------------------------------------------------------\n"

	var results string
	for request, cascadeDeletes := range detector.cascadeDeletes {
		results += fmt.Sprintf("entry request: %s\n", request.entry.String())
		for i, cascadeDelete := range cascadeDeletes {
			results += fmt.Sprintf("\tdelete #%d: %s\n", i, cascadeDelete.op.call.String())
			/* for _, op := range cascadeDelete.cascadingOps {
				results += fmt.Sprintf("\t\tcascading delete: %s\n", op.call.String())
			} */
			for _, pendingField := range cascadeDelete.pendingDBFields {
				results += fmt.Sprintf("\t\tMISSING delete on field: %s\n", pendingField.GetPath())
			}
		}
		results += "\n"
	}
	detector.summary = header + results
}
