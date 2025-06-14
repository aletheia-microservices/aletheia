package app

import (
	"go/ast"
	golangtypes "go/types"
	"path/filepath"
	"slices"

	"golang.org/x/tools/go/packages"

	"analyzer/pkg/controlflow"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/lookup"
	"analyzer/pkg/types"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

func (app *App) createPackage(goPackage *packages.Package) *types.Package {
	modulePath := ""
	if goPackage.Module != nil {
		modulePath = goPackage.Module.Path
	}

	newPackage := &types.Package{
		Name:                    goPackage.Name,    // e.g. models
		PackagePath:             goPackage.PkgPath, // e.g. github.com/blueprint-uservices/blueprint/examples/postnotification/workflow/postnotification/models
		Module:                  modulePath,        // e.g. github.com/blueprint-uservices/blueprint/examples/postnotification/workflow
		ImportedPackages:        make(map[string]*types.Package),
		ImportedPackagesByAlias: make(map[string]*types.Package),
		DeclaredTypes:           make(map[string]gotypes.Type),
		ServiceTypes:            make(map[string]*gotypes.ServiceType),
		DatastoreTypes:          make(map[string]gotypes.Type),
		DeclaredConstants:       make(map[string]objects.Object),
		DeclaredVariables:       make(map[string]objects.Object),
		ImportedTypes:           make(map[string]gotypes.Type),
		ImportedConstants:       make(map[string]objects.Object),
		ImportedVariables:       make(map[string]objects.Object),
		TypesInfo:               goPackage.TypesInfo,
		Block:                   &types.Block{},
	}

	if utils.IsAppPackage(app.PackagePath, goPackage.PkgPath) {
		logger.Logger.Infof("[APP PACKAGE FINDER] found internal package %s", goPackage.String())
		newPackage.Type = types.APP
		app.AddAppPackage(newPackage.Name, newPackage)
	} else if utils.IsBlueprintBackendPath(goPackage.PkgPath) {
		logger.Logger.Infof("[APP PACKAGE FINDER] found blueprint package %s", goPackage.String())
		newPackage.Type = types.BLUEPRINT
		app.AddBlueprintPackage(newPackage.Name, newPackage)
	} else {
		logger.Logger.Infof("[APP PACKAGE FINDER] found external package %s", goPackage.String())
		newPackage.Type = types.EXTERNAL
		app.AddExternalPackage(newPackage.Name, newPackage)
	}

	return newPackage
}

func parseBlueprintPackage(bpPackage *types.Package) {
	packagesPattern := bpPackage.PackagePath
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedModule | packages.NeedFiles | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedImports | packages.NeedDeps}, packagesPattern)
	if err != nil {
		logger.Logger.Fatalf("[APP PACKAGE BLUEPRINT] error loading packages from %s: %s", bpPackage.PackagePath, err.Error())
	}
	// should be only one package
	if len(pkgs) != 1 {
		logger.Logger.Fatalf("[APP PACKAGE BLUEPRINT] unexpected number of packages (%d) in blueprint backend %s", len(pkgs), bpPackage.PackagePath)
	}
	goPackage := pkgs[0]
	visitedNamedTypes := make(map[string]bool)
	for _, def := range goPackage.TypesInfo.Defs {
		if def == nil {
			continue
		}
		if _, isInterface := def.Type().Underlying().(*golangtypes.Interface); isInterface {
			if namedGoType, ok := def.Type().(*golangtypes.Named); ok {
				if visitedNamedTypes[namedGoType.String()] {
					continue
				}
				visitedNamedTypes[namedGoType.String()] = true

				typeName := namedGoType.Obj().Name()
				objectPackage := namedGoType.Obj().Pkg()
				if objectPackage == nil || objectPackage.Path() != bpPackage.PackagePath {
					continue
				}

				if blueprint.IsBackendComponent(typeName) {
					backendType := &blueprint.BlueprintBackendType{
						Name:    typeName,
						Package: objectPackage.Path(),
					}
					backendType.Methods = blueprint.BuildBackendComponentMethods(backendType.Name)
					bpPackage.AddDeclaredType(backendType)
					bpPackage.AddDatastoreType(backendType)
					logger.Logger.Infof("[APP PACKAGE BLUEPRINT] added blueprint backend %s", backendType.String())

				} else {
					logger.Logger.Tracef("[APP PACKAGE BLUEPRINT] ignoring type %s (%s): %v", utils.GetType(def.Type()), utils.GetType(def.Type().Underlying()), def.Type().String())
				}
			}
		}
	}
}

type packageGoInfo struct {
	ParsedPackage *types.Package
	GoPackage     *packages.Package
}

func (app *App) saveNewPackageGoInfo(packagesInfo map[string]*packageGoInfo, goPackage *packages.Package) *packageGoInfo {
	parsedPackage := app.createPackage(goPackage)
	packageInfo := &packageGoInfo{
		ParsedPackage: parsedPackage,
		GoPackage:     goPackage,
	}
	packagesInfo[goPackage.PkgPath] = packageInfo
	return packageInfo
}

