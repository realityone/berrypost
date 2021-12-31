package management

import (
	"github.com/realityone/berrypost/pkg/server"
)

type Method struct {
	Name           string
	GRPCMethodName string
	InputSchema    string
	ServiceMethod  string
	PreferTarget   string
}

type Service struct {
	Name     string
	FileName string
	Methods  []*Method
}

type InvokePage struct {
	Meta              server.ServerMeta
	ServiceIdentifier string
	PackageName       string
	PreferTarget      string
	DefaultTarget     string
	Services          []*Service
	ProtoFiles        []*ProtoFileMeta
	Link              string
	Blueprints        []string
}

type BlueprintMeta struct {
	blueprintIdentifier string
	Methods             []*BlueprintMethodInfo
}

type BlueprintMethodInfo struct {
	Filename   string `json:"filename"`
	MethodName string `json:"methodName"`
}
