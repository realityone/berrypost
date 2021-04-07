package proxy

import (
	"context"

	"google.golang.org/grpc"
)

// Context is
type Context struct {
	context.Context

	srv           interface{}
	serverStream  grpc.ServerStream
	serviceMethod string
}

// Srv is
func (ctx *Context) Srv() interface{} {
	return ctx.srv
}

// ServerStream is
func (ctx *Context) ServerStream() grpc.ServerStream {
	return ctx.serverStream
}

// ServiceMethod is
func (ctx *Context) ServiceMethod() string {
	return ctx.serviceMethod
}
