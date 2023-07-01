package graph

import service "github.com/LabbJoil/Chat/Services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ServiceChat service.ChatInteraction
}
