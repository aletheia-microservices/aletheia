package tainter

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

func propagateTaintNearbyFromNodeOnCall(graph *ssagraph.SSAGraph, node *ssagraph.SSANode, toNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, upwards bool, call *ssa.Call) {
	if !call.Call.IsInvoke() {
		if builtin, ok := call.Call.Value.(*ssa.Builtin); ok {
			if ok, funcType, _ := utils.SSABuiltinFuncIsDirect(builtin); ok {
				if funcType == utils.FUNC_TYPE_APPEND {
					// NOTE: builtin append() can safely taint its arguments because
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
					propagateTaintNearby(graph, true, call, taintInfo, visited, upwards)
				}
				if funcType == utils.FUNC_TYPE_APPEND || funcType == utils.FUNC_TYPE_TRANSFER {
					// e.g., append(t111, t126...)
					// e.g., copy(...)
					for _, arg := range call.Call.Args {
						argNode := graph.GetNodeByName(arg.Name())
						if argNode != node {
							propagateTaintNearby(graph, true, argNode.GetValue(), taintInfo, visited, upwards)
						}
					}
				} else if funcType == utils.FUNC_TYPE_MAP_ELEMS {
					for idxToTaint, edge := range graph.GetEdgesToNode(toNode) {
						funcArg := edge.GetFromNode()
						if funcArg == node {
							continue
						}

						switch call.Call.Value.Name() {
						case "delete":
							// e.g., delete(m map[Type]Type1, key Type)
							// e.g., delete(t0, t59), where t0 is a map and t59 is a map-value
							logrus.WithField("call", call.String()).
								WithField("curr/from", node.String()).
								WithField("other arg", edge.GetFromNode().String()).
								WithField("taint_info", taintInfo.String()).
								Infof("[TAINT NEARBY] [FROM] found call")

							if idxToTaint == 0 {
								taintInfoTmp := taintInfo.updateObjectPathPrefix(DYNAMIC_MAP_KEY + ".Key")
								taintInfoTmp = taintInfoTmp.disableObjectRoot()
								propagateTaintNearby(graph, true, funcArg.GetValue(), taintInfoTmp, visited, upwards)
							} else if idxToTaint == 1 {
								taintInfoTmp, ok := taintInfo.cutObjectPathPrefix(DYNAMIC_MAP_KEY + ".Key")
								taintInfoTmp = taintInfoTmp.tryEnableObjectRoot()
								if ok {
									propagateTaintNearby(graph, true, funcArg.GetValue(), taintInfoTmp, visited, upwards)
								}
							}
						default:
							logrus.WithField("call", call.String()).Fatalf("[TAINT NEARBY] [FROM] unexpected call")
						}
					}
				}
			}
		} else if fn, ok := call.Call.Value.(*ssa.Function); ok && fn.Package() != nil {
			switch fn.Package().Pkg.Name() {
			case "strconv":
				switch fn.Name() {
				case "FormatInt", "Itoa", "ParseInt", "ParseFloat":
					// e.g., myval2 := strconv.ParseInt(myval, 10, 64)
					// e.g., myval2 := strconv.FormatInt(val, 10)
					// e.g., myval2 := strconv.Itoa(myval)

					// propagate from argument value to call value
					if node == graph.GetNodeByName(call.Call.Args[0].Name()) {
						callNode := graph.GetNodeByName(call.Name())
						propagateTaintNearby(graph, true, callNode.GetValue(), taintInfo, visited, upwards)
					}
				default:
					logrus.WithField("call", call.String()).Fatalf("[TAINT NEARBY] [FROM] unexpected call on (strconv) package")
				}
			}
		}
	}
}

func propagateTaintNearbyToNodeOnCall(graph *ssagraph.SSAGraph, taintInfo TaintInfo, visited map[ssa.Value]bool, upwards bool, call *ssa.Call) {
	if !call.Call.IsInvoke() {
		if builtin, ok := call.Call.Value.(*ssa.Builtin); ok {
			if ok, funcType, _ := utils.SSABuiltinFuncIsDirect(builtin); ok {
				if funcType == utils.FUNC_TYPE_APPEND {
					// NOTE: builtin append() can safely taint its arguments because
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
						propagateTaintNearby(graph, true, argNode.GetValue(), taintInfo, visited, upwards)
					}
				}
			}
		} else if fn, ok := call.Call.Value.(*ssa.Function); ok && fn.Package() != nil {
			switch fn.Package().Pkg.Name() {
			case "strconv":
				switch fn.Name() {
				case "FormatInt", "Itoa", "ParseInt", "ParseFloat":
					// e.g., myval2 := strconv.ParseInt(myval, 10, 64)
					// e.g., myval2 := strconv.FormatInt(val, 10)
					// e.g., myval2 := strconv.Itoa(myval)

					// propagate from call value to argument value
					argNode := graph.GetNodeByName(call.Call.Args[0].Name())
					propagateTaintNearby(graph, true, argNode.GetValue(), taintInfo, visited, upwards)
				default:
					logrus.WithField("call", call.String()).Fatalf("[TAINT NEARBY] [FROM] unexpected call on (strconv) package")
				}
			}
		}
	}
}
