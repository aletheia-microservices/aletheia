package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/examples/dsb_sn/workflow/socialnetwork"
	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()

	var postCache backend.Cache
	var postDB backend.NoSQLDatabase
	postStorageService, _ := socialnetwork.NewPostStorageServiceImpl(ctx, postCache, postDB)

	var reqID int64
	var post socialnetwork.Post
	postStorageService.StorePost(ctx, reqID, post)
}
