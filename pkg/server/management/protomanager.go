package management

import (
	"context"

	"github.com/jhump/protoreflect/desc"
)

type ProtoManager interface {
	ListPackages(context.Context) ([]*PackageMeta, error)
	GetPackage(context.Context, *GetPackageRequest) (*ProtoPackage, error)
	ListServiceAlias(context.Context) ([]*ServiceAlias, error)
}

type ServiceAlias struct {
	Package string   `json:"package"`
	Alias   []string `json:"alias"`
}

type GetPackageRequest struct {
	PackageName string
}

type PackageMeta struct {
	Meta    ProtoMeta `json:"meta"`
	Package string    `json:"package"`
}

type ProtoMeta struct {
	ProtoPath  string `json:"proto_path"`
	ImportPath string `json:"import_path"`
}

type ProtoPackage struct {
	Meta           ProtoMeta            `json:"meta"`
	Common         Common               `json:"common"`
	FileDescriptor *desc.FileDescriptor `json:"file_descriptor"`
}

type defaultProtoManager struct{}

func (dpm defaultProtoManager) ListPackages(context.Context) ([]*PackageMeta, error) {
	return []*PackageMeta{}, nil
}

func (dpm defaultProtoManager) GetPackage(context.Context, *GetPackageRequest) (*ProtoPackage, error) {
	return &ProtoPackage{}, nil
}

func (dpm defaultProtoManager) ListServiceAlias(context.Context) ([]*ServiceAlias, error) {
	return []*ServiceAlias{}, nil
}
