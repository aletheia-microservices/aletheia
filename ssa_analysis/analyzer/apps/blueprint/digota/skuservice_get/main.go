package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/examples/digota/workflow/digota"
	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()

	var db backend.NoSQLDatabase
	service, _ := digota.NewProductServiceImpl(ctx, db)

	var id string
	service.Get(ctx, id)
}
