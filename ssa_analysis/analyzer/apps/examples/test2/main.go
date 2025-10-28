package main

import (
	"context"
	"fmt"
)

type Inventory struct {
	Amount int
}

type Product struct {
	ID           int
	Name         string
	InventoryPtr *Inventory
	InventoryVal Inventory
	Nested1      Nested1
	Nested1Ptr   *Nested1
	Metadata     map[string]string
}

type MongoInserter interface {
	Insert(ctx context.Context, document interface{}) error
}

type MongoDB struct{}

func (m *MongoDB) Insert(ctx context.Context, database string, collection string, document interface{}) error {
	fmt.Printf("[INFO] inserted document: %v\n", document)
	return nil
}

type Nested3 struct {
	Val3 int
}

type Nested2 struct {
	Val2 int
	Nested3
}

type Nested1 struct {
	Val1 int
	Nested2
}

func Test(ctx context.Context, postid int, metadata map[string]string) {
	db := &MongoDB{}

	inv1 := &Inventory{
		Amount: 50,
	}

	inv2 := Inventory{
		Amount: postid,
	}

	db.Insert(ctx, "db1", "ids", postid)

	prod1 := Product{
		ID:           1,
		Name:         "prod1",
		Metadata:     metadata,
		InventoryPtr: inv1,
		InventoryVal: inv2,
	}

	metadata["alice"] = "bob"

	inv2 = Inventory{}

	inv1.Amount = 10
	inv2.Amount = 20
	prod1.InventoryPtr.Amount = 30
	prod1.InventoryVal.Amount = 40

	db.Insert(ctx, "db2", "inventory", prod1.InventoryPtr)
	db.Insert(ctx, "db2", "inventory", prod1.InventoryVal)
	db.Insert(ctx, "db3", "product", prod1)

	prod2 := prod1
	prod2.Name = "MyName"
	db.Insert(ctx, "db3", "product", prod2)

	prod1.Metadata["alice2"] = "bob2"

	nested1 := Nested1{
		Val1: 1,
		Nested2: Nested2{
			Val2: 2,
			Nested3: Nested3{
				Val3: 3,
			},
		},
	}
	prod1.Nested1 = nested1

	nested1.Nested2.Nested3.Val3 = 4

	nested1ptr := &Nested1{
		Val1: 11,
		Nested2: Nested2{
			Val2: 22,
			Nested3: Nested3{
				Val3: 33,
			},
		},
	}
	prod1.Nested1Ptr = nested1ptr
	nested1ptr.Nested2.Nested3.Val3 = 44

	db.Insert(ctx, "db3", "product", prod1)
}

func main() {
	ctx := context.Background()
	Test(ctx, 1, make(map[string]string))
}
