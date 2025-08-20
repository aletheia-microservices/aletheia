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
		for _, cascadeDelete := range cascadeDeletes {
			results += fmt.Sprintf("\tDELETE: %s\n", cascadeDelete.op.call.String())
			/* for _, op := range cascadeDelete.cascadingOps {
				results += fmt.Sprintf("\t\tcascading delete: %s\n", op.call.String())
			} */
			for _, pendingField := range cascadeDelete.pendingFields {
				results += fmt.Sprintf("\t\tMISSING DELETE on field: %s\n", pendingField.GetPath())
				for i, constraint := range pendingField.GetConstraints() {
					results += fmt.Sprintf("\t\t\t- affected constraint #%d: %s\n", i, constraint.String())
				}
			}
		}
		results += "\n"
	}
	detector.summary = header + results
}
