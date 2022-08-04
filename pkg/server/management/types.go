package management

import (
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

type InvokePage struct {
	Meta                 server.ServerMeta
	ServiceIdentifier    string
	PackageName          string
	PreferTarget         string
	Services             []*Service
	ProtoFiles           []*ProtoFileMeta
	InvokePageURLBuilder func(*ProtoFileMeta) string
}
