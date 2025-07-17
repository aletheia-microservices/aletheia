package main

import (
	"context"
	"fmt"
)

type Product struct {
	ID   int64
	Name string
}

type MongoInserter interface {
	Insert(ctx context.Context, document interface{}) error
}

type MongoDB struct{}

func (m *MongoDB) Insert(ctx context.Context, document interface{}) error {
	fmt.Printf("[INFO] inserted document: %v\n", document)
	return nil
}

func StoreCart(ctx context.Context, i int) {
	prod1 := &Product{ID: 1}
	prod2 := &Product{ID: 2}
	cart := []*Product{prod1, prod2}

	prod1.ID = 99

	if i < 10 {
		prod := &Product{ID: 5}
		cart[0] = prod
	}

	for _, p := range cart {
		p.ID = 100
	}

	db := &MongoDB{}
	db.Insert(ctx, cart)
}

func main() {
	ctx := context.Background()
	StoreCart(ctx, 1)
}
