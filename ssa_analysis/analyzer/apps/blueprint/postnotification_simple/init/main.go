package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple"
)

func main() {
	ctx := context.Background()
	
	var postsDB backend.NoSQLDatabase
	storageService, _ := postnotification_simple.NewStorageServiceImpl(ctx, postsDB)

	var reqID int64
	var text string
	storageService.StorePost(ctx, reqID, text)
}
