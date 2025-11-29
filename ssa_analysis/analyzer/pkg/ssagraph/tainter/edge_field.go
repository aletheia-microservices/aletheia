package tainter

import (
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
)

func propagateTaintNearbyFromNodeOnField(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, toNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
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
		return
	}

	var taintInfoTmp TaintInfo
	if taintInfo.isObjectRoot() {
		taintInfoTmp = taintInfo.updateCallPathSuffix("." + edge.GetParam())
	} else if edge.GetParam() == taintInfo.objpath { // FIXME: NOT SURE ABOUT THIS!
		taintInfoTmp = taintInfo.enableObjectRoot()
		taintInfoTmp.dbTaint.dbpath, _ = strings.CutSuffix(taintInfoTmp.objpath, taintInfo.objpath)
	} else {
		logrus.Tracef("[TAINT NEARBY] [PART_1] [FIELD] unexpected conditions // FROM_NODE=%s // TAINT INFO (_obj%s, %s)\n", toNode.String(), taintInfo.getObjectPath(), taintInfo.getDatabasePath())
		return
	}
	propagateTaintNearby(graph, true, toNode.GetValue(), taintInfoTmp, visited, checkTaintInfo, upwards)
}

func propagateTaintNearbyToNodeOnField(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, fromNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	var taintInfoTmp TaintInfo
	if !taintInfo.isObjectRoot() && edge.GetParam() == taintInfo.objpath {
		taintInfoTmp = taintInfo.enableObjectRoot()
	} else {
		taintInfoTmp = taintInfo.updateObjectPathPrefix("." + edge.GetParam())
	}
	propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, true)
}
