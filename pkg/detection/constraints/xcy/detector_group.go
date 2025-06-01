package xcy

import (
	"fmt"

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

func (dg *DetectorGroup) GetAnalysisTypeString() string {
	return "xcy"
}
