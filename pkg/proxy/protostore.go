package proxy

import (
	"context"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

type RuntimeProtoStore interface {
	GetMethodMessage(context.Context, string, string) (proto.Message, proto.Message, error)
	GetMessage(context.Context, string) (proto.Message, error)
}

type Metadata struct {
	ProtoRevision string
}

func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return Metadata{}, false
	}

	meta := Metadata{}
	if vs := md.Get("x-proto-revision"); len(vs) > 0 {
		meta.ProtoRevision = vs[0]
	}
	return meta, true
}

type defaultRuntimeProtoStore struct{}

func (defaultRuntimeProtoStore) GetMethodMessage(context.Context, string, string) (proto.Message, proto.Message, error) {
	return nil, nil, errors.New("unimpl")
}

func (defaultRuntimeProtoStore) GetMessage(context.Context, string) (proto.Message, error) {
	return nil, errors.New("unimpl")
}

type ctxedAnyResolver struct {
	RuntimeProtoStore
	ctx context.Context
}

func trimAnyTypePrefix(typeURL string) string {
	pos := strings.Index(typeURL, "/")
	if pos == -1 {
		return typeURL
	}
	return typeURL[pos+1:]
}

func (a ctxedAnyResolver) resolveTo(typeURL string, dstMessage *proto.Message, dstError *error) <-chan struct{} {
	doneSig := make(chan struct{})
	go func() {
		defer func() {
			doneSig <- struct{}{}
		}()
		m, err := a.GetMessage(a.ctx, trimAnyTypePrefix(typeURL))
		if err != nil {
			*dstError = err
			return
		}
		*dstMessage = m
	}()
	return doneSig
}

func (a ctxedAnyResolver) Resolve(typeURL string) (proto.Message, error) {
	var (
		dstMessage   proto.Message
		resolveError error
	)
	select {
	case <-a.resolveTo(typeURL, &dstMessage, &resolveError):
		if resolveError != nil {
			return nil, resolveError
		}
		return dstMessage, nil
	case <-a.ctx.Done():
		return nil, errors.WithStack(a.ctx.Err())
	}
}

func AsContextedAnyResolver(ctx context.Context, in RuntimeProtoStore) jsonpb.AnyResolver {
	return ctxedAnyResolver{
		RuntimeProtoStore: in,
		ctx:               ctx,
	}
}
