package types

import (
	"go/ast"
)

type File struct {
	ast     *ast.File
	pkg     *Package
	name    string
	absPath string
}

func NewFile(ast *ast.File, pkg *Package, name string, absPath string) *File {
	return &File{
		ast:     ast,
		pkg:     pkg,
		name:    name,
		absPath: absPath,
	}
}

func (f *File) GetAst() *ast.File {
	return f.ast
}

func (f *File) String() string {
	return f.absPath
}

func (f *File) GetPackage() *Package {
	return f.pkg
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) SetPackage(newPkg *Package) {
	f.pkg = newPkg
}

func (f *File) HasAbsolutePath(other string) bool {
	return f.absPath == other
}
