package xcy

import (
	"fmt"
	"slices"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/datastores"
	"analyzer/pkg/detection/detector"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

type XCYDetector struct {
	detector.Detector

	currentRequest *Request
	requests       []*Request
	detectionMode  DetectionMode
	datastoresOps  map[*datastores.Datastore][]*Operation
	entryNode      abstractgraph.AbstractNode

	inconsistencies int
}

func NewDetector(entryNode abstractgraph.AbstractNode, mode DetectionMode) *XCYDetector {
	fmt.Println()
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println(" ------------------------------------------- INITIALIZING XCY DETECTOR -------------------------------------------- ")
	fmt.Println(" ------------------------------------------------------------------------------------------------------------------ ")
	fmt.Println()

	return &XCYDetector{
		datastoresOps: make(map[*datastores.Datastore][]*Operation),
		entryNode:     entryNode,
		detectionMode: mode,
	}
}

func (detector *XCYDetector) OnNewRun(app *app.App) {
	//no-op
}

func (detector *XCYDetector) OnEndRun(app *app.App) {
	//no-op
}

func (detector *XCYDetector) OnNewRequest(entryNode *abstractgraph.AbstractServiceCall) {
	//no-op
}

func (detector *XCYDetector) OnEndRequest(app *app.App) {
	//no-op
}

func (detector *XCYDetector) OnNewNode(app *app.App, node abstractgraph.AbstractNode) {
	if DetectionModeUsesLineages(detector) {
		if _, ok := node.(*abstractgraph.AbstractQueueHandler); ok {
			detector.currentRequest.PushLineage()
		}
	}
}

func (detector *XCYDetector) OnEndNode(app *app.App, node abstractgraph.AbstractNode) {
	// FIXME: maybe we should link the two stateful nodes (push and pop) in the abstract graph
	// instead of placing a stateless queue handler between both
	if DetectionModeUsesLineages(detector) {
		if _, ok := node.(*abstractgraph.AbstractQueueHandler); ok {
			detector.currentRequest.PopLineage()
		}
	}
}

func (detector *XCYDetector) OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	if backend, ok := dbCall.ParsedCall.Method.(*blueprint.BackendMethod); ok {
		operation := detector.currentRequest.SaveReadOperation(dbCall, backend)
		detector.attachOperationToDatastore(operation)
		logger.Logger.Infof("[XCY] saved read %s", operation.String())

		if !detector.HasDetectionMode(DEBUG_LINEAGES) {
			inconsistencies := detector.searchInconsistencies(detector.currentRequest, operation, dbCall)
			if detector.HasDetectionMode(DEBUG_XCY_MISSING_DEPENDENCIES) && len(inconsistencies) > 0 {
				detectMissingDependencies(detector.currentRequest, inconsistencies)
			}
		}
		detector.afterAnyOperation(operation)
	}
}

func (detector *XCYDetector) OnWrite(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	detector.onWriteOrUpdate(dbCall)
}

func (detector *XCYDetector) OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	detector.onWriteOrUpdate(dbCall)
}

func (detector *XCYDetector) onWriteOrUpdate(dbCall *abstractgraph.AbstractDatabaseCall) {
	if backend, ok := dbCall.ParsedCall.Method.(*blueprint.BackendMethod); ok {
		operation := detector.currentRequest.SaveWriteOperation(dbCall, backend)
		detector.attachOperationToDatastore(operation)
		logger.Logger.Infof("[XCY] saved write/update %s", operation.String())
		detector.afterAnyOperation(operation)
	}
}

func (detector *XCYDetector) afterAnyOperation(operation *Operation) {
	if detector.HasDetectionMode(DEBUG_XCY_MINIMIZE_DEPENDENCIES) {
		detector.attachOperationToDatastore(operation)
	}
}

func (detector *XCYDetector) OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int) {
	// no-op
	logger.Logger.Warnf("[XCY] ignoring operation: %s", dbCall.String())
}

func (detector *XCYDetector) HasInconsistencies() bool {
	return detector.inconsistencies > 0
}

