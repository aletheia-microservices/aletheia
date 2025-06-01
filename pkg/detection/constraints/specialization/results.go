package specialization

import "fmt"

func (detector *SpecializationDetector) GetResults() string {
	return detector.results
}

func (detector *SpecializationDetector) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "----------------------- SPECIALIZATION ANALYSIS ---------------------\n"
	header += "---------------------------------------------------------------------\n"
	var numRemovedMandatoryFields int
	if len(detector.rmes) > 0 {
		detector.results += "removed mandatory entities:\n"
	}
	for i, rme := range detector.rmes {
		detector.results += fmt.Sprintf("- (#%d) %s", i, rme.String())
		for _, mandatoryField := range rme.mandatoryFields { // AT THE MOMENT MANDATORY FIELDS IS ALWAYS NIL SO WE NEVER PRINT THIS
			detector.results += fmt.Sprintf("\t\t %s REFERENCES %s * {MANDATORY}", mandatoryField.field.GetFullName(), mandatoryField.mandatoryRef.GetFullName())
			numRemovedMandatoryFields++
		}
		if i < len(detector.rmes)-1 {
			detector.results += "\n" // enforce empty line between each foreign key read result
		}
	}

	header += fmt.Sprintf(">> (# REMOVED MANDATORY OBJECTS; # REFERENCES OF OBJECTS):\n>> (%d;%d)\n", len(detector.rmes), numRemovedMandatoryFields)
	detector.results = header + "---------------------------------------------------------------------\n" + detector.results
}
