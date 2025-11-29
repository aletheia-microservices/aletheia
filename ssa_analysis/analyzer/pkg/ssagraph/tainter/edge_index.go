package tainter

import (
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
)

func propagateTaintNearbyFromNodeOnIndex(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, toNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
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
		return
	}
	var taintInfoTmp TaintInfo
	if taintInfo.isObjectRoot() {
		taintInfoTmp = taintInfo.updateCallPathSuffix("[" + edge.GetParam() + "]")
	} else {
		var ok bool
		taintInfoTmp = taintInfo.enableObjectRoot()
		taintInfoTmp.objpath, ok = strings.CutSuffix(taintInfoTmp.objpath, "[*]")
		if !ok {
			logrus.Fatalf("[TAINT NEARBY] [PART_1] [INDEX] could not cut suffix [*] for objpath = (%s)\n", taintInfoTmp.objpath)
		}
	}

	propagateTaintNearby(graph, true, toNode.GetValue(), taintInfoTmp, visited, checkTaintInfo, upwards)
}

func propagateTaintNearbyToNodeOnIndex(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, fromNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	var taintInfoTmp TaintInfo
	if taintInfo.isObjectRoot() {
		taintInfoTmp = taintInfo.updateObjectPathPrefix("[" + edge.GetParam() + "]")
	} else {
		taintInfoTmp = taintInfo.enableObjectRoot()
	}

	propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, true)
}
