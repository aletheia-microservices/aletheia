package tainter

import (
	"fmt"
	"strings"

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
	dbpath  string                 // with full path: <database>.<table>.<fieldname>[.<any sub path>]
	dbcall  *ssagraph.DatabaseCall //currently this is only used to get the ID later in the parser of abstractcallgraph
	readKey bool                   // aka filter key
	readVal bool                   // aka retrived value
}

type SVTaint struct {
	svpath string                // with full path: <method param name>[.<any sub path>]
	svcall *ssagraph.ServiceCall //currently this is only used to get the ID later in the parser of abstractcallgraph
}

type TaintInfoData struct {
	objpath  string
	objval   ssa.Value
	infoType TaintInfoType
	dbTaint  DBTaint
	svTaint  SVTaint
}

type TaintInfo struct {
	TaintInfoData
	prevval ssa.Value // debug purposes
	objroot bool
	callerT string // managed at combiner.go
}

func (ti TaintInfo) String() string {
	if ti.getDatabasePath() != "" {
		return fmt.Sprintf("(_obj%s, %s) (DB)", ti.getObjectPath(), ti.getDatabasePath())
	} else if ti.getServicePath() != "" {
		return fmt.Sprintf("(_obj%s, %s) (SV)", ti.getObjectPath(), ti.getServicePath())
	}
	return fmt.Sprintf("(_obj%s)", ti.getObjectPath())
}

func NewTaintInfoDatabase(dbpath string, path string, val ssa.Value, dbcall *ssagraph.DatabaseCall, readKey bool, readVal bool) TaintInfo {
	return TaintInfo{
		TaintInfoData: TaintInfoData{
			objpath:  path,
			objval:   val,
			infoType: TAINT_INFO_DATABASE,
			dbTaint: DBTaint{
				dbpath:  dbpath,
				dbcall:  dbcall,
				readKey: readKey,
				readVal: readVal,
			},
		},
		objroot: true,
	}
}

func NewTaintInfoService(svpath string, path string, val ssa.Value, svcall *ssagraph.ServiceCall) TaintInfo {
	return TaintInfo{
		TaintInfoData: TaintInfoData{
			objpath:  path,
			objval:   val,
			infoType: TAINT_INFO_SERVICE,
			svTaint: SVTaint{
				svpath: svpath,
				svcall: svcall,
			},
		},
		objroot:  true,
	}
}

func (t TaintInfo) getCallerT() string {
	return t.callerT
}

func (t TaintInfo) isReadKey() bool {
	return t.dbTaint.readKey
}

func (t TaintInfo) isReadValue() bool {
	return t.dbTaint.readVal
}

func (t TaintInfo) setReadKey(readKey bool) {
	t.dbTaint.readKey = readKey
}

func (t TaintInfo) setReadValue(readValue bool) {
	t.dbTaint.readVal = readValue
}

func (t TaintInfo) enableObjectRoot() TaintInfo {
	t.objroot = true
	return t
}

func (t TaintInfo) tryEnableObjectRoot() TaintInfo {
	if t.objpath == "" {
		t.objroot = true
	}
	return t
}

func (t TaintInfo) disableObjectRoot() TaintInfo {
	t.objroot = false
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
	return "_obj" + t.objpath
}

func (t TaintInfo) getObjectPath() string {
	return t.objpath
}

func (t TaintInfo) getDatabasePath() string {
	return t.dbTaint.dbpath
}

func (t TaintInfo) getDatabaseCall() *ssagraph.DatabaseCall {
	return t.dbTaint.dbcall
}

func (t TaintInfo) getServicePath() string {
	return t.svTaint.svpath
}

func (t TaintInfo) getServiceCall() *ssagraph.ServiceCall {
	return t.svTaint.svcall
}

func (t TaintInfo) updateValue(val ssa.Value) TaintInfo {
	t.objval = val
	return t
}

func (t TaintInfo) cutObjectPathSuffix(suffix string) (TaintInfo, bool) {
	var ok bool
	t.objpath, ok = strings.CutSuffix(t.objpath, suffix)
	if !ok {
		// EVAL: fmt.Printf("[TAINTINFO] [WARNING] objectpath (%s) does not contain suffix (%s)\n", t.objpath, suffix)
	}
	return t, ok
}

func (t TaintInfo) cutObjectPathPrefix(prefix string) (TaintInfo, bool) {
	var ok bool
	t.objpath, ok = strings.CutPrefix(t.objpath, prefix)
	if !ok {
		// EVAL: fmt.Printf("[TAINTINFO] [WARNING] objectpath (%s) does not contain suffix (%s)\n", t.objpath, suffix)
	}
	return t, ok
}

func (t TaintInfo) updateObjectPathPrefix(prefix string) TaintInfo {
	t.objpath = prefix + t.objpath
	return t
}

func (t TaintInfo) updateObjectPathSuffix(suffix string) TaintInfo {
	t.objpath = t.objpath + suffix
	return t
}

func (t TaintInfo) setObjectPath(new string) TaintInfo {
	t.objpath = new
	return t
}

func (t TaintInfo) updateCallPathSuffix(suffix string) TaintInfo {
	if t.infoType == TAINT_INFO_DATABASE {
		t.dbTaint.dbpath = t.dbTaint.dbpath + suffix
	} else if t.infoType == TAINT_INFO_SERVICE {
		t.svTaint.svpath = t.svTaint.svpath + suffix
	}
	return t
}

// full path: <database>.<table>.<fieldname>[.<any sub path> or [<any sub path>]
func (t TaintInfo) updateCallPathPrefix(prefix string) TaintInfo {
	if t.infoType == TAINT_INFO_DATABASE {
		t.dbTaint.dbpath = insertAfterFieldName(t.dbTaint.dbpath, prefix)
	} else if t.infoType == TAINT_INFO_SERVICE {
		t.svTaint.svpath = insertAfterFieldName(t.svTaint.svpath, prefix)
	}
	return t
}

// full path: <database>.<table>.<fieldname>[.<any sub path> or [<any sub path>]
func insertAfterFieldName(path, prefix string) string {
	// find first '.'
	first := strings.IndexByte(path, '.')
	// find second '.' (after first)
	second := strings.IndexByte(path[first+1:], '.') + first + 1

	// fieldname starts at second+1
	// ends at next '.' or '[' or end
	i := second + 1
	for i < len(path) && path[i] != '.' && path[i] != '[' {
		i++
	}
	// insert prefix right after fieldname
	return path[:i] + prefix + path[i:]
}
