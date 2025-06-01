package specialization

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

type SpecializationDetector struct {
	results string
	summary string
	rmes    []*RemovedMandatoryEntity
}

func (detector *SpecializationDetector) addRemovedMandatoryEntity(rme *RemovedMandatoryEntity) {
	detector.rmes = append(detector.rmes, rme)
}

func (detector *SpecializationDetector) GetSummary() string {
	return detector.summary
}

func (detector *SpecializationDetector) SetSummary(summary string) {
	detector.summary = summary
}

func NewDetector() *SpecializationDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" -------------------------------------- INITIALIZING SPECIALIZATION DETECTOR -------------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()
	return &SpecializationDetector{}
}

func (detector *SpecializationDetector) OnNewRun(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *SpecializationDetector) OnEndRun(app *app.App) {
	//no-op
}

func (detector *SpecializationDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	//no-op
}

func (detector *SpecializationDetector) OnEndRequest(app *app.App) {
	app.ResetAllDataflows()
}

func (detector *SpecializationDetector) OnNewNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *SpecializationDetector) OnEndNode(app *app.App, node abstractgraph.AbstractNode) {
	//no-op
}

func (detector *SpecializationDetector) OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *SpecializationDetector) OnWrite(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op

	/* datastore := dbCall.DbInstance.GetDatastore()
	params := dbCall.GetParams()
	switch datastore.Type {
	case datastores.Queue:
		msg := params[1]
		abstractgraph.TaintDataflowWrite(app, msg, dbCall, datastore, "", true, child_idx)

	case datastores.NoSQL:
		doc := params[1]
		abstractgraph.TaintDataflowWrite(app, doc, dbCall, datastore, "", true, child_idx)

	case datastores.Cache:
		key, value := params[1], params[2]
		abstractgraph.TaintDataflowWrite(app, key, dbCall, datastore, datastores.ROOT_FIELD_NAME_CACHE_KEY, false, child_idx)
		abstractgraph.TaintDataflowWrite(app, value, dbCall, datastore, datastores.ROOT_FIELD_NAME_CACHE_VALUE, false, child_idx)

	default:
		logger.Logger.Fatalf("[SPECIALIZATION > WRITE] TODO: %s", dbCall.String())
	} */
}

func (detector *SpecializationDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	//no-op
}

func (detector *SpecializationDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	datastore := dbCall.DbInstance.GetDatastore()
	switch datastore.Type {
	case datastores.NoSQL:
		/* doc := params[1]
		abstractgraph.TaintDataflowWrite(detector.app, doc, dbCall, datastore, "", true, requestIdx) */

		dbMandatoryConstraints := datastore.GetSchema().GetConstraints(datastores.ConstraintFilter{Mandatory: utils.BoolPtr(true)})
		if len(dbMandatoryConstraints) > 0 {
			detector.addRemovedMandatoryEntity(newRemovedMandatoryEntity(dbCall.ParsedCall, nil)) // FIXME: IN THE FUTURE, REPLACE NIL MANDATORY FIELDS
		}
	default:
		logger.Logger.Fatalf("[SPECIALIZATION > DELETE] TODO: %s", dbCall.String())
	}
}

func (detector *SpecializationDetector) GetAnalysisTypeString() string {
	return "specialization"
}
