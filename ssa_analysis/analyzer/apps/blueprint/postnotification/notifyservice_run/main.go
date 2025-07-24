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

	var notificationsQueue backend.Queue
	notifyService, _ := postnotification_simple.NewNotifyServiceImpl(ctx, storageService, notificationsQueue)
	
	notifyService.Run(ctx)
}
