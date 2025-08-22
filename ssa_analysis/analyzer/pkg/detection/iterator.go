package detection

import (
	"fmt"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/common"
	"analyzer/pkg/utils"
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
			taintMapping := abstractgraph.NewTaintMapping()

			// --------------------------------------------------
			// --------------------------------------------------

			// propagate taints across services (forward): args (from) >>> params (to)
			for i, toParam := range toNode.GetParams() {
				fmt.Printf("debug toParam: %s\n", toParam.String())
				fromArg := edge.GetArgumentAt(i)
				taintMappingTmp := abstractgraph.MergeTaints(toParam, fromArg.GetPrimaryTaints(), false, false)
				taintMapping.Merge(taintMappingTmp)
			}

			// propagate taints across services (backwards): args (from) <<< params (to)
			for i, fromArg := range edge.GetArguments() {
				fmt.Printf("debug fromArg: %s\n", fromArg.String())
				toParam := toNode.GetParamAt(i)
				taintMappingTmp := abstractgraph.MergeTaints(fromArg, toParam.GetPrimaryTaints(), false, false)
				taintMapping.Merge(taintMappingTmp)
			}

			// propagate taints across services (forward): args (from) >>> params (to)
			for i, toParam := range toNode.GetParams() {
				fmt.Printf("debug toParam: %s\n", toParam.String())
				fromArg := edge.GetArgumentAt(i)
				taintMappingTmp := abstractgraph.MergeTaints(toParam, fromArg.GetPrimaryTaints(), false, false)
				taintMapping.Merge(taintMappingTmp)
			}

			// propagate taints across services (backwards): rets (from) <<< rets (to)
			for i, fromRet := range edge.GetReturns() {
				fmt.Printf("\t\t[ITERATOR] [RETS] [index=%d] taint mapping for ret (%s)\n", i, fromRet.String())
				toRet := toNode.GetReturnAt(i)
				taintMappingTmp := abstractgraph.MergeTaints(fromRet, toRet.GetPrimaryTaints(), false, false)
				taintMapping.Merge(taintMappingTmp)
			}

			// update object taints to the next node
			for _, edge2 := range iterator.graph.GetEdgesFromNode(toNode) {
				fmt.Printf("\t\t[ITERATOR] propagate new taints to objects args on edge: %s\n", edge2.String())
				for i, obj := range edge2.GetArguments() {
					fmt.Printf("\t\t\t[ITERATOR] [ARG %d] > ENTER update object (%s) taints with new taint mapping\n", i, obj.String())
					propagateNewTaintsToObject(obj, taintMapping)
					fmt.Printf("\t\t\t[ITERATOR] [ARG %d] < EXIT update object (%s) taints with new taint mapping\n", i, obj.String())
				}
			}

			abstractgraph.PropagateNewTaintsToDatabases(iterator.graph, taintMapping)

			// --------------------------------------------------
			// --------------------------------------------------

			for _, otherEdge := range iterator.graph.GetEdgesFromNode(node) {
				if otherEdge == edge || otherEdge.GetEdgeType() != abstractgraph.EDGE_SERVICE_RPC {
					continue
				}

				// propagate taints across objects using trace info (args)
				for _, arg := range edge.GetArguments() {
					fmt.Printf("[TRACE] [ARG] arg={%s} // edge={%s} // otherEdge={%s} // taintMapping={%s}\n", arg.String(), edge.String(), otherEdge.String(), taintMapping.String())
					taintTracedObjects(arg, otherEdge, taintMapping)
				}

				// propagate taints across objects using trace info (rets)
				for _, ret := range edge.GetReturns() {
					fmt.Printf("[TRACE] [RET] ret={%s} // edge={%s} // otherEdge={%s} // taintMapping={%s}\n", ret.String(), edge.String(), otherEdge.String(), taintMapping.String())
					taintTracedObjects(ret, otherEdge, taintMapping)
				}
			}

			fmt.Printf("\t[TRACE ITERATOR] final taint mapping: %s\n", taintMapping.String())

			fmt.Println("[1] ============================================================ [1]")

			// update object taints to the next node
			for _, edge2 := range iterator.graph.GetEdgesFromNode(toNode) {
				fmt.Printf("\t\t[ITERATOR] propagate new taints to objects args on edge: %s\n", edge2.String())
				for i, obj := range edge2.GetArguments() {
					fmt.Printf("\t\t\t[ITERATOR] [ARG %d] > ENTER update object (%s) taints with new taint mapping\n", i, obj.String())
					propagateNewTaintsToObject(obj, taintMapping)
					fmt.Printf("\t\t\t[ITERATOR] [ARG %d] < EXIT update object (%s) taints with new taint mapping\n", i, obj.String())
				}
			}

			// update object taints within all the following services
			for _, edge2 := range iterator.graph.GetEdgesFromNode(node) {
				if edge2.GetEdgeType() == abstractgraph.EDGE_SERVICE_RPC {
					toNode := edge2.GetToNode()
					for _, edge3 := range iterator.graph.GetEdgesFromNode(toNode) {
						fmt.Printf("\t\t[TRACE ITERATOR] propagate new taints to objects args on edge: %s\n", edge2.String())
						for i, obj := range edge3.GetArguments() {
							fmt.Printf("\t\t\t[TREACE ITERATOR] [ARG %d] > ENTER update object (%s) taints with new taint mapping\n", i, obj.String())
							propagateNewTaintsToObject(obj, taintMapping)
							fmt.Printf("\t\t\t[TRACE ITERATOR] [ARG %d] < EXIT update object (%s) taints with new taint mapping\n", i, obj.String())
						}
					}
				}
			}

			fmt.Println("[2] ============================================================ [2]")
			abstractgraph.PropagateNewTaintsToDatabases(iterator.graph, taintMapping)

			fmt.Println("[3] ============================================================ [3]")

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
								taintMappingTmp := abstractgraph.MergeTaints(otherArg, arg.GetPrimaryTaints(), false, false)
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

func taintTracedObjects(obj *abstractgraph.AbstractObject, otherEdge *abstractgraph.AbstractEdge, taintMapping *abstractgraph.TaintMapping) {
	for objpath, tracesLst := range obj.GetTraces() {
		// e.g.,
		// MediaMicroservices in APIService.ReadPage(...)
		//
		// movieId := movieIdService.ReadMovieId(title)
		// movieInfo := movieInfoService.ReadMovieInfo(movieId.ID)
		//
		// t4 = ReadMovieId(..) => objpath 		 (@ t4.MovieID) = _obj.MovieID
		// ReadMovieInfo(t10) 	=> tracedObjPath (@ t10) 		= _obj
		//
		// REMINDER: traceObjPath is simply the objpath of the traced object
		for _, trace := range tracesLst {
			if trace.GetServiceCallID() != otherEdge.GetID() {
				continue
			}
			// e.g., SockShop3 @ Frontend.AddItem
			// AddItem(ctx, sessionID, Item{ID: itemID, Quantity: 1, UnitPrice: sock.Price})
			// <=> AddItem(ctx, sessionID, t18)
			// ------------------------------------
			// 		t12: local Item (complit)
			// ------------------------------------
			// 		t13: &t12.ID (#0)
			// ==== tainted ====
			// 		_obj
			// [rpc] @ CatalogueService.Get.itemID
			// [rpc] @ CatalogueService.AddItem.t18.ID
			// 	  	   ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
			// ------------------------------------
			// 		t18: *t12
			// 		^^^^^^^^^
			// ==== tainted ====
			//		 _obj
			// [rpc] @ CartService.AddItem.t18
			// 		_obj.ID
			// 		^^^^^^^^
			// [rpc] @ CartService.Get.itemID
			// [rpc] @ CartService.AddItem.t18.ID
			// ------------------------------------
			// 
			// CURRENT OBJECT is t13 w/ objpath = _obj
			// TRACED OBJECT is t18 w/ traceobjpath = _obj.ID
			//
			// we want to get the taints of t13 at _obj and propagate them to t18 on _obj.ID
			// REMINDER: we just simply associate the SAME dbfield to t18 on _obj.ID
			fmt.Printf("[TRACE] [OBJ=%s // OBJPATH=%s] trace={%s}\n", obj.String(), objpath, trace.LongString())
			tracedObj := otherEdge.GetArgumentByNameIfExists(trace.GetArgumentName())
			if tracedObj != nil {
				tracedObjPath := trace.GetArgumentPath()
				fmt.Printf("[TRACE] [OBJ=%s // OBJPATH=%s] corresponding trace obj (path=%s): %s\n", obj.String(), objpath, tracedObjPath, tracedObj.String())
				var selectedTaints = make(map[string][]*abstractgraph.AbstractTaint)

				var ok = true
				var subpath = ""
				// if there is no taint for current objpath then it is possible that there are upper taints
				// so we go up, create a new subtaint and save to the selected taints of the traced object
				// e.g., MediaMicroservices: APIService.ReadPage():
				//
				// 				CURRENT OBJECT BELOW
				// ------------------------------------------
				// ==== arg 1 (movieID) tainted ====
				// 			_obj
				// [read, secondary] @ movieid_db.movieid
				// 			_obj.ID
				// [rpc] @ MovieIdService.ReadMovieId.movieID
				// ------------------------------------------
				// after going up, we get a new potential subtaint
				// (that we don't save for the current obj but only for the traced obj)
				// ------------------------------------------
				// 			_obj.ID
				// [read, secondary] @ movieid_db.movieid.ID
				// ------------------------------------------
				for ok {
					for _, taint := range obj.GetTaintsForObjectPath(objpath) {
						subTaint := taint.Copy()
						subTaint.AddSuffixToDatabasePath(subpath)
						selectedTaints[tracedObjPath] = append(selectedTaints[tracedObjPath], subTaint)
					}
					objpath, subpath, ok = utils.ExtractUpperPath(objpath)
				}

				taintMappingTmp := abstractgraph.MergeTaints(tracedObj, selectedTaints, false, true)
				fmt.Printf("[TRACE] taint mapping tmp = %s\n", taintMappingTmp.String())
				taintMapping.Merge(taintMappingTmp)
			}

		}
	}
}
