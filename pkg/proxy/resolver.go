package proxy

import "context"

type RuntimeResolver interface {
	Resolve(context.Context, string) (string, error)
}
