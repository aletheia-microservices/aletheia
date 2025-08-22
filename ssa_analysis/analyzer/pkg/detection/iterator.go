package detection

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/app/backends"
	"analyzer/pkg/common"
)

type Iterator struct {
	app       *app.App
	graph     *abstractgraph.AbstractCallGraph
	detectors []Detector
}

func NewIterator(app *app.App, graph *abstractgraph.AbstractCallGraph, detectors ...Detector) *Iterator {
	return &Iterator{
		app:       app,
		graph:     graph,
		detectors: detectors,
	}
}

func (iterator *Iterator) Run() {
	for _, detector := range iterator.detectors {
		detector.OnNewRun(iterator.app)
	}

	clientNode := iterator.graph.GetNodeByName("client")

	for _, edge := range iterator.graph.GetEdgesFromNode(clientNode) {
		toNode := edge.GetToNode()

		for _, detector := range iterator.detectors {
			detector.OnNewRequest(toNode)
		}

		// FIXME: skip for now
		// maybe for the future we can ensure we do not append the Run to the nodes list
		// but let it be attached to edges
		if toNode.GetMethod() == "Run" {
			continue
		}

		iterator.transverse(toNode)

		for _, detector := range iterator.detectors {
			detector.OnEndRequest(iterator.app)
		}

		iterator.clean(toNode)
	}

	for _, detector := range iterator.detectors {
		detector.OnEndRun(iterator.app)
	}
}

func (iterator *Iterator) clean(node *abstractgraph.AbstractNode) {
	for _, param := range node.GetParams() {
		param.CleanSecondaryTaints()
	}
	for _, edge := range iterator.graph.GetEdgesFromNode(node) {
		iterator.clean(edge.GetToNode())
	}
}

func (iterator *Iterator) transverse(node *abstractgraph.AbstractNode) {
	for _, detector := range iterator.detectors {
		detector.OnNewNode(iterator.app, node)
	}

	for _, edge := range iterator.graph.GetEdgesFromNode(node) {
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
				taintMapping.Merge(taintMappingTmp)
			}

			// update future propagation with taints received from caller args to current params
			abstractgraph.PropagateNewTaintsToDatabaseCallObjects(iterator.graph, toNode, taintMapping)
			abstractgraph.PropagateNewTaintsToTracedObjects(iterator.graph, toNode, nil, nil, true)

			// finalize phase by propagating to database schemas
			abstractgraph.PropagateNewTaintsToDatabaseSchemas(iterator.graph, taintMapping)

			// --------------------------
			// PHASE 2: transverse caller
			// --------------------------

			iterator.transverse(edge.GetToNode())

			// -----------------------------------
			// PHASE 3: propagate taints to callee
			// -----------------------------------

			taintMapping = abstractgraph.NewTaintMapping()
			// propagate taints across services (backwards): args (from) <<< params (to)
			for i, fromArg := range edge.GetArguments() {
				toParam := toNode.GetParamAt(i)
				fmt.Printf("[TRANSVERSE] [ARG << PARAM] fromArg=%s // toParam=%s\n", fromArg.String(), toParam.String())
				taintMappingTmp := abstractgraph.MergeTaints(fromArg, toParam.GetPrimaryTaints(), false, false)
				taintMapping.Merge(taintMappingTmp)
			}
			// propagate taints across services (backwards): rets (from) <<< rets (to)
			for i, fromRet := range edge.GetReturns() {
				toRet := toNode.GetReturnAt(i)
				fmt.Printf("[TRANSVERSE] [RET << RET] fromRet=%s // toRet=%s\n", fromRet.String(), toRet.String())
				taintMappingTmp := abstractgraph.MergeTaints(fromRet, toRet.GetPrimaryTaints(), false, false)
				taintMapping.Merge(taintMappingTmp)
			}

			abstractgraph.PropagateNewTaintsToTracedObjects(iterator.graph, node, taintMapping, edge, false)
			abstractgraph.PropagateNewTaintsToDatabaseSchemas(iterator.graph, taintMapping)
		}

		if edge.GetEdgeType() == abstractgraph.EDGE_DATABASE_CALL {
			// ===================
			// DATABASE OPERATIONS
			// ===================
			for _, detector := range iterator.detectors {
				switch edge.GetOpType() {
				case common.OP_READ:
					detector.OnRead(iterator.app, edge)
				case common.OP_WRITE:
					detector.OnWrite(iterator.app, edge)
				case common.OP_UPDATE:
					detector.OnUpdate(iterator.app, edge)
				case common.OP_DELETE:
					detector.OnDelete(iterator.app, edge)
				}
			}
			// FIXME: this is a bit hardcoded for now
			currDB := iterator.app.GetDatabaseByName(edge.GetToNode().GetDatabaseName())
			if currDB.IsQueue() && edge.GetOpType() == common.OP_WRITE {
				iterator.transverseQueue(node, currDB, edge)
			}
		}
	}

	for _, detector := range iterator.detectors {
		detector.OnEndNode(iterator.app, node)
	}

}

func (iterator *Iterator) transverseQueue(node *abstractgraph.AbstractNode, currDB *backends.Database, edge *abstractgraph.AbstractEdge) {
	for _, queueReadEdge := range iterator.graph.GetEdges() {
		if queueReadEdge.GetEdgeType() == abstractgraph.EDGE_DATABASE_CALL && queueReadEdge.GetOpType() == common.OP_READ {
			otherDB := iterator.app.GetDatabaseByName(queueReadEdge.GetToNode().GetDatabaseName())
			if otherDB == currDB {
				taintMapping := abstractgraph.NewTaintMapping()
				for i, arg := range edge.GetArguments() {
					otherArg := queueReadEdge.GetArgumentAt(i)
					// FIXME: maybe we should also propagate secondary taints?
					taintMappingTmp := abstractgraph.MergeTaints(otherArg, arg.GetPrimaryTaints(), false, false)
					taintMapping.Merge(taintMappingTmp)
				}

				abstractgraph.PropagateNewTaintsToTracedObjects(iterator.graph, node, taintMapping, queueReadEdge, true)
				abstractgraph.PropagateNewTaintsToDatabaseSchemas(iterator.graph, taintMapping)

				callerNode := queueReadEdge.GetFromNode()
				iterator.transverse(callerNode)
			}
		}
	}
}
