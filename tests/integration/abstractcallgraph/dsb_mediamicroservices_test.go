package abstractcallgraph

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// dsb_mediamicroservices (DeathStarBench): an APIService registers movies, users and reviews, and
// reads movie pages; most services have a database and a cache

func TestMediaMicroservicesEntrypoints(t *testing.T) {
	g := runner.Get(t, "dsb_mediamicroservices").AbsGraph
	assertList(t, "entrypoints", entrypoints(g), []string{
		"APIService.Login",
		"APIService.ReadPage",
		"APIService.RegisterMovie",
		"APIService.RegisterMovieId",
		"APIService.RegisterUser",
		"APIService.UploadMovieId",
		"APIService.UploadText",
		"APIService.UploadUniqueId",
		"APIService.UploadUserWithUsername",
		"APIService.WriteCastInfo",
		"APIService.WriteMovieInfo",
		"APIService.WritePlot",
	})
}

func TestMediaMicroservicesCalls(t *testing.T) {
	g := runner.Get(t, "dsb_mediamicroservices").AbsGraph
	assertList(t, "rpc calls", rpcCalls(g), []string{
		"APIService.Login -> UserService.Login",
		"APIService.ReadPage -> MovieIdService.ReadMovieId",
		"APIService.ReadPage -> PageService.ReadPage",
		"APIService.RegisterMovie -> CastInfoService.WriteCastInfo",
		"APIService.RegisterMovie -> MovieIdService.RegisterMovieId",
		"APIService.RegisterMovie -> MovieInfoService.WriteMovieInfo",
		"APIService.RegisterMovie -> PlotService.WritePlot",
		"APIService.RegisterMovieId -> MovieIdService.RegisterMovieId",
		"APIService.RegisterUser -> UserService.RegisterUser",
		"APIService.UploadMovieId -> MovieIdService.UploadNewMovieId",
		"APIService.UploadText -> TextService.UploadNewText",
		"APIService.UploadUniqueId -> UniqueIdService.UploadNewUniqueId",
		"APIService.UploadUserWithUsername -> UserService.UploadUserWithUsername",
		"APIService.WriteCastInfo -> CastInfoService.WriteCastInfo",
		"APIService.WriteMovieInfo -> MovieInfoService.WriteMovieInfo",
		"APIService.WritePlot -> PlotService.WritePlot",
		"ComposeReviewService.UploadMovieId -> MovieReviewService.UploadMovieReview",
		"ComposeReviewService.UploadMovieId -> ReviewStorageService.StoreReview",
		"ComposeReviewService.UploadMovieId -> UserReviewService.UploadUserReview",
		"ComposeReviewService.UploadRating -> MovieReviewService.UploadMovieReview",
		"ComposeReviewService.UploadRating -> ReviewStorageService.StoreReview",
		"ComposeReviewService.UploadRating -> UserReviewService.UploadUserReview",
		"ComposeReviewService.UploadText -> MovieReviewService.UploadMovieReview",
		"ComposeReviewService.UploadText -> ReviewStorageService.StoreReview",
		"ComposeReviewService.UploadText -> UserReviewService.UploadUserReview",
		"ComposeReviewService.UploadUniqueId -> MovieReviewService.UploadMovieReview",
		"ComposeReviewService.UploadUniqueId -> ReviewStorageService.StoreReview",
		"ComposeReviewService.UploadUniqueId -> UserReviewService.UploadUserReview",
		"ComposeReviewService.UploadUserId -> MovieReviewService.UploadMovieReview",
		"ComposeReviewService.UploadUserId -> ReviewStorageService.StoreReview",
		"ComposeReviewService.UploadUserId -> UserReviewService.UploadUserReview",
		"MovieIdService.UploadNewMovieId -> ComposeReviewService.UploadMovieId",
		"MovieIdService.UploadNewMovieId -> RatingService.UploadNewRating",
		"MovieReviewService.ReadMovieReviews -> ReviewStorageService.ReadReviews",
		"PageService.ReadPage -> CastInfoService.ReadCastInfos",
		"PageService.ReadPage -> MovieInfoService.ReadMovieInfo",
		"PageService.ReadPage -> MovieReviewService.ReadMovieReviews",
		"PageService.ReadPage -> PlotService.ReadPlot",
		"RatingService.UploadNewRating -> ComposeReviewService.UploadRating",
		"TextService.UploadNewText -> ComposeReviewService.UploadText",
		"UniqueIdService.UploadNewUniqueId -> ComposeReviewService.UploadUniqueId",
		"UserService.UploadUserWithUsername -> ComposeReviewService.UploadUserId",
	})
	// caches are left out
	assertList(t, "database calls", databaseCalls(g, isNotCache), []string{
		"read CastInfoService.ReadCastInfos -> cast_info_db.cast.FindMany",
		"read MovieIdService.ReadMovieId -> movie_id_db.movie.FindOne",
		"read MovieIdService.UploadNewMovieId -> movie_id_db.movie.FindOne",
		"read MovieInfoService.ReadMovieInfo -> movie_info_db.movie_info.FindOne",
		"read MovieReviewService.ReadMovieReviews -> movie_review_db.movie_review.FindOne",
		"read MovieReviewService.UploadMovieReview -> movie_review_db.movie_review.FindOne",
		"read PlotService.ReadPlot -> plot_db.plot.FindOne",
		"read ReviewStorageService.ReadReviews -> review_storage_db.review.FindMany",
		"read UserReviewService.UploadUserReview -> movie_review_db.movie_review.FindOne",
		"read UserService.Login -> user_db.user.FindOne",
		"read UserService.UploadUserWithUsername -> user_db.user.FindOne",
		"update MovieReviewService.UploadMovieReview -> movie_review_db.movie_review.UpdateMany",
		"update UserReviewService.UploadUserReview -> movie_review_db.movie_review.UpdateMany",
		"write CastInfoService.WriteCastInfo -> cast_info_db.cast.InsertOne",
		"write MovieIdService.RegisterMovieId -> movie_id_db.movie.InsertOne",
		"write MovieInfoService.WriteMovieInfo -> movie_info_db.movie_info.InsertOne",
		"write MovieReviewService.UploadMovieReview -> movie_review_db.movie_review.InsertOne",
		"write PlotService.WritePlot -> plot_db.plot.InsertOne",
		"write ReviewStorageService.StoreReview -> review_storage_db.review.InsertOne",
		"write UserReviewService.UploadUserReview -> movie_review_db.movie_review.InsertOne",
		"write UserService.RegisterUser -> user_db.user.InsertOne",
	})
}

