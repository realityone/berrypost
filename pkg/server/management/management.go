package management

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/pkg/errors"
	"github.com/realityone/berrypost/pkg/server"
	"github.com/sirupsen/logrus"
	"k8s.io/kube-openapi/pkg/util/sets"
)

type Option func(*Management)

func SetProtoManager(in ProtoManager) Option {
	return func(m *Management) {
		m.protoManager = in
	}
}

type Management struct {
	server       *server.Server
	protoManager ProtoManager
}

func New(opts ...Option) *Management {
	m := &Management{
		protoManager: defaultProtoManager{},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m Management) intro(ctx *gin.Context) {
	introSchema := struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}{
		Name:    "berrypost-management-server",
		Version: "0.0.1",
	}
	ctx.JSON(http.StatusOK, introSchema)
}

func (m Management) listPackages(ctx *gin.Context) {
	packages, err := m.protoManager.ListPackages(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, packages)
}

func (m Management) getPackage(ctx *gin.Context) {
	packageProfile, err := m.protoManager.GetPackage(ctx, &GetPackageRequest{
		PackageName: ctx.Param("package_name"),
	})
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, packageProfile)
}

func (m Management) listServiceAlias(ctx *gin.Context) {
	alias, err := m.protoManager.ListServiceAlias(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, alias)
}

func (m Management) findPackageNameByServiceIdentifier(ctx context.Context, serviceIdentifier string) (string, bool) {
	alias, err := m.protoManager.ListServiceAlias(ctx)
	if err != nil {
		return "", false
	}
	for _, sa := range alias {
		names := sets.NewString(sa.Package)
		names.Insert(sa.Alias...)
		if names.Has(serviceIdentifier) {
			return sa.Package, true
		}
	}
	return "", false
}

func (m Management) makeInvokePage(ctx context.Context, serviceIdentifier string) (*InvokePage, error) {
	packageName, ok := m.findPackageNameByServiceIdentifier(ctx, serviceIdentifier)
	if !ok {
		return nil, errors.Errorf("Failed to find package from service identifier: %q", serviceIdentifier)
	}
	page := &InvokePage{
		Meta:              m.server.Meta(),
		ServiceIdentifier: serviceIdentifier,
		PackageName:       packageName,
		PreferTarget:      serviceIdentifier,
		ServiceDropdown:   m.allServiceAlias(ctx),
	}

	proto, err := m.protoManager.GetPackage(ctx, &GetPackageRequest{
		PackageName: packageName,
	})
	if err != nil {
		return nil, err
	}

	for _, pkg := range proto.ProtoPackages {
		page.Services = make([]*Service, 0, len(pkg.FileDescriptor.GetServices()))
		for _, s := range pkg.FileDescriptor.GetServices() {
			ps := &Service{
				Name:               s.GetName(),
				FullyQualifiedName: s.GetFullyQualifiedName(),
			}
			ps.Methods = make([]*Method, 0, len(s.GetMethods()))
			for _, m := range s.GetMethods() {
				pm := &Method{
					Name:               m.GetName(),
					GRPCMethodName:     fmt.Sprintf("/%s/%s", s.GetFullyQualifiedName(), m.GetName()),
					FullyQualifiedName: m.GetFullyQualifiedName(),
				}
				descMarshaler := jsonpb.Marshaler{
					EmitDefaults: true,
					Indent:       "    ",
				}
				inputSchema, err := descMarshaler.MarshalToString(dynamic.NewMessage(m.GetInputType()))
				if err != nil {
					logrus.Warn("Failed to marshal method: %q input type as string: %+v", m.GetFullyQualifiedName(), err)
				}
				pm.InputSchema = inputSchema
				ps.Methods = append(ps.Methods, pm)
			}
			page.Services = append(page.Services, ps)
		}
	}

	return page, nil
}

func (m Management) firstServiceAlias(ctx context.Context) string {
	serviceAlias, err := m.protoManager.ListServiceAlias(ctx)
	if err != nil {
		return ""
	}
	for _, sa := range serviceAlias {
		for _, a := range sa.Alias {
			return a
		}
	}
	return ""
}

func (m Management) redirectToFirstService(ctx *gin.Context) {
	serviceIdentifier := m.firstServiceAlias(ctx)
	ctx.Redirect(http.StatusTemporaryRedirect, path.Join("/management/invoke", serviceIdentifier))
}

func (m Management) allServiceAlias(ctx context.Context) []string {
	alias, err := m.protoManager.ListServiceAlias(ctx)
	if err != nil {
		logrus.Error("Failed to list service alias: %+v", err)
		return nil
	}
	out := make([]string, 0, len(alias))
	for _, a := range alias {
		if len(a.Alias) > 0 {
			out = append(out, a.Alias[0])
		}
	}
	return out
}

func (m Management) invoke(ctx *gin.Context) {
	serviceIdentifier := ctx.Param("service-identifier")
	page, err := m.makeInvokePage(ctx, serviceIdentifier)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.HTML(http.StatusOK, "invoke.html", page)
}

func (m Management) emptyInvoke(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "invoke.html", &InvokePage{
		Meta:            m.server.Meta(),
		ServiceDropdown: m.allServiceAlias(ctx),
	})
}

func (m Management) Setup(s *server.Server) error {
	m.server = s

	r := s.Group("/management")
	r.GET("/rediect-to-example", m.redirectToFirstService)
	r.GET("/invoke", m.emptyInvoke)
	r.GET("/invoke/:service-identifier", m.invoke)

	rAPI := s.Group("/management/api")
	rAPI.GET("/_intro", m.intro)
	rAPI.GET("/packages", m.listPackages)
	rAPI.GET("/packages/:package_name", m.getPackage)
	rAPI.GET("/service-alias", m.listServiceAlias)
	return nil
}

func (m Management) Name() string {
	return "management-server"
}

func (m Management) Meta() map[string]string {
	return nil
}
