package tainter

import (
	"go/token"
	"log"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

// REMINDER:
// sometimes there are can be taints such as: _obj[*][*], usertimeline_db.usertimeline.Posts[*].PostID
// this may happen when any[] type is using, for example, //EVAL - fmt.Print calls
// e.g., dsb_sn2 at UserTimelineService.ReadUserTimeline:
// > post_ids = append(new_post_ids, post_ids...)
// > //EVAL - fmt.Println(post_ids)
// ----------------------------
// t56: append(t55, t31...)
// t59: make any <- []int64 (t56)
// *t58 = t59
// t58: &t57[0:int]
// t57: new [1]any (vargs)
// t60: slice t57[:]
// t61: //EVAL - fmt.Println(t60...)
// ----------------------------
// t56 taint is:		 _obj[*], usertimeline_db.usertimeline.Posts[*].PostID
// t57 taint becomes: _obj[*][*], usertimeline_db.usertimeline.Posts[*].PostID
// ----------------------------
func doTaintNode(node *ssagraph.SSANode, taintInfo TaintInfo, taintMode TaintMode) {
	if strings.Contains(taintInfo.objpath, "_obj.Username[*]") {
		log.Fatalf("[TAINT] unexpected taint info path (%s): %v\n", taintInfo.objpath, taintInfo)
	}

	if strings.Contains(taintInfo.objpath, "[*][*][*]") {
		log.Fatalf("[TAINT] unexpected taint info path (%s): %v\n", taintInfo.objpath, taintInfo)
	}

	// sanity check for dsb_hotel2 app
	if strings.Contains(taintInfo.objpath, ".HId[*]") && taintInfo.dbTaint.dbcall.GetDatabaseName() == "recommendation_db" {
		log.Fatalf("[TAINT] [DSB_HOTEL2] unexpected taint info path (%s): %v\n", taintInfo.objpath, taintInfo)
	}

	if taintInfo.isTypeDatabase() && taintInfo.getDatabaseCall() == nil {
		// FIXME: verify this
		//EVAL - fmt.Printf("[TAINT] [4] nil db call for taint info: %v\n", taintInfo)
		return
	}
	if taintInfo.isTypeService() && taintInfo.getServiceCall() == nil {
		// FIXME: verify this
		//EVAL - fmt.Printf("[TAINT] [4] nil sv call for taint info: %v\n", taintInfo)
		return
	}

	var ok bool
	switch taintMode {
	case TAINT_MODE_NEARBY:
		if taintInfo.isTypeDatabase() {
			//EVAL - fmt.Printf("[TAINT NEARBY] [DATABASE] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath())
			ok = node.AddDatabaseTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getDatabasePath(), taintInfo.getDatabaseCall())
		} else if taintInfo.isTypeService() {
			//EVAL - fmt.Printf("[TAINT NEARBY] [SERVICE] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath())
			ok = node.AddServiceTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getServicePath(), taintInfo.getServiceCall())
		}
	case TAINT_MODE_FETCH_UPWARDS:
		if taintInfo.isTypeDatabase() {
			//EVAL - fmt.Printf("[TAINT FETCH] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath()+taintInfo.getObjectPath())
			ok = node.AddDatabaseTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getDatabasePath()+taintInfo.getObjectPath(), taintInfo.getDatabaseCall())
		} else if taintInfo.isTypeService() {
			//EVAL - fmt.Printf("[TAINT FETCH] [SERVICE] tainting node (%s) for objpath (%s) and dbfield (%s)\n", node.String(), taintInfo.getObjectFullPath(), taintInfo.getDatabasePath())
			ok = node.AddServiceTaintIfNotExists(taintInfo.getObjectFullPath(), taintInfo.getServicePath(), taintInfo.getServiceCall())
		}
	}
	if ok {
		//EVAL - fmt.Printf("\t[TAINT] OK!\n")
	}
}

func getObjectPathDiff(longPath1 string, shortPath2 string) string {
	longPath1 = strings.TrimPrefix(longPath1, "_obj")
	shortPath2 = strings.TrimPrefix(shortPath2, "_obj")
	// i.e., pathTop - pathBottomRel
	return strings.TrimPrefix(longPath1, shortPath2)
}

func nodeHasTaintInfo(node *ssagraph.SSANode, objpath string, info TaintInfo) bool {
	for _, taint := range node.GetTaintsForPath(objpath) {
		if taint.IsDatabaseTaint() && info.isTypeDatabase() && taint.GetDatabaseCall() == info.getDatabaseCall() {
			return true
		} else if taint.IsServiceTaint() && info.isTypeService() && taint.GetServiceCall() == info.getServiceCall() {
			return true
		}
	}
	return false
}

func generateRootTaintInfoFromTaint(node *ssagraph.SSANode, taint *ssagraph.SSATaint) TaintInfo {
	if taint.IsDatabaseTaint() {
		return NewTaintInfoDatabase(taint.GetDatabasePath(), "", node.GetValue(), taint.GetDatabaseCall())
	} else if taint.IsServiceTaint() {
		return NewTaintInfoService(taint.GetServicePath(), "", node.GetValue(), taint.GetServiceCall())
	}
	log.Panicf("[TAINT INFO FROM TAINT] unexpected type for taint: %s\n", taint.String())
	return TaintInfo{}
}

