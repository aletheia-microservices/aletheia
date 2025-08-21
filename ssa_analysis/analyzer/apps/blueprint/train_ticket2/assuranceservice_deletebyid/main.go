package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"github.com/blueprint-uservices/blueprint/examples/train_ticket2/workflow/train_ticket2"
)

func main() {
	ctx := context.Background()
	
	var assuranceDB backend.NoSQLDatabase
	assuranceService, _ := train_ticket2.NewAssuranceServiceImpl(ctx, assuranceDB)

	var id string
	assuranceService.DeleteById(ctx, id)
}
