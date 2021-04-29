package management

import (
	"github.com/realityone/berrypost/pkg/server"
)

type Method struct {
	Name               string
	FullyQualifiedName string
	InputSchema        string
}

type Service struct {
	Name               string
	FullyQualifiedName string
	Methods            []*Method
}

type InvokePage struct {
	Meta         server.ServerMeta
	PackageName  string
	PreferTarget string

	Services []*Service
}
