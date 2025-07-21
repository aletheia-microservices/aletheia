package tainter

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/graph"
	"analyzer/pkg/utils"
)

type TaintMode int

const (
	TAINT_MARK_DIRECT TaintMode = iota
	TAINT_CHECK_INDIRECT
)

type TaintInfo struct {
	dbfield string
	path    string
	val     ssa.Value
}

func NewTaintInfo(database string, collection string) TaintInfo {
	return TaintInfo{
		dbfield: database + "." + collection,
		path:    "",
	}
}

func (t TaintInfo) objectFullPath() string {
	return "_obj" + t.path
}

func (t TaintInfo) objectPath() string {
	return t.path
}

func (t TaintInfo) databaseField() string {
	return t.dbfield
}

func (t TaintInfo) updateValue(val ssa.Value) TaintInfo {
	t.val = val
	return t
}

func (t TaintInfo) updatePathPrefix(prefix string) TaintInfo {
	t.path = prefix + t.path
	return t
}

type CheckTaintInfo struct {
	dbFieldsDirect   []string
	dbFieldsIndirect []string
}

func (t *CheckTaintInfo) addDbFieldsDirectIfNotExists(fields []string) {
	for _, field := range fields {
		if !slices.Contains(t.dbFieldsDirect, field) {
			t.dbFieldsDirect = append(t.dbFieldsDirect, field)
		}
	}
}

func (t *CheckTaintInfo) addDbFieldIndirectIfNotExists(field string) {
	if !slices.Contains(t.dbFieldsIndirect, field) {
		t.dbFieldsIndirect = append(t.dbFieldsIndirect, field)
	}
}

func doTaintNode(node *graph.Node, taintInfo TaintInfo, taintMode TaintMode) string {
	switch taintMode {
	case TAINT_MARK_DIRECT:
		// note that objfields/dbfields already have "." before them
		node.AddTaintIfNotExists(taintInfo.objectFullPath(), taintInfo.databaseField())
		return taintInfo.databaseField()
	case TAINT_CHECK_INDIRECT:
		node.AddTaintIfNotExists(taintInfo.objectFullPath(), taintInfo.databaseField()+taintInfo.objectPath())
		return taintInfo.databaseField() + taintInfo.objectPath()
	}
	return ""
}

func doTaintPointerToSets(g *graph.Graph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool) {
	node := g.GetNodeByName(val.Name())
	for _, edge := range g.GetEdgesFromNode(node) {
		if edge.GetType() == graph.EDGE_POINTS_TO {
			if edge.GetPath() != "" {
				// add before
				// note that both edge.path and objfields/dbfields already have "." before them
				taintInfo = taintInfo.updatePathPrefix(edge.GetPath())
			}
			doTaintNode(edge.GetToNode(), taintInfo, TAINT_MARK_DIRECT)

			doTaintBackwards(g, edge.GetToNode().GetValue(), taintInfo, visited, TAINT_MARK_DIRECT, nil)
		}
	}
}