func (detector *XCYDetector) GetActiveDetectionModeIndex() int {
	for i, activeMode := range GetActiveDetectionModes() {
		if detector.detectionMode == activeMode {
			return i
		}
	}
	logger.Logger.Fatalf("[XCY DETECTOR] could not find current detection mode (%s) in active mode list: %v", DetectionModeName(detector), GetActiveDetectionModes())
	return -1
}

func (detector *XCYDetector) GetDatastoreOps() map[*datastores.Datastore][]*Operation {
	return detector.datastoresOps
}

func (detector *XCYDetector) UpdateDatastoreOps(ops map[*datastores.Datastore][]*Operation) {
	detector.datastoresOps = ops
}

func (detector *XCYDetector) InitRequest(cumulativeDatastoreOps map[*datastores.Datastore][]*Operation) {
	if cumulativeDatastoreOps != nil {
		detector.datastoresOps = cumulativeDatastoreOps //FIXME: this is not copying and instead it is just using the pointer for the cumulativeDatastoreOps
	}
	request := &Request{
		EntryNode:     detector.entryNode,
		LineagesStack: stack.New(),
	}
	fmt.Printf("\n\n-------------------- ENTRY NODE = %s.%s --------------------\n", request.EntryNode.GetCallee(), request.EntryNode.GetName())
	fmt.Printf("-------------------- (%d/%d) XCY DETECTOR MODE = %s --------------------\n\n", detector.GetActiveDetectionModeIndex()+1, len(GetActiveDetectionModes()), DetectionModeName(detector))
	request.PushLineage()
	detector.requests = append(detector.requests, request)
	detector.currentRequest = request
}

func (detector *XCYDetector) searchInconsistencies(request *Request, read *Operation, readCall *abstractgraph.AbstractDatabaseCall) []*XCYInconsistency {
	var inconsistencies []*XCYInconsistency
	if readCall.DbInstance.IsQueue() {
		// skip inconsistencies between a write and a read in the same queue
		return inconsistencies
	}

	var writes []*Operation
	if detector.HasDetectionMode(DEBUG_XCY_MISSING_DEPENDENCIES) {
		writes = request.GetXCYDependencies()
	} else {
		writes = request.CurrentLineage().GetXCYDependenciesByMostRecent()
	}
	logger.Logger.Tracef("[XCY LINEAGES] got dependencies: %v", writes)

	if detector.HasDetectionMode(XCY_ALL_DATASTORES) {
		for _, write := range writes {
			inconsistency := &XCYInconsistency{
				Write: write,
				Read:  read,
			}
			request.AddInconsistency(inconsistency)
			detector.inconsistencies += 1
		}
		return inconsistencies
	}

	if detector.HasDetectionMode(XCY_EQUAL_DATASTORES) {
		for _, write := range writes {
			if readCall.DbInstance.GetDatastore() == write.Datastore {
				inconsistency := &XCYInconsistency{
					Write: write,
					Read:  read,
				}
				request.AddInconsistency(inconsistency)
				detector.inconsistencies += 1
			}
		}
		return inconsistencies
	}

	return detector.checkXCYDependencyRead(writes, request, read, readCall)
}

