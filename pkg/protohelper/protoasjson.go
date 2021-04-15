package protohelper

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var (
	emptyMessageMarshaler = &jsonpb.Marshaler{
		OrigName:     true,
		EmitDefaults: true,
	}
)

func AsEmptyMessageJSON(in proto.Message) string {
	out, _ := emptyMessageMarshaler.MarshalToString(in)
	return out
}