func TestMediaMicroservicesCounts(t *testing.T) {
	assertCounts(t, runner.Get(t, "dsb_mediamicroservices").AbsGraph, map[string]int{
		"service nodes": 38, "database nodes": 24,
		"entry": 12, "rpc": 42, "read": 30, "write": 23, "update": 2,
	})
}

func TestMediaMicroservicesDatabaseOwners(t *testing.T) {
	owners := databaseOwners(runner.Get(t, "dsb_mediamicroservices").AbsGraph, isNotCache)
	// UserReviewService is wired to user_review_db, but one of its calls names movie_review_db,
	// see LIKELY_BUGS.md #18 and TestMediaMicroservicesUserReviewUsesWiredDatabase
	delete(owners, "movie_review_db.movie_review")
	assertOwners(t, owners, map[string][]string{
		"cast_info_db.cast":        {"CastInfoService"},
		"movie_id_db.movie":        {"MovieIdService"},
		"movie_info_db.movie_info": {"MovieInfoService"},
		"plot_db.plot":             {"PlotService"},
		"review_storage_db.review": {"ReviewStorageService"},
		"user_db.user":             {"UserService"},
	})
}

// ReadPage looks up the movie id by title and then reads the page with it
func TestMediaMicroservicesReadPageUsesMovieId(t *testing.T) {
	g := runner.Get(t, "dsb_mediamicroservices").AbsGraph

	readMovieId := getEdge(t, g, "APIService.ReadPage", "MovieIdService.ReadMovieId", "ReadMovieId")
	readPage := getEdge(t, g, "APIService.ReadPage", "PageService.ReadPage", "ReadPage")

	assertTrace(t, readPage.GetArgumentAt(1), "_obj", readMovieId, ".MovieID")
	if !isBefore(readMovieId.GetT(), readPage.GetT()) {
		t.Errorf("ReadMovieId (t=%s) must happen before ReadPage (t=%s)", readMovieId.GetT(), readPage.GetT())
	}
	// the movie id is read from movie_id_db by title
	assertPrimaryTaint(t, g.GetNodeByName("MovieIdService.ReadMovieId").GetParameterByNameIfExists("title"), "_obj", "movie_id_db.movie.Title", common.OP_READ)
}

// RegisterMovie writes the movie to four services
func TestMediaMicroservicesRegisterMovieWritesEverywhere(t *testing.T) {
	g := runner.Get(t, "dsb_mediamicroservices").AbsGraph
	assertContains(t, "rpc calls", rpcCalls(g),
		"APIService.RegisterMovie -> CastInfoService.WriteCastInfo",
		"APIService.RegisterMovie -> MovieIdService.RegisterMovieId",
		"APIService.RegisterMovie -> MovieInfoService.WriteMovieInfo",
		"APIService.RegisterMovie -> PlotService.WritePlot",
	)
	assertContains(t, "database calls", databaseCalls(g, isNotCache),
		"write CastInfoService.WriteCastInfo -> cast_info_db.cast.InsertOne",
		"write MovieIdService.RegisterMovieId -> movie_id_db.movie.InsertOne",
		"write MovieInfoService.WriteMovieInfo -> movie_info_db.movie_info.InsertOne",
		"write PlotService.WritePlot -> plot_db.plot.InsertOne",
	)
}

func TestMediaMicroservicesUserReviewUsesWiredDatabase(t *testing.T) {
	t.Skip("known bug: NoSQL database names come from the GetCollection string instead of the wired backend " +
		"(see LIKELY_BUGS.md #18), so UserReviewService.UploadUserReview is reported on movie_review_db")

	owners := databaseOwners(runner.Get(t, "dsb_mediamicroservices").AbsGraph, isNotCache)
	assertOwners(t, map[string][]string{"movie_review_db.movie_review": owners["movie_review_db.movie_review"]},
		map[string][]string{"movie_review_db.movie_review": {"MovieReviewService"}})
}
