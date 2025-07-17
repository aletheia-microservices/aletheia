package main

import (
	"context"
	"fmt"
)

type Product struct {
	ID   int
	Name string
}

type Notification struct {
	ID int
	ProductVal Product
	ProductPtr *Product
}

type MongoInserter interface {
	Insert(ctx context.Context, document interface{}) error
	Insert2(ctx context.Context, product Product) error
	Insert3(ctx context.Context, product *Product) error
}

type MongoDB struct{}

func (m *MongoDB) Insert(ctx context.Context, document interface{}) error {
	fmt.Printf("[INFO] inserted document: %v\n", document)
	return nil
}

func (m *MongoDB) InsertProduct(ctx context.Context, product Product) error {
	fmt.Printf("[INFO] inserted product: %v\n", product)
	return nil
}

func (m *MongoDB) InsertProductPtr(ctx context.Context, product *Product) error {
	fmt.Printf("[INFO] inserted product: %v\n", product)
	return nil
}

func Test(ctx context.Context, i int) {
	prod1 := &Product{ID: 1}
	prod2 := &Product{ID: 2}
	cart := []*Product{prod1, prod2}

	db := &MongoDB{}
	db.Insert(ctx, cart)
}

func StoreCart(ctx context.Context, i int) {
	prod1 := Product{ID: 1}
	prod2 := Product{ID: 2}
	cart := []*Product{&prod1, &prod2}
	prod1.ID = 99

	if i < 10 {
		prod := Product{ID: 5}
		cart[1] = &prod
	}

	prod3 := prod2
	prod4ptr := &Product{ID: 4}
	prod3ptr := prod4ptr
	prod3ptr.ID = -4
	prod4_doubleptr := &prod4ptr

	notif := Notification{ProductVal: prod1, ProductPtr: &prod1}
	notif.ProductVal.Name = "hello_world_1"
	notif.ProductPtr.Name = "hello_world_2"

	for _, p := range cart {
		p.ID = 100
	}

	db := &MongoDB{}
	db.Insert(ctx, cart)
	db.Insert(ctx, cart[0])
	db.Insert(ctx, prod3)
	db.Insert(ctx, prod3ptr)
	db.Insert(ctx, prod4_doubleptr)
	db.Insert(ctx, notif)
}

func main() {
	ctx := context.Background()
	StoreCart(ctx, 1)
}
