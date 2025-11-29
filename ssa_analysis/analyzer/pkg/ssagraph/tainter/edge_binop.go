package tainter

import (
	"go/token"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
)

func propagateTaintNearbyFromNodeOnBinOp(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, toNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	// e.g. t89 = t85 * t88
	// - allow propagation of taints from binOp.X to binOp
	binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
	if binOp.Op >= token.EQL && binOp.Op <= token.GTR || binOp.Op >= token.NEQ && binOp.Op <= token.GEQ {
		// skip (if conditions)
	} else if binOp.Op >= token.ADD && binOp.Op <= token.REM {
		// operations and delimiters
		propagateTaintNearby(graph, true, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
	} else {
		logrus.Fatalf("CONFIRM! binop.Op = %s", binOp.String())
		propagateTaintNearby(graph, true, toNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
	}

}

func propagateTaintNearbyToNodeOnBinOp(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, fromNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	// e.g. t89 = t85 * t88
	// skip
	// - do not propagate taints from binOp to binOp.X
	// - only allow the inverse
	binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
	if binOp.Op >= token.EQL && binOp.Op <= token.GTR || binOp.Op >= token.NEQ && binOp.Op <= token.GEQ {
		// skip (if conditions)
	} else if binOp.Op >= token.ADD && binOp.Op <= token.AND_NOT {
		// skip operations and delimiters
	} else {
		logrus.Fatalf("TO BE REMOVED ONCE REACHED (%s, OP=%s)", binOp.String(), binOp.Op)
		propagateTaintNearby(graph, true, fromNode.GetValue(), taintInfo, visited, checkTaintInfo, upwards)
	}
}