func doTaintBackwards(g *graph.Graph, val ssa.Value, taintInfo TaintInfo, visited map[TaintInfo]bool, taintMode TaintMode, checkTaintInfo *CheckTaintInfo) {
	taintInfo = taintInfo.updateValue(val)

	fmt.Printf("visiting value %s: %s // %v // TAINT INFO = %v\n", val.Name(), val.String(), val, taintInfo)
	if visited[taintInfo] {
		fmt.Printf("\tskipping value %s: %s\n", val.Name(), val.String())
		return
	}
	visited[taintInfo] = true

	node := g.GetNodeByName(val.Name())

	switch taintMode {
	case TAINT_MARK_DIRECT:
		doTaintNode(node, taintInfo, taintMode)
	case TAINT_CHECK_INDIRECT:
		// e.g.,
		// existing taint: 	_obj @ order_db.order
		// potential taint: _obj.Shipping @ order_db.order.Shipping
		// NOTE: in this case, _obj may be a prefix of _obj.Shipping
		for objPath, dbFields := range node.GetTaints() {
			if strings.HasPrefix(taintInfo.objectFullPath(), objPath) {
				checkTaintInfo.addDbFieldsDirectIfNotExists(dbFields)
			}
		}
		// potential taint: _obj.Shipping @ order_db.order.shipping
		// 									^^^^^^^^^^^^^^^^^^^^^^^
		// 									this is the dbfield (in TAINT_CHECK_INDIRECT mode, it is to be set here)
		for _, dbfield := range checkTaintInfo.dbFieldsDirect {
			taintInfoTmp := taintInfo
			taintInfoTmp.dbfield = dbfield
			
			// to keep track of the dbfield we got here
			dbFieldIndirect := doTaintNode(node, taintInfoTmp, taintMode)
			checkTaintInfo.addDbFieldIndirectIfNotExists(dbFieldIndirect)
		}
	}

	switch t := val.(type) {
	case *ssa.MakeInterface:
		doTaintBackwards(g, t.X, taintInfo, visited, taintMode, checkTaintInfo)
	case *ssa.UnOp:
		doTaintBackwards(g, t.X, taintInfo, visited, taintMode, checkTaintInfo)
	case *ssa.Phi:
		// includes values in t.Edges + other nodes pointing to
		for _, edge := range g.GetEdgesFromNode(g.GetNodeByName(t.Name())) {
			// in case it points to an instruction like store we need to fetch the value
			// (in this case, this corresponds to the variable where something is being stored, and NOT the value being stored)
			if edge.GetToNode().GetInstruction() != nil && edge.GetToNode().GetValue() == nil {
				if taintMode == TAINT_MARK_DIRECT {
					doTaintNode(edge.GetToNode(), taintInfo, taintMode)
				}
				for _, edge2 := range g.GetEdgesToNode(edge.GetToNode()) {
					doTaintBackwards(g, edge2.GetFromNode().GetValue(), taintInfo, visited, taintMode, checkTaintInfo)
				}
			}
		}
	case *ssa.FieldAddr:
		fieldName := utils.FieldIndexToName(t)
		fmt.Printf("[INFO] field addr %s, tainting %s\n", fieldName, t.X.String())
		// add after
		taintInfoTmp := taintInfo
		taintInfoTmp = taintInfoTmp.updatePathPrefix("." + fieldName)
		doTaintBackwards(g, t.X, taintInfoTmp, visited, taintMode, checkTaintInfo)
	case *ssa.IndexAddr:
		// add after
		fmt.Printf("[INFO] index addr %s, tainting %s\n", t.Index.String(), t.X.String())
		taintInfoTmp := taintInfo
		taintInfoTmp = taintInfoTmp.updatePathPrefix("[*]")
		doTaintBackwards(g, t.X, taintInfoTmp, visited, taintMode, checkTaintInfo)
	case *ssa.Parameter, *ssa.Alloc:
		if taintMode == TAINT_MARK_DIRECT {
			doTaintNode(node, taintInfo, taintMode)
		}
	default:
		fmt.Printf("[INFO] ignoring value: [%T] %v\n", val, val)
	}

	if taintMode == TAINT_MARK_DIRECT {
		// if its fieldaddr then we use the objfield and dbfield
		// from the parameters and not the updated ones
		doTaintPointerToSets(g, val, taintInfo, visited)
	}
}

func RunTaint(g *graph.Graph) {
	for _, node := range g.GetNodes() {
		if call, _, ok := isDatabaseCall(node.GetInstruction()); ok {
			var database, collection string

			valDatabase := call.Call.Args[2]
			if c, ok := valDatabase.(*ssa.Const); ok {
				database = strings.Trim(c.Value.ExactString(), "\"")
			}

			valCollection := call.Call.Args[3]
			if c, ok := valCollection.(*ssa.Const); ok {
				collection = strings.Trim(c.Value.ExactString(), "\"")
			}

			valDocument := call.Call.Args[4]

			fmt.Printf("[INFO] tainting document %s.%s: %s\n", database, collection, valDocument.String())
			visited := make(map[TaintInfo]bool)
			taintInfo := NewTaintInfo(database, collection)
			doTaintBackwards(g, valDocument, taintInfo, visited, TAINT_MARK_DIRECT, nil)
		}
		if call, _, ok := isServiceCall(node.GetInstruction()); ok {
			for _, arg := range call.Call.Args[2:] {
				fmt.Printf("[INFO] checking taint for service call with arg: %s\n", arg.String())
				visited := make(map[TaintInfo]bool)
				taintInfo := TaintInfo{}
				CheckTaintInfo := &CheckTaintInfo{}
				doTaintBackwards(g, arg, taintInfo, visited, TAINT_CHECK_INDIRECT, CheckTaintInfo)
				fmt.Printf("GOT DB FIELDS DIRECT: %v\n", CheckTaintInfo.dbFieldsDirect)
				fmt.Printf("GOT DB FIELDS INDIRECT: %v\n", CheckTaintInfo.dbFieldsIndirect)

				node := g.GetNodeByName(arg.Name())
				for _, dbfield := range CheckTaintInfo.dbFieldsIndirect {
					taintInfo := TaintInfo{
						dbfield: dbfield,
						val:     arg,
					}
					doTaintNode(node, taintInfo, TAINT_MARK_DIRECT)
				}
			}
		}
	}
}

func isDatabaseCall(instr ssa.Instruction) (*ssa.Call, *ssa.Function, bool) {
	if call, ok := instr.(*ssa.Call); ok {
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			maybeRcv := fn.Params[0]
			if maybeRcv.Type().String() == "*main.MongoDB" && fn.Name() == "Insert" || fn.Name() == "Find" {
				return call, fn, true
			}
		}
	}
	return nil, nil, false
}

func isServiceCall(instr ssa.Instruction) (*ssa.Call, *ssa.Function, bool) {
	if call, ok := instr.(*ssa.Call); ok {
		if fn, ok := call.Call.Value.(*ssa.Function); ok && len(fn.Params) > 0 {
			maybeRcv := fn.Params[0]
			if maybeRcv.Type().String() == "*main.ShippingService" && fn.Name() == "NewShipment" {
				return call, fn, true
			}
		}
	}
	return nil, nil, false
}
