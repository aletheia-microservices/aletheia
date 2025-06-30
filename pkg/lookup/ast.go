package lookup

import (
	"go/ast"

	"analyzer/pkg/logger"
	"analyzer/pkg/types"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/utils"
)

func GetAllSelectorIdentsForAstExpr(expr ast.Expr) ([]*ast.Ident, string) {
	if callExpr, ok := expr.(*ast.CallExpr); ok {
		return GetAllSelectorIdentsForAstExpr(callExpr.Fun)
	}
	if selectorExpr, ok := expr.(*ast.SelectorExpr); ok {
		r1, r2 := GetAllSelectorIdentsForAstExpr(selectorExpr.X)
		return append(r1, selectorExpr.Sel), r2 + "." + selectorExpr.Sel.Name
	}
	if arrayType, ok := expr.(*ast.ArrayType); ok {
		expr = arrayType.Elt
	}
	if ident, ok := expr.(*ast.Ident); ok {
		return []*ast.Ident{ident}, ident.Name
	}
	logger.Logger.Fatalf("[LOOKUP SELECTOR IDENTS] unexpected expression %v (type = %s)", expr, utils.GetType(expr))
	return nil, ""
}

func ComputeTypeForAstExpr(pkg *types.Package, typeExpr ast.Expr) gotypes.Type {
	logger.Logger.Debugf("[LOOKUP - COMPUTE TYPE AST] (%s) visiting type expr (%v)", utils.GetType(typeExpr), typeExpr)
	switch e := typeExpr.(type) {
	case *ast.Ident:
		if utils.IsBuiltInGoType(e.Name) {
			return &gotypes.BasicType{
				Name: e.Name,
			}
		}
		if namedType, ok := pkg.GetNamedType(e.Name); ok {
			logger.Logger.Debugf("[LOOKUP AST IDENT] got named type (%s) (type = %s)", namedType.String(), utils.GetType(namedType))
			return namedType.DeepCopy()
		}

		logger.Logger.Fatalf("[LOOKUP AST IDENT] cannot compute type for ident (%s)", e)
	case *ast.SelectorExpr:
		if _, ok := e.X.(*ast.Ident); ok {
			t := LookupTypeFromImportsForGoTypes(pkg, pkg.GetTypeInfo(typeExpr))
			if t != nil {
				return t
			}

			goType := pkg.GetTypeInfo(e.Sel)
			return ResolveTypeAndAddToPackage(pkg, nil, goType, nil, nil, nil)
		}

		logger.Logger.Fatalf("[LOOKUP AST SELECTOR] cannot parse selector expr (%v)", e)
		return nil
	case *ast.ChanType:
		return &gotypes.ChanType{
			ChanType: ComputeTypeForAstExpr(pkg, e.Value),
		}
	case *ast.MapType:
		return &gotypes.MapType{
			KeyType:   ComputeTypeForAstExpr(pkg, e.Key),
			ValueType: ComputeTypeForAstExpr(pkg, e.Value),
		}
	case *ast.InterfaceType:
		return &gotypes.InterfaceType{Methods: make(map[string]string)}
	case *ast.ArrayType:
		return &gotypes.ArrayType{
			ElementsType: ComputeTypeForAstExpr(pkg, e.Elt),
		}
	case *ast.StructType:
		structType := &gotypes.StructType{Methods: make(map[string]string)}
		for i, f := range e.Fields.List {
			if len(f.Names) != 1 {
				logger.Logger.Fatalf("[LOOKUP AST STRUCT] unexpected number of fields (%d) for %s", len(f.Names), typeExpr)
			}
			name := f.Names[0].Name
			fieldType := &gotypes.FieldType{
				Origin:      structType,
				WrappedType: ComputeTypeForAstExpr(pkg, f.Type),
				StructField: true,
				FieldName:   name,
				FieldTag:    f.Tag.Value,
				Index:       i,
			}
			if _, ok := fieldType.WrappedType.(*gotypes.StructType); ok {
				fieldType.SetEmbedded()
			}
			structType.AddFieldType(fieldType)
		}
		return structType
	case *ast.StarExpr:
		return &gotypes.PointerType{
			PointerTo: ComputeTypeForAstExpr(pkg, e.X),
		}
	case *ast.Ellipsis:
		logger.Logger.Fatalf("[LOOKUP AST] could not compute type for expr %v (type = %s) \n\npkg: %s", typeExpr, utils.GetType(e.Elt), pkg.Name)
	}
	logger.Logger.Fatalf("[LOOKUP AST] could not compute type for expr %v (type = %s) \n\npkg: %s", typeExpr, utils.GetType(typeExpr), pkg.Name)
	return nil
}

func ComputeFieldsForAstFuncDecl(pkg *types.Package, funcDecl *ast.FuncDecl) ([]*types.MethodField, []*types.MethodField, *types.MethodField) {
	parser := func(fieldsList []*ast.Field) []*types.MethodField {
		var params []*types.MethodField
		for _, field := range fieldsList {
			paramType := ComputeTypeForAstExpr(pkg, field.Type)
			// returns with types only, which is usually the most frequent scenario
			if len(field.Names) == 0 {
				param := &types.MethodField{
					FieldInfo: types.FieldInfo{
						Type: paramType,
					},
				}
				params = append(params, param)
			}
			for _, ident := range field.Names {
				param := &types.MethodField{
					FieldInfo: types.FieldInfo{
						Type: paramType,
						Name: ident.Name,
					},
				}
				params = append(params, param)
			}
		}
		return params
	}
	var params []*types.MethodField
	if funcDecl.Type.Params != nil {
		params = parser(funcDecl.Type.Params.List)
	}
	var results []*types.MethodField
	if funcDecl.Type.Results != nil {
		results = parser(funcDecl.Type.Results.List)
	}
	var receiver *types.MethodField
	if funcDecl.Recv != nil {
		receiverLst := parser(funcDecl.Recv.List)
		if len(receiverLst) > 0 {
			receiver = receiverLst[0]
		}
	}
	return params, results, receiver
}
