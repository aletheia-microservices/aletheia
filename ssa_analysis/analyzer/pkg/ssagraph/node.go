package ssagraph

import (
	"log"
	"sort"
	"strings"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/common"
)

type TaintType int

const (
	TAINT_DATABASE TaintType = iota
	TAINT_SERVICE
)

type SSATaint struct {
	taintType TaintType
	dbpath    string
	dbcall    *DatabaseCall
	svpath    string
	svcall    *ServiceCall
}

func (taint *SSATaint) IsDatabaseTaint() bool {
	return taint.taintType == TAINT_DATABASE
}

func (taint *SSATaint) IsServiceTaint() bool {
	return taint.taintType == TAINT_SERVICE
}

func (taint *SSATaint) GetDatabasePath() string {
	return taint.dbpath
}

func (taint *SSATaint) GetDatabaseCall() *DatabaseCall {
	return taint.dbcall
}

func (taint *SSATaint) GetServicePath() string {
	return taint.svpath
}

func (taint *SSATaint) GetServiceCall() *ServiceCall {
	return taint.svcall
}

func (taint *SSATaint) String() string {
	if taint.taintType == TAINT_DATABASE {
		return taint.dbpath
	}
	return taint.svpath
}

type SSANode struct {
	id    string
	name  string
	val   ssa.Value
	instr ssa.Instruction
	isdef bool

	// maps object to database field, e.g.:
	// key: Product    // SSATaint.dbfield: prod_db.Product
	// key: Product.ID // SSATaint.dbfield: prod_db.Product.ID
	// key: Product.ID // SSATaint.dbfield: sku_db.Sku.ProductID
	taints map[string][]*SSATaint
}

func (node *SSANode) GetID() string {
	return node.id
}

func (node *SSANode) GetName() string {
	return node.name
}

func (node *SSANode) GetInstruction() ssa.Instruction {
	return node.instr
}

func (node *SSANode) GetValue() ssa.Value {
	return node.val
}

func (node *SSANode) GetValueLookup() *ssa.Lookup {
	lookup, ok := node.val.(*ssa.Lookup)
	if !ok {
		log.Panicf("[SSA NODE] unexpected type for node value: [%T] %v\n", node.val, node.val)
	}
	return lookup
}

func (node *SSANode) IsTainted() bool {
	return len(node.taints) > 0
}

func (node *SSANode) GetTaintsForPath(path string) []*SSATaint {
	if t, ok := node.taints[path]; ok { // avoid creating new key entry
		return t
	}
	return nil
}

func (node *SSANode) GetTaints() map[string][]*SSATaint {
	return node.taints
}

func (node *SSANode) AddDatabaseTaintIfNotExists(objpath string, dbpath string, dbcall *DatabaseCall) bool {
	lstTaints := node.taints[objpath]
	for _, taint := range lstTaints {
		if taint.dbpath == dbpath && taint.dbcall.opType == dbcall.opType {
			return false // already exists
		}
	}
	node.taints[objpath] = append(lstTaints, &SSATaint{
		taintType: TAINT_DATABASE,
		dbpath:    dbpath,
		dbcall:    dbcall,
	})
	return true
}

func (node *SSANode) AddServiceTaintIfNotExists(objpath string, svpath string, svcall *ServiceCall) bool {
	lstTaints := node.taints[objpath]
	for _, taint := range lstTaints {
		if taint.svpath == svpath {
			return false // already exists
		}
	}
	node.taints[objpath] = append(lstTaints, &SSATaint{
		taintType: TAINT_SERVICE,
		svpath:    svpath,
		svcall:    svcall,
	})
	return true
}

// same logic as AbstractGraph Object
func (node *SSANode) taintString() string {
	if len(node.taints) == 0 {
		return ""
	}

	var objpaths []string
	for objpath := range node.taints {
		objpaths = append(objpaths, objpath)
	}
	sort.Strings(objpaths)

	var builder strings.Builder
	for _, objpath := range objpaths {
		taints := node.taints[objpath]
		builder.WriteString(objpath)
		builder.WriteByte('\n')
		for _, taint := range taints {
			builder.WriteString("[")
			if taint.taintType == TAINT_DATABASE {
				builder.WriteString(common.OperationTypeToString(taint.GetDatabaseCall().GetOpType()))
			} else if taint.taintType == TAINT_SERVICE {
				builder.WriteString("rpc")
			}
			builder.WriteString("]")

			builder.WriteString(" @ ")
			builder.WriteString(taint.String())
			builder.WriteByte('\n')
		}
	}
	return builder.String()
}

func (node *SSANode) String() string {
	if node.val != nil {
		return node.name + ": " + node.val.String()
	}
	return node.instr.String()
}

func (node *SSANode) colorForSSA() string {
	switch node.instr.(type) {
	case *ssa.Store:
		return "blue"
	case *ssa.Alloc:
		return "orange"
	case *ssa.Return:
		return "yellow"
	case *ssa.Call:
		return "yellow"
	case *ssa.UnOp:
		return "red"
	case *ssa.FieldAddr, *ssa.IndexAddr:
		return "green"
	}
	return "black"
}

func RegisterNewNodeValue(graph *SSAGraph, instr ssa.Instruction, val ssa.Value, id string) *SSANode {
	node := &SSANode{
		name:   val.Name(),
		val:    val,
		instr:  instr,
		isdef:  true,
		id:     id,
		taints: make(map[string][]*SSATaint),
	}
	graph.AddNode(node)
	graph.AddNodeDef(node)
	return node
}

func RegisterNewNode(graph *SSAGraph, instr ssa.Instruction, id string) *SSANode {
	node := &SSANode{
		id:     id,
		instr:  instr,
		taints: make(map[string][]*SSATaint),
	}
	graph.AddNode(node)
	graph.nodes = append(graph.nodes, node)
	return node
}
