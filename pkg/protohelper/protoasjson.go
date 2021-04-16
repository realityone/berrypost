package protohelper

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

var (
	emptyMessageMarshaler = &jsonpb.Marshaler{
		OrigName:     true,
		EmitDefaults: true,
		AnyResolver:  EmptyAnyResolver{},
	}
)

func AsEmptyMessageJSON(in proto.Message) string {
	out, _ := emptyMessageMarshaler.MarshalToString(in)
	return out
}

type WrappedAnyResolver struct {
	jsonpb.AnyResolver
}

func (w WrappedAnyResolver) Resolve(typeURL string) (proto.Message, error) {
	m, err := w.AnyResolver.Resolve(typeURL)
	if err != nil {
		logrus.Warnf("Failed to resolve type: %q, using dummy message directly: %+v.", typeURL, err)
		return &DummyMessage{}, nil
	}
	logrus.Infof("Succeeded to resolve type: %q: %+v", typeURL, m)
	return m, nil
}

type EmptyAnyResolver struct{}

func (EmptyAnyResolver) Resolve(typeURL string) (proto.Message, error) {
	return &DummyMessage{}, nil
}
