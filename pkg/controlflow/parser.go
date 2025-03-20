package controlflow

import (
	"fmt"
	"go/ast"
	"go/token"
	golangtypes "go/types"

	"golang.org/x/tools/go/cfg"

	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/lookup"
	"analyzer/pkg/service"
	"analyzer/pkg/types"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

func ParseServiceMethodCFG(service *service.Service, method *types.ParsedMethod) {
	if method.IsParsed() {
		logger.Logger.Warnf("[CFG] [%s] method ignoring parsed method: %s", service.GetName(), method.String())
		return
	}
	method.SetParsed()

	var blocksStr string
	for _, block := range method.ParsedCfg.GetParsedBlocks() {
		blocksStr += "\t\t\t\t\t - " + block.AstInfoString() + "\n"
	}

	/* if method.GetName() == "StorePostNoSQL" {
		logger.Logger.Fatal("EXIT!")
	} */

	logger.Logger.Debugf("[CFG] [%s] parsing service method cfg for (%s):\n%s", service.GetName(), method.String(), blocksStr)

	entryBlock := method.ParsedCfg.GetEntryParsedBlock()
	visitBasicBlock(service, method, entryBlock)

	/* if method.GetName() == "StorePostNoSQL" {
		logger.Logger.Fatal("EXIT!")
	} */
}

func visitBasicBlock(service *service.Service, method *types.ParsedMethod, block *types.Block) {
	if block.Visited {
		return
	}

	switch block.Block.Stmt.(type) {
	case *ast.BlockStmt:
		logger.Logger.Infof("\n--------------------------------------------------------------------------------------------\nVISITING BlockStmt: %s \n--------------------------------------------------------------------------------------------", block.String())
	case *ast.ForStmt:
		logger.Logger.Infof("\n--------------------------------------------------------------------------------------------\nVISITING ForStmt: %s \n--------------------------------------------------------------------------------------------", block.String())
	case *ast.ReturnStmt:
		logger.Logger.Infof("\n--------------------------------------------------------------------------------------------\nVISITING ReturnStmt: %s \n--------------------------------------------------------------------------------------------", block.String())
	default:
		logger.Logger.Infof("\n--------------------------------------------------------------------------------------------\nVISITING UNKNOWN (%s) BLOCK: %s \n--------------------------------------------------------------------------------------------", utils.GetType(block.Block.Stmt), block.String())

	}

	var deferStmts []*ast.DeferStmt
	block.Visited = true

	var visitedRangeType bool
	var visitedRangeElem bool
	var rangeKeyType gotypes.Type
	var rangeValueType gotypes.Type
	var rangeObj objects.Object

	fmt.Println(utils.TEXT_BOLD_LIGHT_BLUE + "[BEFORE] ---------------------\n" +
		service.GetName() + "." + method.GetName() + "()\n" + block.ListObjectsString() +
		"\n------------------------------" + utils.TEXT_RESET_COLOR)

	for i, node := range block.GetNodes() { //FIXME????
		/* if block.Block.Kind == cfg.KindBody && i == len(block.GetNodes())-1 {
			for _, succ := range block.GetSuccs() {
				// skip last node if we have successors for (i) if branches
				// in these cases, the last node usually corresponds to the conditional expression
				// do not skip for forloops since this is the only way to capture the initial declaration (e.g. for i := 0)
				if succ.Kind == cfg.KindIfThen {
					break
				}
			}
		} */

		initialObjsStr := ""
		for i, obj := range block.Objs {
			initialObjsStr += fmt.Sprintf("\t (#%d) %s", i, obj.String())
			if i < len(block.Objs)-1 {
				initialObjsStr += "\n"
			}
		}

		//logger.Logger.Warnf("\n----------------------------------------------\nPARSING BLOCK [%d] W/ KIND = %s; NODE [%d]: %v \n\t@ METHOD: %s.%s\n%s\n----------------------------------------------", block.Block.Index, block.Block.Kind.String(), i, node, service.Name, method.Name, initialObjsStr)

		var parsingLoop bool
		parsingLoop, visitedRangeType, visitedRangeElem, rangeObj, rangeKeyType, rangeValueType = visitBasicBlockRangeHelper(service, method, block, node, visitedRangeType, visitedRangeElem, rangeObj, rangeKeyType, rangeValueType)

		if !parsingLoop {
			stmts := parseNodeBody(service, method, block, node)
			deferStmts = append(deferStmts, stmts...)
		}

		fmt.Println(utils.TEXT_BOLD_LIGHT_GREEN + "------------------------------\n" +
			service.GetName() + "." + method.GetName() + "()\n" +
			fmt.Sprintf("[#%d, %T] NODE:\n", i, node) +
			block.ListObjectsString() +
			"\n------------------------------" + utils.TEXT_RESET_COLOR)
	}

	for _, deferStmt := range deferStmts {
		parseAndSaveCall(service, method, block, deferStmt.Call)
	}

	if block.Block.Kind == cfg.KindForPost { // FIXME: skip going again for loop
		return
	}

	for _, succ := range block.GetSuccs() {
		parsedSucc := method.ParsedCfg.GetParsedBlockAtIndex(succ.Index)
		parsedSucc.AppendVarsFromPredecessor(block)
		parsedSucc.AppendInlineFuncsFromPredecessor(block)
		logger.Logger.Debugf("\n\n----------------------------------------------\nFOUND BLOCK SUCC [%d -> %d]: %s \n----------------------------------------------", block.Block.Index, parsedSucc.Block.Index, parsedSucc.String())

		visitBasicBlock(service, method, parsedSucc)
	}
}

func visitBasicBlockRangeHelper(service *service.Service, method *types.ParsedMethod, block *types.Block, node ast.Node, visitedRangeType bool, visitedRangeElem bool, rangeObj objects.Object, rangeKeyType gotypes.Type, rangeValueType gotypes.Type) (bool, bool, bool, objects.Object, gotypes.Type, gotypes.Type) {
	logger.Logger.Debugf("[CFG - VISIT BASIC BLOCK] [%s.%s] visiting block [%T]: %v", service.GetName(), method.GetName(), block, block)

	succ := block.GetNextSuccessorIfExists()

	var parsingLoop bool
	if succ != nil && succ.Block.Kind == cfg.KindRangeLoop { // as soon as we see an ident then we are "preparing" for the succeeding range loop
		logger.Logger.Warnf("RANGE AHEAD (%t, %t)! %v; ELEMS TYPE = %v", visitedRangeType, visitedRangeElem, succ.Block.Succs, rangeValueType)

		if !visitedRangeType { // range object
			if expr, ok := node.(ast.Expr); ok {
				rangeObj, _ = lookupObjectFromAstExpression(service, method, block, expr, nil, false)

				visitedRangeType = true
				if rangeObjSlice, ok := rangeObj.(*objects.SliceObject); ok {
					rangeValueType = rangeObjSlice.GetSliceType().UnderlyingType
				} else if rangeObjArray, ok := rangeObj.(*objects.ArrayObject); ok {
					rangeValueType = rangeObjArray.GetElementsType() //FIXME: for some reason the type is SliceType and not ArrayType
				} else if mapObjArray, ok := rangeObj.(*objects.MapObject); ok {
					rangeValueType = mapObjArray.GetMapType().ValueType
					rangeKeyType = mapObjArray.GetMapType().KeyType
				} else if ifaceObjArray, ok := rangeObj.(*objects.InterfaceObject); ok {
					rangeValueType = ifaceObjArray.GetType()
				} else {
					logger.Logger.Fatalf("[VISITOR BLOCK] unexpected type [%s] for range ident object: %v", utils.GetType(rangeObj), rangeObj)
				}
				parsingLoop = true
				logger.Logger.Debugf("saved type (%s) for range ahead: %s", utils.GetType(rangeValueType), rangeValueType.String())

			} else {
				// we are still in the expr for the previous block and not on the expr for the range object
				logger.Logger.Debugf("[VISITOR BLOCK] skipping ast type (%s) for node: %v", utils.GetType(node), node)
			}

		} else {
			ident, ok := node.(*ast.Ident)
			if !ok {
				logger.Logger.Fatalf("[CFG - VISIT BASIC BLOCK] unexpected type (%s) for node: %v", utils.GetType(node), node)
			}
			if visitedRangeType && !visitedRangeElem && ident.Name != "_" { // element ident
				visitedRangeElem = true
				obj := lookup.CreateObjectFromType(ident.Name, rangeValueType)
				logger.Logger.Debugf("[CFG - VISIT BASIC BLOCK] adding range obj (%s) to dependencies of range elem: %s", rangeObj.String(), obj.String())
				obj.GetVariableInfo().AddDependency(rangeObj)
				block.AddObject(obj)
				parsingLoop = true
			} else if ident.Name != "_" { // index ident
				if rangeKeyType == nil {
					obj := wrapValueInBasicVariable("0", "int", ident.Name)
					block.AddObject(obj)
				}
				parsingLoop = true
			}
		}
	}
	return parsingLoop, visitedRangeType, visitedRangeElem, rangeObj, rangeKeyType, rangeValueType
}

func getAssignmentRightObjects(service *service.Service, method *types.ParsedMethod, block *types.Block, rightExprs []ast.Expr) []objects.Object {
	var robjs []objects.Object
	for _, rvalue := range rightExprs {
		obj, _ := lookupObjectFromAstExpression(service, method, block, rvalue, nil, true)
		if tupleVariable, ok := obj.(*objects.TupleObject); ok {
			robjs = append(robjs, tupleVariable.Objects...)
		} else {
			robjs = append(robjs, obj)
		}
	}
	return robjs
}

func declareLeftIdents(service *service.Service, block *types.Block, leftIdents []*ast.Ident, t ast.Expr) {
	for _, ident := range leftIdents {
		t := lookup.ComputeTypeForAstExpr(service.File, t)
		declaredObject := lookup.CreateObjectFromType(ident.Name, t)
		logger.Logger.Warnf("[CFG - PARSE EXPR] VARIABLE IS DECLARED: %s", declaredObject.String())
		block.AddObject(declaredObject)
	}
}

func assignLeftIdents(service *service.Service, method *types.ParsedMethod, block *types.Block, leftIdents []*ast.Ident, rightExprs []ast.Expr) {
	rightObjects := getAssignmentRightObjects(service, method, block, rightExprs)
	for i, rObj := range rightObjects {
		leftIdent := leftIdents[i]
		rObj.GetVariableInfo().SetUnassigned()
		rObj.GetVariableInfo().SetName(leftIdent.Name)
		block.AddObject(rObj)
	}
}

func incDecLeftValues(service *service.Service, method *types.ParsedMethod, block *types.Block, leftExpr ast.Expr, tok token.Token) {
	lvariable, _ := lookupObjectFromAstExpression(service, method, block, leftExpr, nil, true)
	switch tok {
	case token.INC:
		lvariable.GetType().AddValue("1")
	case token.DEC:
		lvariable.GetType().AddValue("-1")
	default:
		logger.Logger.Fatalf("[CFG - INC/DEC LEFT] unknown token: %s", tok.String())
	}
}

func parseAssignmentStatement(service *service.Service, method *types.ParsedMethod, block *types.Block, assignStmt *ast.AssignStmt) {
	logger.Logger.Debugf("[CFG - ASSIGN LEFT] [%s] visiting stmt (%s): %v", service.GetName(), utils.GetType(assignStmt), assignStmt)
	rvariables := getAssignmentRightObjects(service, method, block, assignStmt.Rhs)
	for i, rObj := range rvariables {
		lvalue := assignStmt.Lhs[i]
		logger.Logger.Debugf("[CFG - ASSIGN LEFT] [%s] assigning left value: \n\t\t\t\t\t\t - (lvalue) [%T] %v \n\t\t\t\t\t\t - (rvalue) [%T] %s", service.GetName(), lvalue, lvalue, rObj, rObj.LongString())
		switch e := lvalue.(type) {
		case *ast.Ident:
			if assignStmt.Tok == token.DEFINE || assignStmt.Tok == token.ASSIGN { // := OR =
				var newObj objects.Object

				if rObj.GetVariableInfo().GetName() == "" { // (token.DEFINE) defining for the first time
					newObj = rObj
				} else { // (token.ASSIGN) already exists
					newObj = rObj.NewObject()
				}

				newObj.GetVariableInfo().SetUnassigned()
				newObj.GetVariableInfo().SetName(e.Name)
				block.AddObject(newObj)

				if funcVar, ok := rObj.(*objects.FuncObject); ok {
					funcVar.GetFuncType().Name = e.Name
					parseInlineFuncDeclaration(block, funcVar.GetFuncType().Body, e.Name)
				}
			} else if assignStmt.Tok == token.ADD_ASSIGN { // +=
				lObj := block.GetLatestObjectByName(e.Name)
				lObj.GetType().AddValue(rObj.GetType().GetBasicValue())
			} else if assignStmt.Tok == token.SHL_ASSIGN { // <<=
				lObj := block.GetLatestObjectByName(e.Name)
				logger.Logger.Warnf("[CFG - ASSIGN LEFT] ignoring token (%v) for lobj (%s) in assignment: %v", assignStmt.Tok, lObj.String(), assignStmt)
			} else {
				logger.Logger.Fatalf("[CFG - ASSIGN LEFT] [%s] unexpected token (%v) for assignment: %v", service.GetName(), assignStmt.Tok, assignStmt)
			}
		case *ast.SelectorExpr:
			lObj, _ := lookupObjectFromAstExpression(service, method, block, e, nil, true)
			switch ee := lObj.(type) {
			case *objects.FieldObject:
				logger.Logger.Debugf("[CFG - ASSIGN LEFT] got lvariable (%s) in assignStmt: %v", lObj.String(), assignStmt)
				//newLeftVariable := lvariable.NewVersion()
				lObj.AssignVariable(rObj)
			default:
				logger.Logger.Fatalf("[CFG - ASSIGN LEFT] [%s] unsupported left variable type (%s): %v", service.GetName(), utils.GetType(ee), lObj.String())
			}
		case *ast.IndexExpr: // e.g. res[rt] = pc
			lvariable, _ := lookupObjectFromAstExpression(service, method, block, e.X, nil, true)
			//newLeftVariable := lvariable.NewVersion()
			switch ee := lvariable.(type) {
			case *objects.MapObject:
				keyVariable, _ := lookupObjectFromAstExpression(service, method, block, e.Index, nil, true)
				if basicObj, ok := getUnderlyingBasicObjectIfExists(keyVariable); ok {
					ee.AddKeyValue(basicObj, rObj)
				} else {
					ee.AddDynamicKeyValue(keyVariable, rObj)
				}
			case *objects.ArrayObject:
				idxVariable, _ := lookupObjectFromAstExpression(service, method, block, e.Index, nil, true)
				idx, ok := computeArrayIndexFromObject(idxVariable)
				if ok {
					ee.SetElementAt(idx, rObj)
				} else {
					ee.AddDynamicElement(rObj)
					rObj.GetVariableInfo().SetDynamic()
				}
			case *objects.SliceObject:
				idxVariable, _ := lookupObjectFromAstExpression(service, method, block, e.Index, nil, true)
				idx, ok := computeArrayIndexFromObject(idxVariable)
				if ok {
					ee.SetElementAt(idx, rObj)
				} else {
					ee.AddDynamicElement(rObj)
					rObj.GetVariableInfo().SetDynamic()
				}
			default:
				logger.Logger.Fatalf("[CFG - ASSIGN LEFT] [%s] unsupported left variable type (%s): %v", service.GetName(), utils.GetType(ee), lvariable.String())
			}
		default:
			logger.Logger.Fatalf("[CFG - ASSIGN LEFT] [%s] unexpected type (%s) for left value (%v) in assignment with token (%v): %v", method.Name, utils.GetType(lvalue), lvalue, assignStmt.Tok, assignStmt)
		}
	}
}

func parseInlineFuncDeclaration(block *types.Block, body *ast.BlockStmt, name string) *types.CFG {
	inlineCFG := GenerateInlineFuncCFG(body, name)
	if name != "" { // if empty then it is an anonymous function
		block.AddInlineFunc(name, inlineCFG)
	}
	return inlineCFG
}

func parseInlineFuncCall(service *service.Service, method *types.ParsedMethod, block *types.Block, inlineCFG *types.CFG, params *ast.FieldList, args []ast.Expr) []objects.Object {
	entryBlock := inlineCFG.GetEntryParsedBlock()
	entryBlock.AppendVarsFromPredecessor(block)
	entryBlock.AppendInlineFuncsFromPredecessor(block)

	var variables []objects.Object
	for i, arg := range args {
		v, _ := lookupObjectFromAstExpression(service, method, block, arg, nil, true)
		v = v.DeepCopy()
		v.GetVariableInfo().SetName(params.List[i].Names[0].Name)
		variables = append(variables, v)
	}
	entryBlock.AddVariables(variables)

	var blocksStr string
	for _, block := range inlineCFG.GetParsedBlocks() {
		blocksStr += "\t\t\t\t\t - " + block.AstInfoString() + "\n"
	}

	visitBasicBlock(service, method, entryBlock)
	return entryBlock.Results

	//logger.Logger.Fatalf("[CFG] [%s] parsing service method cfg for (%s):\n%s", service.GetName(), method.String(), blocksStr)
}

func parseNodeBody(service *service.Service, method *types.ParsedMethod, block *types.Block, node ast.Node) []*ast.DeferStmt {
	var deferStmts []*ast.DeferStmt
	logger.Logger.Debugf("[CFG - PARSE NODE BODY] (%s) visiting node (%v)", utils.GetType(node), node)
	switch e := node.(type) {
	// ------------
	// Go Routines
	// ------------
	case *ast.GoStmt:
		if funcLit, ok := e.Call.Fun.(*ast.FuncLit); ok {
			cfg := parseInlineFuncDeclaration(block, funcLit.Body, "")
			parseInlineFuncCall(service, method, block, cfg, funcLit.Type.Params, e.Call.Args)
		} else { // e.g. we can have a previous assignment to a variable function and then call it here
			parseAndSaveCall(service, method, block, e.Call)
		}
	// ----------------------------
	// Declarations and Assignments
	// ----------------------------
	case *ast.DeclStmt: // e.g. `foobar := "foobar"`
		for _, spec := range e.Decl.(*ast.GenDecl).Specs {
			parseNodeBody(service, method, block, spec)
		}
	case *ast.ValueSpec: // e.g. var foobar OR var foobar = "foobar"
		logger.Logger.Warnf("[CFG - PARSE EXPR] parsing value spec with names = (%v) and values = (%v)", e.Names, e.Values)
		if len(e.Values) == 0 { // variables are being declared with types e.g., `var foobar string`
			declareLeftIdents(service, block, e.Names, e.Type)
		} else { // variables are being declared and assigned e.g., `var foobar := "foobar"`
			assignLeftIdents(service, method, block, e.Names, e.Values)
		}
	case *ast.AssignStmt: // e.g. `for i := 0`
		parseAssignmentStatement(service, method, block, e)
	// -----------------------------------
	// Calls and Parenthesized Expressions
	// -----------------------------------
	case *ast.CallExpr:
		// a call expr can also happen upon a structure assignment
		// e.g. post := Post { ... }
		parseAndSaveCall(service, method, block, e)
	case *ast.ParenExpr:
		// e.g. when used as a bool value in an if statement (assumes it is inside a parentheses)
		// in this case, the unfolded service from ParenExpr is a CallExpr
		parseNodeBody(service, method, block, e.X)
	// -----------------
	// Other Expressions
	// -----------------
	case *ast.ExprStmt:
		parseNodeBody(service, method, block, e.X)
	case *ast.UnaryExpr: // e.g. <-forever
		logger.Logger.Warnf("[CFG - PARSE EXPR] [%s.%s] ignoring %s: %s", service.GetName(), method.GetName(), utils.GetType(node), node)
	// -------
	// Returns
	// -------
	case *ast.ReturnStmt:
		for _, resultExpr := range e.Results {
			v, _ := lookupObjectFromAstExpression(service, method, block, resultExpr, nil, false)
			logger.Logger.Infof("ADDED RESULT: %s", v.String())
			block.AddResult(v)
		}
	// -------------
	// If Statements
	// -------------
	case *ast.IfStmt: // FIXME: we should not needs this! we are only encountering this because we are fully parsing the GoStmt
		logger.Logger.Warnf("[CFG - PARSE EXPR] [%s.%s] ignoring %s: %s", service.GetName(), method.GetName(), utils.GetType(node), node)
	case *ast.BinaryExpr: // FIXME: same... e.g. if err != nil
		logger.Logger.Warnf("[CFG - PARSE EXPR] [%s.%s] ignoring %s: %s", service.GetName(), method.GetName(), utils.GetType(node), node)
	case *ast.Ident: // FIXME: same... e.g. if flag ...
		logger.Logger.Warnf("[CFG - PARSE EXPR] [%s.%s] ignoring %s: %s", service.GetName(), method.GetName(), utils.GetType(node), node)
	case *ast.SelectorExpr: // FIXME: same... e.g. for ... range userInfo.Followers
		logger.Logger.Warnf("[CFG - PARSE EXPR] [%s.%s] ignoring %s: %s", service.GetName(), method.GetName(), utils.GetType(node), node)

	// ---------
	// For Loops
	// ---------
	case *ast.IncDecStmt: // e.g. increment in ForPost block (e.g. i++)
		incDecLeftValues(service, method, block, e.X, e.Tok)

	// ----------------
	// Other Statements
	// ----------------
	case *ast.DeferStmt:
		deferStmts = append(deferStmts, e)

	// ----------------
	// Ignore Remaining
	// ----------------
	default:
		logger.Logger.Fatalf("[CFG - PARSE EXPR] [%s.%s] unknown type in parseExpressions: %s", service.GetName(), method.GetName(), utils.GetType(node))
	}
	return deferStmts
}

func saveParsedFuncCallParams(service *service.Service, method *types.ParsedMethod, block *types.Block, parsedCall types.Call, args []ast.Expr) {
	if parsedCall.GetName() == "StorePost" {
		logger.Logger.Tracef("(1) FOUND CALL TO SERVICE VAR %s", parsedCall.GetName())
	}
	for i, arg := range args {
		logger.Logger.Tracef("[CFG] inside save func call params")
		param, _ := lookupObjectFromAstExpression(service, method, block, arg, nil, false)

		// upgrade variable with known type from function method
		if _, ok := param.GetType().(*gotypes.GenericType); ok {
			param.GetVariableInfo().Type = parsedCall.GetMethod().GetParams()[i].GetType()
			logger.Logger.Warnf("[CFG] upgrading variable %s with new type %s", param.GetVariableInfo().Name, param.GetType().String())
		}
		parsedCall.AddParam(param)
		logger.Logger.Tracef("ADDED PARAM: %s", param.LongString())
	}
	if parsedCall.GetName() == "StorePost" {
		logger.Logger.Tracef("(1) FOUND CALL TO SERVICE VAR %s", parsedCall.GetName())
	}
}

func getFuncCallDeps(service *service.Service, method *types.ParsedMethod, block *types.Block, callExpr *ast.CallExpr) []objects.Object {
	var deps []objects.Object
	for _, expr := range callExpr.Args {
		logger.Logger.Warnf("[CFG] [%s] searching for function call deps in expression %v", service.GetName(), expr)
		if v, _ := lookupObjectFromAstExpression(service, method, block, expr, nil, false); v != nil {
			deps = append(deps, v)
		}
	}
	return deps
}

// 1. creates new variables from golang types extracted from the original golang expression types
// 2. returns new tuple variable that encompasses all the new created variables
// 3. add the new tuple variable to the return parameter of the call argument if not nil
// (the call argument will be nil when we encounter built-in calls - TO BE FIXED)
func computeInternalFuncCallReturns(service *service.Service, callExpr *ast.CallExpr, call types.Call) *objects.TupleObject {
	tupleVar := &objects.TupleObject{
		ObjectInfo: &objects.ObjectInfo{
			Type: &gotypes.TupleType{},
			Id:   objects.VARIABLE_INLINE_ID,
		},
	}

	if signatureGoType, ok := service.GetPackage().GetTypeInfo(callExpr.Fun).(*golangtypes.Signature); ok {
		if signatureGoType.Results() != nil {
			signatureResults := lookup.LookupTypesForGoTypes(service.GetPackage(), signatureGoType.Results())
			for _, t := range signatureResults.(*gotypes.TupleType).Types {
				newVar := lookup.CreateObjectFromType("", t)
				tupleVar.AddObjectAndType(newVar)
				newVar.GetVariableInfo().AddParent(newVar, tupleVar)

				if call != nil {
					call.AddReturn(newVar)
				}
			}
		}
	} else {
		logger.Logger.Fatalf("[CFG CALLS] unexpected type for imported method call %v", callExpr.Fun)
	}
	return tupleVar
}

func computeExternalFuncCallReturns(service *service.Service, callExpr *ast.CallExpr, deps []objects.Object) *objects.TupleObject {
	tupleType := &gotypes.TupleType{}
	tupleVar := &objects.TupleObject{
		ObjectInfo: &objects.ObjectInfo{
			Type: tupleType,
			Id:   objects.VARIABLE_INLINE_ID,
		},
	}

	if signatureGoType, ok := service.GetPackage().GetTypeInfo(callExpr.Fun).(*golangtypes.Signature); ok {
		if signatureGoType.Results() != nil {
			signatureResults := lookup.LookupTypesForGoTypes(service.GetPackage(), signatureGoType.Results())

			if len(signatureResults.(*gotypes.TupleType).Types) == 1 && len(deps) == 0 {
				newVar := lookup.CreateObjectFromType("", signatureResults.(*gotypes.TupleType).Types[0])
				logger.Logger.Warnf("[FIXMEEEEEEE!!!!!!] (IS THIS EVEN CORRECT???) CREATED VAR FOR RETURNED TUPLE IN EXTERNAL FUNC CALL: %s", newVar.String())
				tupleVar.Objects = append(tupleVar.Objects, newVar)
				newVar.GetVariableInfo().AddParent(newVar, tupleVar)
			} else {
				logger.Logger.Warnf("[CFG CALLS] call returns tuple with len %d and depends on %d variables: %v", len(signatureResults.(*gotypes.TupleType).Types), len(deps), deps)
				for _, t := range signatureResults.(*gotypes.TupleType).Types {
					newVar := lookup.CreateObjectFromType("", t)
					ok := objects.AddUnderlyingDepsFromFuncCall(newVar, deps)
					if !ok {
						logger.Logger.Warnf("[CFG CALLS] cannot keep variable (%s) (%s) for underlying deps list with len (%d): %v", objects.VariableTypeName(newVar), newVar.String(), len(deps), deps)
						/* newVar = &objects.GenericVariable{
							ObjectInfo: &objects.ObjectInfo{
								Type: newVar.GetType(),
								Id:   objects.VARIABLE_UNASSIGNED_ID,
							},
							Params: deps,
						} */
						objects.AddVariableInfoDependencies(newVar, deps)
					}
					tupleVar.Objects = append(tupleVar.Objects, newVar)
					newVar.GetVariableInfo().AddParent(newVar, tupleVar)
				}
				logger.Logger.Warnf("CREATED COMPOSITE VAR FOR (%d) TUPLE: %s", len(deps), tupleVar.String())
			}
		}
	} else {
		logger.Logger.Fatalf("[CFG CALLS] unexpected type for imported method call %v", callExpr.Fun)
	}
	return tupleVar
}

func saveCallToStructOrInterface(service *service.Service, method *types.ParsedMethod, block *types.Block, callExpr *ast.CallExpr, leftVariableTypeName string, methodName string, pkgPath string) (*types.ParsedInternalCall, *objects.TupleObject, *objects.TupleObject) {
	logger.Logger.Debugf("[CFG CALLS] [%s] saving call in (%s) for ast CallExpr: %v", service.GetName(), method.GetName(), callExpr)
	if pkgPath == "" {
		logger.Logger.Debugf("FIX ME: WE ENCOUNTER BUILT-IN PACKAGES e.g. err.Error() -> string")
		tupleVar := computeInternalFuncCallReturns(service, callExpr, nil)
		return nil, nil, tupleVar
	}
	var parsedMethod *types.ParsedMethod

	pkg := service.GetPackage().GetImportedPackageIfExists(pkgPath)

	if pkg != nil {
		parsedMethod = pkg.GetParsedMethodIfExists(methodName, leftVariableTypeName)
	}

	if parsedMethod != nil {
		logger.Logger.Warnf("[CFG CALLS] [%s] !!!!!!!!!!!! GOT PARSED METHOD (%s): %v", service.GetName(), methodName, parsedMethod.GetParsedCfg())
		//logger.Logger.Debugf("[CFG CALLS] got parsed method %s", parsedMethod.String())
		parsedCall := &types.ParsedInternalCall{
			ParsedCall: types.ParsedCall{
				Ast:    callExpr,
				Name:   parsedMethod.Name,
				Method: parsedMethod,
			},
			ServiceTypeName: &gotypes.ServiceType{Name: service.Name, PackagePath: service.GetPackageName()},
		}
		saveParsedFuncCallParams(service, method, block, parsedCall, callExpr.Args)
		method.Calls = append(method.Calls, parsedCall)
		variable := computeInternalFuncCallReturns(service, callExpr, parsedCall)
		logger.Logger.Infof("[CFG CALLS] [%s] found internal call (%s) in package (%s) -- returned tuple: %s", service.GetName(), parsedCall.GetName(), pkg.GetName(), variable.String())
		return parsedCall, variable, nil
	}
	deps := getFuncCallDeps(service, method, block, callExpr)
	tupleVar := computeExternalFuncCallReturns(service, callExpr, deps)
	return nil, nil, tupleVar
}

// (b1) call to variable (including receiver) in block
func parseCallToVariableInBlock(service *service.Service, method *types.ParsedMethod, block *types.Block, callExpr *ast.CallExpr, variable objects.Object, idents []*ast.Ident, identsStr string) *objects.TupleObject {
	funcIdent := idents[len(idents)-1]
	logger.Logger.Infof("[CFG] [%s.%s] parsing call to variable (%s) in block for method (%s): %v", service.GetName(), method.Name, variable.String(), funcIdent, callExpr)

	// check if variable is a declared func
	if funcVar, ok := variable.(*objects.FuncObject); ok {
		inlineFunc := block.GetLatestInlineFunc(funcVar.GetFuncType().GetName())
		tupleVar := &objects.TupleObject{
			ObjectInfo: &objects.ObjectInfo{
				Id:   objects.VARIABLE_INLINE_ID,
				Type: &gotypes.TupleType{},
			},
		}
		results := parseInlineFuncCall(service, method, block, inlineFunc.CFG, funcVar.GetFuncType().Params, callExpr.Args)
		for _, r := range results {
			tupleVar.AddObjectAndType(r)
		}
		return tupleVar
	}

	logger.Logger.Tracef("[CFG CALLS] (1) call to variable: %s", variable.LongString())
	leftVariableTypeName := variable.GetType().GetName()
	variable = objects.UnwrapPointerVariable(variable)
	logger.Logger.Debugf("START ITERATING for variable: %s", variable)
	for i := 1; i < len(idents); i++ {
		logger.Logger.Debugf("ITERATING (%d) for variable: %s", i, variable)
		ident := idents[i]
		variable = objects.UnwrapPointerVariable(variable)

		if structVar, ok := variable.(*objects.StructObject); ok {
			fieldName := ident.Name
			variable = structVar.GetFieldVariableIfExists(fieldName)

			if variable == nil {
				fieldType := structVar.GetStructType().GetFieldTypeByNameIfExists(fieldName)
				if fieldType != nil {
					fieldVar := lookup.CreateObjectFromType(fieldName, fieldType)
					structVar.AddOrGetFieldKeyVariable(fieldName, fieldVar)
				} else {
					methodName := ident.Name
					pkgPath := structVar.GetStructType().GetMethodPackagePath(methodName)

					var externalCallTupleVar *objects.TupleObject
					_, variable, externalCallTupleVar = saveCallToStructOrInterface(service, method, block, callExpr, leftVariableTypeName, methodName, pkgPath)
					if externalCallTupleVar != nil {
						return externalCallTupleVar
					}
				}
			}
		} else if genericVar, ok := variable.(*objects.GenericObject); ok {
			logger.Logger.Warnf("[CFG CALLS] ignoring generic var %s", genericVar.LongString())
			if deps := getFuncCallDeps(service, method, block, callExpr); deps != nil {
				tupleVar := computeExternalFuncCallReturns(service, callExpr, deps)
				return tupleVar
			}
			parsedCall := &types.ParsedInternalCall{
				ParsedCall: types.ParsedCall{
					Ast:     callExpr,
					CallStr: identsStr,
					Name:    funcIdent.Name,
					Pos:     callExpr.Pos(),
				},
			}
			tupleVar := computeInternalFuncCallReturns(service, callExpr, parsedCall)
			return tupleVar
		} else if interfaceVar, ok := variable.(*objects.InterfaceObject); ok {
			logger.Logger.Debugf("[CFG CALLS] call to interface variable %s: %v", interfaceVar.LongString(), service.GetPackage().GetTypeInfo(callExpr.Fun))
			methodName := ident.Name
			pkgPath := interfaceVar.GetInterfaceType().GetMethodPackagePath(methodName)
			var externalCallTupleVar *objects.TupleObject
			_, variable, externalCallTupleVar = saveCallToStructOrInterface(service, method, block, callExpr, leftVariableTypeName, methodName, pkgPath)
			if externalCallTupleVar != nil {
				logger.Logger.Debugf("[CFG CALLS] EXTERNAL TUPLE VAR: %s (%s)", externalCallTupleVar.String(), utils.GetType(externalCallTupleVar.Objects[0]))
				return externalCallTupleVar
			}

			logger.Logger.Fatalf("[CFG CALLS] call to interface variable (%s) for method (%s) in package path (%s)", interfaceVar.LongString(), methodName, pkgPath)
			return nil
		} else if fieldVariable, ok := variable.(*objects.FieldObject); ok {
			variable = fieldVariable.GetWrappedVariable()
			logger.Logger.Debugf("GOT VAR: %s", variable.String())
		} else if _, ok := variable.(*objects.ServiceObject); ok {
			break
		} else if _, ok := variable.(*blueprint.BlueprintBackendObject); ok {
			break
		} else if _, ok := variable.(*objects.TupleObject); ok {
			logger.Logger.Fatalf("[CFG CALLS] [TODO] nested calls!")
		} else {
			logger.Logger.Fatalf("[CFG CALLS] unexpected call for variable (%s) with type (%s) (%s)", variable.String(), utils.GetType(variable), utils.GetType(variable.GetType()))
		}
	}
	// e.g. from internal calls
	if tupleVar, ok := variable.(*objects.TupleObject); ok {
		return tupleVar
	}
	if serviceVar, ok := variable.(*objects.ServiceObject); ok {
		// store function call either as service call or database call
		// if the field corresponds to a service field
		// 1. extract the service field from the current service
		// 2. get the target service service for the type
		// 3. add the targeted method of the other service for the current call expression
		targetService := service.Services[serviceVar.GetServiceName()]
		targetMethod := targetService.GetExportedMethod(funcIdent.Name)
		parsedCall := &types.ParsedServiceCall{
			ParsedCall: types.ParsedCall{
				Ast:     callExpr,
				CallStr: identsStr,
				Name:    funcIdent.Name,
				Pos:     callExpr.Pos(),
				Method:  targetMethod,
			},
			CallerTypeName: &gotypes.ServiceType{Name: service.Name, PackagePath: service.GetPackageName()},
			CalleeTypeName: serviceVar.GetType(),
		}
		saveParsedFuncCallParams(service, method, block, parsedCall, callExpr.Args)
		tupleVar := computeInternalFuncCallReturns(service, callExpr, parsedCall)
		method.Calls = append(method.Calls, parsedCall)
		logger.Logger.Infof("[CFG CALLS] [%s] found service call (%s) -- returned tuple: %s", service.GetName(), parsedCall.Name, tupleVar.String())
		return tupleVar
	}
	if blueprintBackendVar, ok := variable.(*blueprint.BlueprintBackendObject); ok {
		logger.Logger.Warnf("GOT BLUEPRINT BACKEND VAR: %s", blueprintBackendVar.String())
		blueprintBackendType := blueprintBackendVar.GetBlueprintBackendType()
		blueprintMethod := blueprintBackendType.GetMethod(funcIdent.Name)
		blueprintMethod.SetCalledBackendType(blueprintBackendType)
		parsedCall := &types.ParsedDatabaseCall{
			ParsedCall: types.ParsedCall{
				Ast:     callExpr,
				CallStr: identsStr,
				Name:    funcIdent.Name,
				Pos:     callExpr.Pos(),
				Method:  blueprintMethod,
			},
			CallerTypeName: &gotypes.ServiceType{Name: service.Name, PackagePath: service.GetPackageName()},
			DbInstance:     blueprintBackendType.DatastoreInstance,
		}

		if blueprintBackendType.DatastoreInstance == nil {
			logger.Logger.Fatalf("[CFG CALLS] unexpected nil db instance for backend type (%s)", blueprintBackendType.String())
		}

		saveParsedFuncCallParams(service, method, block, parsedCall, callExpr.Args)
		tupleVar := computeInternalFuncCallReturns(service, callExpr, parsedCall)

		// maybe user is just getting the collection
		if blueprintMethod.IsNoSQLBackendCall() {
			if ok, index := blueprintMethod.ReturnsNoSQLCollection(); ok {
				databaseName := parsedCall.Params[1].GetType().GetBasicValue()
				collectionName := parsedCall.Params[2].GetType().GetBasicValue()

				blueprintMethod.SetNoSQLDatabaseCollection(databaseName, collectionName, blueprintBackendType.DatastoreInstance)
				collectionType := tupleVar.GetVariableAt(index).(*blueprint.BlueprintBackendObject).GetBlueprintBackendType()
				collectionType.SetNoSQLDatabaseCollection(databaseName, collectionName, blueprintBackendType.DatastoreInstance)
				logger.Logger.Infof("[CFG CALLS] found NoSQLDatabase call (%s) to instance (%s) -- returned tuple: %s", parsedCall.Name, parsedCall.DbInstance.String(), tupleVar.String())
				return tupleVar
			} else {
				logger.Logger.Fatalf("[CFG CALLS] method (%s) must return NoSQL collection", blueprintMethod.String())
			}
		}
		if blueprintMethod.IsNoSQLCollectionCall() && blueprintBackendType.IsNoSQLCollection() {
			databaseName := blueprintBackendType.NoSQLComponent.Database
			collectionName := blueprintBackendType.NoSQLComponent.Collection

			if ok, index := blueprintMethod.ReturnsNoSQLCursor(); ok {
				fmt.Printf("[%s] HERE: %s (DATABASE = %s, COLLECTION = %s)", service.GetName(), blueprintMethod.String(), databaseName, collectionName)
				collectionType := tupleVar.GetVariableAt(index).(*blueprint.BlueprintBackendObject).GetBlueprintBackendType()
				collectionType.SetNoSQLDatabaseCursor(databaseName, collectionName, blueprintBackendType.DatastoreInstance)

				method.Calls = append(method.Calls, parsedCall)
				logger.Logger.Infof("[CFG CALLS] found (NoSQLDatabase.NoSQLCollection -> NoSQLCursor) call (%s) to instance (%s) -- returned tuple: %s", parsedCall.Name, parsedCall.DbInstance.String(), tupleVar.String())
				return tupleVar
			} else {
				// e.g. NoSQLDatabase.NoSQLCollection.InsertOne(), NoSQLDatabase.NoSQLCollection.InsertMany()
				method.Calls = append(method.Calls, parsedCall)
				logger.Logger.Infof("[CFG CALLS] found (NoSQLDatabase.NoSQLCollection -> _ ) call (%s) to instance (%s) -- returned tuple: %s", parsedCall.Name, parsedCall.DbInstance.String(), tupleVar.String())
				return tupleVar
			}
		}
		if blueprintMethod.IsNoSQLCursorCall() {
			// e.g. NoSQLDatabase.NoSQLCursor.One(), NoSQLDatabase.NoSQLCursor.All()
			// cursor is tainted from operation NoSQLDatabase.FindOne() when building the schema later
			// this way, we also need to attach the target variable that fetches the result from the cursor to its dependencies
			// so that we can latter taint the target variable as well
			targetVariable := parsedCall.GetParam(1)
			blueprintBackendVar.SetTargetVariable(targetVariable)
			logger.Logger.Infof("[CFG CALLS] [%s] found NoSQLCursor call (%s) to instance (%s) -- returned tuple: %s // %s", service.GetName(), parsedCall.Name, parsedCall.DbInstance.String(), tupleVar.String(), parsedCall.Method.LongString())
			return tupleVar
		}
		method.Calls = append(method.Calls, parsedCall)
		logger.Logger.Infof("[CFG CALLS] found Cache or Queue datastore call (%s) to instance (%s) -- returned tuple: %s", parsedCall.Name, parsedCall.DbInstance.String(), tupleVar.String())
		return tupleVar
	}
	logger.Logger.Fatalf("[CFG CALLS] unable to parse call to variable (%s) with type (%s) in call expr fun: %v\nBLOCK VARS: %v", variable.String(), utils.GetType(variable), callExpr.Fun, block.Objs)
	return nil
}

func searchCallToMethodInImportedPackage(service *service.Service, method *types.ParsedMethod, block *types.Block, callExpr *ast.CallExpr, pkg *types.Package, idents []*ast.Ident, identsStr string) (*objects.TupleObject, *types.Package, bool) {
	funcIdent := idents[len(idents)-1]
	logger.Logger.Infof("[CFG CALLS] [%s.%s] searching call to method (%s) in imported package (%s) (type = %d)", service.GetName(), method.Name, funcIdent, pkg.GetName(), pkg.Type)
	switch pkg.Type {
	case types.EXTERNAL:
		if deps := getFuncCallDeps(service, method, block, callExpr); deps != nil {
			tupleVar := computeExternalFuncCallReturns(service, callExpr, deps)
			return tupleVar, nil, false
		}
		parsedCall := &types.ParsedInternalCall{
			ParsedCall: types.ParsedCall{
				Ast:     callExpr,
				CallStr: identsStr,
				Name:    funcIdent.Name,
				Pos:     callExpr.Pos(),
			},
		}
		tupleVar := computeInternalFuncCallReturns(service, callExpr, parsedCall)
		logger.Logger.Infof("[CFG CALLS] [%s.%s] found call (%s) to method in imported external package (%s) -- returned tuple: %s", service.GetName(), method.Name, parsedCall.CallStr, pkg.GetName(), tupleVar.String())
		return tupleVar, nil, false
	case types.BLUEPRINT:
		logger.Logger.Warnf("[CFG CALLS] [%s.%s] ignoring call with idents (%v) in blueprint package", service.GetName(), method.Name, identsStr)
		// ignore direct calls to blueprint package
		// we only care about backend calls to functions of well-defined interfaces (cache, queue, nosqldatabase)
		return nil, nil, true
	case types.APP:
		if parsedMethod := pkg.GetParsedMethodIfExists(funcIdent.Name, ""); parsedMethod != nil {
			return nil, pkg, false

		}
	}
	return nil, nil, false
}

func searchCallToMethodInImport(service *service.Service, method *types.ParsedMethod, block *types.Block, callExpr *ast.CallExpr, impt *types.Import, idents []*ast.Ident, identsStr string) (*objects.TupleObject, *types.Package, bool) {
	funcIdent := idents[len(idents)-1]
	logger.Logger.Infof("[CFG CALLS] [%s.%s] searching call to method (%s) in imported package (%s)", service.GetName(), method.Name, funcIdent, impt.Alias)

	if pkg := service.GetPackage().GetImportedPackage(impt.PackagePath); pkg != nil {
		return searchCallToMethodInImportedPackage(service, method, block, callExpr, pkg, idents, identsStr)
	}
	return nil, nil, false
}

func parseCallToMethodInImportedOrCurrentPackage(service *service.Service, method *types.ParsedMethod, block *types.Block, callExpr *ast.CallExpr, callPkg *types.Package, idents []*ast.Ident, identsStr string) *objects.TupleObject {
	funcIdent := idents[len(idents)-1]
	logger.Logger.Infof("[CFG CALLS] [%s.%s] parsing call to method (%s) in imported or current package (%s): %v", service.GetName(), method.Name, funcIdent, callPkg.Name, callExpr)

	if parsedMethod := callPkg.GetParsedMethodIfExists(funcIdent.Name, ""); parsedMethod != nil {
		parsedCall := &types.ParsedInternalCall{
			ParsedCall: types.ParsedCall{
				Ast:     callExpr,
				CallStr: identsStr,
				Name:    funcIdent.Name,
				Method:  parsedMethod,
			},
			ServiceTypeName: &gotypes.ServiceType{Name: service.Name, PackagePath: service.GetPackageName()},
		}
		saveParsedFuncCallParams(service, method, block, parsedCall, callExpr.Args)
		tupleVar := computeInternalFuncCallReturns(service, callExpr, parsedCall)
		if callPkg == service.GetPackage() {
			logger.Logger.Infof("[CFG CALLS] [%s.%s] found call (%s) to method in current package (%s) -- returned tuple: %s", service.GetName(), method.Name, parsedCall.CallStr, callPkg.Name, tupleVar.String())
			method.Calls = append(method.Calls, parsedCall)
			return tupleVar
		}

		logger.Logger.Infof("[CFG CALLS] [%s.%s] found call (%s) to method in imported app package (%s) -- returned tuple: %s", service.GetName(), method.Name, parsedCall.CallStr, callPkg.Name, tupleVar.String())
		if deps := getFuncCallDeps(service, method, block, callExpr); deps != nil {
			tupleVar := computeExternalFuncCallReturns(service, callExpr, deps)
			return tupleVar
		}
		parsedCall = &types.ParsedInternalCall{
			ParsedCall: types.ParsedCall{
				Ast:     callExpr,
				CallStr: identsStr,
				Name:    funcIdent.Name,
				Pos:     callExpr.Pos(),
			},
		}
		tupleVar = computeInternalFuncCallReturns(service, callExpr, parsedCall)
		return tupleVar
	}
	return nil
}

func wrapInTupleVariable(varsToWrap ...objects.Object) *objects.TupleObject {
	var typesToWrap []gotypes.Type
	for _, v := range varsToWrap {
		typesToWrap = append(typesToWrap, v.GetType())
	}
	return &objects.TupleObject{
		Objects: varsToWrap,
		ObjectInfo: &objects.ObjectInfo{
			Type: &gotypes.TupleType{
				Types: typesToWrap,
			},
			Id: objects.VARIABLE_UNASSIGNED_ID,
		},
	}
}

func wrapInBasicVariable(variable objects.Object, typeName string) *objects.BasicObject {
	var underlyingObjects []objects.Object
	var basicValue string
	if basicObj, ok := getUnderlyingBasicObjectIfExists(variable); ok {
		basicValue = basicObj.GetBasicType().GetBasicValue()
		if basicValue == "" {
			underlyingObjects = append(underlyingObjects, basicObj)
		}
	} else {
		underlyingObjects = append(underlyingObjects, variable)
	}

	return &objects.BasicObject{
		ObjectInfo: &objects.ObjectInfo{
			Type: gotypes.NewBasicType(typeName, basicValue),
			Id:   objects.VARIABLE_UNASSIGNED_ID,
		},
		UnderlyingObjects: underlyingObjects,
	}
}

func wrapValueInBasicVariable(basicValue string, typeName string, objName string) *objects.BasicObject {
	return &objects.BasicObject{
		ObjectInfo: &objects.ObjectInfo{
			Name: objName,
			Type: &gotypes.BasicType{
				Name:  typeName,
				Value: basicValue,
			},
			Id: objects.VARIABLE_UNASSIGNED_ID,
		},
	}
}

// call to golang built-in func or type
func parseBuiltInGoFuncOrTypeCall(service *service.Service, method *types.ParsedMethod, block *types.Block, callExpr *ast.CallExpr, funcIdent *ast.Ident) objects.Object {
	logger.Logger.Infof("[CFG] [%s.%s] parsing built-in go function or type call (%s) in imported or current package: %v", service.GetName(), method.Name, funcIdent.Name, callExpr)

	if utils.IsBuiltInGoFunc(funcIdent.Name) {
		return parseBuiltInGoFuncCall(service, method, block, callExpr, funcIdent)
	} else if utils.IsBuiltInGoType(funcIdent.Name) {
		return parseBuiltInGoTypeCall(service, method, block, callExpr, funcIdent)
	}
	logger.Logger.Fatalf("[CFG] [%s.%s] unexpected built-in go func or type (%s): %v", service.GetName(), method.Name, funcIdent.Name, callExpr)
	return nil
}

// FIXME: we could actually parse the builtin.go file
// call to golang built-in func e.g. make(...), println(...), append(...)
func parseBuiltInGoFuncCall(service *service.Service, method *types.ParsedMethod, block *types.Block, callExpr *ast.CallExpr, funcIdent *ast.Ident) objects.Object {
	logger.Logger.Infof("[CFG] [%s.%s] parsing built-in go function call (%s) in imported or current package: %v", service.GetName(), method.Name, funcIdent.Name, callExpr)

	deps := getFuncCallDeps(service, method, block, callExpr)
	switch funcIdent.Name {
	case "make":
		return wrapInTupleVariable(deps[0])
	case "delete":
		return nil
	case "len":
		return &objects.BasicObject{
			ObjectInfo: &objects.ObjectInfo{
				Type: &gotypes.BasicType{
					Name:  "int",
					Value: fmt.Sprintf("len(%s)", deps[0].String()),
				},
				Id: objects.VARIABLE_UNASSIGNED_ID,
			},
			UnderlyingObjects: deps,
		}
	case "append":
		slice := deps[0]
		elems := deps[1]

		if sliceVariable, ok := slice.(*objects.SliceObject); ok {
			sliceVariable.AppendElements(elems)
		} else if arrayVariable, ok := slice.(*objects.ArrayObject); ok {
			arrayVariable.AppendElement(elems)
		} else {
			logger.Logger.Fatalf("[CFG] [%s] unexpected slice type (%s) in \"append\" call (%v)", service.GetName(), utils.GetType(slice), callExpr)
		}
		return wrapInTupleVariable(slice)
	case "println":
		return nil
	default:
		logger.Logger.Fatalf("[CFG] [%s] unexpected built-in go func (%s) for function call (%v)", service.GetName(), funcIdent.Name, callExpr)
	}
	return nil
}

// call to golang built-in type e.g. []byte(...)
func parseBuiltInGoTypeCall(service *service.Service, method *types.ParsedMethod, block *types.Block, callExpr *ast.CallExpr, funcIdent *ast.Ident) objects.Object {
	logger.Logger.Infof("[CFG] [%s.%s] parsing built-in go type call (%s) in imported or current package: %v", service.GetName(), method.Name, funcIdent.Name, callExpr)

	deps := getFuncCallDeps(service, method, block, callExpr)
	switch funcIdent.Name {
	case "byte", "string", "float32", "int64", "uint16", "uint64", "int32":
		return wrapInBasicVariable(deps[0], funcIdent.Name)
	case "delete":
		return nil
	default:
		logger.Logger.Fatalf("[CFG] [%s] unexpected built-in go type (%s) for function call (%v)", service.GetName(), funcIdent.Name, callExpr)
	}
	return nil
}

// FIXME: this does not support nested calls!!!!
func parseAndSaveCall(service *service.Service, method *types.ParsedMethod, block *types.Block, callExpr *ast.CallExpr) objects.Object {
	logger.Logger.Infof("[CFG CALLS] [%s] parsing call (%s) for args (%v)", service.GetName(), callExpr.Fun, callExpr.Args)
	idents, identsStr := lookup.GetAllSelectorIdents(callExpr.Fun)
	leftIdent := idents[0]
	funcIdent := idents[len(idents)-1]

	var varsStr = ""
	for i, expr := range callExpr.Args {
		v, _ := lookupObjectFromAstExpression(service, method, block, expr, nil, false)
		varsStr += fmt.Sprintf("\t\t\t\t\t\t\t - argument %d: (%s)\n", i, v.String())
	}

	logger.Logger.Infof("[CFG CALLS] [%s] found arguments for call with idents (%s):\n%s", service.GetName(), identsStr, varsStr)

	// call to variable (including receiver) in block
	if variable := block.GetLastestVariableIfExists(leftIdent.Name); variable != nil {
		tupleVar := parseCallToVariableInBlock(service, method, block, callExpr, variable, idents, identsStr)
		if tupleVar != nil {
			return tupleVar
		}
	}

	// call to golang built-in type (e.g. make, println, append)
	if utils.IsBuiltInGoTypeOrFunc(funcIdent.Name) {
		return parseBuiltInGoFuncOrTypeCall(service, method, block, callExpr, funcIdent)
	}

	var callInPackage bool
	var callPkg *types.Package
	var tupleVar *objects.TupleObject

	// call to method in imported package of file
	logger.Logger.Debugf("[CFG CALLS] check if call is to imported package (%s) for package import map:\n%v", leftIdent.Name, service.GetPackage().ImportsByAliasMapStr())
	if pkg := service.GetPackage().GetImportedPackageByAliasIfExists(leftIdent.Name); pkg != nil {
		var isBlueprintCall bool
		tupleVar, callPkg, isBlueprintCall = searchCallToMethodInImportedPackage(service, method, block, callExpr, pkg, idents, identsStr)
		logger.Logger.Warnf("!!!!!!!!!!!!!!! FOUND CALL TO METHOD IN IMPORTED PACKAGE: %v // %v // %v", tupleVar, callPkg, isBlueprintCall)
		if isBlueprintCall { // skip all blueprint calls that are not on backend components - e.g. backend.GetLogger().Info(...)
			return nil
		}
		if tupleVar != nil {
			return tupleVar
		}
		if callPkg != nil {
			callInPackage = true
		}
	} else if impt := service.File.GetImportIfExists(leftIdent.Name); impt != nil {
		var isBlueprintCall bool
		tupleVar, callPkg, isBlueprintCall = searchCallToMethodInImport(service, method, block, callExpr, impt, idents, identsStr)
		if isBlueprintCall { // skip all blueprint calls that are not on backend components - e.g. backend.GetLogger().Info(...)
			return nil
		}
		if tupleVar != nil {
			return tupleVar
		}
		if callPkg != nil {
			callInPackage = true
		}
	}

	// call to method in current package
	logger.Logger.Infof("[CFG CALLS] [%s.%s] try parsing call to method (%s) in current package (%s): %v", service.GetName(), method.Name, funcIdent.Name, service.GetPackageName(), callExpr)
	if parsedMethod := service.GetPackage().GetParsedMethodIfExists(funcIdent.Name, ""); parsedMethod != nil {
		callInPackage = true
		callPkg = service.GetPackage()
	}

	if callInPackage && callPkg != nil {
		tupleVar := parseCallToMethodInImportedOrCurrentPackage(service, method, block, callExpr, callPkg, idents, identsStr)
		if tupleVar != nil {
			return tupleVar
		}
	}

	logger.Logger.Fatalf("[TODO] unexpected call: %v (call in package = %t, call pkg = %s) -- idents types = %v", callExpr.Fun, callInPackage, callPkg, idents)
	return nil
}
