package proxy

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

type RuntimeProtoStore interface {
	GetMethodMessage(context.Context, string, string) (proto.Message, proto.Message, error)
}

type defaultRuntimeProtoStore struct{}

func (defaultRuntimeProtoStore) GetMethodMessage(context.Context, string, string) (proto.Message, proto.Message, error) {
	return nil, nil, errors.New("unimpl")
}
