package detection_test

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/detection/constraints/keycoordination"
	"github.com/aletheia-microservices/aletheia/internal/app"
	"github.com/aletheia-microservices/aletheia/internal/app/backends"
)

// newMediaApp builds the dsb_mediamicroservices schema, where a movie is split across movie_id_db
// and movie_info_db under the same primary key _id
func newMediaApp(primaryKeys bool, foreignKey bool) *app.App {
	a := newApp(map[string]string{"movie_id_db": "NoSQLDatabase", "movie_info_db": "NoSQLDatabase"})
	movieID := field(a, "movie_id_db.movie._id")
	movieInfoID := field(a, "movie_info_db.movie_info._id")
	if primaryKeys {
		addConstraint(backends.CONSTRAINT_PRIMARY, movieID)
		addConstraint(backends.CONSTRAINT_PRIMARY, movieInfoID)
	}
	if foreignKey {
		// both records are written by the same request (0)
		addForeignKey(movieInfoID, movieID, 0)
	}
	return a
}

// readMoviePage returns APIService.ReadPage, which reads both parts of a movie by its _id
func readMoviePage() (*abstractgraph.AbstractNode, []*abstractgraph.AbstractEdge) {
	readInfo := databaseCall("find_movie_info", "FindOne", serviceNode("MovieInfoService", "ReadMovieInfo"), databaseNode("movie_info_db", "movie_info"), common.OP_READ,
		object(
			primary("movie_info_db.movie_info._id", "find_movie_info", common.OP_READ),
			secondary("movie_id_db.movie._id", "find_movie", common.OP_READ),
		))
	readID := databaseCall("find_movie", "FindOne", serviceNode("MovieIdService", "ReadMovieId"), databaseNode("movie_id_db", "movie"), common.OP_READ,
		object(primary("movie_id_db.movie._id", "find_movie", common.OP_READ)))
	return serviceNode("APIService", "ReadPage"), []*abstractgraph.AbstractEdge{readID, readInfo}
}

func TestPrimaryKeyCoordinationReportsReadsBySamePrimaryKey(t *testing.T) {
	a := newMediaApp(true, true)
	d := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_PRIMARY_KEY)
	entry, reads := readMoviePage()

	runRequest(d, a, 1, entry, reads...)

	got := results(d, a)
	wantWarnings(t, got, 1)
	wantLines(t, got,
		"entry request: APIService.ReadPage()",
		"\tPRIMARY KEY READS #1:",
		"\t\tREAD: MovieInfoService.ReadMovieInfo() ... movie_info_db.movie_info.FindOne()",
		"\t\t\t- constraint: PRIMARY KEY (movie_info_db.movie_info._id)",
		"\t\tREAD: MovieIdService.ReadMovieId() ... movie_id_db.movie.FindOne()",
		"\t\t\t- constraint: PRIMARY KEY (movie_id_db.movie._id)",
	)
}

func TestPrimaryKeyCoordinationRequiresPrimaryKeys(t *testing.T) {
	a := newMediaApp(false, true)
	d := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_PRIMARY_KEY)
	entry, reads := readMoviePage()

	runRequest(d, a, 1, entry, reads...)

	wantWarnings(t, results(d, a), 0)
}

func TestPrimaryKeyCoordinationRequiresMandatoryLink(t *testing.T) {
	// with RestrictivePrimaryKeyCoordinationAnalysis, the two parts must be linked by a mandatory foreign key
	a := newMediaApp(true, false)
	d := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_PRIMARY_KEY)
	entry, reads := readMoviePage()

	runRequest(d, a, 1, entry, reads...)

	wantWarnings(t, results(d, a), 0)
}

func TestPrimaryKeyCoordinationTypeString(t *testing.T) {
	if got := keycoordination.NewDetector(keycoordination.DETECTION_TYPE_PRIMARY_KEY).GetTypeString(); got != "primary-key-coordination" {
		t.Errorf("GetTypeString() = %q, want primary-key-coordination", got)
	}
}
