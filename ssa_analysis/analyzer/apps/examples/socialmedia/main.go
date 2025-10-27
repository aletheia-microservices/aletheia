package main

import (
	"context"
)

type User struct {
	ID       int
	Username string
	Friends  []*User
}

type Message struct {
	Sender    *User
	Recipient *User
	Content   string
}

type DB interface {
	Save(ctx context.Context, data interface{}) error
}

type InMemoryDB struct{}

func (db *InMemoryDB) Save(ctx context.Context, data interface{}) error {
	//EVAL - fmt.Printf("[DB] Saved: %+v\n", data)
	return nil
}

func createUser(id int, username string) *User {
	return &User{
		ID:       id,
		Username: username,
	}
}

func createUserCopy(id int, username string) User {
	return User{
		ID:       id,
		Username: username,
	}
}

func addFriend(u *User, friend *User) {
	u.Friends = append(u.Friends, friend)
}

func sendMessage(ctx context.Context, db DB, sender *User, recipient *User, content string) {
	msg := &Message{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
	}
	db.Save(ctx, msg)
}

func main() {
	ctx := context.Background()
	db := &InMemoryDB{}

	alice := createUser(1, "alice")
	bob := createUser(2, "bob")
	charlie := createUser(3, "charlie")

	addFriend(alice, bob)
	addFriend(alice, charlie)

	// testing pointer relationships
	for _, friend := range alice.Friends {
		sendMessage(ctx, db, alice, friend, "Hi "+friend.Username+"!")
	}

	tera := createUserCopy(4, "terra")
	db.Save(ctx, tera)
}
