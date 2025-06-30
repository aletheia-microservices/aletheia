package types

import (
	"fmt"
	"go/ast"
	golangtypes "go/types"
	"slices"
	"strings"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/gotypes"
	"analyzer/pkg/types/objects"
)

type PackageType int

const (
	APP       PackageType = iota
	BLUEPRINT             // blueprint backend package
	EXTERNAL
)

type Package struct {
	Name        string
	Module      string
	PackagePath string
	Files       []*File
	Block       *Block

	// app packages are fully parsed
	// blueprint packages only contain declared types and package metadata (above)
	// external packages are not parsed
	Type PackageType

	ImportedPackages        map[string]*Package // maps import paths appearing in source files to loaded packages - TODO: add support to go modules
	ImportedPackagesByAlias map[string]*Package // maps import alias appearing in source files to loaded packages - TODO: add support to go modules

	TypesInfo         *golangtypes.Info
	DeclaredVariables map[string]objects.Object
	DeclaredConstants map[string]objects.Object
	DeclaredTypes     map[string]gotypes.Type
	ServiceTypes      map[string]*gotypes.ServiceType
	ParsedMethods     []*ParsedMethod

	// contains blueprint backend types for datastores
	// key: type_name (e.g. Cache, Queue, etc)
	DatastoreTypes map[string]gotypes.Type

	// contains all imported types, including blueprint backend types for datastores
	// key: package_path.type_name, where package_path is the real package name
	ImportedTypes     map[string]gotypes.Type
	ImportedConstants map[string]objects.Object
	ImportedVariables map[string]objects.Object
}

func (p *Package) ImportsByAliasMapStr() string {
	str := ""
	for alias, impt := range p.ImportedPackagesByAlias {
		if alias == impt.PackagePath {
			str += "\t\t\t\t\t - " + impt.PackagePath + "\n"
		} else {
			str += "\t\t\t\t\t - (" + alias + ") " + impt.PackagePath + "\n"
		}
	}
	return str
}

