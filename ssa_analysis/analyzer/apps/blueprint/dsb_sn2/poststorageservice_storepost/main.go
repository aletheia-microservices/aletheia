package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/examples/dsb_sn2/workflow/socialnetwork2"
	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()

	var postCache backend.Cache
	var postDB backend.NoSQLDatabase
	postStorageService, _ := socialnetwork2.NewPostStorageServiceImpl(ctx, postCache, postDB)

	var reqID int64
	var post socialnetwork2.Post
	postStorageService.StorePost(ctx, reqID, post)
}
