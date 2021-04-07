package proxy

import "context"

type RuntimeProtoStore interface {
	GetProto(context.Context) ([]byte, error)
}
