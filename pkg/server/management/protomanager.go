package management

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/sirupsen/logrus"
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

func (pp *ProtoPackage) MarshalJSON() ([]byte, error) {
	marshalStruct := struct {
		Meta           ProtoMeta       `json:"meta"`
		FileDescriptor json.RawMessage `json:"file_descriptor"`
	}{
		Meta: pp.Meta,
	}

	descMarshaler := jsonpb.Marshaler{}
	descString, err := descMarshaler.MarshalToString(pp.FileDescriptor.AsProto())
	if err != nil {
		logrus.Warnf("Failed to marshal %+v as json string: %+v", err)
	}
	marshalStruct.FileDescriptor = json.RawMessage(descString)

	return json.Marshal(marshalStruct)
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
