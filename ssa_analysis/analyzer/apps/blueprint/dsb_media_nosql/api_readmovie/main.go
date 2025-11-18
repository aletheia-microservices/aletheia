package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/examples/dsb_media_nosql/workflow/mediamicroservices_nosql"
	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()
	var db backend.NoSQLDatabase
	service, _ := mediamicroservices_nosql.NewPlotServiceImpl(ctx, db)
	var reqID int64
	var plotID string
	var text string
	service.WritePlot(ctx, reqID, plotID, text)
}
