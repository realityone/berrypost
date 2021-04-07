package proxy

import "context"

type RuntimeServiceResolver interface {
	ResolveOnce(context.Context, string) (string, error)
}
