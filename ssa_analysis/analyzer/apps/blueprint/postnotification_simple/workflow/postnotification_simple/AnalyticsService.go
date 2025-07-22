package postnotification_simple

import (
	"context"
	"sync"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple/common"
)

type AnalyticsService interface {
	Run(ctx context.Context) error
	ReadAnalytics(ctx context.Context, postID int64) (Analytics, error)
}

type AnalyticsServiceImpl struct {
	analyticsQueue backend.Queue
	analyticsDb    backend.NoSQLDatabase
	numWorkers     int
}

func NewAnalyticsServiceImpl(ctx context.Context, analyticsDb backend.NoSQLDatabase, analyticsQueue backend.Queue) (AnalyticsService, error) {
	n := &AnalyticsServiceImpl{analyticsDb: analyticsDb, analyticsQueue: analyticsQueue, numWorkers: 4}
	return n, nil
}

func (a *AnalyticsServiceImpl) ReadAnalytics(ctx context.Context, postID int64) (Analytics, error) {
	var analytics Analytics
	collection, err := a.analyticsDb.GetCollection(ctx, "analyticsDb", "analytics_collection")
	if err != nil {
		return analytics, err
	}
	analyticsQuery := bson.D{{Key: "postid", Value: postID}}
	result, err := collection.FindOne(ctx, analyticsQuery)
	if err != nil {
		return analytics, err
	}
	res, err := result.One(ctx, &analytics)
	if !res || err != nil {
		return analytics, err
	}
	return analytics, nil
}

func (a *AnalyticsServiceImpl) handleMessage(ctx context.Context, message TriggerAnalyticsMessage) error {
	postID, err := common.StringToInt64(message.PostID)
	if err != nil {
		return err
	}

	collection, err := a.analyticsDb.GetCollection(ctx, "analyticsDb", "analytics_collection")
	if err != nil {
		return err
	}

	analytics := Analytics{
		PostID: postID,
	}
	return collection.InsertOne(ctx, analytics)
}

func (n *AnalyticsServiceImpl) workerThread(ctx context.Context) error {
	var forever chan struct{}
	go func() {
		var analyticsEvent TriggerAnalyticsMessage
		n.analyticsQueue.Pop(ctx, &analyticsEvent)
		n.handleMessage(ctx, analyticsEvent)
	}()
	<-forever
	return nil
}

func (n *AnalyticsServiceImpl) Run(ctx context.Context) error {
	backend.GetLogger().Info(ctx, "initializing %d workers", n.numWorkers)
	var wg sync.WaitGroup
	wg.Add(n.numWorkers)
	for i := 1; i <= n.numWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			err := n.workerThread(ctx)
			if err != nil {
				backend.GetLogger().Error(ctx, "error in worker thread: %s", err.Error())
				panic(err)
			}
		}(i)
	}
	wg.Wait()
	backend.GetLogger().Info(ctx, "joining %d workers", n.numWorkers)
	return nil
}
