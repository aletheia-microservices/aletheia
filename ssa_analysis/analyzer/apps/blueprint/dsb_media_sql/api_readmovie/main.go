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
	var castInfoDB backend.RelationalDB
	var plotDB backend.RelationalDB

	movieIdService, _ := mediamicroservices_sql.NewMovieIdServiceImpl(ctx, movieIdDB)
	movieInfoService, _ := mediamicroservices_sql.NewMovieInfoServiceImpl(ctx, movieInfoDB)
	castInfoService, _ := mediamicroservices_sql.NewCastInfoServiceImpl(ctx, castInfoDB)
	plotService, _ := mediamicroservices_sql.NewPlotServiceImpl(ctx, plotDB)
	api, _ := mediamicroservices_sql.NewAPIServiceImpl(ctx, movieIdService, movieInfoService, castInfoService, plotService)

	var reqID int64
	var title string
	api.ReadPage(ctx, reqID, title)
}
