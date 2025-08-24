package tainter

import (
	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/ssagraph"
)

type TaintMode int

const (
	TAINT_MODE_NEARBY TaintMode = iota
	TAINT_MODE_FETCH_UPWARDS
)

type TaintInfoType int

const (
	TAINT_INFO_DATABASE TaintInfoType = iota
	TAINT_INFO_SERVICE
)

type DBTaint struct {
	path string                 // with full path: <database>.<table>.<fieldname>[.<any sub path>]
	call *ssagraph.DatabaseCall //currently this is only used to get the ID later in the parser of abstractcallgraph
}

type SVTaint struct {
	path string                // with full path: <method param name>[.<any sub path>]
	call *ssagraph.ServiceCall //currently this is only used to get the ID later in the parser of abstractcallgraph
}

type TaintInfo struct {
	path     string
	val      ssa.Value
	infoType TaintInfoType
	objroot  bool
	dbTaint  DBTaint
	svTaint  SVTaint
}

func NewTaintInfoDatabase(dbpath string, path string, val ssa.Value, dbcall *ssagraph.DatabaseCall) TaintInfo {
	return TaintInfo{
		path:     path,
		val:      val,
		infoType: TAINT_INFO_DATABASE,
		objroot:  true,
		dbTaint: DBTaint{
			path: dbpath,
			call: dbcall,
		},
	}
}

func NewTaintInfoService(svpath string, path string, val ssa.Value, svcall *ssagraph.ServiceCall) TaintInfo {
	return TaintInfo{
		path:     path,
		val:      val,
		infoType: TAINT_INFO_SERVICE,
		svTaint: SVTaint{
			path: svpath,
			call: svcall,
		},
	}
}

func (t TaintInfo) enableObjectRoot() TaintInfo {
	t.objroot = true
	return t
}

func (t TaintInfo) isObjectRoot() bool {
	return t.objroot
}

func (t TaintInfo) isTypeDatabase() bool {
	return t.infoType == TAINT_INFO_DATABASE
}

func (t TaintInfo) isTypeService() bool {
	return t.infoType == TAINT_INFO_SERVICE
}

func (t TaintInfo) getObjectFullPath() string {
	return "_obj" + t.path
}

func (t TaintInfo) getObjectPath() string {
	return t.path
}

func (t TaintInfo) getDatabasePath() string {
	return t.dbTaint.path
}

func (t TaintInfo) getDatabaseCall() *ssagraph.DatabaseCall {
	return t.dbTaint.call
}

func (t TaintInfo) getServicePath() string {
	return t.svTaint.path
}

func (t TaintInfo) getServiceCall() *ssagraph.ServiceCall {
	return t.svTaint.call
}

func (t TaintInfo) updateValue(val ssa.Value) TaintInfo {
	t.val = val
	return t
}

func (t TaintInfo) updateObjectPathPrefix(prefix string) TaintInfo {
	t.path = prefix + t.path
	return t
}

func (t TaintInfo) updateObjectPathSuffix(suffix string) TaintInfo {
	t.path = t.path + suffix
	return t
}

func (t TaintInfo) updateCallPathSuffix(suffix string) TaintInfo {
	if t.infoType == TAINT_INFO_DATABASE {
		t.dbTaint.path = t.dbTaint.path + suffix
	} else if t.infoType == TAINT_INFO_SERVICE {
		t.svTaint.path = t.svTaint.path + suffix
	}
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
		if taint.path == dbfield {
			return
		}
	}
	t.inheritedTaints[objPath] = append(t.inheritedTaints[objPath], DBTaint{
		path: dbfield,
		call: dbcall,
	})
}

func (t *CheckTaintInfo) addToIndirectTaints(dbfield string, dbcall *ssagraph.DatabaseCall) {
	for _, taint := range t.indirectTaints {
		if taint.path == dbfield {
			return
		}
	}
	t.indirectTaints = append(t.indirectTaints, DBTaint{
		path: dbfield,
		call: dbcall,
	})
}
