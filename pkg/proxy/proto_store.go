package proxy

import (
	"context"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RuntimeProtoStore interface {
	GetMethodMessage(context.Context, string, string) (proto.Message, proto.Message, error)
}

type defaultRuntimeProtoStore struct{}

func (defaultRuntimeProtoStore) GetMethodMessage(context.Context, string, string) (proto.Message, proto.Message, error) {
	return nil, nil, errors.New("unimpl")
}

type wrappedAnyResolver struct {
	jsonpb.AnyResolver
}

func (w wrappedAnyResolver) Resolve(typeURL string) (proto.Message, error) {
	m, err := w.AnyResolver.Resolve(typeURL)
	if err != nil {
		logrus.Warnf("Failed to resolve type: %q, using dynamic message directly: %+v.", typeURL, err)
		return &dynamic.Message{}, nil
	}
	return m, nil
}

type emptyAnyResolver struct{}

func (emptyAnyResolver) Resolve(typeURL string) (proto.Message, error) {
	return &dynamic.Message{}, nil
}