func getLastWordAfterSlash(input string) string {
	parts := strings.Split(input, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return input
}

func (p *Package) GetUnparsedMethods() []*ParsedMethod {
	var lst []*ParsedMethod
	for _, m := range p.ParsedMethods {
		if !m.IsParsed() {
			lst = append(lst, m)
		}
	}
	return lst
}

// AddImportedPackageByAliasIfNotExists adds an entry that maps the imported alias to the loaded package
// The ImportedPackagesByAlias can have many different keys from different source files that map to the same package
// e.g. in "encoding/json" the alias is just "json"
func (p *Package) AddImportedPackageByAliasIfNotExists(alias string, pkg *Package) {
	alias = getLastWordAfterSlash(alias)
	if existingPkg, exists := p.ImportedPackagesByAlias[alias]; exists && existingPkg != pkg {
		logger.Logger.Fatalf("[PACKAGE] [%s] a different package (%s) is already imported for alias (%s) (%s)", p.GetName(), existingPkg.GetPackagePath(), alias, pkg.GetPackagePath())
	}
	p.ImportedPackagesByAlias[alias] = pkg
}

func (p *Package) GetImportedPackageByAliasIfExists(alias string) *Package {
	if pkg, exists := p.ImportedPackagesByAlias[alias]; exists {
		return pkg
	}

	// alternative:
	// packages whose name contains "-" (e.g., "github.com/Rhymond/go-money"), can only be callable with the alias after "-"
	for key, pkg := range p.ImportedPackagesByAlias {
		splits := strings.Split(key, "-")
		usableAlias := splits[len(splits)-1]
		if alias == usableAlias {
			return pkg
		}
	}
	logger.Logger.Warnf("[PACKAGE] [%s] package with alias (%s) does not exist for import map: %s", p.GetName(), alias, p.ImportsByAliasMapStr())
	return nil
}

func (p *Package) IsAppPackage() bool {
	return p.Type == APP
}

func (p *Package) IsExternalPackage() bool {
	return p.Type == EXTERNAL
}

func (p *Package) IsBlueprintBackendPackage() bool {
	return p.Type == BLUEPRINT
}

func (p *Package) GetName() string {
	return p.Name
}

func (p *Package) GetPackagePath() string {
	return p.PackagePath
}

func (p *Package) GetDeclaredVariableOrConst(typeNameIdent *ast.Ident) objects.Object {
	if v, ok := p.DeclaredVariables[typeNameIdent.Name]; ok {
		return v
	}
	if v, ok := p.DeclaredConstants[typeNameIdent.Name]; ok {
		return v
	}
	logger.Logger.Fatalf("unknown declared variable or const (%s) in package (%s) -- declared variables: %v", typeNameIdent.Name, p.Name, p.DeclaredVariables)
	return nil
}

func (p *Package) GetDeclaredConstant(typeNameIdent *ast.Ident) objects.Object {
	if v, ok := p.DeclaredConstants[typeNameIdent.Name]; ok {
		return v
	}
	logger.Logger.Fatalf("unknown declared constant (%s) in package (%s) -- constant variables: %v", typeNameIdent.Name, p.Name, p.DeclaredConstants)
	return nil
}

func (p *Package) GetDeclaredVariableOrConstIfExists(name string) objects.Object {
	logger.Logger.Debugf("[PACKAGE DECLS] looking for variable or const named (%s) in package named (%s): %s", name, p.String(), p.PackagePath)
	if v, ok := p.DeclaredVariables[name]; ok {
		return v
	}
	if v, ok := p.DeclaredVariables[p.PackagePath+"."+name]; ok {
		return v
	}
	logger.Logger.Warnf("[BLOCK] variable (%s) does not exist in declared variables list: %v", name, p.DeclaredVariables)
	if v, ok := p.DeclaredConstants[name]; ok {
		return v
	}
	if v, ok := p.DeclaredConstants[p.PackagePath+"."+name]; ok {
		return v
	}
	logger.Logger.Warnf("[BLOCK] variable (%s) does not exist in declared constants list: %v", name, p.DeclaredVariables)
	return nil
}

func (p *Package) RegisterMethodToParse(method *ParsedMethod) {
	p.ParsedMethods = append(p.ParsedMethods, method)
}

func (p *Package) GetAllParsedMethods() []*ParsedMethod {
	return p.ParsedMethods
}

func (p *Package) GetParsedMethod(methodName string, recvTypeName string) *ParsedMethod {
	for _, m := range p.ParsedMethods {
		if m.Name == methodName && (!m.HasReceiver() || m.GetReceiver().GetType() == nil || m.GetReceiver().GetType().GetName() == recvTypeName) {
			if m.ParsedCfg == nil {
				logger.Logger.Warnf("[TYPES PACKAGES] encountered nil cfg - method (%s) is not yet parsed", methodName)
			}
			return m
		}
	}
	logger.Logger.Fatalf("unknown method (%s) with receiver type name (%s) for package (%s)", methodName, recvTypeName, p.Name)
	return nil
}

func (p *Package) GetParsedMethodIfExists(methodName string, recvTypeName string) *ParsedMethod {
	// only internal methods are parsed
	// if method is not found the it must be from an external package
	for _, m := range p.ParsedMethods {
		if m.Name == methodName && (!m.HasReceiver() || m.GetReceiver().GetType() == nil || m.GetReceiver().GetType().GetName() == recvTypeName) {
			lst := ""
			for _, m := range p.ParsedMethods {
				lst += fmt.Sprintf("\t\t\t\t\t\t - %s\n", m)
			}

			logger.Logger.Debugf("[TYPES PACKAGE] found parsed method for name (%s) - parsed methods list in package:\n%s", methodName, lst)
			return m
		}
	}
	logger.Logger.Warnf("unknown method (%s) with receiver type name (%s) in package (%s)", methodName, recvTypeName, p.Name)
	return nil
}

func (p *Package) String() string {
	return p.Name
}

func (p *Package) GetAllDeclaredTypes() map[string]gotypes.Type {
	return p.DeclaredTypes
}

func (p *Package) GetDeclaredType(name string) gotypes.Type {
	decl, ok := p.DeclaredTypes[name]
	if !ok {
		logger.Logger.Fatalf("type %s not found in declared types of package %s", name, p.Name)
	}
	return decl
}

func (p *Package) GetImportedPackageIfExists(pkgPath string) *Package {
	if pkgPath == p.PackagePath {
		return p
	}

	if pkg, ok := p.ImportedPackages[pkgPath]; ok {
		return pkg
	}

	logger.Logger.Warnf("ignoring... package (%s) not found in imports of package (%s): %v", pkgPath, p.PackagePath, p.ImportedPackages)
	return nil
}

func (p *Package) GetImportedPackage(pkgPath string) *Package {
	if pkgPath == p.PackagePath {
		return p
	}
	imptPkg, ok := p.ImportedPackages[pkgPath]
	if !ok {
		logger.Logger.Fatalf("package (%s) not found in imports of package (%s): %v", pkgPath, p.PackagePath, p.ImportedPackages)
	}
	return imptPkg
}

func (p *Package) GetTypeInfo(expr ast.Expr) golangtypes.Type {
	if p.TypesInfo == nil {
		logger.Logger.Fatalf("[PACKAGE] unexpected nil types info for package: %s", p.String())
	}
	return p.TypesInfo.TypeOf(expr)
}

func (p *Package) HasPath(path string) bool {
	return p.PackagePath == path
}

func (p *Package) GetFile(filepath string) *File {
	for _, f := range p.Files {
		if f.HasAbsolutePath(filepath) {
			return f
		}
	}
	logger.Logger.Fatalf("could not find file with path %s for package %s", filepath, p.Name)
	return nil
}

func ImportedTypeKey(packagePath string, typeName string) string {
	return packagePath + "." + typeName
}

func (p *Package) AddConstant(v objects.Object) {
	if _, exists := p.DeclaredConstants[v.GetVariableInfo().GetName()]; exists {
		logger.Logger.Fatalf("package %s already constains declared constant %s", v.GetVariableInfo().GetName(), v.String())
	}
	p.DeclaredConstants[v.GetVariableInfo().GetName()] = v
}

func (p *Package) AddVariable(v objects.Object) {
	if _, exists := p.DeclaredVariables[v.GetVariableInfo().GetName()]; exists {
		logger.Logger.Fatalf("package %s already constains declared variable %s", v.GetVariableInfo().GetName(), v.String())
	}
	p.DeclaredVariables[v.GetVariableInfo().GetName()] = v
}

func (p *Package) AddDeclaredType(e gotypes.Type) {
	if _, exists := p.DeclaredTypes[e.GetName()]; exists {
		logger.Logger.Fatalf("package %s already constains declared type %s", p.Name, e.String())
	}
	p.DeclaredTypes[e.GetName()] = e
}

func (p *Package) AddImportedTypeIfNotExists(e gotypes.Type) {
	key := ImportedTypeKey(e.GetPackage(), e.GetName())
	if _, exists := p.ImportedTypes[key]; exists {
		logger.Logger.Warnf("package %s already constains imported type %s", p.Name, e.String())
		return
	}
	p.ImportedTypes[key] = e
}

func (p *Package) AddImportedConstant(v objects.Object, packagePath string) {
	key := ImportedTypeKey(packagePath, v.GetVariableInfo().GetName())
	if _, exists := p.ImportedConstants[key]; exists {
		logger.Logger.Fatalf("package %s already constains imported constant %s", v.GetVariableInfo().GetName(), v.String())
	}
	p.ImportedConstants[key] = v
}

func (p *Package) AddImportedVariable(v objects.Object, packagePath string) {
	key := ImportedTypeKey(packagePath, v.GetVariableInfo().GetName())
	if _, exists := p.ImportedVariables[key]; exists {
		logger.Logger.Fatalf("package %s already constains imported variable %s", v.GetVariableInfo().GetName(), v.String())
	}
	p.ImportedVariables[key] = v
}

func (p *Package) AddServiceType(e *gotypes.ServiceType) {
	if _, exists := p.ServiceTypes[e.Name]; exists {
		logger.Logger.Fatalf("package %s already constains service type %s", p.Name, e.String())
	}
	p.ServiceTypes[e.Name] = e
}

func (p *Package) AddDatastoreType(e gotypes.Type) {
	if _, exists := p.DatastoreTypes[e.GetName()]; exists {
		logger.Logger.Fatalf("package %s already constains datastore type %s", p.Name, e.String())
	}
	p.DatastoreTypes[e.GetName()] = e
}

func (p *Package) LinkFile(file *File) {
	if slices.Contains(p.Files, file) {
		logger.Logger.Fatalf("package (%s) already constains file (%s)", p.Name, file.GetName())
	}
	p.Files = append(p.Files, file)
	file.SetPackage(p)
}

func (p *Package) getImportedTypesLstStr() string {
	importedTypesStr := ""
	for k, t := range p.ImportedTypes {
		importedTypesStr += "- " + k + ": " + t.String() + "\n"
	}
	return importedTypesStr
}

func (p *Package) getImportedPackagesLstStr() string {
	importedPackagesStr := ""
	for k, t := range p.ImportedPackages {
		importedPackagesStr += "- " + k + ": " + t.String() + "\n"
	}
	return importedPackagesStr
}

func (p *Package) GetImportedTypeIfExists(importPath string) gotypes.Type {
	if e, ok := p.ImportedTypes[importPath]; ok {
		return e
	}
	logger.Logger.Warnf("[PACKAGE] unknown imported type (%s) in package (%s) with import types list and packages:\n%v\n\n%v", importPath, p.Name, p.getImportedTypesLstStr(), p.getImportedPackagesLstStr())
	return nil
}

func (p *Package) GetDeclaredTypeIfExists(typeNameIdent *ast.Ident) gotypes.Type {
	if p.Type == APP {
		if e, ok := p.DeclaredTypes[typeNameIdent.Name]; ok {
			return e
		}
	}
	logger.Logger.Warnf("[PACKAGE] unknown declared type (%s) in package (%s)", typeNameIdent.Name, p.Name)
	return nil
}

func (p *Package) GetDeclaredTypeByName(name string) gotypes.Type {
	if t, ok := p.DeclaredTypes[name]; ok {
		return t
	}
	logger.Logger.Fatalf("[PACKAGE] unknown declared type (%s) for package: %s", name, p.String())
	return nil
}

func (p *Package) GetNamedType(name string) (gotypes.Type, bool) {
	if e, ok := p.ServiceTypes[name]; ok {
		logger.Logger.Tracef("[PACKAGE] found service type for (%s)", name)
		return e, true
	}
	if e, ok := p.DatastoreTypes[name]; ok {
		logger.Logger.Tracef("[PACKAGE] found datastore type for (%s)", name)
		return e, true
	}
	if e, ok := p.DeclaredTypes[name]; ok {
		logger.Logger.Tracef("[PACKAGE] found declared type for (%s)", name)
		return e, true
	}
	logger.Logger.Fatalf("unknown named type (%s) in package (%s)", name, p.Name)
	return nil, false
}
