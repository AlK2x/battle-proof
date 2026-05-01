package graph

import "user/internal/app"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	UserService app.UserService
}
