package proxy

import (
	"context"

	"github.com/pkg/errors"
)

type RuntimeServiceResolver interface {
	ResolveOnce(context.Context, string) (string, error)
}

type defaultRuntimeServiceResolver struct{}

func (defaultRuntimeServiceResolver) ResolveOnce(context.Context, string) (string, error) {
	return "", errors.New("unimpl")
}
