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
	Blueprints        []string
	UserId            string
}

type BlueprintPage struct {
	Meta                server.ServerMeta
	BlueprintIdentifier string
	PreferTarget        string
	DefaultTarget       string
	Services            []*Service
	ProtoFiles          []*ProtoFileMeta
	Blueprints          []string
	UserId              string
}

type LoginPage struct {
	Meta server.ServerMeta
}

type BlueprintMeta struct {
	blueprintIdentifier string
	Methods             []*BlueprintMethodInfo
}

type BlueprintMethodInfo struct {
	Filename   string `json:"filename"`
	MethodName string `json:"methodName"`
}
