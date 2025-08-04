package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/examples/digota/workflow/digota"
	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()

	var skusDB backend.NoSQLDatabase
	skuService, _ := digota.NewSkuServiceImpl(ctx, skusDB)

	var id string
	skuService.Get(ctx, id)
}
