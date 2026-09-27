package detection

import (
	"testing"

	"github.com/aletheia-microservices/aletheia/tests/runner"
)

// EI-1: APIService.ReadPage reads replicated objects by primary key from two services without coordination
func TestMediaMicroservicesPrimaryKeyCoordination(t *testing.T) {
	a := runner.Get(t, "dsb_mediamicroservices")
	assertConstraint(t, a, "PRIMARY KEY (movie_id_db.movie._id)", true)
	assertWarnings(t, a, "primary-key-coordination", 1,
		"entry request: APIService.ReadPage()",
		"READ: MovieInfoService.ReadMovieInfo() ... movie_info_db.movie_info.FindOne()",
		"READ: MovieIdService.ReadMovieId() ... movie_id_db.movie.FindOne()",
	)
}

// Un-1: concurrent RegisterMovie requests may write the same (unique) movie title
func TestMediaMicroservicesUniquenessConcurrency(t *testing.T) {
	a := runner.Get(t, "dsb_mediamicroservices")
	assertConstraint(t, a, "UNIQUE (movie_id_db.movie.Title)", true)
	assertWarnings(t, a, "uniqueness-concurrency", 1,
		"entry request: APIService.RegisterMovie()",
		"- field (constrained): movie_id_db.movie.Title (UNIQUE)",
		"affected write #1: MovieInfoService.WriteMovieInfo() ... movie_info_db.movie_info.InsertOne()",
	)
}

// reviews store the movie and the user they were written for
func TestMediaMicroservicesReviewsReferenceMoviesAndUsers(t *testing.T) {
	a := runner.Get(t, "dsb_mediamicroservices")
	assertConstraint(t, a, "FOREIGN_KEY review_storage_db.review.MovieID REFERENCES movie_id_db.movie._id [T]", true)
	assertConstraint(t, a, "FOREIGN_KEY review_storage_db.review.UserID REFERENCES user_review_db.user_review.UserID [MANDATORY]", true)
}
