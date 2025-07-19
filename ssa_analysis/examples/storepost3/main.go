package main

import (
	"context"
	"fmt"
	"math/rand"
)

type MongoInserter interface {
	Insert(ctx context.Context, document interface{}) error
}

type MongoDB struct{}

func (m *MongoDB) Insert(ctx context.Context, document interface{}) error {
	fmt.Printf("[INFO] inserted document: %v\n", document)
	return nil
}

type RabbitMQPusher interface {
	Push(ctx context.Context, document interface{}) error
}

type RabbitMQ struct{}

func (r *RabbitMQ) Push(ctx context.Context, message interface{}) error {
	fmt.Printf("[INFO] pushed message: %v\n", message)
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
	user   string
}

type Profile struct {
	Bio     string
	Website string
	Age     int
}

func getUserInput() string {
	return "myuser"
}

func getUserWebsite() string {
	return "mywebsite"
}

func sanitize(username string) string {
	return username
}

func processCreator(ctx context.Context, creator *Creator) {
	creator.Profile.Age = -1
	db := &MongoDB{}
	db.Insert(ctx, creator)
}

type Product struct {
	productID int64
}

func StorePost(ctx context.Context, reqID int64, text string, mentions []string) (int64, error) {
	postID := rand.Int63()
	timestamp := rand.Int63()

	username := getUserInput()
	safeUsername := sanitize(username)

	creator := &Creator{
		Username: safeUsername,
		Profile: &Profile{
			Bio:     "researcher",
			Website: getUserWebsite(),
		},
	}

	post := Post{
		reqID:     reqID,
		postID:    postID,
		text:      text,
		mentions:  mentions,
		timestamp: timestamp,
		creator:   creator,
	}
	db := &MongoDB{}
	db.Insert(ctx, post)

	post.reqID = -1
	db.Insert(ctx, post)
	
	post = Post{
		reqID:     reqID,
		postID:    postID,
	}
	db.Insert(ctx, post)

	product := Product{
		productID: 0,
	}

	notification := Notification{
		reqID:  reqID,
		postID: postID,
		user:   creator.Username,
	}

	queue := &RabbitMQ{}
	queue.Push(ctx, notification)
	
	notification.postID = product.productID
	queue.Push(ctx, notification.postID)

	processCreator(ctx, creator)

	return postID, nil
}

func main() {
	ctx := context.Background()
	reqID := int64(0)
	text := "mytext"
	mentions := []string{"mention1", "mention2"}
	StorePost(ctx, reqID, text, mentions)
}