func (detector *XCYDetector) checkXCYDependencyRead(writes []*Operation, request *Request, read *Operation, readCall *abstractgraph.AbstractDatabaseCall) []*XCYInconsistency {
	var inconsistencies []*XCYInconsistency
	for _, write := range writes {
		if readCall.DbInstance.GetDatastore() == write.Datastore {
			logger.Logger.Debugf("[XCY] evaluating XCY violation for read (%s) and write (%s) at (%s)", read.Key.GetVariableInfo().Name, write.Key.GetVariableInfo().Name, write.Datastore.GetName())

			var dataflows []*objects.ObjectDataflow

			writeUnderlyingVars := write.GetObjectUnderlyingVariables(true)
			readUnderlyingVars := read.GetKeyUnderlyingVariables(true)

			logger.Logger.Infof("[WRITE VALUE] writes (xcy deps) for (%s) %s: \n%v", utils.GetType(write.Object), write.Object.String(), objects.GetDependenciesStringLst(writeUnderlyingVars...))
			logger.Logger.Infof("[READ KEY] reads (xcy deps) for (%s) %s: \n%v", utils.GetType(read.Key), read.Key.String(), objects.GetDependenciesStringLst(readUnderlyingVars...))

			for _, readVar := range readUnderlyingVars {
				readDfs := readVar.GetVariableInfo().GetAllDataflowsForDatastore(readCall.DbInstance.GetName())
				for _, readDf := range readDfs {
					for _, writeVar := range writeUnderlyingVars {
						if readDf.HasVariable(writeVar) {
							dataflows = append(dataflows, writeVar.GetVariableInfo().GetAllDataflows()...)
						}
					}
				}
			}

			if dataflows != nil {
				inconsistency := NewXCYInconsistency(write, read, dataflows)
				request.AddInconsistency(inconsistency)
				inconsistencies = append(inconsistencies, inconsistency)
				logger.Logger.Warnf("[XCY] found inconsistency: %s", inconsistency.String())

			}
		}
	}
	detector.inconsistencies += len(inconsistencies)
	return inconsistencies
}

func operationDataDependsOnPrevious(op_1 *Operation, op_2 *Operation) bool {
	datastore_1 := op_1.GetDatastore()
	datastore_2 := op_2.GetDatastore()

	if datastore_1 != datastore_2 {
		opUnderlyingVars_1 := op_1.GetAllUnderlyingVariables(false)
		opUnderlyingVars_2 := op_2.GetKeyUnderlyingVariables(false)

		for _, opVar_2 := range opUnderlyingVars_2 { //FIXME: SHOULD BE ONLY READ DATAFLOWS
			dataflows := opVar_2.GetVariableInfo().GetAllDataflowsForDatastore(datastore_2.GetName())
			for _, dataflow := range dataflows {
				if dataflow.HasAnyVariable(opUnderlyingVars_1) {
					return true
				}
			}
		}
	}
	return false
}

func filterWriteDependencySet(writes_1 []*Operation, read_2 *Operation) {
	for _, write_1 := range writes_1 {
		for _, dep_1 := range write_1.GetVisibleDependenciesForDatastore(read_2.GetDatastore()) {
			write_1.AddMinimizedDependencyIfNotExists(dep_1)
		}
	}
}

func (detector *XCYDetector) MinimizeDependecySets(request *Request) {
	if !detector.HasDetectionMode(DEBUG_XCY_MINIMIZE_DEPENDENCIES) {
		return
	}
	reads := request.GetReadOperations()
	for i, read_1 := range reads {
		writes_1 := detector.getWriteOperationsForDatastore(read_1.GetDatastore())
		for j := i + 1; j < len(reads); j++ {
			read_2 := reads[j]
			if operationDataDependsOnPrevious(read_1, read_2) {
				logger.Logger.Warnf("READ 2 (%s) DEPENDS ON READ 1 (%s)", read_2.String(), read_1.String())
				filterWriteDependencySet(writes_1, read_2)
			}
		}
	}
}

func (detector *XCYDetector) HasDetectionMode(mode DetectionMode) bool {
	return detector.detectionMode == mode
}

func detectMissingDependencies(request *Request, inconsistencies []*XCYInconsistency) {
	for _, inconsistency := range inconsistencies {
		if !slices.Contains(request.CurrentLineage().GetXCYDependenciesByMostRecent(), inconsistency.Write) {
			inconsistency.MissingDependency = true
		}
	}
}

func (detector *XCYDetector) attachOperationToDatastore(operation *Operation) {
	detector.datastoresOps[operation.GetDatastore()] = append(detector.datastoresOps[operation.GetDatastore()], operation)
}

func (detector *XCYDetector) getWriteOperationsForDatastore(datastore *datastores.Datastore) []*Operation {
	var ops []*Operation
	for _, op := range detector.datastoresOps[datastore] {
		if op.Write {
			ops = append(ops, op)
		}
	}
	return ops
}
