package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"github.com/blueprint-uservices/blueprint/examples/large_scale_app_C/workflow/large_scale_app_C"
)

func main() {
	ctx := context.Background()
	
	var db backend.NoSQLDatabase
	service, _ := large_scale_app_C.NewService121Impl(ctx, db)

	var id string
	var data string
	service.Method121(ctx, id, data)
}
