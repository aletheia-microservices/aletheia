package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/examples/dsb_media_sql/workflow/mediamicroservices_sql"
	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()

	var movieIdDB backend.RelationalDB
	var movieInfoDB backend.RelationalDB

	movieIdService, _ := mediamicroservices_sql.NewMovieIdServiceImpl(ctx, movieIdDB)
	movieInfoService, _ := mediamicroservices_sql.NewMovieInfoServiceImpl(ctx, movieInfoDB)
	api, _ := mediamicroservices_sql.NewAPIServiceImpl(ctx, movieIdService, movieInfoService)

	var reqID int64
	var movieID string
	api.ReadMovie(ctx, reqID, movieID)
}