func (app *App) RegisterPackages() {
	logger.Logger.Info("[APP] registering packages")
	packagesPattern := app.PackagePath + "/..."
	logger.Logger.Debug("[APP] loading packages")
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedModule | packages.NeedFiles | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedImports}, packagesPattern)
	if err != nil {
		logger.Logger.Fatalf("[APP PACKAGES] error loading packages from %s: %s", app.PackagePath, err.Error())
	}

	packagesInfo := make(map[string]*packageGoInfo)
	for _, goPackage := range pkgs {
		if _, exists := packagesInfo[goPackage.PkgPath]; !exists {
			packageInfo := app.saveNewPackageGoInfo(packagesInfo, goPackage)
			for k, imptgoPackage := range goPackage.Imports {
				importedPackageInfo, exists := packagesInfo[imptgoPackage.PkgPath]
				if !exists {
					importedPackageInfo = app.saveNewPackageGoInfo(packagesInfo, imptgoPackage)
				}
				packageInfo.ParsedPackage.ImportedPackages[k] = importedPackageInfo.ParsedPackage
				logger.Logger.Debugf("[APP PACKAGES] [%s] added import: %s", packageInfo.ParsedPackage.GetName(), importedPackageInfo.ParsedPackage.PackagePath)
			}
		}
	}

	for _, info := range packagesInfo {
		pkg := info.ParsedPackage
		files := info.GoPackage.Syntax
		for _, file := range files {
			for _, impt := range file.Imports {
				logger.Logger.Debugf("IMPT: %s", impt.Name.String())
				imptAlias := ""
				if impt.Name != nil {
					imptAlias = impt.Name.Name
				}
				imptPath := utils.RemoveQuotesFromPathImport(impt.Path.Value)
				logger.Logger.Warnf("[APP PACKAGES] [%s] added import by alias (%s): %s", pkg.Name, imptAlias, imptPath)

				if impt.Name == nil {
					imptAlias = imptPath
				}

				importedPkg := pkg.GetImportedPackage(imptPath)
				pkg.AddImportedPackageByAliasIfNotExists(imptAlias, importedPkg)
			}
		}
	}

	for _, parsedPackage := range app.GetBlueprintPackages() {
		parseBlueprintPackage(parsedPackage)
	}
	for _, parsedPackage := range app.GetAppPackages() {
		app.ParseAppPackage(parsedPackage, packagesInfo[parsedPackage.PackagePath].GoPackage)
	}

	for _, pkg := range app.GetAppPackages() {
		for _, m := range pkg.GetAllParsedMethods() {
			controlflow.GenerateMethodCFG(m)
		}
	}

	app.dumpYamlPackages()
}


//FIXME: we should actually be careful with methods that call other methods
func (app *App) ParseMethods() {
	for _, pkg := range app.GetAppPackages() {
		// (1) parse top-level methods
		for _, method := range pkg.GetUnparsedMethods() {
			if !method.HasReceiver() {
				controlflow.ParseMethodCFG(pkg, nil, method)
			}
		}
		// (2) parse methods that implement declared types
		for _, method := range pkg.GetUnparsedMethods() {
			if method.HasReceiver() && !method.HasAttachedService() {
				controlflow.ParseMethodCFG(pkg, nil, method)
			}
		}
	}

	// (3) parse all services' exposed methods
	for _, node := range app.Services {
		for _, method := range node.GetMethods() {
			controlflow.ParseMethodCFG(node.GetPackage(), node, method)
		}
	}
}

func parseGolangFuncTypes(pkg *types.Package, def golangtypes.Object, typeNameToFuncs map[string][]*golangtypes.Func) []*golangtypes.Func {
	var funcs []*golangtypes.Func
	if funcGoType, ok := def.(*golangtypes.Func); ok && funcGoType.Pkg().Path() == pkg.PackagePath {
		//logger.Logger.Debugf("FUNC TYPE %v", funcGoType)
		if signatureGoType, ok := funcGoType.Type().Underlying().(*golangtypes.Signature); ok {
			//logger.Logger.Debugf("\t\t SIGNATURE TYPE %v", signatureGoType)
			if signatureGoType.Recv() != nil {
				recv := signatureGoType.Recv().Type()
				if pointerGoType, ok := signatureGoType.Recv().Type().(*golangtypes.Pointer); ok {
					recv = pointerGoType.Elem()
				}
				if namedGoType, ok := recv.(*golangtypes.Named); ok {
					var recvTypeName string
					_, isStruct := recv.Underlying().(*golangtypes.Struct)
					_, isInterface := recv.Underlying().(*golangtypes.Interface)
					if !isStruct && !isInterface {
						logger.Logger.Fatalf("[APP] unexpected type (%s) for recv underlying (%s)", utils.GetType(recv.Underlying()), recv.Underlying())
					}
					recvTypeName = namedGoType.Obj().Name()
					if !slices.Contains(typeNameToFuncs[recvTypeName], funcGoType) {
						typeNameToFuncs[recvTypeName] = append(typeNameToFuncs[recvTypeName], funcGoType)
					}
				} else {
					logger.Logger.Fatalf("[APP] unexpected type (%s) for recv (%s)", utils.GetType(namedGoType), recv.String())
				}
			} else {
				if !slices.Contains(funcs, funcGoType) {
					funcs = append(funcs, funcGoType)
				}
			}
		}
	}
	return funcs
}

