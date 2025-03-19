package controlflow

import (
	"go/ast"

	"golang.org/x/tools/go/cfg"

	"analyzer/pkg/logger"
	"analyzer/pkg/lookup"
	"analyzer/pkg/service"
	"analyzer/pkg/types"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

// if func is anonymous than name is empty
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
		receiver := lookup.CreateObjectFromType(parsedMethod.Receiver.GetName(), parsedMethod.Receiver.GetType())
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

	/* if parsedMethod.GetName() == "ReadPostMedia" {
		logger.Logger.Fatal("EXIT!")
	} */

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

	receiver := lookup.CreateObjectFromType(parsedMethod.Receiver.GetName(), parsedMethod.Receiver.GetType())
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
