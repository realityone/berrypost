package management

import (
	"github.com/realityone/berrypost/pkg/server"
)

type Method struct {
	Name               string
	FullyQualifiedName string
}

type Service struct {
	Name               string
	FullyQualifiedName string
	Methods            []*Method
}

type InvokePage struct {
	Meta        server.ServerMeta
	PackageName string

	Service []*Service
}
