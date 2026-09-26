package detection_test

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/detection/constraints/uniquenessconcurrency"
	"github.com/aletheia-microservices/aletheia/internal/app"
	"github.com/aletheia-microservices/aletheia/internal/app/backends"
)

// newMovieRegistrationApp builds the dsb_mediamicroservices schema written by RegisterMovie, where
// movie_id_db.movie.Title is optionally unique
func newMovieRegistrationApp(unique bool) *app.App {
	a := newApp(map[string]string{"movie_id_db": "NoSQLDatabase", "movie_info_db": "NoSQLDatabase"})
	title := field(a, "movie_id_db.movie.Title")
	field(a, "movie_info_db.movie_info.Title")
	if unique {
		addConstraint(backends.CONSTRAINT_UNIQUE, title)
	}
	return a
}

// registerMovie returns APIService.RegisterMovie, which writes the movie title to movie_id_db and
// then, when related is true, writes the movie info derived from that write to movie_info_db
func registerMovie(related bool) (*abstractgraph.AbstractNode, []*abstractgraph.AbstractEdge) {
	writeID := databaseCall("insert_movie", "InsertOne", serviceNode("MovieIdService", "RegisterMovieId"), databaseNode("movie_id_db", "movie"), common.OP_WRITE,
		object(primary("movie_id_db.movie.Title", "insert_movie", common.OP_WRITE)))
	infoTaints := []*abstractgraph.AbstractTaint{primary("movie_info_db.movie_info.Title", "insert_movie_info", common.OP_WRITE)}
	if related {
		infoTaints = append(infoTaints, secondary("movie_id_db.movie.Title", "insert_movie", common.OP_WRITE))
	}
	writeInfo := databaseCall("insert_movie_info", "InsertOne", serviceNode("MovieInfoService", "WriteMovieInfo"), databaseNode("movie_info_db", "movie_info"), common.OP_WRITE,
		object(infoTaints...))
	return serviceNode("APIService", "RegisterMovie"), []*abstractgraph.AbstractEdge{writeID, writeInfo}
}

func TestUniquenessReportsRelatedWriteToAnotherDatabase(t *testing.T) {
	a := newMovieRegistrationApp(true)
	d := uniquenessconcurrency.NewDetector()
	entry, writes := registerMovie(true)

	runRequest(d, a, 0, entry, writes...)

	got := results(d, a)
	wantWarnings(t, got, 1)
	wantLines(t, got,
		"entry request: APIService.RegisterMovie()",
		"write (origin): APIService.RegisterMovie() ... MovieIdService.RegisterMovieId() ... movie_id_db.movie.InsertOne()",
		"\t\t- field (constrained): movie_id_db.movie.Title (UNIQUE)",
		"\t- affected write #1: MovieInfoService.WriteMovieInfo() ... movie_info_db.movie_info.InsertOne()",
	)
}

func TestUniquenessIgnoresWritesWithoutUniqueField(t *testing.T) {
	a := newMovieRegistrationApp(false)
	d := uniquenessconcurrency.NewDetector()
	entry, writes := registerMovie(true)

	runRequest(d, a, 0, entry, writes...)

	wantWarnings(t, results(d, a), 0)
}

func TestUniquenessIgnoresUnrelatedWrite(t *testing.T) {
	a := newMovieRegistrationApp(true)
	d := uniquenessconcurrency.NewDetector()
	entry, writes := registerMovie(false)

	runRequest(d, a, 0, entry, writes...)

	wantWarnings(t, results(d, a), 0)
}

func TestUniquenessIgnoresRelatedWriteInAnotherRequest(t *testing.T) {
	a := newMovieRegistrationApp(true)
	d := uniquenessconcurrency.NewDetector()
	entry, writes := registerMovie(true)

	runRequest(d, a, 0, entry, writes[0])
	runRequest(d, a, 1, serviceNode("MovieInfoService", "WriteMovieInfo"), writes[1])

	wantWarnings(t, results(d, a), 0)
}

func TestUniquenessTypeString(t *testing.T) {
	if got := uniquenessconcurrency.NewDetector().GetTypeString(); got != "uniqueness-concurrency" {
		t.Errorf("GetTypeString() = %q, want uniqueness-concurrency", got)
	}
}

func TestUniquenessReportsRelatedWriteBeforeUniqueWrite(t *testing.T) {
	t.Skip("known bug: uniquenessconcurrency only reports related writes that come after the unique write " +
		"in the request (see LIKELY_BUGS.md #12)")
	a := newMovieRegistrationApp(true)
	d := uniquenessconcurrency.NewDetector()
	// both writes use the same title, so each argument carries a secondary taint from the other write
	writeInfo := databaseCall("insert_movie_info", "InsertOne", serviceNode("MovieInfoService", "WriteMovieInfo"), databaseNode("movie_info_db", "movie_info"), common.OP_WRITE,
		object(
			primary("movie_info_db.movie_info.Title", "insert_movie_info", common.OP_WRITE),
			secondary("movie_id_db.movie.Title", "insert_movie", common.OP_WRITE),
		))
	writeID := databaseCall("insert_movie", "InsertOne", serviceNode("MovieIdService", "RegisterMovieId"), databaseNode("movie_id_db", "movie"), common.OP_WRITE,
		object(
			primary("movie_id_db.movie.Title", "insert_movie", common.OP_WRITE),
			secondary("movie_info_db.movie_info.Title", "insert_movie_info", common.OP_WRITE),
		))

	runRequest(d, a, 0, serviceNode("APIService", "RegisterMovie"), writeInfo, writeID)

	wantWarnings(t, results(d, a), 1)
}
