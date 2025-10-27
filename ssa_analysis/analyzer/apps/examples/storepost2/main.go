package main

import (
	"context"
	"math/rand"
)

type MongoInserter interface {
	Insert(ctx context.Context, document interface{}) error
}

type MongoDB struct{}

func (m *MongoDB) Insert(ctx context.Context, document interface{}) error {
	//EVAL - fmt.Printf("[INFO] inserted document: %v\n", document)
	return nil
}

type RabbitMQPusher interface {
	Push(ctx context.Context, document interface{}) error
}

type RabbitMQ struct{}

func (r *RabbitMQ) Push(ctx context.Context, message interface{}) error {
	//EVAL - fmt.Printf("[INFO] pushed message: %v\n", message)
	return nil
}

type Creator struct {
	Username string
	Profile  *Profile
}

type Post struct {
	reqID     int64
	postID    int64
	text      string
	mentions  []string
	timestamp int64
	creator   *Creator
}

type Notification struct {
	reqID  int64
	postID int64
}

type Profile struct {
	Bio     string
	Website string
	Age     int
}

func StorePost(ctx context.Context, reqID int64, text string, mentions []string) (int64, error) {
	postID := rand.Int63()
	timestamp := rand.Int63()

	post := Post{
		reqID:     reqID,
		postID:    postID,
		text:      text,
		mentions:  mentions,
		timestamp: timestamp,
		creator: &Creator{
			Username: "some username",
		},
	}
	db := &MongoDB{}
	db.Insert(ctx, post)

	notification := Notification{
		reqID:  reqID,
		postID: postID,
	}
	queue := &RabbitMQ{}
	queue.Push(ctx, notification)

	return postID, nil
}

func main() {
	ctx := context.Background()
	reqID := int64(0)
	text := "mytext"
	mentions := []string{"mention1", "mention2"}
	StorePost(ctx, reqID, text, mentions)
}
