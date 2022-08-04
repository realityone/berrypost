package management

import (
	"github.com/realityone/berrypost/pkg/metadata"
	"github.com/realityone/berrypost/pkg/server"
)

type Method struct {
	Name           string
	GRPCMethodName string
	InputSchema    string
	ServiceMethod  string
}

type Service struct {
	Name    string
	Methods []*Method
}

type MetadataItem struct {
	Key   string
	Value string
}

type InvokePage struct {
	Meta                 server.ServerMeta
	ServiceIdentifier    string
	PackageName          string
	PreferTarget         string
	Services             []*Service
	ProtoFiles           []*ProtoFileMeta
	InvokePageURLBuilder func(string, string) string
	DefaultGRPCMetadata  []*MetadataItem
	Metadata             metadata.Metadata
	KnownReferences      []*ReferenceItem
}
