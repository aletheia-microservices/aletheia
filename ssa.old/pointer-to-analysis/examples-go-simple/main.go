package main

import (
	"context"
	"fmt"
	"math/rand"
)

type Inventory struct {
	amount int
}

type Product struct {
	ID           int
	Name         string
	InventoryPtr *Inventory
	InventoryVal Inventory
}

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

func Test(ctx context.Context, i int) {
	prod1 := &Product{ID: 1}
	prod2 := &Product{ID: 2}
	cart := []*Product{prod1, prod2}

	db := &MongoDB{}
	db.Insert(ctx, cart)

	/* cart2
	prod2 = cart.prod[0]
	svc.Call(cart2) */

	if i < 10 {
		prod1.ID = 50
	} else {
		prod2.ID = 100
	}

	prod2.InventoryPtr = &Inventory{amount: 10}
	prod2.InventoryVal = Inventory{amount: 20}

	db.Insert(ctx, prod2)

	/* prod2.InventoryPtr = &Inventory{amount: 30}
	prod2.ID = prod1.ID
	db.Insert(ctx, prod2) */
}

type Address struct {
	City string
}

type Shipping struct {
	Address *Address
}

type OrderItem struct {
	Type        int32
	Quantity    int64
	Amount      int64
	Currency    int32
	Parent      string
	Description string
}

func (item *OrderItem) IsTypeTax() bool {
	return true
}

func (item *OrderItem) IsTypeSku() bool {
	return true
}

type Sku struct {
	Name     string
	Price    uint64
	Currency int32
}

type OrderItems struct {
	Items []*OrderItem
}

type Order struct {
	Currency int32
	Items    []*OrderItem
	Metadata map[string]string
	Email    string
	Shipping *Shipping
	Amount   int
}

func Get(ctx context.Context, id string) (*Sku, error) {
	return &Sku{}, nil
}

func calculateTotal(currency int32, items []*OrderItem) (int, error) {
	return 0, nil
}

func New(ctx context.Context, currency int32, items []*OrderItem, metadata map[string]string, email string, shipping *Shipping) (*Order, error) {
	order := &Order{
		Currency: currency,
		Items:    items,
		Metadata: metadata,
		Shipping: shipping,
	}

	var orderItems []*OrderItem
	for _, myitem1 := range items {
		if myitem1.IsTypeTax() {
			if myitem1.Quantity <= 0 {
				myitem1.Quantity = 1
			}
			orderItems = append(orderItems, myitem1)
		}
	}
	for _, myitem2 := range orderItems {
		if myitem2.IsTypeSku() {
			item, _ := Get(ctx, myitem2.Parent)
			myitem2.Amount = int64(item.Price)
			myitem2.Currency = item.Currency
			myitem2.Description = item.Name
		}
	}

	order.Shipping = &Shipping{}
	order.Shipping.Address = &Address{}
	order.Currency = 50

	order.Items = orderItems

	db := &MongoDB{}
	db.Insert(ctx, order.Shipping)
	db.Insert(ctx, order.Items)
	db.Insert(ctx, order.Currency)
	db.Insert(ctx, order.Currency)

	/* order.Items[0].Amount = 100
	order.Shipping.Address.City = "myaddress" */
	return order, nil
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

/* func StorePost(ctx context.Context, reqID int64, text string, mentions []string) (int64, error) {
	postID := rand.Int63()
	timestamp := rand.Int63()

	post := Post{
		reqID:     reqID,
		postID:    postID,
		text:      text,
		mentions:  mentions,
		timestamp: timestamp,
		creator: Creator{
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
} */

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

	notification := Notification{
		reqID:  reqID,
		postID: postID,
		user:   creator.Username,
	}

	queue := &RabbitMQ{}
	queue.Push(ctx, notification)

	processCreator(ctx, creator)

	return postID, nil
}

func main() {
	ctx := context.Background()
	//Test(ctx, 1)
	/* currency := int32(0)
	items := []*OrderItem{}
	metadata := make(map[string]string, 0)
	email := ""
	shipping := &Shipping{}
	New(ctx, currency, items, metadata, email, shipping) */
	reqID := int64(0)
	text := "mytext"
	mentions := []string{"mention1", "mention2"}
	StorePost(ctx, reqID, text, mentions)
}
