package tainter

import (
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
	"analyzer/pkg/utils"
)

const (
	DYNAMIC_MAP       = ".*"
	DYNAMIC_MAP_KEY   = ".MapKey"
	DYNAMIC_MAP_VALUE = ".MapVal"
)

func propagateTaintNearbyFromNodeOnMap(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, toNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	switch edge.GetType() {
	case ssagraph.EDGE_MAP_UPDATE:
		if upwards {
			// TODO (similar to FIELD and INDEXES but needs more boring logic)
			break
		}

		// e.g., t0[t6] = t1
		// curr/from: t0
		// to: 		  t0[t6] = t1
		// t0 >>> t6
		// t0 >>> t1

		var keyOk, valOk bool
		taintInfoTmpKey := taintInfo.enableObjectRoot()
		taintInfoTmpKey.objpath, keyOk = strings.CutPrefix(taintInfoTmpKey.objpath, DYNAMIC_MAP_KEY+".Key")
		taintInfoTmpVal := taintInfo.enableObjectRoot()
		taintInfoTmpVal.objpath, valOk = strings.CutPrefix(taintInfoTmpKey.objpath, DYNAMIC_MAP_KEY+".Val")

		logrus.WithField("root", taintInfo.isObjectRoot()).
			WithField("ok_key", keyOk).
			WithField("ok_val", valOk).
			WithField("curr/from", node.String()).
			WithField("to", toNode.String()).
			WithField("taint_info", taintInfo.String()).
			WithField("taint_info_tmp", taintInfoTmpVal.String()).
			Infof("[TAINT NEARBY] [FROM] found EDGE_MAP_UPDATE")

		if keyOk {
			// triggers ssagraph.EDGE_MAP_KEY ahead
			mapKey := toNode.GetInstruction().(*ssa.MapUpdate).Key
			propagateTaintNearby(graph, true, mapKey, taintInfoTmpKey, make(map[ssa.Value]bool), checkTaintInfo, upwards)
		}
		if valOk {
			// triggers ssagraph.EDGE_MAP_VAL ahead
			mapVal := toNode.GetInstruction().(*ssa.MapUpdate).Value
			propagateTaintNearby(graph, true, mapVal, taintInfoTmpVal, make(map[ssa.Value]bool), checkTaintInfo, upwards)
		}

	case ssagraph.EDGE_MAP_KEY:
		// e.g., [ssa.MapUpdate] t0[t27] = t7
		// curr/from: t27
		// to:		  t0[t27] = t7 (INSTRUCTION)
		// t27 >>> t10
		mapInstr := toNode.GetInstructionMapUpdate()
		var prefix string
		keyStr, ok := utils.ExtractStringFromValue(mapInstr.Key)
		if !ok {
			prefix = DYNAMIC_MAP_KEY
		} else {
			prefix = "." + keyStr
		}
		prefix_key := prefix + ".Key"
		taintInfoTmp := taintInfo.updateObjectPathPrefix(prefix_key)
		taintInfoTmp = taintInfoTmp.disableObjectRoot()

		logrus.WithField("prefix_key", prefix_key).
			WithField("curr/from", node).
			WithField("to", toNode).
			WithField("taint_info_tmp", taintInfoTmp.String()).
			Infof("[TAINT NEARBY] [FROM] found EDGE_MAP_KEY")
		taintInfoTmp = taintInfoTmp.enableObjectRoot()

		propagateTaintNearby(graph, true, mapInstr.Map, taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, upwards)

	case ssagraph.EDGE_MAP_VALUE:
		// e.g., [ssa.MapUpdate] t0[t27] = t7
		// curr/from: t7
		// to:		  t0[t27] = t7 (INSTRUCTION)
		// t7 >>> t0
		var prefix string
		instr := toNode.GetInstruction().(*ssa.MapUpdate)
		keyStr, ok := utils.ExtractStringFromValue(instr.Key)
		if !ok {
			prefix = DYNAMIC_MAP_KEY // dynamic
			// logrus.Fatalf("TODO: if one key is dynamic, then all keys must be dynamic, even if a subset is static!")
		} else {
			prefix = "." + keyStr
		}
		prefix += ".Val"
		taintInfoTmp := taintInfo.updateObjectPathPrefix(prefix)
		taintInfoTmp = taintInfoTmp.disableObjectRoot()

		logrus.WithField("to", toNode.String()).
			WithField("curr/from", node.String()).
			WithField("prefix", prefix).
			WithField("taint_info", taintInfo.String()).
			WithField("taint_info_tmp", taintInfoTmp.String()).
			Infof("[TAINT NEARBY] [FROM] found EDGE_MAP_VALUE")

		propagateTaintNearby(graph, true, instr.Map, taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, upwards)

	case ssagraph.EDGE_LOOKUP_MAP:
		// e.g., [ssa.Lookup] t12: t0[t11]
		// curr/from: t0
		// to:		  t12
		// t0 >>> t11
		// t0 >>> t12
		lookupIndex := toNode.GetValueLookup().Index
		var prefix string
		keyStr, ok := utils.ExtractStringFromValue(toNode.GetValueLookup().Index)
		if !ok {
			prefix = DYNAMIC_MAP_KEY
		} else {
			prefix = "." + keyStr
		}
		prefix_key := prefix + ".Key"
		taintInfoTmpKey, keyOk := taintInfo.cutObjectPathPrefix(prefix_key)
		if keyOk {
			taintInfoTmpKey = taintInfoTmpKey.enableObjectRoot()
			propagateTaintNearby(graph, true, lookupIndex, taintInfoTmpKey, make(map[ssa.Value]bool), checkTaintInfo, upwards)
		}

		prefix_val := prefix + ".Val"
		taintInfoTmpVal, valOk := taintInfo.cutObjectPathPrefix(prefix_val)
		if valOk {
			taintInfoTmpVal = taintInfoTmpVal.enableObjectRoot()
			propagateTaintNearby(graph, true, toNode.GetValue(), taintInfoTmpVal, make(map[ssa.Value]bool), checkTaintInfo, upwards)
		}

		if !keyOk && !valOk {
			logrus.WithField("prefix_val", prefix_val).WithField("prefix_key", prefix_key).
				WithField("taint_info", taintInfo.String()).
				WithField("taint_info_tmp_key", taintInfoTmpKey.String()).
				WithField("taint_info_tmp_val", taintInfoTmpVal.String()).
				WithField("graph", graph.String()).
				WithField("map", node.String()).
				Warnf("EDGE_LOOKUP_MAP not ok for KEY and VAL!")
		}

	case ssagraph.EDGE_LOOKUP_MAP_INDEX:
		// e.g., [ssa.Lookup] t12: t0[t11]
		// curr/from: t11
		// to:		  t12
		// t11 >>> t0
		lookupTarget := toNode.GetValueLookup().X
		var prefix string
		keyStr, ok := utils.ExtractStringFromValue(node.GetValue())
		if !ok {
			prefix = DYNAMIC_MAP_KEY // dynamic
		} else {
			prefix = "." + keyStr
		}
		prefix += ".Key"
		taintInfoTmp := taintInfo.updateObjectPathPrefix(prefix)
		taintInfoTmp = taintInfoTmp.disableObjectRoot()

		logrus.WithField("to", toNode.String()).
			WithField("curr/from", node.String()).
			WithField("taint_info", taintInfo.String()).
			WithField("taint_info_tmp", taintInfoTmp.String()).
			Infof("[TAINT NEARBY] [FROM] found EDGE_LOOKUP_MAP_INDEX: %s\n", toNode.String())
		propagateTaintNearby(graph, true, lookupTarget, taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, upwards)
	}
}

