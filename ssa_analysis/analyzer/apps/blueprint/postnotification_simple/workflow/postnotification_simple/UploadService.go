package postnotification_simple

import (
	"context"
	"math/rand"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"

	"github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple/common"
)

type UploadService interface {
	UploadPost(ctx context.Context, username string, text string) (int64, error)
	DeletePost(ctx context.Context, postID int64) error
	ReadPostWithAnalytics(ctx context.Context, reqID int64, postID int64) (Post, Analytics, error)
}

type UploadServiceImpl struct {
	storageService     StorageService
	analyticsService   AnalyticsService
	notificationsQueue backend.Queue
}

func NewUploadServiceImpl(ctx context.Context, storageService StorageService, analyticsService AnalyticsService, notificationsQueue backend.Queue) (UploadService, error) {
	return &UploadServiceImpl{storageService: storageService, analyticsService: analyticsService, notificationsQueue: notificationsQueue}, nil
}

func (u *UploadServiceImpl) DeletePost(ctx context.Context, postID int64) error {
	return u.storageService.DeletePost(ctx, postID)
}

func (u *UploadServiceImpl) ReadPostWithAnalytics(ctx context.Context, reqID int64, postID int64) (Post, Analytics, error) {
	post, err := u.storageService.ReadPost(ctx, reqID, postID)
	if err != nil {
		return Post{}, Analytics{}, err
	}

	analytics, err := u.analyticsService.ReadAnalytics(ctx, postID)
	if err != nil {
		return Post{}, Analytics{}, err
	}

	return post, analytics, err
}

/* func (u *UploadServiceImpl) UploadPost(ctx context.Context, username string, text string) (int64, error) {
	reqID := rand.Int63()
	postID_UploadSVC := rand.Int63()

	timestamp := rand.Int63()
	mentions := []string{"alice", "bob"}
	post := Post{
		ReqID:     reqID,
		PostID:    postID_UploadSVC,
		Text:      text,
		Mentions:  mentions,
		Timestamp: timestamp,
		Creator: Creator{
			Username: "some username",
		},
	}
	u.storageService.StorePost(ctx, post.ReqID, post)

	message := Message{
		ReqID:          post.ReqID,
		PostID_MESSAGE: post.PostID,
	}
	_, err := u.notificationsQueue.Push(ctx, message)
	if err != nil {
		return 0, err
	}
	return post.PostID, nil
} */

func (u *UploadServiceImpl) UploadPost(ctx context.Context, username string, text string) (int64, error) {
	reqID := rand.Int63()

	common.TestFunc()

	post, _ := u.storageService.StorePost(ctx, reqID, text)

	message := Message{
		ReqID:          reqID,
		PostID_MESSAGE: post.PostID,
	}
	_, err := u.notificationsQueue.Push(ctx, message)
	if err != nil {
		return 0, err
	}
	return post.PostID, nil
}
