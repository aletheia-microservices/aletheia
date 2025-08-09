package tainter

import (
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
)

type TaintMode int

const (
	PROPAGATE_TAINT_NEARBY TaintMode = iota
	PROPAGATE_TAINT_FETCH_UPWARDS
)

type DBTaint struct {
	dbfield string
	dbcall  *ssagraph.DatabaseCall //currently this is only used to get the ID later in the parser of abstractcallgraph
}

type TaintInfo struct {
	path    string
	val     ssa.Value
	dbTaint DBTaint
}

func NewTaintInfo(dbfield string, path string, val ssa.Value, dbcall *ssagraph.DatabaseCall) TaintInfo {
	return TaintInfo{
		path: path,
		val:  val,
		dbTaint: DBTaint{
			dbfield: dbfield,
			dbcall:  dbcall,
		},
	}
}

func (t TaintInfo) getObjectFullPath() string {
	return "_obj" + t.path
}

func (t TaintInfo) getObjectPath() string {
	return t.path
}

func (t TaintInfo) getPath() string {
	return t.path
}

func (t TaintInfo) getDatabaseField() string {
	return t.dbTaint.dbfield
}

func (t TaintInfo) getDbCall() *ssagraph.DatabaseCall {
	return t.dbTaint.dbcall
}

func (t TaintInfo) updateValue(val ssa.Value) TaintInfo {
	t.val = val
	return t
}

func (t TaintInfo) updatePathPrefix(prefix string) TaintInfo {
	t.path = prefix + t.path
	return t
}

func (t TaintInfo) updatePathSuffix(suffix string) TaintInfo {
	t.path = t.path + suffix
	return t
}

func (t TaintInfo) updateFieldSuffix(prefix string) TaintInfo {
	t.dbTaint.dbfield = t.dbTaint.dbfield + prefix
	return t
}

func (t TaintInfo) updatePathSuffixAndField(prefix string) TaintInfo {
	t.path = t.path + prefix
	t.dbTaint.dbfield = t.dbTaint.dbfield + prefix
	return t
}

type CheckTaintInfo struct {
	indirectTaints  []DBTaint
	inheritedTaints map[string][]DBTaint
}

func NewCheckTaintInfo() *CheckTaintInfo {
	return &CheckTaintInfo{
		inheritedTaints: make(map[string][]DBTaint),
	}
}

func (t *CheckTaintInfo) addToInheritedTaints(objPath string, dbfield string, dbcall *ssagraph.DatabaseCall) {
	for _, taint := range t.inheritedTaints[objPath] {
		if taint.dbfield == dbfield {
			return
		}
	}
	t.inheritedTaints[objPath] = append(t.inheritedTaints[objPath], DBTaint{
		dbfield: dbfield,
		dbcall:  dbcall,
	})
}

func (t *CheckTaintInfo) addToIndirectTaints(dbfield string, dbcall *ssagraph.DatabaseCall) {
	for _, taint := range t.indirectTaints {
		if taint.dbfield == dbfield {
			return
		}
	}
	t.indirectTaints = append(t.indirectTaints, DBTaint{
		dbfield: dbfield,
		dbcall:  dbcall,
	})
}
