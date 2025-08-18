package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/examples/foobar/workflow/foobar"
	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()

	var barDB backend.NoSQLDatabase
	barService, _ := foobar.NewBarServiceImpl(ctx, barDB)

	var id, text string
	barService.WriteBar(ctx, id, text)
}
