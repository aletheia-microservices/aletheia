package abstractgraph

import (
	"fmt"
	"strings"
)

type AbstractTrace struct {
	svpath   string
	svcallID string
}

func NewAbstractTrace(svpath string, svcallID string) *AbstractTrace {
	return &AbstractTrace{
		svpath:   svpath,
		svcallID: svcallID,
	}
}

// [TO BE IMPROVED]
// format: <service>.<method>.<ssa name>[.<any sub path>]
// e.g., MovieIdService.RegisterMovieId.t4
// e.g., MovieIdService.RegisterMovieId.t4.MovieId
func (trace *AbstractTrace) GetArgumentName() string {
	splits := strings.Split(trace.GetServicePath(), ".")
	return splits[2]
}

func (trace *AbstractTrace) GetArgumentPath() string {
	splits := strings.SplitN(trace.GetServicePath(), ".", 4)
	if len(splits) > 3 {
		return "_obj." + splits[3]
	}
	return "_obj"
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
	fmt.Printf("[ABSTRACT TRACE] [EQUAL] checking if traces are equal:\n\t%s\n\t%s\n", trace.LongString(), other.LongString())
	return trace.svpath == other.svpath && trace.svcallID == other.svcallID
}

func (trace *AbstractTrace) IsUpperPath(other *AbstractTrace) (bool, string) {
	fmt.Printf("[ABSTRACT TRACE] [SUPER] checking if trace is super path:\n\t%s\n\t%s\n", trace.LongString(), other.LongString())
	if trace.svpath != other.svpath && strings.HasPrefix(other.svpath, trace.svpath) {
		var subpath string
		_, subpath, _ = strings.Cut(other.svpath, trace.svpath)
		fmt.Printf("got subpath: %s\n", subpath)
		return trace.svcallID == other.svcallID, subpath
	}
	return false, ""
}
