package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/examples/sockshop3/workflow/sockshop3"
	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

func main() {
	ctx := context.Background()

	var userDB backend.NoSQLDatabase
	userService, _ := sockshop3.NewUserServiceImpl(ctx, userDB)

	var username, password string
	userService.Login(ctx, username, password)
}
