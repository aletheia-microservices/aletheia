package tainter

import (
	"go/token"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
)

func propagateTaintNearbyFromNodeOnBinOp(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, toNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, upwards bool) {
	// e.g. t89 = t85 * t88
	// - allow propagation of taints from binOp.X to binOp
	binOp := edge.GetToNode().GetValue().(*ssa.BinOp)
	if binOp.Op >= token.EQL && binOp.Op <= token.GEQ {
		// skip if conditions
	} else if binOp.Op >= token.ADD && binOp.Op <= token.REM {
		// allow operations and delimiters
		propagateTaintNearby(graph, true, toNode.GetValue(), taintInfo, visited, upwards)
	} else {
		logrus.WithField("bin_op", binOp.String()).Fatalf("[TAINT NEARBY] [FROM] [BINOP] to implement")
	}

}

func propagateTaintNearbyToNodeOnBinOp(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, fromNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, upwards bool) {
	// e.g. t89 = t85 * t88
	// skip
	// - do not propagate taints from binOp to binOp.X
	// - only allow the inverse
	return
}
