package controlflow

import (
	"go/ast"

	"github.com/golang-collections/collections/stack"
	"golang.org/x/tools/go/cfg"

	"analyzer/pkg/logger"
	"analyzer/pkg/lookup"
	"analyzer/pkg/service"
	"analyzer/pkg/types"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

type ControlflowContext struct {
	pkg     *types.Package
	service *service.Service
	cfg     *stack.Stack
}

func NewControlflowContext(pkg *types.Package, service *service.Service, initialCfg *types.CFG) *ControlflowContext {
	ctx := &ControlflowContext{
		pkg:     pkg,
		service: service,
		cfg:     stack.New(),
	}
	ctx.PushCFG(initialCfg)
	return ctx
}

func (ctx *ControlflowContext) String() string {
	out := "{pkg: "
	if ctx.pkg != nil {
		out += ctx.pkg.GetName()
	}
	out += ", svc: "
	if ctx.service != nil {
		out += ctx.service.GetName()
	}
	out += "}"
	return out
}

func (ctx *ControlflowContext) GetPackage() *types.Package {
	if ctx.pkg == nil {
		logger.Logger.Fatalf("[CONTROLFLOW CONTEXT] package not found for controlflow context: %s", ctx.String())
	}
	return ctx.pkg
}

func (ctx *ControlflowContext) GetService() *service.Service {
	if ctx.service == nil {
 		logger.Logger.Fatalf("[CONTROLFLOW CONTEXT] service not found for controlflow context: %s", ctx.String())
	}
	return ctx.service
}

func (ctx *ControlflowContext) CurrentCFG() *types.CFG {
	return ctx.cfg.Peek().(*types.CFG)
}

func (ctx *ControlflowContext) PushCFG(new *types.CFG) {
	ctx.cfg.Push(new)
}

func (ctx *ControlflowContext) PopCFG() *types.CFG {
	return ctx.cfg.Pop().(*types.CFG)
}

// if func is anonymous then name is empty
func GenerateInlineFuncCFG(inlineBlock *ast.BlockStmt, name string) *types.CFG {
	cfg := cfg.New(inlineBlock, mayReturn)
	parsedCfg := types.InitParsedCFG(cfg, name)
	return parsedCfg
}

func GenerateMethodCFG(parsedMethod *types.ParsedMethod) {
	cfg := cfg.New(parsedMethod.GetBody(), mayReturn)
	parsedCfg := types.InitParsedCFG(cfg, parsedMethod.Name)
	parsedMethod.SetParsedCFG(parsedCfg)
	entryBlock := parsedCfg.GetEntryParsedBlock()

	receiver := parsedMethod.GetReceiverIfExists()
	if receiver != nil {
		receiver := lookup.CreateObjectFromType(parsedMethod.GetReceiver().GetName(), parsedMethod.GetReceiver().GetType())
		entryBlock.AddObject(receiver)
		parsedCfg.HasReceiver = true
		parsedCfg.ReceiverType = receiver.GetType()
		logger.Logger.Tracef("[CFG] added receiver (%s) %s", objects.VariableTypeName(receiver), receiver.String())
	}

	for i, param := range parsedMethod.Params {
		v := lookup.CreateObjectFromType(param.GetName(), param.GetType())
		v.GetVariableInfo().IsBlockParam = true
		v.GetVariableInfo().BlockParamIdx = i
		entryBlock.Objs = append(entryBlock.Objs, v)
	}

	for _, ret := range parsedMethod.Returns {
		// return values can also be declared by providing a name when declaring the function
		if ret.GetName() != "" {
			v := lookup.CreateObjectFromType(ret.GetName(), ret.GetType())
			entryBlock.Objs = append(entryBlock.Objs, v)
		}
	}

	// note that parameters also include receiver
	logger.Logger.Infof("[CFG] parsed CFG with receiver (%v) and (%d) initial variables for method (%s)", receiver, len(entryBlock.Objs), parsedMethod.String())
}

func InitServiceReceiverFieldsForParsedCFG(service *service.Service, parsedMethod *types.ParsedMethod) {
	logger.Logger.Debugf("[CFG] [%s] updating service receiver for method: %s", service.GetName(), parsedMethod.String())

	parsedCfg := parsedMethod.GetParsedCfg()
	receiver := parsedCfg.GetEntryParsedBlock().GetReceiver()
	implVariable := service.GetImplVariable()

	if ptrReceiver, ok := receiver.(*objects.PointerObject); ok {
		ptrReceiver.PointerTo = implVariable
	} else {
		logger.Logger.Fatalf("[CFG] TODO!!!!")
	}

	logger.Logger.Warnf("[CFG] [%s] updated service receiver for method (%s):\n\t\t\t\t\t - %s", service.GetName(), parsedMethod.GetName(), implVariable.LongString())

}

func GenerateMethodCFGForService(service *service.Service, parsedMethod *types.ParsedMethod) {
	cfg := cfg.New(parsedMethod.GetBody(), mayReturn)
	parsedCfg := types.InitParsedCFG(cfg, parsedMethod.Name)
	parsedMethod.SetParsedCFG(parsedCfg)
	entryBlock := parsedCfg.GetEntryParsedBlock()

	receiver := lookup.CreateObjectFromType(parsedMethod.GetReceiver().GetName(), parsedMethod.GetReceiver().GetType())
	entryBlock.AddObject(receiver)

	variable := receiver
	if pointerVar, ok := variable.(*objects.PointerObject); ok {
		variable = pointerVar.PointerTo
	}
	if structVar, ok := variable.(*objects.StructObject); ok {
		for name, f := range service.Fields {
			structVar.SetFieldByKey(name, lookup.CreateObjectFromType(name, f.GetType()))
		}
	}
	logger.Logger.Tracef("[CFG] added service receiver %s (%s) (%s)", receiver.String(), utils.GetType(receiver.(*objects.PointerObject).PointerTo), utils.GetType(receiver.(*objects.PointerObject).PointerTo.GetType()))

	for i, param := range parsedMethod.Params {
		v := lookup.CreateObjectFromType(param.GetName(), param.GetType())
		v.GetVariableInfo().IsBlockParam = true
		v.GetVariableInfo().BlockParamIdx = i
		entryBlock.Objs = append(entryBlock.Objs, v)
	}

	// note that parameters also include receiver
	logger.Logger.Infof("[CFG] parsed CFG with (%d) initial variables for method (%s) in service (%s)", len(entryBlock.Objs), parsedCfg.FullMethod, service.GetName())
}

// https://github.com/coder/go-tools/blob/master/go/analysis/passes/ctrlflow/ctrlflow_test.go
func mayReturn(call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return fun.Name != "panic"
	case *ast.SelectorExpr:
		return fun.Sel.Name != "Fatal"
	}
	return true
}
