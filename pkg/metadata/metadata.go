package metadata

import "context"

type Metadata struct {
	ProtoRevision string
}

const ContextKey = "berrypost-metadata-key"

func FromContext(ctx context.Context) (Metadata, bool) {
	meta, ok := ctx.Value(ContextKey).(Metadata)
	return meta, ok
}
