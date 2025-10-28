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
}

type Post struct {
	reqID        int64
	postID       int64
	text         string
	usermentions UserMentions
	timestamp    int64
	creator      *Creator
}

type Notification struct {
	reqID  int64
	postID int64
}
type UserMentions struct {
	mentions map[string]string
}

func Swap(x *int, y *int) int {
	temp := *x
	*x = *y
	*y = temp
	return temp
}

func Test() {
	x := 1
	p := &x
	y := *p
	_ = y // to avoid "declared and not used" error
}

func StorePost(ctx context.Context, reqID int64, text string, someuser string) (int64, error) {
	postID := rand.Int63()
	timestamp := rand.Int63()

	usermentions := UserMentions{
		mentions: make(map[string]string),
	}

	usermentions.mentions["alice"] = "hello alice"
	usermentions.mentions["bob"] = "hello bob"
	usermentions.mentions[someuser] = "hello some user"

	post := Post{
		reqID:     reqID,
		postID:    postID,
		text:      text,
		timestamp: timestamp,
		creator: &Creator{
			Username: "some username",
		},
		usermentions: usermentions,
	}
	db := &MongoDB{}
	db.Insert(ctx, post)

	db.Insert(ctx, usermentions.mentions["alice"])
	db.Insert(ctx, usermentions.mentions[someuser])
	db.Insert(ctx, usermentions.mentions[someuser])

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
	someuser := "some other user"
	StorePost(ctx, reqID, text, someuser)
}
