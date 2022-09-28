package metadata

import "context"

type Metadata struct {
	ProtoRevision string
	ProtoPath     string
}

const ContextKey = "berrypost-metadata-key"

func FromContext(ctx context.Context) (Metadata, bool) {
	meta, ok := ctx.Value(ContextKey).(Metadata)
	return meta, ok
}

const ProtoRevisionGRPCMetadataKey = "x-proto-revision"
const ProtoPathGRPCMetadataKey = "x-proto-path"
