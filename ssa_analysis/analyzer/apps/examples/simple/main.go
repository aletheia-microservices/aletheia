package main

import (
	"context"
)

func main() {
	ctx := context.Background()
	myfunc_ptrs(ctx)
	myfunc_vals(ctx)
	myfunc_channels(ctx)
}

type DB interface {
	Save(ctx context.Context, data interface{}) error
}

type InMemoryDB struct{}

func (db *InMemoryDB) Save(ctx context.Context, data interface{}) error {
	//EVAL - fmt.Printf("[DB] Saved: %+v\n", data)
	return nil
}

type Product struct {
	name string
	desc string
}

func myfunc_ptrs(ctx context.Context) {
	db := InMemoryDB{}

	var prod1 *Product
	prod1 = &Product{
		name: "name",
	}
	var prod11 *Product
	prod11 = &Product{
		name: "name",
	}
	var prod12 *Product
	prod12 = &Product{
		name: "name",
	}
	var prod2 *Product
	prod2 = prod1
	prod2.desc = "desc"

	prod2 = prod11
	db.Save(ctx, prod2)

	prod2 = prod12
	db.Save(ctx, prod2)
}

func myfunc_vals(ctx context.Context) {
	var prod1 Product
	prod1 = Product{
		name: "name",
	}
	var prod2 Product
	prod2 = prod1
	prod2.desc = "desc"

	db := InMemoryDB{}
	db.Save(ctx, prod2)
}

func myfunc_channels(ctx context.Context) {
	var prod1 Product
	prod1 = Product{
		name: "chan",
	}
	prod1.desc = "sent over channel"
	ch := make(chan Product)

	go func(p Product) {
		ch <- p
	}(prod1)

	prod2 := <-ch

	db := InMemoryDB{}
	db.Save(ctx, prod2)
}
