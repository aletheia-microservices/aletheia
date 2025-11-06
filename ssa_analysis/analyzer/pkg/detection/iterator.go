package detection

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/common"
)

type IterationPhase int

const (
	PHASE_1_SCHEMA_BUILDER IterationPhase = iota
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

	for _, detector := range it.detectors {
		detector.OnNewRun(it.app)
	}

	clientNode := it.graph.GetNodeByName("client")

	for _, edge := range it.graph.GetEdgesFromNode(clientNode) {
		toNode := edge.GetToNode()

		it.newReqIdx()
		for _, detector := range it.detectors {
			detector.OnNewRequest(toNode, it.currentReqIdx())
		}

		// FIXME: skip for now
		// maybe for the future we can ensure we do not append the Run to the nodes list
		// but let it be attached to edges
		if toNode.GetMethod() == "Run" {
			continue
		}

		it.transverse(toNode)

		for _, detector := range it.detectors {
			detector.OnEndRequest(it.app)
		}
		it.popReqIdx()

		it.clean(toNode)
	}

	for _, detector := range it.detectors {
		detector.OnEndRun(it.app)
	}
}

func (it *Iterator) clean(node *abstractgraph.AbstractNode) {
	for _, param := range node.GetParams() {
		param.CleanSecondaryTaints()
	}
	for _, edge := range it.graph.GetEdgesFromNode(node) {
		it.clean(edge.GetToNode())
	}
}

func (it *Iterator) transverse(node *abstractgraph.AbstractNode) {
	for _, detector := range it.detectors {
		detector.OnNewNode(it.app, node)
	}

	for _, edge := range it.graph.GetEdgesFromNode(node) {
		if edge.GetEdgeType() == abstractgraph.EDGE_SERVICE_RPC {
			// ============
			// SERVICE RPCs
			// ============
			fmt.Printf("[TRANSVERSE] edge=%s\n", edge.String())
			toNode := edge.GetToNode()

			// -----------------------------------
			// PHASE 1: propagate taints to caller
			// -----------------------------------
			taintMapping := abstractgraph.NewTaintMapping()
			// propagate taints across services (forward): args (from) >>> params (to)
			for i, toParam := range toNode.GetParams() {
				fromArg := edge.GetArgumentAt(i)
				fmt.Printf("[TRANSVERSE] [ARG >> PARAM] fromArg=%s // toParam=%s\n", fromArg.String(), toParam.String())
				taintMappingTmp := abstractgraph.MergeTaints(toParam, fromArg.GetPrimaryTaints(), false, false)
				taintMapping.Merge(taintMappingTmp, true)
			}

			// update future propagation with taints received from caller args to current params
			abstractgraph.PropagateNewTaintsToDatabaseCallObjects(it.graph, toNode, taintMapping)
			abstractgraph.PropagateNewTaintsToTracedObjects(it.graph, toNode, nil, nil, true)

			// finalize phase by propagating to database schemas
			if it.mode == PHASE_1_SCHEMA_BUILDER {
				abstractgraph.PropagateNewTaintsToDatabaseSchemas(it.graph, it.currentReqIdx(), taintMapping)
			}
			/* modified := abstractgraph.PropagateNewTaintsToDatabaseSchemas(it.graph, it.currentReqIdx(), taintMapping)
			if it.mode == PHASE_2_PATTERN_DETECTOR && modified {
				log.Fatalf("HERE!")
			} */

			// --------------------------
			// PHASE 2: transverse caller
			// --------------------------

			it.transverse(edge.GetToNode())

			// -----------------------------------
			// PHASE 3: propagate taints to callee
			// -----------------------------------

			taintMapping = abstractgraph.NewTaintMapping()
			// propagate taints across services (backwards): args (from) <<< params (to)
			for i, fromArg := range edge.GetArguments() {
				toParam := toNode.GetParamAt(i)
				fmt.Printf("[TRANSVERSE] [ARG << PARAM] fromArg=%s // toParam=%s\n", fromArg.String(), toParam.String())
				taintMappingTmp := abstractgraph.MergeTaints(fromArg, toParam.GetPrimaryTaints(), false, false)
				taintMapping.Merge(taintMappingTmp, true)
			}
			// propagate taints across services (backwards): rets (from) <<< rets (to)
			for i, fromRet := range edge.GetReturns() {
				toRet := toNode.GetReturnAt(i)
				fmt.Printf("[TRANSVERSE] [RET << RET] fromRet=%s // toRet=%s\n", fromRet.String(), toRet.String())
				taintMappingTmp := abstractgraph.MergeTaints(fromRet, toRet.GetPrimaryTaints(), false, false)
				taintMapping.Merge(taintMappingTmp, true)
			}

			abstractgraph.PropagateNewTaintsToTracedObjects(it.graph, node, taintMapping, edge, false)

			if it.mode == PHASE_1_SCHEMA_BUILDER {
				abstractgraph.PropagateNewTaintsToDatabaseSchemas(it.graph, it.currentReqIdx(), taintMapping)
			}
		}

		if edge.GetEdgeType() == abstractgraph.EDGE_DATABASE_CALL {
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
			// FIXME: this is a bit hardcoded for now
			currDB := it.app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
			if currDB.IsQueue() && edge.GetOpType() == common.OP_WRITE {
				it.transverseQueue(node, currDB, edge)
			}
		}
	}

	for _, detector := range it.detectors {
		detector.OnEndNode(it.app, node)
	}

}

func (it *Iterator) transverseQueue(node *abstractgraph.AbstractNode, currDB *backends.Database, edge *abstractgraph.AbstractEdge) {
	for _, queueReadEdge := range it.graph.GetEdges() {
		if queueReadEdge.GetEdgeType() == abstractgraph.EDGE_DATABASE_CALL && queueReadEdge.GetOpType() == common.OP_READ {
			otherDB := it.app.GetDatabaseByName(queueReadEdge.GetToNode().GetDatabaseName())
			if otherDB == currDB {
				//log.Fatalf("HERE!")
				taintMapping := abstractgraph.NewTaintMapping()
				for i, arg := range edge.GetArguments() {
					otherArg := queueReadEdge.GetArgumentAt(i)
					// FIXME: maybe we should also propagate secondary taints?
					taintMappingTmp := abstractgraph.MergeTaints(otherArg, arg.GetPrimaryTaints(), false, false)
					taintMapping.Merge(taintMappingTmp, true)
				}

				abstractgraph.PropagateNewTaintsToTracedObjects(it.graph, node, taintMapping, queueReadEdge, true)

				if it.mode == PHASE_1_SCHEMA_BUILDER {
					abstractgraph.PropagateNewTaintsToDatabaseSchemas(it.graph, it.currentReqIdx(), taintMapping)
				}

				callerNode := queueReadEdge.GetFromNode()
				it.newReqIdx()
				it.transverse(callerNode)
				it.popReqIdx()
			}
		}
	}
}
