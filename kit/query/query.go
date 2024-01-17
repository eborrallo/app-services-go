package query

import "context"

// Bus defines the expected behaviour from a query  bus.
type Bus interface {
	// Ask is the method used to dispatch new query s.
	Ask(context.Context, Query) (interface{}, error)
	// Register is the method used to register a new query  controllers.
	Register(Type, Handler)
}

//go:generate mockery --case=snake --outpkg=querymocks --output=querymocks --name=Bus

// Type represents an application command type.
type Type string

// Command represents an application command.
type Query interface {
	Type() Type
}

// Handler defines the expected behaviour from a query controllers.
type Handler interface {
	Handle(context.Context, Query) (interface{}, error)
}
