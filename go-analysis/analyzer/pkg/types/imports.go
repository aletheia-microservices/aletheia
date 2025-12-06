package types

import "fmt"

// ---------------------
// ------- Import ------
// ---------------------
type Import struct {
	Alias       string
	ImportPath  string
	PackageName string
	// real pack of the imported package
	// FIXME: if dev uses go.mod to replace paths then we have to translate to the real path
	PackagePath string
}

func (i *Import) IsBlueprintPackage() bool {
	return i.PackagePath == "github.com/blueprint-uservices/blueprint/runtime/core/backend"
}
func (i *Import) String() string {
	return fmt.Sprintf("%s \"%s\" (package: %s)", i.Alias, i.ImportPath, i.PackagePath)
}
func (i *Import) GetPackagePath() string {
	return i.PackagePath
}
