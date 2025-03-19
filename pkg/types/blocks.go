package types

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/cfg"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

type InlineFunc struct {
	Name string
	CFG  *CFG
}

type Block struct {
	Block      *cfg.Block
	Index      int
	Successors []*Block

	Objs        []objects.Object
	Results     []objects.Object
	InlineFuncs []*InlineFunc // declared inline blocks (NOT ANONYMOUS)
	Info        *BlockInfo
	Visited     bool
}

func (block *Block) AddInlineFunc(name string, cfg *CFG) {
	block.InlineFuncs = append(block.InlineFuncs, &InlineFunc{
		Name: name,
		CFG:  cfg,
	})
	logger.Logger.Infof("[BLOCK] added inline func (%s) for cfg: %s", name, block.String())
}

func (block *Block) GetLatestInlineFunc(name string) *InlineFunc {
	for i := 0; i < len(block.InlineFuncs); i++ {
		if block.InlineFuncs[i].Name == name {
			return block.InlineFuncs[i]
		}
	}
	logger.Logger.Warnf("[BLOCK] could not find latest inline func (%s) in list: %v", name, block.InlineFuncs)
	return nil
}

func (block *Block) HasVariable(name string) bool {
	for i := len(block.Objs) - 1; i >= 0; i-- {
		if block.Objs[i].GetVariableInfo().Name == name {
			return true
		}
	}
	return false
}

func (block *Block) GetObject(name string) objects.Object {
	for i := len(block.Objs) - 1; i >= 0; i-- {
		if block.Objs[i].GetVariableInfo().Name == name {
			return block.Objs[i]
		}
	}
	logger.Logger.Fatalf("[BLOCK] object with name (%s) not found: %v", name, block.Objs)
	return nil
}

func (block *Block) GetVars() []objects.Object {
	return block.Objs
}

func (block *Block) VarsString() string {
	s := "variables = ["
	for i, v := range block.Objs {
		s += v.String()
		if i < len(block.Objs)-1 {
			s += ", "
		}

	}
	return s + "]"
}

func (block *Block) ListObjectsString() string {
	s := fmt.Sprintf("[#%d, %s] BLOCK:\n", block.GetIndex(), block.Block.Kind)
	for i, v := range block.Objs {
		deps := v.GetDependencies()
		depsStr := "None"
		if len(deps) > 0 {
			depStrings := make([]string, len(deps))
			for j, d := range deps {
				depStrings[j] = fmt.Sprintf("%v", d)
			}
			depsStr = strings.Join(depStrings, ", ")
		}
		
		s += fmt.Sprintf("\t [%d] %s ---> depends on: %v", i, v.String(), depsStr)
		if i < len(block.Objs)-1 {
			s += "\n"
		}
	}
	return s
}

func (block *Block) DeepCopy() *Block {
	var copyVars []objects.Object
	for _, v := range block.Objs {
		copyVars = append(copyVars, v.DeepCopy())
	}
	return &Block{
		Block: block.Block,
		Objs:  copyVars,
		Info:  block.Info,
	}
}

func (block *Block) AddResult(variable objects.Object) {
	logger.Logger.Debugf("[BLOCK (%d)] added new result: \n\t\t\t\t\t\t - (%s) %s", block.Block.Index, objects.VariableTypeName(variable), variable.LongString())
	block.Results = append(block.Results, variable)
}

func (block *Block) GetVariableAt(index int) objects.Object {
	if index < len(block.Objs) {
		return block.Objs[index]
	}
	logger.Logger.Fatalf("[BLOCK] received index (%d) but block vars (%s) is length (%d)", index, block.String(), len(block.Objs))
	return nil
}

func (block *Block) LongString() string {
	str := fmt.Sprintf("Block %d", block.Block.Index)
	for i, v := range block.Block.Nodes {
		str += fmt.Sprintf("\t Node %d: %s", i, utils.GetType(v))
	}
	for _, v := range block.Block.Succs {
		str += fmt.Sprintf("\t Succ: %d", v.Index)
	}
	str += "\n"
	return str
}

type BlockInfo struct {
	Gen []objects.Object
	In  []objects.Object
	Out []objects.Object
}

func (block *Block) GetVariables() []objects.Object {
	return block.Objs
}

func (block *Block) GetResults() []objects.Object {
	return block.Results
}

