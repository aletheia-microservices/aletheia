package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"github.com/blueprint-uservices/blueprint/examples/large_scale_app/workflow/large_scale_app"
)

func main() {
	ctx := context.Background()
	
	var service100DB backend.NoSQLDatabase
	service100, _ := large_scale_app.NewService500Impl(ctx, service100DB)

	var id string
	var data string
	var datatwo string
	service100.Method500(ctx, id, data, datatwo)
}