func propagateTaintNearby(graph *ssagraph.SSAGraph, recurse bool, val ssa.Value, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	if val == nil {
		log.Panicf("[TAINT NEARBY] unexpected nil val // TAINT INFO (_obj%s, %s)\n", taintInfo.getObjectPath(), taintInfo.getDatabasePath())
	}

	/* //var prevVal ssa.Value
	var prevValName string
	if taintInfo.prevval == nil {
		//prevVal = nil
		prevValName = "<nil>"
		taintInfo.prevval = val
	} else {
		//prevVal = taintInfo.prevval
		prevValName = taintInfo.prevval.Name()
	} */

	taintInfo = taintInfo.updateValue(val)

	//EVAL - fmt.Printf("[TAINT NEARBY] visiting %s: %s // TAINT INFO (_obj%s, %s)\n", val.Name(), val.String(), taintInfo.getObjectPath(), taintInfo.getDatabasePath())
	if visited[val] {
		//EVAL - fmt.Printf("\t[TAINT NEARBY] skipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[val] = true

	if ssaValueIsUsedInMongoBsonFilter(graph, val) {
		return // skip
	}

	node := graph.GetNodeByName(val.Name())

	// avoid infinite recursion
	if nodeHasTaintInfo(node, "_obj"+taintInfo.objpath, taintInfo) {
		return
	}

	doTaintNode(node, taintInfo, TAINT_MODE_NEARBY)

	//EVAL - fmt.Printf("[TAINT NEARBY] [PART_1] [ROOT=%t] [RECURSE=%t] [PREV=%s] current node: %v\n", taintInfo.objroot, recurse, prevValName, node)
	taintInfo.prevval = node.GetValue()

	for _, edge := range graph.GetEdgesFromNode(node) {
		toNode := edge.GetToNode()
		//EVAL - fmt.Printf("\t[TAINT NEARBY] [PART_1] [r=%t] [FromNode %s] edge (%s) to node: %v\n", recurse, node.GetValue().Name(), edge.GetTypeString(), toNode)

		switch edge.GetType() {

		case ssagraph.EDGE_FIELD:
			if upwards {
				// TODO: EXTEND TO TYPE DATABASE??
				// FOR NOW WE ONLY NEED TO SERVICE TYPE
				// BECAUSE WE NEED TO SPREAD THE TRACES FROM UPPER STRUCTS TO LOWER FIELDS
				// E.G. TRAINTICKET PRESERVESERVICE.PRESERVE() WHERE TRIPALLINFO AND ORDER USE THE OTI.TRIPID PARAMETER
				// we don't need to do for type database because we already to "check upper taints"
				if taintInfo.isTypeService() || taintInfo.isTypeDatabase() {
					// found field corresponding to upper taintinfo objpath
					if taintInfo.objpath == "."+edge.GetParam() {
						for _, upperTaint := range node.GetTaintsForPath("_obj" + taintInfo.objpath) {
							taintInfoTmp := generateRootTaintInfoFromTaint(toNode, upperTaint)
							/* if taintInfo.isTypeDatabase() && graph.GetService() == "UserTimelineService" && graph.GetMethodName() == "ReadUserTimeline" {
								if taintInfo.isTypeDatabase() && strings.Contains(upperTaint.GetDatabasePath(), edge.GetParam()) {
									if field, ok := prevVal.(*ssa.FieldAddr); ok {
										log.Fatalf("HERE!! (PREV_VAL=%v // TO_NODE=%v // OBJPATH=%s // DBPATH=%s)\n", field.Name(), toNode.String(), taintInfoTmp.objpath, upperTaint.GetDatabasePath())
									}
								}
							} */
							// node has taint info if it was the previous node calling propagateTaintNearby
							// (e.g., lower field propagating to upper struct)
							// we need to avoid visiting it again otherwise we will have infinite recursion!
							if !nodeHasTaintInfo(toNode, "_obj", taintInfoTmp) {
								propagateTaintNearby(graph, true, toNode.GetValue(), taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, false)
							}
						}
					}
				} else {
					//TODO
				}
				break
			}

			var taintInfoTmp TaintInfo
			if taintInfo.isObjectRoot() {
				taintInfoTmp = taintInfo.updateCallPathSuffix("." + edge.GetParam())
			} else if edge.GetParam() == taintInfo.objpath { // FIXME: NOT SURE ABOUT THIS!
				taintInfoTmp = taintInfo.enableObjectRoot()
				taintInfoTmp.dbTaint.dbpath, _ = strings.CutSuffix(taintInfoTmp.objpath, taintInfo.objpath)
			} else {
				//EVAL - fmt.Printf("[TAINT NEARBY] [PART_1] [FIELD] unexpected conditions // FROM_NODE=%s // TAINT INFO (_obj%s, %s)\n", toNode.String(), taintInfo.getObjectPath(), taintInfo.getDatabasePath())
				continue
			}

			propagateTaintNearby(graph, true, toNode.GetValue(), taintInfoTmp, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_INDEX:
			if upwards {
				if taintInfo.isTypeService() {
					// found index corresponding to upper taintinfo objpath
					if taintInfo.objpath == "["+edge.GetParam()+"]" {
						for _, upperTaint := range node.GetTaintsForPath("_obj" + taintInfo.objpath) {
							taintInfoTmp := generateRootTaintInfoFromTaint(toNode, upperTaint)
							// node has taint info if it was the previous node calling propagateTaintNearby
							// (e.g., lower index propagating to upper slice)
							// we need to avoid visiting it again otherwise we will have infinite recursion!
							if !nodeHasTaintInfo(toNode, "_obj", taintInfoTmp) {
								propagateTaintNearby(graph, true, toNode.GetValue(), taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, false)
							}
						}
					}
				} else {
					// TODO
				}
				break
			}
			var taintInfoTmp TaintInfo
			if taintInfo.isObjectRoot() {
				taintInfoTmp = taintInfo.updateCallPathSuffix("[" + edge.GetParam() + "]")
			} else {
				var ok bool
				taintInfoTmp = taintInfo.enableObjectRoot()
				taintInfoTmp.objpath, ok = strings.CutSuffix(taintInfoTmp.objpath, "[*]")
				if !ok {
					log.Fatalf("[TAINT NEARBY] [PART_1] [INDEX] could not cut suffix [*] for objpath = (%s)\n", taintInfoTmp.objpath)
				}
			}

			propagateTaintNearby(graph, true, toNode.GetValue(), taintInfoTmp, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_MAP_UPDATE:
			if upwards {
				// TODO (similar to FIELD and INDEXES but needs more boring logic)
				break
			}

			var taintInfoTmpKey, taintInfoTmpVal TaintInfo
			var keyOk, valOk bool
			if taintInfo.isObjectRoot() {
				taintInfoTmpKey = taintInfo.updateCallPathSuffix(".MapKey")
				taintInfoTmpVal = taintInfo.updateCallPathSuffix(".MapVal")
				keyOk = true
				valOk = true
			} else {
				taintInfoTmpKey = taintInfo.enableObjectRoot()
				taintInfoTmpKey.objpath, keyOk = strings.CutSuffix(taintInfoTmpKey.objpath, ".MapKey")
				taintInfoTmpVal = taintInfo.enableObjectRoot()
				taintInfoTmpVal.objpath, valOk = strings.CutSuffix(taintInfoTmpKey.objpath, ".MapVal")
			}

			for _, edge := range graph.GetEdgesToNode(toNode) {
				if edge.GetFromNode() != node {
					if edge.IsType(ssagraph.EDGE_MAP_KEY) && keyOk {
						propagateTaintNearby(graph, true, edge.GetFromNode().GetValue(), taintInfoTmpKey, visited, checkTaintInfo, upwards)
					} else if edge.IsType(ssagraph.EDGE_MAP_VALUE) && valOk {
						propagateTaintNearby(graph, true, edge.GetFromNode().GetValue(), taintInfoTmpVal, visited, checkTaintInfo, upwards)
					}
				}
			}

		case ssagraph.EDGE_LOOKUP_INDEX:
			// PROPAGATE DOWNWARDS (to lookup target)
			//
			// e.g., in dsb_sn2 @ UserTimelineService.ReadUserTimeline
			// seen_posts := make(map[int64]bool)
			// [...]
			// if _, ok := seen_posts[post.PostID]; ok {
			// 	 continue
			// }
			// ------------------------
			// t124: &t114.PostID [#0]
			// t125: *t124
			// t126: t23[t125], ok
			// t127: extract t126 #0
			// t128: extract t126 #1
			// ------------------------
			lookupTarget := toNode.GetValueLookup().X
			taintInfoTmp := taintInfo.updateObjectPathSuffix(".MapKey")
			taintInfoTmp = taintInfoTmp.disableObjectRoot()
			propagateTaintNearby(graph, true, lookupTarget, taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, upwards)
			// TODO: propagate to nodes that extract the current val

		case ssagraph.EDGE_LOOKUP_TARGET:
			// PROPAGATE UPWARDS (to lookup index)
			lookupIndex := toNode.GetValueLookup().Index
			// ignore if it actually has the suffix or not:
			// - if there is suffix then we propagate the corresponding taint
			// - otherwise (meaning that the entire map was passed in the call), we just propagate the generic top taint
			taintInfoTmp, _ := taintInfo.cutObjectPathSuffix(".MapKey")
			taintInfoTmp = taintInfoTmp.enableObjectRoot()
			propagateTaintNearby(graph, true, lookupIndex, taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, upwards)
			// TODO: propagate to nodes that extract the current val

		case ssagraph.EDGE_MAP_KEY, ssagraph.EDGE_MAP_VALUE:
			// skip for now

		case ssagraph.EDGE_STORE_ADDRESS:
			val := toNode.GetInstruction().(*ssa.Store).Val
			valNode := graph.GetNodeByName(val.Name())
			propagateTaintNearby(graph, true, valNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		case ssagraph.EDGE_STORE_VALUE:
			addr := toNode.GetInstruction().(*ssa.Store).Addr
			addrNode := graph.GetNodeByName(addr.Name())
			propagateTaintNearby(graph, true, addrNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		case ssagraph.EDGE_POINTS_TO:
			// ignore for now
		case ssagraph.EDGE_ARG_ON_CALL:
			call := toNode.GetValue().(*ssa.Call)
			if call == nil {
				log.Fatalf("[TAINT NEARBY] [PART_1] unexpected type for node: [%T] %s\n", toNode.GetValue(), toNode.GetValue())
			}
			if !call.Call.IsInvoke() {
				if builtin, ok := call.Call.Value.(*ssa.Builtin); ok {
					if ok, funcType, _ := utils.SSABuiltinFuncIsDirect(builtin); ok {
						if funcType == utils.FUNC_TYPE_APPEND {
							// REMINDER:
							// builtin append() can safely taint its arguments because
							// the Go SSA transforms the second argument into a slice
							// e.g. in dsb_sn2 at UserTimelineService.ReadUserTimeline:
							// > new_post_ids = append(new_post_ids, post.PostID)
							// ----------------------------
							// t107: local PostInfo (post)
							// t122: &t107.PostID [#0]
							// t123: *t122
							// t124: new[1]int64 (varargs)
							// t125: &t124[0:int]
							// *t125 = t123
							// t126: slice t124[:]
							// t127: append(t111, t126...)
							// ----------------------------
							propagateTaintNearby(graph, true, call, taintInfo, visited, checkTaintInfo, upwards)
						}
						if funcType == utils.FUNC_TYPE_APPEND || funcType == utils.FUNC_TYPE_TRANSFER {
							// e.g.
							// append(t111, t126...)
							// copy(...)
							for _, arg := range call.Call.Args {
								argNode := graph.GetNodeByName(arg.Name())
								if argNode != node {
									propagateTaintNearby(graph, true, argNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
								}
							}
						} else if funcType == utils.FUNC_TYPE_MAP_ELEMS {
							for idxToTaint, edge := range graph.GetEdgesToNode(toNode) {
								// e.g., delete(m map[Type]Type1, key Type)
								// propagate from map to key
								funcArg := edge.GetFromNode()
								if edge.IsType(ssagraph.EDGE_POINTS_TO) || funcArg == node {
									continue
								}

								// TODO!!!!!
								if idxToTaint == 0 {
									taintInfoTmp := taintInfo.updateObjectPathSuffix(".MapKey")
									taintInfoTmp = taintInfoTmp.disableObjectRoot()
									propagateTaintNearby(graph, true, funcArg.GetValue(), taintInfoTmp, visited, checkTaintInfo, upwards)
								} else if idxToTaint == 1 {
									propagateTaintNearby(graph, true, funcArg.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
								} else {
									log.Fatalf("[TAINT NEARBY] [PART_1] unexpected index (%d)\n", idxToTaint)
								}
							}
						}
					}
				} else if fn, ok := call.Call.Value.(*ssa.Function); ok {
					if fn.Package() != nil {
						if fn.Package().Pkg.Name() == "strconv" {
							if fn.Name() == "FormatInt" || fn.Name() == "Itoa" {
								// e.g.,
								// (1) id_str := strconv.FormatInt(id, 10)
								// 		- taint: id_str <<< id
								// (2) amount, err := strconv.ParseFloat(declineOverAmount, 32)
								// 		- taint: amount <<< declineOverAmount
								argNode := graph.GetNodeByName(call.Call.Args[0].Name())
								// check if its the current node
								if node == argNode {
									callNode := graph.GetNodeByName(call.Name())
									propagateTaintNearby(graph, true, callNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
								}
							} else if fn.Name() == "ParseFloat" {
								// e.g., amount, err := strconv.ParseFloat(declineOverAmount, 32)
								// taint: amount <<< declineOverAmount
								// TODO (needs SSA EXTRACT logic)
							} else {
								log.Fatalf("[TAINT NEARBY] [PART_1] skipping function (%s) for call: [%T] %s\n", fn.Name(), call.Call.Value, call.Call.Value)
							}
						} else {
							//EVAL - fmt.Printf("[TAINT NEARBY] [PART_1] skipping package (%s) for call: [%T] %s\n", fn.Package().Pkg.Name(), call.Call.Value, call.Call.Value)
						}
					} else {
						//EVAL - fmt.Printf("[TAINT NEARBY] [PART_1] skipping nil package for call: [%T] %s\n", call.Call.Value, call.Call.Value)
					}
				} else {
					log.Fatalf("[TAINT NEARBY] [PART_1] unexpected type for call value: [%T] %s\n", call.Call.Value, call.Call.Value)
				}
			}
		case ssagraph.EDGE_RETURN_ON, ssagraph.EDGE_EXTRACT:
			// skip
		case ssagraph.EDGE_BINOP_X:
			binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
			if binOp.Op >= token.EQL && binOp.Op <= token.GTR || binOp.Op >= token.NEQ && binOp.Op <= token.GEQ {
				// skip (if conditions)
			} else {
				propagateTaintNearby(graph, true, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			}
		case ssagraph.EDGE_BINOP_Y:
			binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
			if binOp.Op >= token.EQL && binOp.Op <= token.GTR || binOp.Op >= token.NEQ && binOp.Op <= token.GEQ {
				// skip (if conditions)
			} else {
				propagateTaintNearby(graph, true, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			}

		default:
			propagateTaintNearby(graph, true, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			/* if taintInfo.isObjectRoot() {
				propagateTaintNearby(graph, true, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			} */
		}
	}

	//EVAL - fmt.Printf("[TAINT NEARBY] [PART_2] [ROOT=%t] [RECURSE=%t] [PREV=%s] current node: %v\n", taintInfo.objroot, recurse, prevValName, node)
	for _, edge := range graph.GetEdgesToNode(node) {
		fromNode := edge.GetFromNode()
		//EVAL - fmt.Printf("\t[TAINT NEARBY] [PART_2] [r=%t] [ToNode %s] edge (%s) from node: %v\n", recurse, node.GetValue().Name(), edge.GetTypeString(), fromNode)

		if ssaValueIsUsedInMongoBsonFilter(graph, fromNode.GetValue()) {
			continue // skip
		}

		switch edge.GetType() {

		case ssagraph.EDGE_LOOKUP_INDEX:
			// (COPY PASTE FROM PART 1)
			lookupTarget := node.GetValueLookup().X
			taintInfoTmp := taintInfo.updateObjectPathSuffix(".MapKey")
			taintInfoTmp = taintInfoTmp.disableObjectRoot()
			propagateTaintNearby(graph, true, lookupTarget, taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, upwards)

		case ssagraph.EDGE_LOOKUP_TARGET:
			// (COPY PASTE FROM PART 1)
			lookupIndex := node.GetValueLookup().Index
			taintInfoTmp, _ := taintInfo.cutObjectPathSuffix(".MapKey")
			taintInfoTmp = taintInfoTmp.enableObjectRoot()
			propagateTaintNearby(graph, true, lookupIndex, taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, upwards)

		case ssagraph.EDGE_FIELD:
			visitedTmp := make(map[ssa.Value]bool)

			var taintInfoTmp TaintInfo
			if taintInfo.isObjectRoot() {
				taintInfoTmp = taintInfo.updateObjectPathPrefix("." + edge.GetParam())
			} else if edge.GetParam() == taintInfo.objpath { // FIXME: NOT SURE ABOUT THIS!
				taintInfoTmp = taintInfo.enableObjectRoot()
			} else {
				//EVAL - fmt.Printf("[TAINT NEARBY] [PART_2] [FIELD] unexpected conditions // FROM_NODE=%s // TAINT INFO (_obj%s, %s)\n", fromNode.String(), taintInfo.getObjectPath(), taintInfo.getDatabasePath())
				continue
			}

			/* taintInfoTmp := taintInfo.updateObjectPathPrefix("." + edge.GetParam()) */
			propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfoTmp, visitedTmp, checkTaintInfo, true)

		case ssagraph.EDGE_INDEX:
			visitedTmp := make(map[ssa.Value]bool)

			var taintInfoTmp TaintInfo
			if taintInfo.isObjectRoot() {
				taintInfoTmp = taintInfo.updateObjectPathPrefix("[" + edge.GetParam() + "]")
			} else {
				taintInfoTmp = taintInfo.enableObjectRoot()
			}

			propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfoTmp, visitedTmp, checkTaintInfo, true)

		case ssagraph.EDGE_MAP_UPDATE:
			/* mapUpdateNode := edge.GetToNode().GetInstruction() */

		case ssagraph.EDGE_MAP_KEY, ssagraph.EDGE_MAP_VALUE:
			// skip for now

		case ssagraph.EDGE_STORE_ADDRESS:
			val := fromNode.GetInstruction().(*ssa.Store).Val
			valNode := graph.GetNodeByName(val.Name())
			propagateTaintNearby(graph, true, valNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		case ssagraph.EDGE_STORE_VALUE:
			addr := fromNode.GetInstruction().(*ssa.Store).Addr
			addrNode := graph.GetNodeByName(addr.Name())
			propagateTaintNearby(graph, true, addrNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		case ssagraph.EDGE_RETURN_ON, ssagraph.EDGE_EXTRACT:
			// skip
		case ssagraph.EDGE_BINOP_X:
			binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
			if binOp.Op >= token.EQL && binOp.Op <= token.GTR || binOp.Op >= token.NEQ && binOp.Op <= token.GEQ {
				// skip (if conditions)
			} else {
				propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			}
		case ssagraph.EDGE_BINOP_Y:
			binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
			if binOp.Op >= token.EQL && binOp.Op <= token.GTR || binOp.Op >= token.NEQ && binOp.Op <= token.GEQ {
				// skip (if conditions)
			} else {
				propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			}

		case ssagraph.EDGE_USAGE, ssagraph.EDGE_PHI_ON:
			propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)

		case ssagraph.EDGE_POINTS_TO:
		// ignore for now

		case ssagraph.EDGE_ARG_ON_CALL:
			call := node.GetValue().(*ssa.Call)
			if call == nil {
				log.Fatalf("[TAINT NEARBY] [PART_2] unexpected type for node: [%T] %s\n", fromNode.GetValue(), fromNode.GetValue())
			}
			if !call.Call.IsInvoke() {
				if builtin, ok := call.Call.Value.(*ssa.Builtin); ok {
					if ok, funcType, _ := utils.SSABuiltinFuncIsDirect(builtin); ok {
						if funcType == utils.FUNC_TYPE_APPEND {
							// REMINDER:
							// builtin append() can safely taint its arguments because
							// the Go SSA transforms the second argument into a slice
							// e.g. in dsb_sn2 at UserTimelineService.ReadUserTimeline:
							// > new_post_ids = append(new_post_ids, post.PostID)
							// ----------------------------
							// t107: local PostInfo (post)
							// t122: &t107.PostID [#0]
							// t123: *t122
							// t124: new[1]int64 (varargs)
							// t125: &t124[0:int]
							// *t125 = t123
							// t126: slice t124[:]
							// t127: append(t111, t126...)
							// ----------------------------
							for _, arg := range call.Call.Args {
								argNode := graph.GetNodeByName(arg.Name())
								propagateTaintNearby(graph, true, argNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
							}
						}
					}
				} else if fn, ok := call.Call.Value.(*ssa.Function); ok && fn.Package() != nil {
					if fn.Package().Pkg.Name() == "strconv" {
						if fn.Package() != nil {
							if fn.Name() == "FormatInt" || fn.Name() == "Itoa" {
								// e.g.,
								// (1) id_str := strconv.FormatInt(id, 10)
								// 		- taint: id_str >>> id
								// (2) amount, err := strconv.ParseFloat(declineOverAmount, 32)
								// 		- taint: amount >>> declineOverAmount
								arg := call.Call.Args[0]
								argNode := graph.GetNodeByName(arg.Name())
								propagateTaintNearby(graph, true, argNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
							} else if fn.Name() == "ParseFloat" {
								// e.g., amount, err := strconv.ParseFloat(declineOverAmount, 32)
								// taint: amount >>> declineOverAmount
								// TODO (needs SSA EXTRACT logic)
							} else {
								log.Fatalf("[TAINT NEARBY] [PART_2] skipping function (%s) for call: [%T] %s\n", fn.Name(), call.Call.Value, call.Call.Value)
							}
						} else {
							//EVAL - fmt.Printf("[TAINT NEARBY] [PART_2] skipping nil package for call: [%T] %s\n", call.Call.Value, call.Call.Value)
						}
					} else {
						//EVAL - fmt.Printf("[TAINT NEARBY] [PART_2] skipping package (%s) for call: [%T] %s\n", fn.Package().Pkg.Name(), call.Call.Value, call.Call.Value)
					}
				} else {
					log.Fatalf("[TAINT NEARBY] [PART_2] unexpected type for call value: [%T] %s\n", call.Call.Value, call.Call.Value)
				}
			}

		default:
			/* if taintInfo.isObjectRoot() {
				propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
			} */
			propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
		}
	}
	//EVAL - fmt.Printf("\t[TAINT NEARBY] [PART_2] exiting %s: %s\n", val.Name(), val.String())
}

func propagateTaintFetchUpwards(graph *ssagraph.SSAGraph, val ssa.Value, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	taintInfo = taintInfo.updateValue(val)

	//EVAL - fmt.Printf("[TAINT FETCH] visiting %s: %s // TAINT INFO (%s, %s)\n", val.Name(), val.String(), taintInfo.getObjectPath(), taintInfo.getDatabasePath())
	if visited[val] {
		//EVAL - fmt.Printf("\t[TAINT FETCH] skipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[val] = true

	node := graph.GetNodeByName(val.Name())
	//EVAL - fmt.Printf("\t[TAINT FETCH] checking upper taints: %v\n", node.GetTaints())

	if ssaValueIsUsedInMongoBsonFilter(graph, val) {
		return // skip
	}

	// 1. taint "subpaths" for current variable and save to later taint the corresponding "subobjects" that requested the upper taint
	for objPath, taints := range node.GetTaints() {

		//EVAL - fmt.Printf("\t[TAINT FETCH] comparing prefixes:\n\t - tainted obj path:\t %s\n\t - bottom to upper:\t %s\n", objPath, taintInfo.getObjectFullPath())

		if strings.HasPrefix(taintInfo.getObjectFullPath(), objPath) && taintInfo.getObjectFullPath() != objPath {
			// e.graph.,
			// existing path: 	_obj
			// current path: 	_obj.Shipping
			//
			// in this case, '_obj.Shipping' has prefix '_obj'
			// as result, we may get:
			//
			// existing taint: 	_obj			@ order_db.order
			// potential taint: _obj.Shipping 	@ order_db.order.Shipping
			for _, taint := range taints {
				// save the taint in the upper node
				taintInfoTmp := taintInfo
				taintInfoTmp.dbTaint.dbpath = taint.GetDatabasePath()
				taintInfoTmp.dbTaint.dbcall = taint.GetDatabaseCall()
				doTaintNode(node, taintInfoTmp, TAINT_MODE_FETCH_UPWARDS)

				// so that we can later taint the bottom node
				dbFieldIndirect := taintInfoTmp.getDatabasePath() + taintInfo.getObjectPath()
				if taintInfoTmp.getDatabaseCall() == nil {
					// FIXME: verify this
					//EVAL - fmt.Printf("[TAINT FETCH] [4] nil db call for taint info tmp: %v\n", taintInfoTmp)
				} else {
					checkTaintInfo.addToIndirectTaints(dbFieldIndirect, taintInfoTmp.getDatabaseCall())
				}
			}
			break
		} else if strings.HasPrefix(objPath, taintInfo.getObjectFullPath()) { // also true if strings are equal
			// e.graph.,
			// upper's taint: 		_obj.PostID @ posts_db.post.PostID
			// bottom's path: 		_obj.PostID
			// => bottom's taint: 	_obj		@ posts_db.post.PostID

			pathDiff := getObjectPathDiff(objPath, taintInfo.getObjectFullPath())
			for _, taint := range taints {
				checkTaintInfo.addToInheritedTaints(pathDiff, taint.GetDatabasePath(), taint.GetDatabaseCall())
			}
		}
	}
	// 2. taint forward propagation
	// propagate/update the taints for/of objects that use the current one
	// e.g., in postnotification.StorageService.ReadPost():
	// after "t8: new primitive.E (slicelit)" taint is updated with a subpath for the postid value of the bson filter
	// then we need to go forward again and propagate the taints for "t13: slice t8[...]"
	for _, edge := range graph.GetEdgesFromNode(node) {
		toVal := edge.GetToNode().GetValue()
		switch toVal.(type) {
		case *ssa.MakeInterface, *ssa.Slice:
			propagateTaintFetchUpwards(graph, toVal, taintInfo, visited, checkTaintInfo, upwards)
			// TODO: maybe we also need to do this for:
			// (i) nodes whose pointerto set have the current node
			// (ii) nodes within the pointerto set of the current node
			// (iii) load and store instrs?
		}

	}

	//EVAL - fmt.Printf("\t[TAINT FETCH] exiting %s: %s\n", val.Name(), val.String())
}

func RunTainter(graph *ssagraph.SSAGraph) {
	// for simplicity, always run taint on database calls before service calls due to cases where
	// 1st there is a service call that returns some value
	// 2nd that value is then used in a database call
	//
	// the taint is detected when analyzing the service call
	// only if the database call was analyzed before
	// because of the way we are checking upper taints with the "spreading of taints in store points"
	// at the service call
	//
	// this logic is close to the tainting process of returns and parameters
	// where we only taint them after all database calls have been processed
	runTainterOnDatabaseCalls(graph)
	runTainterOnServiceCalls(graph)
	runTainterOnParameters(graph)
	runTainterOnReturns(graph)
}

func runTainterOnParameters(graph *ssagraph.SSAGraph) {
	paramNodes := graph.GetFuncParametersExceptMemberAndContext()

	/* service := graph.GetService()
	method := graph.GetMethodName()
	funcShortPath := graph.GetFunctionShortPath()
	callId := graph.GetServiceWithMethod()

	call := ssagraph.NewServiceCall(callId, nil, paramNodes, nil, service, method, funcShortPath)
	call.EnableDummy()
	graph.SetMethodCall(call)

	for _, paramNode := range paramNodes {
		arg := paramNode.GetValue()
		svpath := call.String() + "." + arg.Name()
		visited := make(map[ssa.Value]bool)
		taintInfo := NewTaintInfoService(svpath, "", nil, call)
		propagateTaintNearby(graph, false, arg, taintInfo, visited, nil, false)
	} */

	checkUpperTaintsForObjects(graph, paramNodes)
}

func runTainterOnReturns(graph *ssagraph.SSAGraph) {
	var nodesToVisit []*ssagraph.SSANode
	for _, node := range graph.GetNodes() {
		// keep track of arguments returned in the current function possibly invoked by other services
		// so that we can mark their indirect taints
		if ret, ok := node.GetInstruction().(*ssa.Return); ok {
			var rets []*ssagraph.SSANode
			for _, res := range ret.Results {
				resNode := graph.GetNodeByName(res.Name())
				nodesToVisit = append(nodesToVisit, resNode)
				rets = append(rets, resNode)
			}
			graph.AddReturnsToLst(rets)
		}

		checkUpperTaintsForObjects(graph, nodesToVisit)
	}
}

func runTainterOnDatabaseCalls(graph *ssagraph.SSAGraph) {
	for _, node := range graph.GetNodes() {
		var nodesToVisit []*ssagraph.SSANode
		if database, collectionOrTopic, method, opType, valFieldPathLst, ok := isDatabaseCall(graph, node.GetValue()); ok {
			var argNodes []*ssagraph.SSANode

			for _, valFieldPath := range valFieldPathLst {
				argNodes = append(argNodes, graph.GetNodeByName(valFieldPath.val.Name()))
			}

			callId := ssagraph.ComputeCallID(graph, node)
			dbCall := ssagraph.NewDatabaseCall(callId, node, argNodes, database, collectionOrTopic, method, opType)
			graph.AddDatabaseCall(dbCall)

			for _, valFieldPath := range valFieldPathLst {
				dbfield := valFieldPath.fieldpath
				obj := valFieldPath.val

				visited := make(map[ssa.Value]bool)

				taintInfo := NewTaintInfoDatabase(dbfield, "", nil, dbCall)

				// e.g. dsb_hotel2: RecommendationService.LoadRecommendations()
				// where &hotels is the destination
				// -----------------------------------------------
				// var hotels []Hotel
				// cursor, err := collection.FindMany(ctx, filter)
				// cursor.All(ctx, &hotels)
				// -----------------------------------------------
				// REMINDER:
				// cannot do default propagateTaintNearby (root=true) for read operations fetched
				// into arrays/slices (e.g., FindMany) because the default taintinfo is the root path
				// when propagating to other objects
				if valFieldPath.bsonCursorMany || valFieldPath.bsonFilterIn {
					taintInfo.objpath += "[*]"
					taintInfo.objroot = false
				}

				propagateTaintNearby(graph, false, obj, taintInfo, visited, nil, false)
				objNode := graph.GetNodeByName(obj.Name())
				nodesToVisit = append(nodesToVisit, objNode)
			}

			checkUpperTaintsForObjects(graph, nodesToVisit)
		}
	}
}

func runTainterOnServiceCalls(graph *ssagraph.SSAGraph) {
	for _, node := range graph.GetNodes() {
		var nodesToVisit []*ssagraph.SSANode
		// keep track of arguments passed in service RPCs
		// so that we can mark their indirect taints
		if service, method, funcShortPath, args, call, ok := isServiceCall(graph, node.GetValue()); ok {
			// keep track of objects passed as arguments
			var argNodes []*ssagraph.SSANode
			//EVAL - fmt.Printf("[TAINT] added service call (%s) --> (%s)\n", graph.GetFunctionShortPath(), funcShortPath)
			for _, arg := range args {
				argNodes = append(argNodes, graph.GetNodeByName(arg.Name()))
				//EVAL - fmt.Printf("[TAINT] checking taint for service call with arg: %s\n", arg.String())
				node := graph.GetNodeByName(arg.Name())
				nodesToVisit = append(nodesToVisit, node)
			}

			// keep track of objects extracted from returns
			var retNodes []*ssagraph.SSANode
			callNode := graph.GetNodeByName(call.Name())
			if call.Call.Signature().Results().Len() > 1 {
				for _, edge := range graph.GetEdgesFromNode(callNode) {
					if edge.GetType() == ssagraph.EDGE_EXTRACT {
						nodesToVisit = append(nodesToVisit, edge.GetToNode())
						retNodes = append(retNodes, edge.GetToNode())
					}
				}
			} else {
				// when there is only one return value then there
				// are no extract instructions and the value is just
				// the one declared when invoking the function
				nodesToVisit = append(nodesToVisit, callNode)
				retNodes = append(retNodes, callNode)
			}
			//EVAL - fmt.Printf("[TAINT] [%s] got ret nodes: %v\n", node.String(), retNodes)

			callId := ssagraph.ComputeCallID(graph, node)
			svcCall := ssagraph.NewServiceCall(callId, node, argNodes, retNodes, service, method, funcShortPath)
			graph.AddServiceCall(svcCall)

			for _, argNode := range argNodes {
				arg := argNode.GetValue()
				svpath := svcCall.String() + "." + arg.Name()
				visited := make(map[ssa.Value]bool)
				taintInfo := NewTaintInfoService(svpath, "", nil, svcCall)
				propagateTaintNearby(graph, false, arg, taintInfo, visited, nil, false)
			}

			for _, retNode := range retNodes {
				ret := retNode.GetValue()
				svpath := svcCall.String() + "." + ret.Name()
				visited := make(map[ssa.Value]bool)
				taintInfo := NewTaintInfoService(svpath, "", nil, svcCall)
				propagateTaintNearby(graph, false, ret, taintInfo, visited, nil, false)
			}

			//EVAL - fmt.Printf("[TAINT] visiting nodes for call (%s) --> (%s)\n", graph.GetFunctionShortPath(), funcShortPath)
			checkUpperTaintsForObjects(graph, nodesToVisit)
		}

	}
}

func checkUpperTaintsForObjects(graph *ssagraph.SSAGraph, nodesToVisit []*ssagraph.SSANode) {
	// check for upper taints affecting the current database/service calls objects
	for _, originNode := range nodesToVisit {
		//EVAL - fmt.Println()
		//EVAL - fmt.Printf("[TAINT] check upper taints for node: %v\n", originNode.String())
		visited := make(map[ssa.Value]bool)
		taintInfo := NewTaintInfoDatabase("", "", nil, nil)
		checkTaintInfo := NewCheckTaintInfo()
		propagateTaintFetchUpwards(graph, originNode.GetValue(), taintInfo, visited, checkTaintInfo, false)
		node := graph.GetNodeByName(originNode.GetValue().Name())

		// indirect taints
		for _, taint := range checkTaintInfo.indirectTaints {
			if taint.dbcall == nil {
				log.Fatalf("[1] nil db call for taint: %v\n", taint)
			}
			taintInfo := NewTaintInfoDatabase(taint.dbpath, "", originNode.GetValue(), taint.dbcall)
			doTaintNode(node, taintInfo, TAINT_MODE_NEARBY)
		}

		// inherited taints
		for objpath, taints := range checkTaintInfo.inheritedTaints {
			//EVAL - fmt.Printf("[TAINT] check inherited taints for objpath (%s): %v\n", objpath, taints)
			for _, taint := range taints {
				if taint.dbcall == nil {
					// FIXME: verify this
					//EVAL - fmt.Printf("[2] nil db call for taint: %v\n", taint)
				} else {
					taintInfo := NewTaintInfoDatabase(taint.dbpath, objpath, originNode.GetValue(), taint.dbcall)
					doTaintNode(node, taintInfo, TAINT_MODE_NEARBY)
				}
			}
		}
	}
}
