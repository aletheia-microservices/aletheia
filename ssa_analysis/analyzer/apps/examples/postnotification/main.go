package main

import (
	"context"
	"fmt"
)

type Post struct {
	ReqID    string
	PostID   string
	Text     string
	Username string
	Mentions []string
}

type Notification struct {
	ReqID  string
	PostID string
}

type Analytics struct {
	ReqID    string
	PostID   string
	Username string
}

type AnalyticsService struct{}

func (s *AnalyticsService) UpdateAnalytics(ctx context.Context, reqID string, postID string) error {
	analytics := Analytics{
		ReqID:  reqID,
		PostID: postID,
	}

	orderDB := &MongoDB{}
	orderDB.Insert(ctx, "analytics_db", "analytics", analytics)

	return nil
}

type MongoDB struct{}

func (m *MongoDB) Insert(ctx context.Context, database string, collection string, document interface{}) error {
	fmt.Printf("[INFO] inserted document: %v\n", document)
	return nil
}

func (m *MongoDB) Find(ctx context.Context, database string, collection string, id string) interface{} {
	fmt.Printf("[INFO] found document for id: %v\n", id)
	return nil
}

type RabbitMQ struct{}

func (r *RabbitMQ) Push(ctx context.Context, database string, topic string, message interface{}) error {
	fmt.Printf("[INFO] pushed message: %v\n", message)
	return nil
}

func StorePost(ctx context.Context, reqID string, postID string, text string, username string, mentions []string) (Post, error) {
	post := Post{
		ReqID:    reqID,
		PostID:   postID,
		Text:     text,
		Username: username,
		Mentions: mentions,
	}

	orderDB := &MongoDB{}
	orderDB.Insert(ctx, "posts_db", "post", post)

	analyticsService := &AnalyticsService{}
	analyticsService.UpdateAnalytics(ctx, reqID, post.PostID)

	notification := Notification{
		ReqID:  reqID,
		PostID: post.PostID,
	}

	rabbitMQ := &RabbitMQ{}
	rabbitMQ.Push(ctx, "notif_queue", "notification", notification)

	return post, nil
}

func main() {
	ctx := context.Background()
	var reqID, postID, text, username string
	var mentions []string
	StorePost(ctx, reqID, postID, text, username, mentions)
}
