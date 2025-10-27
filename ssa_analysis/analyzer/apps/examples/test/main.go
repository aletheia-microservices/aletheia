package main

import (
	"context"
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
	//EVAL - fmt.Printf("[INFO] inserted document: %v\n", document)
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

func main() {
	ctx := context.Background()
	Test(ctx, 1)
}