func propagateTaintNearbyToNodeOnMap(graph *ssagraph.SSAGraph, edge *ssagraph.SSAEdge, node *ssagraph.SSANode, fromNode *ssagraph.SSANode, taintInfo TaintInfo, visited map[ssa.Value]bool, checkTaintInfo *CheckTaintInfo, upwards bool) {
	switch edge.GetType() {
	case ssagraph.EDGE_MAP_UPDATE:
		logrus.WithField("curr/to", node.String()).
			WithField("from", fromNode.String()).
			Fatalf("[TAINT NEARBY] [TO] ignoring EDGE_MAP_UPDATE")

	case ssagraph.EDGE_MAP_KEY:
		logrus.WithField("curr/to", node.String()).
			WithField("from", fromNode.String()).
			Fatalf("[TAINT NEARBY] [TO] ignoring EDGE_MAP_KEY")

	case ssagraph.EDGE_MAP_VALUE:
		// e.g., [ssa.MapUpdate] t0[t27] = t7
		logrus.WithField("curr/to", node.String()).
			WithField("from", fromNode.String()).
			Fatalf("[TAINT NEARBY] [TO] ignoring EDGE_MAP_VALUE")

	case ssagraph.EDGE_LOOKUP_MAP:
		// e.g., t41: t6[t27]
		// curr/to: t41
		// from: 	t6
		// t41 >>> t6

		lookupTarget := node.GetValueLookup().X
		var prefix string
		keyStr, ok := utils.ExtractStringFromValue(node.GetValueLookup().Index)
		if !ok {
			prefix = DYNAMIC_MAP_KEY // dynamic
			// logrus.Fatalf("TODO: if one key is dynamic, then all keys must be dynamic, even if a subset is static!")
		} else {
			prefix = "." + keyStr
		}
		prefix += ".Val"

		logrus.WithField("prefix", prefix).WithField("from", fromNode.String()).WithField("curr/to", node.String()).
			Infof("[TAINT NEARBY] [TO] found EDGE_LOOKUP_MAP (%s: %s)\n", node.GetValueLookup().Name(), node.GetValueLookup().String())

		taintInfoTmp := taintInfo.updateObjectPathPrefix(prefix)
		taintInfoTmp = taintInfoTmp.disableObjectRoot()
		propagateTaintNearby(graph, true, lookupTarget, taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, upwards)

	case ssagraph.EDGE_LOOKUP_MAP_INDEX:
		// e.g., t41: t6[t27]
		// curr/to: t41
		// from: 	t27
		// t41 >>> t6

		lookupTarget := node.GetValueLookup().X
		var prefix string
		keyStr, ok := utils.ExtractStringFromValue(node.GetValueLookup().Index)
		if !ok {
			prefix = DYNAMIC_MAP_KEY // dynamic
			// logrus.Fatalf("TODO: if one key is dynamic, then all keys must be dynamic, even if a subset is static!")
		} else {
			prefix = "." + keyStr
		}
		prefix += ".Val"
		taintInfoTmp := taintInfo.updateObjectPathPrefix(prefix)
		taintInfoTmp = taintInfoTmp.disableObjectRoot()

		logrus.WithField("prefix", prefix).
			WithField("from", fromNode.String()).
			WithField("curr/to", node.String()).
			WithField("taint_info", taintInfo.String()).
			WithField("taint_info_tmp", taintInfoTmp.String()).
			Infof("[TAINT NEARBY] [TO] found EDGE_LOOKUP_MAP_INDEX (%s: %s)\n", node.GetValueLookup().Name(), node.GetValueLookup().String())

		propagateTaintNearby(graph, true, lookupTarget, taintInfoTmp, make(map[ssa.Value]bool), checkTaintInfo, upwards)
	}
}
