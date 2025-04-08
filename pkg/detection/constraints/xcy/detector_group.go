package xcy

import (
	"fmt"
	"slices"

	"gopkg.in/yaml.v2"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection/detection"
)

type DetectorGroup struct {
	detection.Detector
	detectors []*XCYDetector
	results   string
	summary   string
}

func (detector *DetectorGroup) GetSummary() string {
	return detector.summary
}

func (detector *DetectorGroup) SetSummary(summary string) {
	detector.summary = summary
}

func NewDetectorGroup(entryNodes []abstractgraph.AbstractNode) *DetectorGroup {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" ----------------------------------------- INITIALIZING XCY DETECTOR GROUP ---------------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	dg := &DetectorGroup{}
	for _, entryNode := range entryNodes {
		for _, mode := range GetActiveDetectionModes() {
			detector := NewDetector(entryNode, mode)
			dg.detectors = append(dg.detectors, detector)
		}
	}
	return dg
}

func (dg *DetectorGroup) GetAllDetectors() []*XCYDetector {
	return dg.detectors
}

func (dg *DetectorGroup) OnNewRun(app *app.App)                                     { /* no-op */ }
func (dg *DetectorGroup) OnEndRun(app *app.App)                                     { /* no-op */ }
func (dg *DetectorGroup) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) { /* no-op */ }
func (dg *DetectorGroup) OnEndRequest(app *app.App)                                 { /* no-op */ }
func (dg *DetectorGroup) OnNewNode(app *app.App, node abstractgraph.AbstractNode)   { /* no-op */ }
func (dg *DetectorGroup) OnEndNode(app *app.App, node abstractgraph.AbstractNode)   { /* no-op */ }
func (dg *DetectorGroup) OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) { /* no-op */
}
func (dg *DetectorGroup) OnWrite(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) { /* no-op */
}
func (dg *DetectorGroup) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) { /* no-op */
}

func (dg *DetectorGroup) ComputeResults() {
	header := "---------------------------------------------------------------------\n"
	header += "---------------------------- XCY ANALYSIS ---------------------------\n"
	header += "---------------------------------------------------------------------\n"

	var numInconsistencies int
	for _, detector := range dg.detectors {
		if detector.HasInconsistencies() {
			for _, request := range detector.requests {
				data, _ := yaml.Marshal(detector.dumpRequestYaml(request, false))
				dg.results += string(data)
				dg.results += "----------------------------------------------------------\n"
			}
			numInconsistencies += detector.inconsistencies
		}
	}
	header += fmt.Sprintf(">> (# CROSS-SERVICE INCONSISTENCIES):\n>> (%d)\n", numInconsistencies)
	dg.results = header + "---------------------------------------------------------------------\n" + dg.results
}

func (dg *DetectorGroup) GetAnalysisTypeString() string {
	return "xcy"
}

func (dg *DetectorGroup) GetResults() string {
	return dg.results
}

func (detector *XCYDetector) dumpOperationYaml(operation *Operation) map[string]interface{} {
	data := make(map[string]interface{})
	data["_operation"] = operation.String()
	data["visible_dependency_set"] = operation.GetVisibleDependenciesString()

	if detector.HasDetectionMode(DEBUG_XCY_MINIMIZE_DEPENDENCIES) && operation.Write {
		lenVisible := len(operation.GetVisibleDependencies())
		lenMinimized := len(operation.GetMinimizedDependencySet())
		data[fmt.Sprintf("minimized_dependency_set (%d --> %d)", lenVisible, lenMinimized)] = operation.GetMinimizedDependencySetString()
	}
	return data
}

func (detector *XCYDetector) dumpInconsistencyYaml(inconsistency *XCYInconsistency) map[string]interface{} {
	data := make(map[string]interface{})
	data["write"] = inconsistency.Write.String()
	data["read"] = inconsistency.Read.String()
	if detector.HasDetectionMode(DEBUG_XCY_MISSING_DEPENDENCIES) { // include missing dependencies
		data["missing_dependency"] = inconsistency.MissingDependency
		data["visible_dependency_set"] = inconsistency.Read.GetVisibleDependenciesString()
	}
	return data
}

func (detector *XCYDetector) dumpLineageYaml(request *Request, lineage *Lineage) map[string]interface{} {
	data := make(map[string]interface{})
	data["_id"] = lineage.ID
	var dataOperations []interface{}
	for _, op := range lineage.GetOperations() {
		if detector.HasDetectionMode(DEBUG_LINEAGES) || detector.HasDetectionMode(DEBUG_XCY_MINIMIZE_DEPENDENCIES) { // include visible dependencies
			dataOperations = append(dataOperations, detector.dumpOperationYaml(op))
		} else {
			dataOperations = append(dataOperations, op.String())
		}
	}
	data["_operations"] = dataOperations

	if detector.HasDetectionMode(DEBUG_LINEAGES) {
		var dataDependencies []string
		lineageDependencies := lineage.GetXCYDependenciesByMostRecent()
		for _, op := range lineageDependencies {
			dataDependencies = append(dataDependencies, op.String())
		}
		data["lineage_dependencies"] = dataDependencies

		var dataMissingDependencies []string
		allReqOps := request.Operations
		if len(lineage.Operations) > 0 {
			firstLineageOp := lineage.Operations[0]
			for _, op := range allReqOps {
				if op.LineageID < firstLineageOp.LineageID && !slices.Contains(lineageDependencies, op) {
					dataMissingDependencies = append(dataMissingDependencies, op.String())
				}
			}
		}
		data["missing_dependencies"] = dataMissingDependencies
	}

	return data
}

func (detector *XCYDetector) dumpRequestYaml(request *Request, includeLineages bool) map[string]interface{} {
	data := make(map[string]interface{})
	data["entry"] = detector.entryNode.ShortString()
	data["mode"] = DetectionModeName(detector)
	data["number_inconsistencies"] = len(request.Inconsistencies)

	var dataInconsistencies []map[string]interface{}
	for _, inconsistency := range request.Inconsistencies {
		result := detector.dumpInconsistencyYaml(inconsistency)
		//result["number"] = i+1
		dataInconsistencies = append(dataInconsistencies, result)
	}
	data["xcy_inconsistencies"] = dataInconsistencies

	if includeLineages {
		var dataLineages []map[string]interface{}
		for _, lineage := range request.Lineages {
			dataLineages = append(dataLineages, detector.dumpLineageYaml(request, lineage))
		}
		data["xcy_lineages"] = dataLineages
	}
	return data
}
