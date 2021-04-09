package proxy

import (
	"context"
	"net/http"
)

// Context is
type Context struct {
	context.Context

	req    *http.Request
	writer http.ResponseWriter

	serviceMethod string
}
