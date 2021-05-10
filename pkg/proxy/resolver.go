package proxy

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ToNextResolver = errors.New("to next resolver")
)

type RuntimeServiceResolver interface {
	ResolveOnce(context.Context, *ResolveOnceRequest) (string, error)
	Name() string
}

type ResolveOnceRequest struct {
	ServiceFullyQualifiedName string
	UserDefinedTarget         string
}

type defaultRuntimeServiceResolver struct{}

func (defaultRuntimeServiceResolver) ResolveOnce(ctx context.Context, req *ResolveOnceRequest) (string, error) {
	if req.UserDefinedTarget == "" {
		return "", ToNextResolver
	}
	parsed, err := url.Parse(req.UserDefinedTarget)
	if err != nil {
		return "", ToNextResolver
	}
	switch parsed.Scheme {
	case "tcp", "udp":
		return req.UserDefinedTarget, nil
	default:
		return "", ToNextResolver
	}
}

func (defaultRuntimeServiceResolver) Name() string {
	return "default-resolver"
}

type chainedRuntimeResolver struct {
	all []RuntimeServiceResolver
}

func (crr chainedRuntimeResolver) ResolveOnce(ctx context.Context, req *ResolveOnceRequest) (string, error) {
	for _, r := range crr.all {
		addr, err := r.ResolveOnce(ctx, req)
		if err != nil {
			logrus.Warn("Failed to resolve %+v with resolver %q: %+v", req, r.Name(), err)
			continue
		}
		return addr, nil
	}
	return "", errors.Errorf("Could not resolve service: %+v", req)
}

func (crr chainedRuntimeResolver) Name() string {
	names := make([]string, 0, len(crr.all))
	for _, r := range crr.all {
		names = append(names, r.Name())
	}
	return fmt.Sprintf("chained-resolver:%q", strings.Join(names, ">"))
}

func ChainDefaultResolver(in ...RuntimeServiceResolver) RuntimeServiceResolver {
	resolvers := []RuntimeServiceResolver{
		defaultRuntimeServiceResolver{},
	}
	resolvers = append(resolvers, in...)
	return chainedRuntimeResolver{all: resolvers}
}
