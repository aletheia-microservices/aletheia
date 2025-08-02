package foreignkeycoordination

import (
	"fmt"

	//"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/detection/detector"
)

type ForeignKeyCoordinationDetector struct {
	detector.Detector
	//requestInfoStack *stack.Stack

	keyType string // 'primary_key' or 'foreign_key'

	// results
	results string
	summary string
	//reads   []*ForeignKeyRead
}

func NewDetector(keyType string) *ForeignKeyCoordinationDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" --------------------------------------- INITIALIZING KEY_COORD DETECTOR ---------------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	return &ForeignKeyCoordinationDetector{
		keyType:          keyType,
		//requestInfoStack: stack.New(),
	}
}

func (detector *ForeignKeyCoordinationDetector) GetResults() string {
	// TODO
	return ""
}

func (detector *ForeignKeyCoordinationDetector) GetSummary() string {
	// TODO
	return ""
}

func (detector *ForeignKeyCoordinationDetector) SetSummary(summary string) {
	// TODO
}

func (detector *ForeignKeyCoordinationDetector) ComputeResults() {
	// TODO
}

func (detector *ForeignKeyCoordinationDetector) GetAnalysisTypeString() string {
	// TODO
	return ""
}

func (detector *ForeignKeyCoordinationDetector) OnNewRun(app *app.App) {
	// TODO
}

func (detector *ForeignKeyCoordinationDetector) OnEndRun(app *app.App) {
	// TODO
}

func (detector *ForeignKeyCoordinationDetector) OnNewRequest(node *abstractgraph.AbstractNode) {
	// TODO
}

func (detector *ForeignKeyCoordinationDetector) OnEndRequest(app *app.App) {
	// TODO
}

func (detector *ForeignKeyCoordinationDetector) OnNewNode(app *app.App, node *abstractgraph.AbstractNode) {
	// TODO
}

func (detector *ForeignKeyCoordinationDetector) OnEndNode(app *app.App, node *abstractgraph.AbstractNode) {
	// TODO
}

func (detector *ForeignKeyCoordinationDetector) OnRead(app *app.App, edge *abstractgraph.AbstractEdge) {
	fmt.Printf("[DETECTOR] [READ] operation for arguments:\n")
	for i, arg := range edge.GetArguments() {
		fmt.Printf("\t[arg %d] %s:\n", i, arg.String())
		for objpath, taintsLst := range arg.GetAllTaints() {
			for _, taint := range taintsLst {
				fmt.Printf("%s @ %s\n", objpath, taint.LongString())
			}
		}
	}
	fmt.Println()
	for i, param := range edge.GetToNode().GetParams() {
		fmt.Printf("\t[param %d] %s:\n", i, param.String())
		for objpath, taintsLst := range param.GetAllTaints() {
			for _, taint := range taintsLst {
				fmt.Printf("%s @ %s\n", objpath, taint.LongString())
			}
		}
	}
	fmt.Println()
}

func (detector *ForeignKeyCoordinationDetector) OnWrite(app *app.App, edge *abstractgraph.AbstractEdge) {
	fmt.Printf("[DETECTOR] [WRITE] operation for arguments:\n")
	for i, arg := range edge.GetArguments() {
		fmt.Printf("\t[arg %d] %s:\n", i, arg.String())
		for objpath, taintsLst := range arg.GetAllTaints() {
			for _, taint := range taintsLst {
				fmt.Printf("%s @ %s\n", objpath, taint.LongString())
			}
		}
	}
	fmt.Println()
	for i, param := range edge.GetToNode().GetParams() {
		fmt.Printf("\t[param %d] %s:\n", i, param.String())
		for objpath, taintsLst := range param.GetAllTaints() {
			for _, taint := range taintsLst {
				fmt.Printf("%s @ %s\n", objpath, taint.LongString())
			}
		}
	}
	fmt.Println()
}

func (detector *ForeignKeyCoordinationDetector) OnUpdate(app *app.App, edge *abstractgraph.AbstractEdge) {
	// TODO
}

func (detector *ForeignKeyCoordinationDetector) OnDelete(app *app.App, edge *abstractgraph.AbstractEdge) {
	// TODO
}
