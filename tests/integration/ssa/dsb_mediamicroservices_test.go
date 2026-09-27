package ssa

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// MovieIdService.RegisterMovieId builds a movie from its parameters and inserts it: each parameter
// is tainted with the field it was assigned to, and the request id is not written at all
func TestMediaMicroservicesWriteFieldsFromParameters(t *testing.T) {
	a := runner.Get(t, "dsb_mediamicroservices")
	graph := getSSAGraph(t, a, "mediamicroservices.MovieIdService.RegisterMovieId")

	call := getDatabaseCall(t, graph, "movie_id_db.movie", "InsertOne")
	assertDatabaseTaint(t, getParam(t, graph, "movieID"), dbTaint{objpath: "_obj", dbpath: "movie_id_db.movie.MovieID", op: common.OP_WRITE, t: call.GetT()})
	assertDatabaseTaint(t, getParam(t, graph, "title"), dbTaint{objpath: "_obj", dbpath: "movie_id_db.movie.Title", op: common.OP_WRITE, t: call.GetT()})
	assertNotTaintedBy(t, getParam(t, graph, "reqID"), "movie_id_db")

	movie := call.GetArguments()[0]
	assertDatabaseTaint(t, movie, dbTaint{objpath: "_obj", dbpath: "movie_id_db.movie", op: common.OP_WRITE})
	assertDatabaseTaint(t, movie, dbTaint{objpath: "_obj.Title", dbpath: "movie_id_db.movie.Title", op: common.OP_WRITE})
}

// PlotService.ReadPlot first looks up the plot in the cache and, on a miss, reads it from the database:
// the plot id is the read key of both calls, each with the t of its call
func TestMediaMicroservicesCacheThenDatabaseRead(t *testing.T) {
	a := runner.Get(t, "dsb_mediamicroservices")
	graph := getSSAGraph(t, a, "mediamicroservices.PlotService.ReadPlot")

	get := getDatabaseCall(t, graph, "plot_cache.*", "Get")
	find := getDatabaseCall(t, graph, "plot_db.plot", "FindOne")
	plotID := getParam(t, graph, "plotID")
	assertDatabaseTaint(t, plotID, dbTaint{objpath: "_obj", dbpath: "plot_cache.*.Key", op: common.OP_READ, readKey: true, t: get.GetT()})
	assertDatabaseTaint(t, plotID, dbTaint{objpath: "_obj", dbpath: "plot_db.plot.PlotID", op: common.OP_READ, readKey: true, t: find.GetT()})
	assertNotTaintedBy(t, getParam(t, graph, "reqID"), "plot")
}

// ComposeReviewService.UploadRating stores the rating in the cache under the request id: the request id
// is the written key and the rating the written value
func TestMediaMicroservicesCacheWriteKeyAndValue(t *testing.T) {
	a := runner.Get(t, "dsb_mediamicroservices")
	graph := getSSAGraph(t, a, "mediamicroservices.ComposeReviewService.UploadRating")

	put := getDatabaseCall(t, graph, "compose_review_cache.rating", "Put")
	if put.GetOpType() != common.OP_WRITE {
		t.Fatalf("Put op type = %s, want write", common.OperationTypeToString(put.GetOpType()))
	}
	assertDatabaseTaint(t, getParam(t, graph, "reqID"), dbTaint{objpath: "_obj", dbpath: "compose_review_cache.rating.Key", op: common.OP_WRITE, readKey: true, t: put.GetT()})
	assertDatabaseTaint(t, getParam(t, graph, "rating"), dbTaint{objpath: "_obj", dbpath: "compose_review_cache.rating.Value", op: common.OP_WRITE, readVal: true, t: put.GetT()})
}
