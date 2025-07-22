package postnotification_simple

import (
	"context"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()
	
	var postsDB backend.NoSQLDatabase
	var analyticsQueue backend.Queue
	storageService, _ := NewStorageServiceImpl(ctx, postsDB, analyticsQueue)

	var reqID int64
	var text string
	storageService.StorePost(ctx, reqID, text)
}
