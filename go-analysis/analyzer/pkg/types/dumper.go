package types

import (
	"fmt"
	"slices"
	"sort"

	"analyzer/pkg/logger"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

// --------------
// PARSED METHODS
// --------------

func (f *ParsedMethod) Yaml() (interface{}, string) {
	logger.Logger.Tracef("[YAML] dumping yaml data for method (%s)", f.Name)
	if f.ParsedCfg != nil {
		return f.ParsedCfg.Yaml()
	}
	return nil, ""
}

func (f *ParsedMethod) YamlCalls() []string {
	logger.Logger.Tracef("[YAML CALLS] dumping yaml calls for method (%s)", f.Name)
	var lst []string
	for _, call := range f.Calls {
		logger.Logger.Tracef("\t\t\t - call: %s, method: %v", call.GetName(), call.GetMethod())
		var callStr string
		if dbCall, ok := call.(*ParsedDatabaseCall); ok {
			callStr = "DATASTORE CALL @ " + dbCall.DbInstance.GetName() + " >> "
		} else if svcCall, ok := call.(*ParsedServiceCall); ok {
			callStr = "SERVICE CALL @ " + svcCall.Name + " >> "
		} else {
			callStr = "INTERNAL CALL >> "
		}
		callStr += call.GetMethod().LongString()
		lst = append(lst, callStr)
	}
	return lst
}

// ----
// CFGS
// ----

func (cfg *CFG) Yaml() (map[string][]string, string) {
	data := make(map[string][]string)
	block := cfg.GetLastLiveBlock()
	var blockVarsStr string
	data[block.Block.String()], blockVarsStr = block.Yaml()
	return data, blockVarsStr
}

// ------
// BLOCKS
// ------

func getTaintedFieldsListAndString(dfs []*objects.ObjectDataflow) ([]string, string) {
	var fieldsLst []string
	fieldsStr := ""
	for _, df := range dfs {
		if df.Field == nil {
			logger.Logger.Fatalf("GOT NIL FIELD FOR DF (%s, %s). %s", df.Datastore, df.Service, df.Variable.String())
		}
		n := df.Datastore + "." + df.Field.GetName()
		if !slices.Contains(fieldsLst, n) {
			fieldsLst = append(fieldsLst, n)
		}
	}
	for i, db := range fieldsLst {
		fieldsStr += db
		if i < len(fieldsLst)-1 {
			fieldsStr += ", "
		}
	}
	return fieldsLst, fieldsStr
}

func getIndirectDependencies(obj objects.Object, i int) ([]objects.Object, string, []int) {
	blockVarsStr := ""

	padding := ""
	paddingSpace := "  "
	for j := 0; j < i; j++ {
		padding += "_"
		paddingSpace += " "
	}

	padding += fmt.Sprintf("%d", i)

	opStr := ""

	visitedOp := make(map[string]bool)
	for _, df := range obj.GetVariableInfo().GetAllDataflows() {

		s := df.GetOpString() + "(" + df.Datastore + "), "
		if ok := visitedOp[s]; !ok {
			visitedOp[s] = true
			opStr += s
		}
	}

	/* if opStr != "" {
		opStr = "// " + opStr
	} */

	dfsWriteOps := obj.GetVariableInfo().GetAllWriteDataflows()
	dfsReadOps := obj.GetVariableInfo().GetAllReadDataflows()
	if len(dfsWriteOps) > 0 || len(dfsReadOps) > 0 {
		writeOpsLst, writeOpsStr := getTaintedFieldsListAndString(dfsWriteOps)
		readOpsLst, readOpsStr := getTaintedFieldsListAndString(dfsReadOps)
		if len(dfsWriteOps) > 0 {
			blockVarsStr += fmt.Sprintf("%s %s --> w-tainted: write(%s) {%d}", paddingSpace, blockVarsStr, writeOpsStr, len(writeOpsLst))
		}
		if len(dfsReadOps) > 0 {
			blockVarsStr += fmt.Sprintf("%s %s --> r-tainted: read(%s) {%d}", paddingSpace, blockVarsStr, readOpsStr, len(readOpsLst))
		}
		// the '{}' symbol resets the red color caused by the from the '-->' symbol
		blockVarsStr += "\n"
	}

	blockVarsStr += fmt.Sprintf("[%s] (%s) %s\n", padding, objects.VariableTypeName(obj), obj.String())

	/* if obj.GetVariableInfo().GetName() == "items" && (len(dfsWriteOps) > 0 || len(dfsReadOps) > 0) {
		for _, df := range obj.GetVariableInfo().GetAllDataflows() {
			logger.Logger.Warnf("DF: %s", df.String())
		}
		logger.Logger.Fatalf("HERE: %s\n%s", obj.String(), blockVarsStr)
	} */

	var deps = []objects.Object{obj}
	var depth = []int{i}

	// indirect dependencies from reference
	for _, ref := range obj.GetVariableInfo().GetReferences() {
		indirectDeps, indirectStr, indirectDepth := getIndirectDependencies(ref, i+1)
		blockVarsStr += indirectStr
		deps = append(deps, indirectDeps...)
		depth = append(depth, indirectDepth...)
	}
	// direct dependencies
	for _, dep := range obj.GetDependencies() {
		directDeps, indirectStr, indirectDepth := getIndirectDependencies(dep, i+1)
		blockVarsStr += indirectStr
		deps = append(deps, directDeps...)
		depth = append(depth, indirectDepth...)
	}

	/* for _, refBy := range v.GetVariableInfo().ReferencedBy {
		opStr := ""
		for _, df := range refBy.GetVariableInfo().GetAllDataflows() {
			s := df.GetOpString() + "(" + df.Datastore + "), "
			if ok := visitedOp[s]; !ok {
				visitedOp[s] = true
				opStr += s
			}
		}
		if opStr != "" {
			opStr = "// " + opStr
		}
		blockVarsStr += fmt.Sprintf("[%s] (%d) (%s) ref_by <%s> %s\n", padding, refBy.GetId(), objects.VariableTypeName(refBy), refBy.String(), opStr)
	} */

	return deps, blockVarsStr, depth
}

func (block *Block) Yaml() ([]string, string) {
	data := []string{}
	visited := make(map[objects.Object]bool)
	allBlockVarsStr := ""
	for _, obj := range block.Objs {
		deps, blockVarsStr, depths := getIndirectDependencies(obj, 0)
		allBlockVarsStr += blockVarsStr + "\n"
		/* slices.Reverse(deps) */
		for i, v := range deps {
			/* lastIndex := len(deps) - 1 */
			if _, isVisited := visited[v]; isVisited {
				continue
			}
			visited[v] = true

			variableString := v.String()
			depth := depths[i]

			padding := ""
			for j := 0; j < depth; j++ {
				padding += "_"
			}

			padding += fmt.Sprintf("%d", i)

			/* if i != lastIndex {
				// last index corresponds to the original variabl from where we got the parameters
				// note that it is the last since the deps slice was reversed
				variableString = fmt.Sprintf("(inline) (%s) %s", objects.VariableTypeName(v), variableString)
			} else {
				variableString = fmt.Sprintf("(%s) %s", objects.VariableTypeName(v), variableString)
			} */
			variableString = fmt.Sprintf("%s depth %d (%s) %s", padding, depth, objects.VariableTypeName(v), variableString)

			dfsWriteOps := v.GetVariableInfo().GetAllWriteDataflows()
			dfsReadOps := v.GetVariableInfo().GetAllReadDataflows()
			if len(dfsWriteOps) > 0 || len(dfsReadOps) > 0 {
				writeOpsLst, writeOpsStr := getTaintedFieldsListAndString(dfsWriteOps)
				readOpsLst, readOpsStr := getTaintedFieldsListAndString(dfsReadOps)
				if len(dfsWriteOps) > 0 && len(dfsReadOps) == 0 {
					data = append(data, fmt.Sprintf("%s --> w-tainted %dx: write(%s)", variableString, len(writeOpsLst), writeOpsStr))
				} else if len(dfsWriteOps) == 0 && len(dfsReadOps) > 0 {
					data = append(data, fmt.Sprintf("%s --> r-tainted %dx: read(%s)", variableString, len(readOpsLst), readOpsStr))
				} else {
					data = append(data, fmt.Sprintf("%s --> w-tainted %dx: write(%s) --> r-tainted %dx: read(%s)", variableString, len(writeOpsLst), writeOpsStr, len(readOpsLst), readOpsStr))
				}
			} else {
				data = append(data, variableString)
			}
		}

	}
	return data, allBlockVarsStr
}

// --------
// PACKAGES
// --------

// for external packages
func (p *Package) DumpExternalYaml() utils.OrderedProperties {
	propsData := utils.NewOrderedPropertyList()

	// metadata
	propsData.AddOrderedProperty("package", p.Name)
	propsData.AddOrderedProperty("package path", p.PackagePath)
	propsData.AddOrderedProperty("module", p.Module)

	// save final data
	return propsData.Result()
}

// for blueprint packages
func (p *Package) DumpBlueprintYaml() utils.OrderedProperties {
	propsData := utils.NewOrderedPropertyList()

	// metadata
	propsData.AddOrderedProperty("package", p.Name)
	propsData.AddOrderedProperty("package path", p.PackagePath)
	propsData.AddOrderedProperty("module", p.Module)

	// declared types
	declaredTypes := []string{}
	for _, e := range p.DeclaredTypes {
		declaredTypes = append(declaredTypes, e.LongString())
	}
	sort.Strings(declaredTypes)
	propsData.AddOrderedProperty("declared types", declaredTypes)

	// datastore types
	datastoreTypes := []string{}
	for _, e := range p.DatastoreTypes {
		datastoreTypes = append(datastoreTypes, e.LongString())
	}
	sort.Strings(datastoreTypes)
	propsData.AddOrderedProperty("datastore types", datastoreTypes)

	// save final data
	return propsData.Result()
}

func (p *Package) DumpYaml() utils.OrderedProperties {
	propsData := utils.NewOrderedPropertyList()

	// metadata
	propsData.AddOrderedProperty("package", p.Name)
	propsData.AddOrderedProperty("package path", p.PackagePath)
	propsData.AddOrderedProperty("module", p.Module)
	// metadata > files
	files := []string{}
	for _, f := range p.Files {
		files = append(files, f.String())
	}
	sort.Strings(files)
	propsData.AddOrderedProperty("files", files)

	// metadata > imports (extracted from package)
	imports := []string{}
	for key := range p.ImportedPackages {
		imports = append(imports, key)
	}
	sort.Strings(imports)
	propsData.AddOrderedProperty("imports (extracted from package)", imports)

	// metadata > imports (extracted from source files)
	imports = []string{}
	for alias, pkg := range p.ImportedPackagesByAlias {
		if alias == pkg.GetPackagePath() {
			imports = append(imports, alias)
		} else {
			imports = append(imports, fmt.Sprintf("(%s) %s", alias, pkg.GetPackagePath()))
		}
	}
	sort.Strings(imports)
	propsData.AddOrderedProperty("imports (extracted from source files)", imports)

	// imported types
	importedTypes := []string{}
	for _, e := range p.ImportedTypes {
		importedTypes = append(importedTypes, e.LongString())
	}
	sort.Strings(importedTypes)
	propsData.AddOrderedProperty("imported types", importedTypes)

	// imported constants
	importedConstants := []string{}
	for _, e := range p.ImportedConstants {
		importedConstants = append(importedConstants, e.LongString())
	}
	sort.Strings(importedConstants)
	propsData.AddOrderedProperty("imported constants", importedConstants)

	// imported types
	importedVariables := []string{}
	for _, e := range p.ImportedVariables {
		importedVariables = append(importedVariables, e.LongString())
	}
	sort.Strings(importedVariables)
	propsData.AddOrderedProperty("imported variables", importedVariables)

	// declared types
	declaredTypes := []string{}
	for _, e := range p.DeclaredTypes {
		declaredTypes = append(declaredTypes, e.LongString())
	}
	sort.Strings(declaredTypes)
	propsData.AddOrderedProperty("declared types", declaredTypes)

	// declared variables
	declaredVariables := []string{}
	for _, v := range p.DeclaredVariables {
		declaredVariables = append(declaredVariables, v.String())
	}
	sort.Strings(declaredVariables)
	propsData.AddOrderedProperty("declared variables", declaredVariables)

	// declared constants
	declaredConstants := []string{}
	for _, v := range p.DeclaredConstants {
		declaredConstants = append(declaredConstants, v.String())
	}
	sort.Strings(declaredConstants)
	propsData.AddOrderedProperty("declared constants", declaredConstants)

	// service types
	serviceTypes := []string{}
	for _, e := range p.ServiceTypes {
		serviceTypes = append(serviceTypes, e.LongString())
	}
	sort.Strings(serviceTypes)
	propsData.AddOrderedProperty("service types", serviceTypes)

	// datastore types
	datastoreTypes := []string{}
	for _, e := range p.DatastoreTypes {
		datastoreTypes = append(datastoreTypes, e.LongString())
	}
	sort.Strings(datastoreTypes)
	propsData.AddOrderedProperty("datastore types", datastoreTypes)

	// parsed methods
	parsedMethods := []string{}
	for _, m := range p.ParsedMethods {
		parsedMethods = append(parsedMethods, m.LongString())
	}
	sort.Strings(parsedMethods)
	propsData.AddOrderedProperty("parsed methods", parsedMethods)

	// save final data
	return propsData.Result()
}
