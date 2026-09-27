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

// same format as GetArgumentName; we want the path after the ssa name, with _obj in its place
// examples:
// - MovieIdService.RegisterMovieId.t4 => _obj
// - MovieIdService.RegisterMovieId.t4.MovieId => _obj.MovieId
// - CastInfoService.ReadCastInfos.t17[*] => _obj[*]
// - CastInfoService.ReadCastInfos.t17[*].CastInfoID => _obj[*].CastInfoID
func (trace *AbstractTrace) ExtractTracedObjectPath() string {
	splits := strings.SplitN(trace.GetServicePath(), ".", 4)
	path := "_obj"

	// handle array case if it exists
	// e.g., CastInfoService.ReadCastInfos.t17[*]...
	// we want to extract [*]...
	arraySplits := strings.SplitN(splits[2], "[*]", 2)
	if len(arraySplits) > 1 {
		path += "[*]" + arraySplits[1]
	}

	// handle sub path if it exists (everything after the ssa name)
	// e.g., CastInfoService.ReadCastInfos.t17[*].CastInfoID
	// splits[3] is CastInfoID (SplitN drops the dots), so we append "." + splits[3] = .CastInfoID
	if len(splits) > 3 {
		path += "." + splits[3]
	}
	return path
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
	return trace.svpath == other.svpath && trace.svcallID == other.svcallID
}

func (trace *AbstractTrace) IsUpperPath(other *AbstractTrace) (bool, string) {
	if trace.svpath != other.svpath && strings.HasPrefix(other.svpath, trace.svpath) {
		var subpath string
		_, subpath, _ = strings.Cut(other.svpath, trace.svpath)
		return trace.svcallID == other.svcallID, subpath
	}
	return false, ""
}
