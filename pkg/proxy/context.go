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

func GetUserDefinedTarget(ctx context.Context) (string, bool) {
	proxyCtx, ok := ctx.(*Context)
	if !ok {
		return "", false
	}
	target := proxyCtx.req.Header.Get("X-Berrypost-Target")
	return target, target != ""
}
