package detection

import (
	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/common"
	"analyzer/pkg/config"
)

type IterationPhase int

const (
	PHASE_0_DEBUG IterationPhase = iota
	PHASE_1_SCHEMA_BUILDER
	PHASE_1_SCHEMA_BUILDER_READ_ONLY
	PHASE_2_PATTERN_DETECTOR
)

type Iterator struct {
	app       *app.App
	graph     *abstractgraph.AbstractCallGraph
	detectors []Detector
	mode      IterationPhase
	// need to simulate stack because of queue writes triggering new requests
	reqIdxStack []int
	nextReqIdx  int
}

func NewIterator(app *app.App, graph *abstractgraph.AbstractCallGraph, detectors ...Detector) *Iterator {
	return &Iterator{
		app:       app,
		graph:     graph,
		detectors: detectors,
	}
}

func (it *Iterator) newReqIdx() {
	it.reqIdxStack = append(it.reqIdxStack, it.nextReqIdx)
	it.nextReqIdx++
}

func (it *Iterator) popReqIdx() int {
	if len(it.reqIdxStack) == 0 {
		return -1
	}
	top := it.reqIdxStack[len(it.reqIdxStack)-1]
	it.reqIdxStack = it.reqIdxStack[:len(it.reqIdxStack)-1]
	return top
}

func (it *Iterator) currentReqIdx() int {
	if len(it.reqIdxStack) == 0 {
		return -1
	}
	return it.reqIdxStack[len(it.reqIdxStack)-1]
}

func (it *Iterator) Run(mode IterationPhase) {
	it.mode = mode
	// reset for the next phase (at the end the indexes will be the same)
	it.nextReqIdx = 0
	it.reqIdxStack = []int{}

	if it.mode == PHASE_2_PATTERN_DETECTOR {
		for _, detector := range it.detectors {
			detector.OnNewRun(it.app)
		}
	}

	clientNode := it.graph.GetNodeByName("client")

	for _, edge := range it.graph.GetEdgesFromNode(clientNode) {
		toNode := edge.GetToNode()

		it.newReqIdx()

		if it.mode == PHASE_2_PATTERN_DETECTOR {
			for _, detector := range it.detectors {
				detector.OnNewRequest(toNode, it.currentReqIdx())
			}
		}

		// FIXME: skip for now
		// maybe for the future we can ensure we do not append the Run to the nodes list
		// but let it be attached to edges
		if toNode.GetMethod() == "Run" {
			continue
		}

		it.transverse(toNode)

		if it.mode == PHASE_2_PATTERN_DETECTOR {
			for _, detector := range it.detectors {
				detector.OnEndRequest(it.app)
			}
		}
		it.popReqIdx()

		it.clean(toNode)
	}

	if it.mode == PHASE_2_PATTERN_DETECTOR {
		for _, detector := range it.detectors {
			detector.OnEndRun(it.app)
		}
	}
}

func (it *Iterator) clean(node *abstractgraph.AbstractNode) {
	if it.mode == PHASE_0_DEBUG {
		return
	}

	for _, param := range node.GetParams() {
		param.CleanSecondaryTaints()
	}
	for _, ret := range node.GetReturns() {
		ret.CleanSecondaryTaints()
	}
	for _, edge := range it.graph.GetEdgesFromNode(node) {
		for _, arg := range edge.GetArguments() {
			arg.CleanSecondaryTaints()
		}
		for _, ret := range edge.GetReturns() {
			ret.CleanSecondaryTaints()
		}
		it.clean(edge.GetToNode())

		if edge.GetEdgeType() == abstractgraph.EDGE_DATABASE_CALL {
			currDB := it.app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
			if currDB.IsQueue() && edge.GetOpType() == common.OP_WRITE {
				_, callerNode, ok := it.findMatchingQueuePop(node, currDB, edge)
				if !ok {
					continue
				}
				it.clean(callerNode)
			}
		}
	}
}

