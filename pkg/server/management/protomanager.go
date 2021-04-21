package management

import (
	"context"

	"github.com/jhump/protoreflect/desc"
)

type ProtoManager interface {
	ListPackages(context.Context) ([]*PackageMeta, error)
	GetPackage(context.Context, *GetPackageRequest) (*ProtoPackageProfile, error)
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

type ProtoPackageProfile struct {
	Common        Common          `json:"common"`
	ProtoPackages []*ProtoPackage `json:"proto_packages"`
}

type ProtoPackage struct {
	Meta           ProtoMeta            `json:"meta"`
	FileDescriptor *desc.FileDescriptor `json:"file_descriptor"`
}

type defaultProtoManager struct{}

func (dpm defaultProtoManager) ListPackages(context.Context) ([]*PackageMeta, error) {
	return []*PackageMeta{}, nil
}

func (dpm defaultProtoManager) GetPackage(context.Context, *GetPackageRequest) (*ProtoPackageProfile, error) {
	return &ProtoPackageProfile{}, nil
}

func (dpm defaultProtoManager) ListServiceAlias(context.Context) ([]*ServiceAlias, error) {
	return []*ServiceAlias{}, nil
}