func (app *App) ParseAppPackage(parsedPackage *types.Package, goPackage *packages.Package) {
	logger.Logger.Infof("[APP PACKAGE PARSER] parsing named types for app go package (%s)", goPackage)

	visitedNamedTypes := make(map[*golangtypes.Named]bool)
	typeNameToFuncs := make(map[string][]*golangtypes.Func)
	var funcs []*golangtypes.Func

	for _, obj := range goPackage.TypesInfo.Defs {
		if obj == nil {
			continue
		}
		logger.Logger.Warnf("- [OBJ] (pkg-level=%t) [%s %s] %s", obj.Parent() == goPackage.Types.Scope(), utils.GetType(obj), utils.GetType(obj.Type()), obj.Name())

		if funcObj, ok := obj.(*golangtypes.Func); ok {
			funcs = append(funcs, parseGolangFuncTypes(parsedPackage, funcObj, typeNameToFuncs)...)
			continue
		} else if namedType, ok := obj.Type().(*golangtypes.Named); ok {
			// built-in (e.g. error, make, println) golang types have no package
			t := lookup.FindDefTypesAndAddToPackage(parsedPackage, obj, namedType, visitedNamedTypes, typeNameToFuncs, app.ServiceTypes)
			if t != nil {
				continue
			}
		}

		if obj.Parent() == goPackage.Types.Scope() {
			// if it is not a named type, then we want to save package-level objects
			// whose object type is either a CONST or VAR
			// note that objects whose object.Type() is Named are also VARs
			// and that's why this 'else if' condition needs to be after the previous
			lookup.SaveObjectToPackage(parsedPackage, obj)
		}
	}

	for i, fileAst := range goPackage.Syntax {
		file := &types.File{
			Ast:     fileAst,
			Name:    filepath.Base(goPackage.GoFiles[i]),
			AbsPath: goPackage.GoFiles[i],
			Package: parsedPackage,
			Imports: make(map[string]*types.Import),
		}
		lookup.ParseFileImports(file)
		parsedPackage.LinkFile(file)

		ast.Inspect(fileAst, func(n ast.Node) bool {
			if funcDecl, ok := n.(*ast.FuncDecl); ok {
				parseFuncDecl(funcDecl, typeNameToFuncs, parsedPackage, file, funcs)
			}
			return true
		})
	}
	logger.Logger.Tracef("[APP PACKAGE PARSER] parsed named types app package %s", parsedPackage.Name)
}

func parseFuncDecl(funcDecl *ast.FuncDecl, typeNameToFuncs map[string][]*golangtypes.Func, parsedPackage *types.Package, file *types.File, funcs []*golangtypes.Func) {
	var recvTypeIdent *ast.Ident
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		recvField := funcDecl.Recv.List[0]
		if startExpr, ok := recvField.Type.(*ast.StarExpr); ok {
			if ident, ok := startExpr.X.(*ast.Ident); ok {
				recvTypeIdent = ident
			}
		} else if ident, ok := recvField.Type.(*ast.Ident); ok {
			recvTypeIdent = ident
		}
	}

	if recvTypeIdent != nil {
		if structTypeFuncs, ok := typeNameToFuncs[recvTypeIdent.Name]; ok {
			createAndSaveMethodForFuncDecl(parsedPackage, file, funcDecl, structTypeFuncs)
		} else {
			logger.Logger.Fatalf("[APP PACKAGE PARSER] cannot find func for type name (%s)", recvTypeIdent.Name)
		}
	} else {
		createAndSaveMethodForFuncDecl(parsedPackage, file, funcDecl, funcs)
	}
}

func createAndSaveMethodForFuncDecl(pkg *types.Package, file *types.File, funcDecl *ast.FuncDecl, funcGoTypes []*golangtypes.Func) {
	for _, f := range funcGoTypes {
		if f.Name() == funcDecl.Name.Name {
			params, returns, receiver := lookup.ComputeFuncDeclFields(file, pkg, funcDecl)
			if receiver != nil {
				logger.Logger.Warnf("COMPUTED RECEIVER (%s)", receiver.GetName())
			}
			parsedMethod := &types.ParsedMethod{
				Ast:      funcDecl,
				Name:     funcDecl.Name.Name,
				Package:  pkg,
				Exported: f.Exported(),
				Params:   params,
				Returns:  returns,
				Receiver: receiver,
			}
			pkg.RegisterMethodToParse(parsedMethod)
			logger.Logger.Warnf("[APP] registered method (%s) to be parsed in package (%s)", parsedMethod.Name, pkg.Name)
		}
	}
}
