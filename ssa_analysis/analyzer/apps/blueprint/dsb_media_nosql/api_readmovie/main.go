package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/examples/dsb_media_nosql/workflow/mediamicroservices_nosql"
	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()
	var database backend.NoSQLDatabase
	var cache backend.Cache
	service, _ := mediamicroservices_nosql.NewPlotServiceImpl(ctx, database, cache)
	var reqID int64
	var plotID int64
	var text string
	service.WritePlot(ctx, reqID, plotID, text)
}
