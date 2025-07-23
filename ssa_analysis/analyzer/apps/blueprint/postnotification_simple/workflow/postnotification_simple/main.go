package postnotification_simple

import (
	"context"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()
	
	var postsDB backend.NoSQLDatabase
	storageService, _ := NewStorageServiceImpl(ctx, postsDB)

	var reqID int64
	var text string
	storageService.StorePost(ctx, reqID, text)
}
