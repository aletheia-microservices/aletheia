package postnotification_simple

import (
	"context"
	"math/rand"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"go.mongodb.org/mongo-driver/bson"

	//"github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple/common"
)

type StorageService interface {
	//StorePost(ctx context.Context, reqID int64, post Post) error
	StorePost(ctx context.Context, reqID int64, text string) (*Post, error)
	ReadPost(ctx context.Context, reqID int64, postID int64) (Post, error)
	DeletePost(ctx context.Context, postID int64) error
}

type StorageServiceImpl struct {
	postsDb        backend.NoSQLDatabase
}

func NewStorageServiceImpl(ctx context.Context, postsDb backend.NoSQLDatabase) (StorageService, error) {
	s := &StorageServiceImpl{postsDb: postsDb}
	return s, nil
}

/* func (s *StorageServiceImpl) StorePost(ctx context.Context, reqID int64, post Post) error {
	collection, err := s.postsDb.GetCollection(ctx, "post", "post")
	if err != nil {
		return err
	}
	err = collection.InsertOne(ctx, post)
	return err
} */

func (s *StorageServiceImpl) DeletePost(ctx context.Context, postID int64) error {
	collection, err := s.postsDb.GetCollection(ctx, "post", "post")
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "postid", Value: postID}}
	err = collection.DeleteOne(ctx, filter)
	return err
}

func (s *StorageServiceImpl) StorePost(ctx context.Context, reqID int64, text string) (*Post, error) {
	postID_STORAGE_SVC := rand.Int63()
	timestamp := rand.Int63()
	mentions := []string{"alice", "bob"}

	post := &Post{
		ReqID:     reqID,
		PostID:    postID_STORAGE_SVC,
		Text:      text,
		Mentions:  mentions,
		Timestamp: timestamp,
		Creator: Creator{
			Username: "some username",
		},
	}

	/* myval := 0
	var mymentions []string
	for idx, mention := range mentions {
		myval += idx
		mymentions = append(mymentions, mention)
	} */

	collection, err := s.postsDb.GetCollection(ctx, "post", "post")
	if err != nil {
		return nil, err
	}
	err = collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}

	/* message := TriggerAnalyticsMessage{
		PostID: common.Int64ToString(post.PostID),
	}
	_, err = s.analyticsQueue.Push(ctx, message) */

	return post, err
}

func (s *StorageServiceImpl) ReadPost(ctx context.Context, reqID int64, postID int64) (Post, error) {
	var post Post
	collection, err := s.postsDb.GetCollection(ctx, "post", "post")
	if err != nil {
		return post, err
	}
	query := bson.D{{Key: "postid", Value: postID}}
	result, err := collection.FindOne(ctx, query)
	if err != nil {
		return post, err
	}
	res, err := result.One(ctx, &post)
	if !res || err != nil {
		return post, err
	}
	
	return post, err
}