func (it *Iterator) transverse(node *abstractgraph.AbstractNode) {
	// EVAL: logrus.Tracef("[TRAVERSE] traversing node: %s\n", node.String())
	if it.mode == PHASE_2_PATTERN_DETECTOR {
		for _, detector := range it.detectors {
			detector.OnNewNode(it.app, node)
		}
	}

	for _, edge := range it.graph.GetEdgesFromNode(node) {
		if edge.GetEdgeType() == abstractgraph.EDGE_SERVICE_RPC {
			// EVAL: logrus.Tracef("[TRAVERSE] visiting service call edge: %s\n", edge.String())
			// ============
			// SERVICE RPCs
			// ============
			toNode := edge.GetToNode()

			// -----------------------------------
			// PHASE 1: propagate taints to caller
			// -----------------------------------
			taintMapping := abstractgraph.NewTaintMapping()

			// propagate taints across services (forward): args (from) >>> params (to)
			for i, toParam := range toNode.GetParams() {
				fromArg := edge.GetArgumentAt(i)
				// EVAL: logrus.Tracef("[TRANSVERSE] [ARG >> PARAM] fromArg=%s // toParam=%s\n", fromArg.String(), toParam.String())
				taintMappingTmp := abstractgraph.MergeTaints(toParam, fromArg.GetAllTaintsBeforeT(edge.GetT()), nil, abstractgraph.MERGE_MODE_TAINT, config.MIN_T, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)
				taintMapping.Join(taintMappingTmp, true)
				taintMappingTmp = abstractgraph.MergeTaints(toParam, fromArg.GetAllTaintsAfterT(edge.GetT()), nil, abstractgraph.MERGE_MODE_TAINT, config.MAX_T, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)
				taintMapping.Join(taintMappingTmp, true)
			}

			// update taints on future node
			abstractgraph.PropagateTaintsToServiceCallObjects(it.graph, toNode, taintMapping, nil, true, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)
			abstractgraph.PropagateNewTaintsToDatabaseCallObjects(it.graph, toNode, taintMapping, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)

			// finalize phase by propagating to database schemas
			if it.mode == PHASE_1_SCHEMA_BUILDER || it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY {
				abstractgraph.PropagateNewTaintsToDatabaseSchemas(it.graph, it.currentReqIdx(), taintMapping, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)
			}

			// --------------------------
			// PHASE 2: transverse caller
			// --------------------------

			it.transverse(edge.GetToNode())

			// -----------------------------------
			// PHASE 3: propagate taints to callee
			// -----------------------------------

			taintMapping.Clear()

			// propagate taints across services (backwards): args (from) <<< params (to)
			for i, fromArg := range edge.GetArguments() {
				toParam := toNode.GetParamAt(i)
				// EVAL: logrus.Tracef("[TRANSVERSE] [ARG << PARAM] fromArg=%s // toParam=%s\n", fromArg.String(), toParam.String())
				taintMappingTmp := abstractgraph.MergeTaints(fromArg, toParam.GetAllTaints(), nil, abstractgraph.MERGE_MODE_TAINT, edge.GetT(), it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)
				taintMapping.Join(taintMappingTmp, true)
			}

			// propagate taints across services (backwards): rets (from) <<< rets (to)
			for i, fromRet := range edge.GetReturns() {
				toRet := toNode.GetReturnAt(i)
				// EVAL: logrus.Tracef("[TRANSVERSE] [RET << RET] fromRet=%s // toRet=%s\n", fromRet.String(), toRet.String())
				taintMappingTmp := abstractgraph.MergeTaints(fromRet, toRet.GetAllTaints(), nil, abstractgraph.MERGE_MODE_TAINT, edge.GetT(), it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)
				taintMapping.Join(taintMappingTmp, true)
			}

			// update taints on current node
			abstractgraph.PropagateTaintsToServiceCallObjects(it.graph, node, taintMapping, edge, false, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)
			abstractgraph.PropagateNewTaintsToDatabaseCallObjects(it.graph, node, taintMapping, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)

			// finalize phase by propagating to database schemas
			if it.mode == PHASE_1_SCHEMA_BUILDER || it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY {
				abstractgraph.PropagateNewTaintsToDatabaseSchemas(it.graph, it.currentReqIdx(), taintMapping, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)
			}
		}

		if edge.GetEdgeType() == abstractgraph.EDGE_DATABASE_CALL {
			// EVAL: logrus.Tracef("[TRAVERSE] visiting database call edge: %s\n", edge.String())
			// ===================
			// DATABASE OPERATIONS
			// ===================
			if it.mode == PHASE_2_PATTERN_DETECTOR {
				for _, detector := range it.detectors {
					switch edge.GetOpType() {
					case common.OP_READ, common.OP_READ_MANY:
						detector.OnRead(it.app, it.currentReqIdx(), edge)
					case common.OP_WRITE:
						detector.OnWrite(it.app, it.currentReqIdx(), edge)
					case common.OP_UPDATE:
						detector.OnUpdate(it.app, it.currentReqIdx(), edge)
					case common.OP_DELETE:
						detector.OnDelete(it.app, it.currentReqIdx(), edge)
					}
				}
			}
			// TODO: this is a bit hardcoded for now but can definitely be improved
			currDB := it.app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
			if currDB.IsQueue() && edge.GetOpType() == common.OP_WRITE {
				it.transverseQueue(node, currDB, edge)
			}
		}
	}

	if it.mode == PHASE_2_PATTERN_DETECTOR {
		for _, detector := range it.detectors {
			detector.OnEndNode(it.app, node)
		}
	}

}

func (it *Iterator) findMatchingQueuePop(node *abstractgraph.AbstractNode, currDB *backends.Database, edge *abstractgraph.AbstractEdge) (*abstractgraph.AbstractEdge, *abstractgraph.AbstractNode, bool) {
	if !currDB.IsQueue() {
		return nil, nil, false
	}

	for _, edge := range it.graph.GetEdges() {
		if edge.GetEdgeType() == abstractgraph.EDGE_DATABASE_CALL && edge.GetOpType() == common.OP_READ {
			otherDB := it.app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
			if otherDB.IsQueue() && otherDB == currDB {
				callerNode := edge.GetFromNode()
				return edge, callerNode, true
			}
		}
	}
	return nil, nil, false
}

func (it *Iterator) transverseQueue(node *abstractgraph.AbstractNode, currDB *backends.Database, edge *abstractgraph.AbstractEdge) {
	queueReadEdge, callerNode, ok := it.findMatchingQueuePop(node, currDB, edge)
	if !ok {
		return
	}

	if config.Global.PropagateTaintsAcrossQueueOperations {
		taintMapping := abstractgraph.NewTaintMapping()
		for i, arg := range edge.GetArguments() {
			otherArg := queueReadEdge.GetArgumentAt(i)
			taintMappingTmp := abstractgraph.MergeTaints(otherArg, arg.GetPrimaryTaints(), nil, abstractgraph.MERGE_MODE_TAINT, config.MIN_T, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)
			taintMapping.Join(taintMappingTmp, true)
		}
		abstractgraph.PropagateTaintsToServiceCallObjects(it.graph, node, taintMapping, queueReadEdge, true, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)

		if it.mode == PHASE_1_SCHEMA_BUILDER || it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY {
			abstractgraph.PropagateNewTaintsToDatabaseSchemas(it.graph, it.currentReqIdx(), taintMapping, it.mode == PHASE_1_SCHEMA_BUILDER_READ_ONLY)
		}
	}

	it.newReqIdx()
	it.transverse(callerNode)
	it.popReqIdx()
}
