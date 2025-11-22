package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/examples/foobar2/workflow/foobar2"
	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()

	var db backend.NoSQLDatabase
	service, _ := foobar2.NewRouteServiceImpl(ctx, db)

	var id string
	service.ReadRoute(ctx, id)
}
