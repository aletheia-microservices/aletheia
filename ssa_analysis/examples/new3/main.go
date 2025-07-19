package main

import (
	"context"
	"fmt"
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
	Insert(ctx context.Context, db string, document interface{}) error
	Find(ctx context.Context, db string, id string) error
}

type MongoDB struct{}

func (m *MongoDB) Insert(ctx context.Context, db string,  document interface{}) error {
	fmt.Printf("[INFO] inserted document: %v\n", document)
	return nil
}

func (m *MongoDB) Find(ctx context.Context, db string, id string) Sku {
	fmt.Printf("[INFO] found document for id: %v\n", id)
	return Sku{}
}

type RabbitMQPusher interface {
	Push(ctx context.Context, document interface{}) error
}

type RabbitMQ struct{}

func (r *RabbitMQ) Push(ctx context.Context, message interface{}) error {
	fmt.Printf("[INFO] pushed message: %v\n", message)
	return nil
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
			// logic of sku.Get
			skuDB := &MongoDB{}
			sku := skuDB.Find(ctx, "sku", myitem2.Parent)
			myitem2.Amount = int64(sku.Price)
			myitem2.Currency = sku.Currency
			myitem2.Description = sku.Name
		}
	}

	order.Items = orderItems
	order.Amount = 100

	orderDB := &MongoDB{}
	orderDB.Insert(ctx, "order", order)

	return order, nil
}

func main() {
	ctx := context.Background()
	currency := int32(0)
	items := []*OrderItem{}
	metadata := make(map[string]string, 0)
	email := ""
	shipping := &Shipping{}
	New(ctx, currency, items, metadata, email, shipping)
}
