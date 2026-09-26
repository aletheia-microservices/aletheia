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
	// --- input for abstract call graph ---
	taintType TaintType
	callerTS  string // originated at combiner.go
	path      string // database path for TAINT_DATABASE; service path for TAINT_SERVICE
	call      Call   // *DatabaseCall for TAINT_DATABASE; *ServiceCall for TAINT_SERVICE
	// database-specific info
	readKey   bool // aka filter key
	readValue bool // aka retrived value
}

func newSSATaint(taintType TaintType, path string, call Call, readKey bool, readValue bool, callerT string) *SSATaint {
	return &SSATaint{
		taintType: taintType,
		path:      path,
		call:      call,
		readKey:   readKey,
		readValue: readValue,
		callerTS:  callerT,
	}
}

func (taint *SSATaint) SimpleCopy() *SSATaint {
	return &SSATaint{
		taintType: taint.taintType,
		path:      taint.path,
		readKey:   taint.readKey,
		readValue: taint.readValue,
		callerTS:  taint.callerTS,
	}
}

func (taint *SSATaint) SetCallerT(callerT string) {
	if taint.callerTS != "" {
		logrus.Fatalf("callerT already existis for taint (existing_callerT=%s) (new_callerT=%s) (taint=%s)", taint.callerTS, callerT, taint.String())
	}

	taint.callerTS = callerT
}

func (taint *SSATaint) GetCallerT() string {
	return taint.callerTS
}

func (taint *SSATaint) GetT() string {
	var prefix string
	if taint.callerTS != "" {
		prefix = taint.callerTS + "."
	}
	return prefix + taint.call.GetT()
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

func (taint *SSATaint) GetPath() string {
	return taint.path
}

func (taint *SSATaint) GetCall() Call {
	return taint.call
}

func (taint *SSATaint) GetDatabasePath() string {
	if !taint.IsDatabaseTaint() {
		return ""
	}
	return taint.path
}

func (taint *SSATaint) GetDatabaseCall() *DatabaseCall {
	dbcall, _ := taint.call.(*DatabaseCall)
	return dbcall
}

func (taint *SSATaint) GetServicePath() string {
	if !taint.IsServiceTaint() {
		return ""
	}
	return taint.path
}

func (taint *SSATaint) GetServiceCall() *ServiceCall {
	svcall, _ := taint.call.(*ServiceCall)
	return svcall
}

func (taint *SSATaint) String() string {
	return taint.path
}

type SSANode struct {
	// --- input for abstract call graph ---
	name   string
	taints map[string][]*SSATaint // key format: _obj.<object path>; e.g., _obj.ID

	// --- extra info for SSA graph ---
	id         string //format: val_<SSA value name>; not used anywhere (nor in abstract call graph), but useful for debugging and testing
	val        ssa.Value
	instr      ssa.Instruction
	inDefs     bool
	usedInBson bool
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
		if taint.IsDatabaseTaint() && taint.path == dbpath && taint.GetDatabaseCall().opType == dbcall.opType {
			return false // already exists
		}
	}
	node.taints[objpath] = append(lstTaints, newSSATaint(TAINT_DATABASE, dbpath, dbcall, readKey, readVal, callerT))
	return true
}

func (node *SSANode) AddServiceTaintIfNotExists(objpath string, svpath string, svcall *ServiceCall, callerT string) bool {
	lstTaints := node.taints[objpath]
	for _, taint := range lstTaints {
		if taint.IsServiceTaint() && taint.path == svpath {
			return false // already exists
		}
	}
	node.taints[objpath] = append(lstTaints, newSSATaint(TAINT_SERVICE, svpath, svcall, false, false, callerT))
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
			switch taint.taintType {
			case TAINT_DATABASE:
				builder.WriteString(common.OperationTypeToString(taint.GetDatabaseCall().GetOpType()))
			case TAINT_SERVICE:
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
