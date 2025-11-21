package abstractgraph

import (
	"fmt"
	"strings"
)

type AbstractTrace struct {
	t        string // format: <ssa_variable_name> (only when primary!!)
	svpath   string
	svcallID string
}

func NewAbstractTrace(t string, svpath string, svcallID string) *AbstractTrace {
	return &AbstractTrace{
		t:        t,
		svpath:   svpath,
		svcallID: svcallID,
	}
}

// format: <service>.<method>.<ssa name>[.<any sub path>]
// examples:
// - MovieIdService.RegisterMovieId.t4
// - MovieIdService.RegisterMovieId.t4.MovieId
// - MovieInfoService.ReadMovieInfo.t4.Casts[*].CastInfoID
// we want to extract t4
func (trace *AbstractTrace) GetArgumentName() string {
	splits := strings.Split(trace.GetServicePath(), ".")

	// handle array case if it exists
	// e.g., CastInfoService.ReadCastInfos.t17[*]...
	// we want to extract t17
	arraySplits := strings.Split(splits[2], "[*]")

	return arraySplits[0]
}

func (trace *AbstractTrace) GetArgumentPath() string {
	splits := strings.SplitN(trace.GetServicePath(), ".", 4)
	if len(splits) > 3 {
		return "_obj." + splits[3]
	}

	// handle array case if it exists
	// e.g., CastInfoService.ReadCastInfos.t17[*]...
	// we want to extract [*]...
	arraySplits := strings.SplitN(splits[2], "[*]", 2)
	if len(arraySplits) > 1 {
		return "_obj[*]" + arraySplits[1]
	}
	return "_obj"
}

func (trace *AbstractTrace) GetT() string {
	return trace.t
}

func (trace *AbstractTrace) GetServicePath() string {
	return trace.svpath
}

func (trace *AbstractTrace) GetServiceCallID() string {
	return trace.svcallID
}

func (trace *AbstractTrace) String() string {
	return trace.svpath
}

func (trace *AbstractTrace) LongString() string {
	return fmt.Sprintf("{%s, %s, rpc}", trace.svpath, trace.svcallID)
}

func (trace *AbstractTrace) Equals(other *AbstractTrace) bool {
	// EVAL: fmt.Printf("[ABSTRACT TRACE] [EQUAL] checking if traces are equal:\n\t%s\n\t%s\n", trace.LongString(), other.LongString())
	return trace.svpath == other.svpath && trace.svcallID == other.svcallID
}

func (trace *AbstractTrace) IsUpperPath(other *AbstractTrace) (bool, string) {
	// EVAL: fmt.Printf("[ABSTRACT TRACE] [SUPER] checking if trace is super path:\n\t%s\n\t%s\n", trace.LongString(), other.LongString())
	if trace.svpath != other.svpath && strings.HasPrefix(other.svpath, trace.svpath) {
		var subpath string
		_, subpath, _ = strings.Cut(other.svpath, trace.svpath)
		// EVAL: fmt.Printf("got subpath: %s\n", subpath)
		return trace.svcallID == other.svcallID, subpath
	}
	return false, ""
}
