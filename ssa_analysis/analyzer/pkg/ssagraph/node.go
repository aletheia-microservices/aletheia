package ssagraph

import (
	"fmt"

	"golang.org/x/tools/go/ssa"
)

type SSATaint struct {
	dbfield string
	dbcall  *DatabaseCall
}

func (taint *SSATaint) GetDbField() string {
	return taint.dbfield
}

func (taint *SSATaint) GetDbCall() *DatabaseCall {
	return taint.dbcall
}

func (taint *SSATaint) String() string {
	return taint.dbfield
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

func (node *SSANode) AddTaintIfNotExists(objPath string, dbField string, dbCall *DatabaseCall) bool {
	lstTaints := node.taints[objPath]
	for _, taint := range lstTaints {
		if taint.dbfield == dbField {
			return false // already exists
		}
	}
	node.taints[objPath] = append(lstTaints, &SSATaint{
		dbfield: dbField,
		dbcall:  dbCall,
	})
	return true
}

func (node *SSANode) taintString() string {
	if len(node.taints) == 0 {
		return ""
	}
	var taintStr string
	for obj, taints := range node.taints {
		taintStr += fmt.Sprintf("\n%s\n", obj)
		for _, taint := range taints {
			taintStr += fmt.Sprintf("@ %s\n", taint.String())
		}
	}
	return taintStr
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
		name:    val.Name(),
		val:     val,
		instr:   instr,
		isdef:   true,
		id:      id,
		taints: make(map[string][]*SSATaint),
	}
	graph.nodes = append(graph.nodes, node)
	graph.defs[node.name] = node
	return node
}

func RegisterNewNode(graph *SSAGraph, instr ssa.Instruction, id string) *SSANode {
	node := &SSANode{
		id:      id,
		instr:   instr,
		taints: make(map[string][]*SSATaint),
	}
	graph.nodes = append(graph.nodes, node)
	return node
}