func (block *Block) GetFirstResult() objects.Object {
	if len(block.Results) == 0 {
		logger.Logger.Fatalf("[BLOCKS (%d)] block has no results: %s", block.Block.Index, block.String())
	}
	return block.Results[0]
}

func (block *Block) GetReceiver() objects.Object {
	if len(block.Objs) == 0 {
		logger.Logger.Fatalf("[BLOCKS] unexpected empty CFG block (%v)", block)
	}
	return block.Objs[0]
}

func (block *Block) UpdateReceiver(v objects.Object) {
	if len(block.Objs) == 0 {
		logger.Logger.Fatalf("[BLOCKS] unexpected empty CFG block (%v)", block)
	}
	block.Objs[0] = v
}

func (block *Block) IsReceiverVariable(v objects.Object) bool {
	if len(block.Objs) > 0 {
		return block.Objs[0] == v
	}
	return false
}

func (block *Block) GetLatestObjectByName(name string) objects.Object {
	for i := len(block.Objs) - 1; i >= 0; i-- {
		v := block.Objs[i]
		if name == v.GetVariableInfo().GetName() {
			return v
		}
	}
	logger.Logger.Fatalf("[BLOCK] variable (%s) not found in block: %s)", name, block.VarsString())
	return nil
}

func (block *Block) GetLastestVariableIfExists(name string) objects.Object {
	for i := len(block.Objs) - 1; i >= 0; i-- {
		v := block.Objs[i]
		if name == v.GetVariableInfo().GetName() {
			return v
		}
	}
	logger.Logger.Debugf("[BLOCK] variable (%s) does not exist in block with vars lst: %v", name, block.Objs)
	return nil
}

func (block *Block) AddVariables(objects []objects.Object) {
	logger.Logger.Tracef("[BLOCK] added variables to block: %s", objects)
	block.Objs = append(block.Objs, objects...)
}

func (block *Block) AddObject(variable objects.Object) {
	if variable.GetVariableInfo().Id == objects.VARIABLE_INLINE_ID {
		logger.Logger.Warnf("[BLOCK] skipping addition of new inline variable (%s) to block", variable.String())
		return
	}

	logger.Logger.Debugf("[BLOCK (%d)] added new variable: \n\t\t\t\t\t\t - (%s) %s", block.Block.Index, objects.VariableTypeName(variable), variable.LongString())
	block.Objs = append(block.Objs, variable)
}

func (block *Block) GetPosition() token.Pos {
	return block.Block.Stmt.Pos()
}

func (block *Block) GetIndex() int32 {
	return block.Block.Index
}

func (block *Block) GetNodes() []ast.Node {
	return block.Block.Nodes
}

func (block *Block) GetSuccs() []*cfg.Block {
	return block.Block.Succs
}

func (block *Block) GetNextSuccessorIfExists() *Block {
	if len(block.Successors) > 0 {
		return block.Successors[0]
	}
	return nil
}

func (block *Block) AppendVarsFromPredecessor(predecessor *Block) {
	logger.Logger.Tracef("[BLOCK] copying vars %d -> %d: %v -> %v", predecessor.GetIndex(), block.GetIndex(), predecessor.Objs, block.Objs)
	block.Objs = append(predecessor.Objs, block.Objs...)
}

func (block *Block) AppendInlineFuncsFromPredecessor(predecessor *Block) {
	logger.Logger.Tracef("[BLOCK] copying inline funcs %d -> %d: %v -> %v", predecessor.GetIndex(), block.GetIndex(), predecessor.InlineFuncs, block.InlineFuncs)
	block.InlineFuncs = append(predecessor.InlineFuncs, block.InlineFuncs...)
}

func (block *Block) String() string {
	str := fmt.Sprintf("Block %d [%s] (", block.Block.Index, block.Block.Kind)
	for i, v := range block.Objs {
		str += v.GetVariableInfo().GetName()
		if i < len(block.Objs)-1 {
			str += ", "
		}
	}
	str += ")"
	return str
}

func (block *Block) AstInfoString() string {
	var nodesStr []string
	for _, node := range block.Block.Nodes {
		nodesStr = append(nodesStr, utils.GetType(node))
	}
	return fmt.Sprintf("Block %d [%s] (stmt = %s) with %d nodes %v and succs %v", block.Block.Index, block.Block.Kind, utils.GetType(block.Block.Stmt), len(nodesStr), nodesStr, block)
}
