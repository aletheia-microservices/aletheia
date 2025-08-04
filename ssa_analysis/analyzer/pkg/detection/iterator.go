package detection

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
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
		iterator.clean(toNode)

		for _, detector := range iterator.detectors {
			detector.OnEndRequest(iterator.app)
		}
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
	fmt.Printf("[ITERATOR] visiting node: %s\n", node.String())

	for _, detector := range iterator.detectors {
		detector.OnNewNode(iterator.app, node)
	}

	fmt.Printf("[ITERATOR] on new node (%s) with edges:\n", node.String())
	for _, edge := range iterator.graph.GetEdgesFromNode(node) {
		fmt.Printf("\t[ITERATOR] (ID = %s) %s\n", edge.GetID(), edge.String())
	}
	fmt.Println()

	for i, edge := range iterator.graph.GetEdgesFromNode(node) {
		if edge.GetEdgeType() == abstractgraph.EDGE_SERVICE_RPC {
			fmt.Printf("\t[ITERATOR] visiting service call: %s\n", edge.String())
			toNode := edge.GetToNode()

			// propagate taints across services (forward): args (from) >>> params (to)
			taintMapping := abstractgraph.NewTaintMapping()
			for i, toParam := range toNode.GetParams() {
				fmt.Printf("debug to param: %s\n", toParam.String())
				fromArg := edge.GetArgumentAt(i)
				taintMappingTmp := abstractgraph.MergeTaints(toParam, fromArg.GetPrimaryTaints(), false)
				taintMapping.Merge(taintMappingTmp)
			}
			for _, edge := range iterator.graph.GetEdgesFromNode(toNode) {
				fmt.Printf("\t\t[ITERATOR] propagate new taints to objects args on edge: %s\n", edge.String())
				for i, obj := range edge.GetArguments() {
					fmt.Printf("\t\t\t[ITERATOR] [ARG %d] > ENTER update object (%s) taints with new taint mapping\n", i, obj.String())
					propagateNewTaintsToObject(obj, taintMapping)
					fmt.Printf("\t\t\t[ITERATOR] [ARG %d] < EXIT update object (%s) taints with new taint mapping\n", i, obj.String())
				}
			}
			abstractgraph.PropagateNewTaintsToDatabases(iterator.graph, taintMapping)

			// propagate taints across services (backwards): rets (from) <<< rets (to)
			taintMapping = abstractgraph.NewTaintMapping()
			for i, fromRet := range edge.GetReturns() {
				toRet := toNode.GetReturnAt(i)
				taintMappingTmp := abstractgraph.MergeTaints(fromRet, toRet.GetPrimaryTaints(), false)
				taintMapping.Merge(taintMappingTmp)
				fmt.Printf("\t\t[ITERATOR] [RETS] [index=%d] taint mapping for ret (%s): %s\n", i, fromRet.String(), taintMappingTmp.String())
			}
			abstractgraph.PropagateNewTaintsToDatabases(iterator.graph, taintMapping)

			fmt.Printf("\t[ITERATOR] final taint mapping: %s\n", taintMapping.String())

			for i, obj := range node.GetParams() {
				fmt.Printf("\t\t[ITERATOR] [PARAM %d] > ENTER update object (%s) taints with new taint mapping\n", i, obj.String())
				propagateNewTaintsToObject(obj, taintMapping)
				fmt.Printf("\t\t[ITERATOR] [PARAM %d] < EXIT update object (%s) taints with new taint mapping\n", i, obj.String())
			}
			for i, obj := range node.GetReturns() {
				fmt.Printf("\t\t[ITERATOR] [RET %d] > ENTER update object (%s) taints with new taint mapping\n", i, obj.String())
				propagateNewTaintsToObject(obj, taintMapping)
				fmt.Printf("\t\t[ITERATOR] [RET %d] < EXIT update object (%s) taints with new taint mapping\n", i, obj.String())
			}
			for _, edge := range iterator.graph.GetEdgesFromNode(node)[i+1:] {
				fmt.Printf("\t\t[ITERATOR] propagate new taints to objects args on edge: %s\n", edge.String())
				for i, obj := range edge.GetArguments() {
					fmt.Printf("\t\t\t[ITERATOR] [ARG %d] > ENTER update object (%s) taints with new taint mapping\n", i, obj.String())
					propagateNewTaintsToObject(obj, taintMapping)
					fmt.Printf("\t\t\t[ITERATOR] [ARG %d] < EXIT update object (%s) taints with new taint mapping\n", i, obj.String())
				}
			}

			iterator.transverse(edge.GetToNode())
		}

		if edge.GetEdgeType() == abstractgraph.EDGE_DATABASE_CALL {
			fmt.Printf("\t[ITERATOR] visiting database call: %s\n", edge.String())
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
				for _, otherEdge := range iterator.graph.GetEdges() {
					if otherEdge.GetEdgeType() == abstractgraph.EDGE_DATABASE_CALL && otherEdge.GetOpType() == common.OP_READ {
						otherDB := iterator.app.GetDatabaseByName(otherEdge.GetToNode().GetDatabaseName())
						if otherDB == currDB {
							taintMapping := abstractgraph.NewTaintMapping()
							for i, arg := range edge.GetArguments() {
								otherArg := otherEdge.GetArgumentAt(i)
								// FIXME: maybe we should also propagate secondary taints?
								taintMappingTmp := abstractgraph.MergeTaints(otherArg, arg.GetPrimaryTaints(), false)
								taintMapping.Merge(taintMappingTmp)
							}

							for _, edge := range iterator.graph.GetEdgesFromNode(otherEdge.GetFromNode()) {
								if otherEdge == edge {
									continue
								}
								fmt.Printf("\t\t[ITERATOR] propagate new taints to objects args on edge: %s\n", edge.String())
								for i, obj := range edge.GetArguments() {
									fmt.Printf("\t\t\t[ITERATOR] [ARG %d] > ENTER update object (%s) taints with new taint mapping\n", i, obj.String())
									propagateNewTaintsToObject(obj, taintMapping)
									fmt.Printf("\t\t\t[ITERATOR] [ARG %d] < EXIT update object (%s) taints with new taint mapping\n", i, obj.String())
								}
							}

							abstractgraph.PropagateNewTaintsToDatabases(iterator.graph, taintMapping)
							callerNode := otherEdge.GetFromNode()
							iterator.transverse(callerNode)
						}
						
					}
				}
			}
		}
	}

	for _, detector := range iterator.detectors {
		detector.OnEndNode(iterator.app, node)
	}
}

func propagateNewTaintsToObject(obj *abstractgraph.AbstractObject, taintMapping *abstractgraph.TaintMapping) {
	for currTaint, otherTaintsLst := range taintMapping.GetMapping() {
		objpath, found := obj.FindObjectPathWithEqualOrUpperTaint(currTaint)
		for _, otherTaint := range otherTaintsLst {
			if found {
				obj.AddTaintIfNotExists(objpath, otherTaint)
			}
		}
	}
}
