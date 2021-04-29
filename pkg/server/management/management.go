package management

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/realityone/berrypost/pkg/server"

	"github.com/gin-gonic/gin"
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
		for _, a := range sa.Alias {
			if a == serviceIdentifier {
				return sa.Package, true
			}
		}
	}
	return "", false
}

func (m Management) invoke(ctx *gin.Context) {
	page := &InvokePage{
		Meta: m.server.Meta(),
	}
	serviceIdentifier := ctx.Param("service-identifier")
	packageName, ok := m.findPackageNameByServiceIdentifier(ctx, serviceIdentifier)
	if !ok {
		ctx.Error(errors.Errorf("Failed to find package from service identifier: %q", serviceIdentifier))
		return
	}
	page.PackageName = packageName

	proto, err := m.protoManager.GetPackage(ctx, &GetPackageRequest{
		PackageName: packageName,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	for _, pkg := range proto.ProtoPackages {
		page.Service = make([]*Service, 0, len(pkg.FileDescriptor.GetServices()))
		for _, s := range pkg.FileDescriptor.GetServices() {
			ps := &Service{
				Name:               s.GetName(),
				FullyQualifiedName: s.GetFullyQualifiedName(),
			}
			ps.Methods = make([]*Method, 0, len(s.GetMethods()))
			for _, m := range s.GetMethods() {
				ps.Methods = append(ps.Methods, &Method{
					Name:               m.GetName(),
					FullyQualifiedName: m.GetFullyQualifiedName(),
				})
			}
			page.Service = append(page.Service, ps)
		}
	}

	ctx.HTML(http.StatusOK, "invoke.html", page)
}

func (m Management) Setup(s *server.Server) error {
	m.server = s

	r := s.Group("/management")
	r.GET("/invoke", m.invoke)
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
