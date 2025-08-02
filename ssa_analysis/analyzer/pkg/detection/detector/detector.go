package detector

import (
	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
)

type Detector interface {
	GetResults() string
	GetSummary() string
	SetSummary(summary string)
	ComputeResults()
	GetAnalysisTypeString() string

	OnNewRun(app *app.App)
	OnEndRun(app *app.App)
	OnNewRequest(node *abstractgraph.AbstractNode)
	OnEndRequest(app *app.App)
	OnNewNode(app *app.App, node *abstractgraph.AbstractNode)
	OnEndNode(app *app.App, node *abstractgraph.AbstractNode)

	// database calls
	OnRead(app *app.App, edge *abstractgraph.AbstractEdge)
	OnWrite(app *app.App, edge *abstractgraph.AbstractEdge)
	OnUpdate(app *app.App, edge *abstractgraph.AbstractEdge)
	OnDelete(app *app.App, edge *abstractgraph.AbstractEdge)
}
