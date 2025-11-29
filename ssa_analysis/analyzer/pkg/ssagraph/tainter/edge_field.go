package tainter

import (
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
)

func propagateTaintNearbyFromNodeOnField(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, toNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, upwards bool) {
	if upwards {
		if taintInfo.isTypeService() || taintInfo.isTypeDatabase() {
			// found field corresponding to upper taintinfo objpath
			if taintInfo.objpath == "."+edge.GetParam() {
				for _, upperTaint := range node.GetTaintsForPath("_obj" + taintInfo.objpath) {
					taintInfoTmp := generateRootTaintInfoFromTaint(toNode, upperTaint)
					// node has taint info if it was the previous node calling propagateTaintNearby
					// (e.g., lower field propagating to upper struct)
					// we need to avoid visiting it again otherwise we will have infinite recursion!
					propagateTaintNearby(graph, true, toNode.GetValue(), taintInfoTmp, make(map[ssa.Value]bool), false)
				}
			}
		}
		return
	}
	var taintInfoTmp TaintInfo
	if taintInfo.isObjectRoot() {
		taintInfoTmp = taintInfo.updateCallPathSuffix("." + edge.GetParam())
	} else if edge.GetParam() == taintInfo.objpath {
		taintInfoTmp = taintInfo.enableObjectRoot()
		taintInfoTmp.dbTaint.dbpath, _ = strings.CutSuffix(taintInfoTmp.objpath, taintInfo.objpath)
	} else {
		logrus.WithField("graph", graph.String()).WithField("curr/from", node.String()).WithField("to", toNode.String()).
		WithField("taint_info", taintInfo.String()).WithField("taint_info_tmp", taintInfoTmp.String()).
		Fatalf("[TAINT NEARBY] [FROM] [FIELD] unexpected conditions")
	}
	propagateTaintNearby(graph, true, toNode.GetValue(), taintInfoTmp, visited, upwards)
}

func propagateTaintNearbyToNodeOnField(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, fromNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, upwards bool) {
	var taintInfoTmp TaintInfo
	if !taintInfo.isObjectRoot() && edge.GetParam() == taintInfo.objpath {
		taintInfoTmp = taintInfo.enableObjectRoot()
	} else {
		taintInfoTmp = taintInfo.updateObjectPathPrefix("." + edge.GetParam())
	}
	propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfoTmp, make(map[ssa.Value]bool), true)
}
