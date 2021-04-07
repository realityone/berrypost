package proxy

import (
	"context"

	"github.com/golang/protobuf/proto"
)

type RuntimeProtoStore interface {
	GetMethodMessage(context.Context, string, string) (proto.Message, proto.Message, error)
}
