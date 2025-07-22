package postnotification_simple

import (
	"context"
	"sync"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"

	//"github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple/common"
)

// does not expose any methods to other services
// it defines Run that runs workers that pull messages from the notificationsQueue
type NotifyService interface {
	Run(ctx context.Context) error
	/* Notify(ctx context.Context, message Message) error */
}

type NotifyServiceImpl struct {
	storageService     StorageService
	notificationsQueue backend.Queue
	num_workers         int
}

func NewNotifyServiceImpl(ctx context.Context, storageService StorageService, notificationsQueue backend.Queue) (NotifyService, error) {
	n := &NotifyServiceImpl{storageService: storageService, notificationsQueue: notificationsQueue, num_workers: 4}
	return n, nil
}

func (n *NotifyServiceImpl) workerThread(ctx context.Context) error {
	var forever chan struct{}
	go func() {
		var workerMessage Message
		n.notificationsQueue.Pop(ctx, &workerMessage)
		n.storageService.ReadPost(ctx, workerMessage.ReqID, workerMessage.PostID_MESSAGE)
	}()
	<-forever
	return nil
}

/*
- Block 0 [Body] (stmt = BlockStmt) with 4 nodes [ExprStmt ValueSpec ExprStmt AssignStmt] and succs [block 3 (ForLoop)]
- Block 1 [ForBody] (stmt = ForStmt) with 1 nodes [GoStmt] and succs [block 4 (ForPost)]
- Block 2 [ForDone] (stmt = ForStmt) with 3 nodes [ExprStmt ExprStmt ReturnStmt] and succs []
- Block 3 [ForLoop] (stmt = ForStmt) with 1 nodes [BinaryExpr] and succs [block 1 (ForBody) block 2 (ForDone)]
- Block 4 [ForPost] (stmt = ForStmt) with 1 nodes [IncDecStmt] and succs [block 3 (ForLoop)]
- Block 5 [Unreachable] (stmt = ReturnStmt) with 0 nodes [] and succs []
*/

/*
- Block 0 [Body] (stmt = BlockStmt) with 4 nodes [ExprStmt ValueSpec ExprStmt AssignStmt] and succs [block 3 (ForLoop)]
- Block 1 [ForBody] (stmt = ForStmt) with 1 nodes [GoStmt] and succs [block 4 (ForPost)]
- Block 2 [ForDone] (stmt = ForStmt) with 3 nodes [ExprStmt ExprStmt ReturnStmt] and succs []
- Block 3 [ForLoop] (stmt = ForStmt) with 1 nodes [BinaryExpr] and succs [block 1 (ForBody) block 2 (ForDone)]
- Block 4 [ForPost] (stmt = ForStmt) with 1 nodes [IncDecStmt] and succs [block 3 (ForLoop)]
- Block 5 [Unreachable] (stmt = ReturnStmt) with 0 nodes [] and succs []
*/


/* 
- Block 0 [Body] (stmt = BlockStmt) with 4 nodes [ExprStmt ValueSpec ExprStmt BinaryExpr] and succs [block 1 (IfThen) block 3 (IfElse)]
- Block 1 [IfThen] (stmt = IfStmt) with 1 nodes [IncDecStmt] and succs [block 2 (IfDone)]
- Block 2 [IfDone] (stmt = IfStmt) with 1 nodes [AssignStmt] and succs [block 6 (ForLoop)]
- Block 3 [IfElse] (stmt = IfStmt) with 1 nodes [IncDecStmt] and succs [block 2 (IfDone)]
- Block 4 [ForBody] (stmt = ForStmt) with 1 nodes [GoStmt] and succs [block 7 (ForPost)]
- Block 5 [ForDone] (stmt = ForStmt) with 3 nodes [ExprStmt ExprStmt ReturnStmt] and succs []
- Block 6 [ForLoop] (stmt = ForStmt) with 1 nodes [BinaryExpr] and succs [block 4 (ForBody) block 5 (ForDone)]
- Block 7 [ForPost] (stmt = ForStmt) with 1 nodes [IncDecStmt] and succs [block 6 (ForLoop)]
- Block 8 [Unreachable] (stmt = ReturnStmt) with 0 nodes [] and succs [] 
*/

func (n *NotifyServiceImpl) Run(ctx context.Context) error {
	backend.GetLogger().Info(ctx, "initializing %d workers", n.num_workers)
	var wg sync.WaitGroup
	wg.Add(n.num_workers)
	/* if n.num_workers > 4 {
		n.num_workers--
	} else {
		n.num_workers++
	} */
	
	fn := func(i int) {
		defer wg.Done()
		err := n.workerThread(ctx)
		if err != nil {
			backend.GetLogger().Error(ctx, "error in worker thread (%d): %s", i, err.Error())
			panic(err)
		}
	}
	
	for i := 1; i <= n.num_workers; i++ { // for 
		go fn(i)
	}
	wg.Wait()
	backend.GetLogger().Info(ctx, "joining %d workers", n.num_workers)
	return nil
}
