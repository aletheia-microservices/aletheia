package ssagraph

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
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

	readKey   bool // aka filter key
	readValue bool // aka retrived value

	callerT string // originated at combiner.go
}

func NewSSATaintDB(dbpath string, dbcall *DatabaseCall, readKey bool, readValue bool, callerT string) *SSATaint {
	return &SSATaint{
		taintType: TAINT_DATABASE,
		dbpath:    dbpath,
		dbcall:    dbcall,
		readKey:   readKey,
		readValue: readValue,
		callerT:   callerT,
	}
}

func NewSSATaintSV(svpath string, svcall *ServiceCall, callerT string) *SSATaint {
	return &SSATaint{
		taintType: TAINT_SERVICE,
		svpath:    svpath,
		svcall:    svcall,
		callerT:   callerT,
	}
}

func (taint *SSATaint) SimpleCopy() *SSATaint {
	return &SSATaint{
		taintType: taint.taintType,
		dbpath:    taint.dbpath,
		svpath:    taint.svpath,
		readKey:   taint.readKey,
		readValue: taint.readValue,
		callerT:   taint.callerT,
	}
}

func (taint *SSATaint) SetCallerT(callerT string) {
	if taint.callerT != "" {
		logrus.Fatalf("callerT already existis for taint (existing_callerT=%s) (new_callerT=%s) (taint=%s)", taint.callerT, callerT, taint.String())
	}

	taint.callerT = callerT
}

func (taint *SSATaint) GetCallerT() string {
	return taint.callerT
}

func (taint *SSATaint) GetT() string {
	var prefix string
	if taint.callerT != "" {
		prefix = taint.callerT + "."
	}
	if taint.IsDatabaseTaint() {
		return prefix + taint.dbcall.GetT()
	}
	return prefix + taint.svcall.GetT()
}

func (taint *SSATaint) IsDatabaseTaint() bool {
	return taint.taintType == TAINT_DATABASE
}

func (taint *SSATaint) IsServiceTaint() bool {
	return taint.taintType == TAINT_SERVICE
}

func (taint *SSATaint) IsReadKey() bool {
	return taint.readKey
}

func (taint *SSATaint) IsReadValue() bool {
	return taint.readValue
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
	id         string
	name       string
	val        ssa.Value
	instr      ssa.Instruction
	inDefs     bool
	usedInBson bool

	// maps object to database field, e.g.:
	// key: Product    // SSATaint.dbfield: prod_db.Product
	// key: Product.ID // SSATaint.dbfield: prod_db.Product.ID
	// key: Product.ID // SSATaint.dbfield: sku_db.Sku.ProductID
	taints map[string][]*SSATaint
}

func RegisterNewNodeVal(graph *SSAGraph, instr ssa.Instruction, val ssa.Value, id string) *SSANode {
	node := &SSANode{
		name:   val.Name(),
		val:    val,
		instr:  instr,
		inDefs: true,
		id:     id,
		taints: make(map[string][]*SSATaint),
	}
	graph.AddNode(node)
	graph.AddNodeDef(node)
	return node
}

func RegisterNewNodeInstr(graph *SSAGraph, instr ssa.Instruction, id string) *SSANode {
	node := &SSANode{
		id:     id,
		instr:  instr,
		taints: make(map[string][]*SSATaint),
	}
	graph.AddNode(node)
	graph.nodes = append(graph.nodes, node)
	return node
}

func (node *SSANode) SimpleCopy() *SSANode {
	return &SSANode{
		id:         node.id,
		name:       node.name,
		val:        node.val,
		instr:      node.instr,
		inDefs:     node.inDefs,
		usedInBson: node.usedInBson,
		taints:     make(map[string][]*SSATaint),
	}
}

func (node *SSANode) CombineTaints(new map[string][]*SSATaint) {
	for newPath, newTaintsLst := range new {
		if taintLst, ok := node.taints[newPath]; ok {
			node.taints[newPath] = append(taintLst, newTaintsLst...)
		} else {
			node.taints[newPath] = newTaintsLst
		}
	}
}

func (node *SSANode) EnableUsedInBson() {
	node.usedInBson = true
}

func (node *SSANode) IsUsedInBson() bool {
	return node.usedInBson
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

func (node *SSANode) GetInstructionMapUpdate() *ssa.MapUpdate {
	mapupdate, ok := node.instr.(*ssa.MapUpdate)
	if !ok {
		log.Panicf("[SSA NODE] unexpected type for node value: [%T] %v\n", node.val, node.val)
	}
	return mapupdate
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

func (node *SSANode) AddDatabaseTaintIfNotExists(objpath string, dbpath string, dbcall *DatabaseCall, readKey bool, readVal bool, callerT string) bool {
	lstTaints := node.taints[objpath]
	for _, taint := range lstTaints {
		if taint.dbpath == dbpath && taint.dbcall.opType == dbcall.opType {
			return false // already exists
		}
	}
	taint := NewSSATaintDB(dbpath, dbcall, readKey, readVal, callerT)
	// EVAL: logrus.Tracef("added new taint: %s\n", taint.String())
	node.taints[objpath] = append(lstTaints, taint)
	return true
}

func (node *SSANode) AddServiceTaintIfNotExists(objpath string, svpath string, svcall *ServiceCall, callerT string) bool {
	lstTaints := node.taints[objpath]
	for _, taint := range lstTaints {
		if taint.svpath == svpath {
			return false // already exists
		}
	}
	node.taints[objpath] = append(lstTaints, NewSSATaintSV(svpath, svcall, callerT))
	return true
}

// same logic as AbstractGraph Object
func (node *SSANode) TaintAndTraceString() string {
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

			if taint.IsReadKey() {
				builder.WriteString(" [K]")
			} else if taint.IsReadValue() {
				builder.WriteString(" [V]")
			}

			builder.WriteString(fmt.Sprintf(" [%s]", taint.GetT()))

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

func (node *SSANode) LabelsString() string {
	var lbls []string
	var ssaTypeStr string
	if node.val != nil {
		ssaTypeStr = fmt.Sprintf("%T", node.val)
	} else {
		ssaTypeStr = fmt.Sprintf("%T", node.instr)
	}
	ssaTypeStr, _ = strings.CutPrefix(ssaTypeStr, "*ssa.")
	ssaTypeStr = strings.ToLower(ssaTypeStr)
	lbls = append(lbls, fmt.Sprintf("[ssa: %s]", ssaTypeStr))
	if node.IsUsedInBson() {
		lbls = append(lbls, "[bson]")
	}
	return strings.Join(lbls, " ")
}
