package management

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/pkg/errors"
	"github.com/realityone/berrypost/pkg/metadata"
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
	packages, err := m.resolveProtoManager(ctx).ListPackages(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, packages)
}

func (m Management) getPackage(ctx *gin.Context) {
	packageProfile, err := m.resolveProtoManager(ctx).GetPackage(ctx, &GetPackageRequest{
		PackageName: ctx.Param("package_name"),
	})
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, packageProfile)
}

func (m Management) listServiceAlias(ctx *gin.Context) {
	alias, err := m.resolveProtoManager(ctx).ListServiceAlias(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, alias)
}

// find proto file by service identifier in this order:
// - proto file identifier
// - proto file name
// - proto package name
// - service alias
func (m Management) findProtoFileByServiceIdentifier(ctx context.Context, serviceIdentifier string) (*ProtoFileProfile, bool) {
	files, err := m.resolveProtoManager(ctx).ListProtoFiles(ctx)
	if err != nil {
		logrus.Error("Failed to list proto files: %+v", err)
	}

	fromProtoImportPath := func() (string, bool, error) {
		for _, f := range files {
			if f.Meta.ImportPath == serviceIdentifier {
				return f.Meta.ImportPath, true, nil
			}
		}
		return "", false, nil
	}

	fromProtoFilename := func() (string, bool, error) {
		for _, f := range files {
			if f.Filename == serviceIdentifier {
				return f.Meta.ImportPath, true, nil
			}
		}
		return "", false, nil
	}

	fromProtoPackageName := func() (string, bool, error) {
		packages, err := m.resolveProtoManager(ctx).ListPackages(ctx)
		if err != nil {
			return "", false, err
		}
		for _, p := range packages {
			if p.Package == serviceIdentifier {
				return p.Meta.ImportPath, true, nil
			}
		}
		return "", false, nil
	}

	fromServiceAlias := func() (string, bool, error) {
		alias, err := m.resolveProtoManager(ctx).ListServiceAlias(ctx)
		if err != nil {
			return "", false, nil
		}
		packages, err := m.resolveProtoManager(ctx).ListPackages(ctx)
		if err != nil {
			return "", false, err
		}
		packageGroupByName := map[string]*PackageMeta{}
		for _, p := range packages {
			packageGroupByName[p.Package] = p
		}

		for _, sa := range alias {
			names := sets.NewString(sa.Package)
			names.Insert(sa.Alias...)
			if names.Has(serviceIdentifier) {
				p, ok := packageGroupByName[sa.Package]
				if ok {
					return p.Meta.ImportPath, true, nil
				}
			}
		}
		return "", false, nil
	}

	for _, fn := range []func() (string, bool, error){
		fromProtoImportPath,
		fromProtoFilename,
		fromProtoPackageName,
		fromServiceAlias,
	} {
		importPath, ok, err := fn()
		if err != nil {
			logrus.Error("Failed to find proto import: %+v", err)
			continue
		}
		if !ok {
			logrus.Errorf("Failed to find proto file desctrption with given service identifier: %q", serviceIdentifier)
			continue
		}
		profile, err := m.resolveProtoManager(ctx).GetProtoFile(ctx, &GetProtoFileRequest{
			ImportPath: importPath,
		})
		if err != nil {
			logrus.Error("Failed to get proto file by import path: %q: %+v", importPath, err)
			continue
		}
		return profile, ok
	}

	return nil, false
}

func (m Management) resolveProtoManager(ctx context.Context) ProtoManager {
	meta, ok := metadata.FromContext(ctx)
	if !ok {
		return m.protoManager
	}
	if meta.ProtoRevision == "" {
		return m.protoManager
	}
	rm, ok := m.protoManager.(RevisionManager)
	if !ok {
		logrus.Warnf("Proto manager %T does not support revision management", m.protoManager)
		return m.protoManager
	}
	pm, err := rm.ResolveRevision(ctx, meta.ProtoRevision)
	if err != nil {
		logrus.Warnf("Failed to resolve proto manager on revision: %s: %+v", meta.ProtoRevision, err)
		return m.protoManager
	}
	return pm
}

func (m Management) makeInvokePage(ctx context.Context, serviceIdentifier string) (*InvokePage, error) {
	fileProfile, ok := m.findProtoFileByServiceIdentifier(ctx, serviceIdentifier)
	if !ok {
		return nil, errors.Errorf("Failed to find package profile from service identifier: %q", serviceIdentifier)
	}
	page := &InvokePage{
		Meta:              m.server.Meta(),
		ServiceIdentifier: serviceIdentifier,
		PackageName:       fileProfile.ProtoPackage.FileDescriptor.GetFullyQualifiedName(),
		PreferTarget:      serviceIdentifier,
		ProtoFiles:        m.allProtoFiles(ctx),
		InvokePageURLBuilder: func(pf *ProtoFileMeta) string {
			dst := fmt.Sprintf("/management/invoke/%s", pf.Meta.ImportPath)
			q := url.Values{}
			meta, _ := metadata.FromContext(ctx)
			if meta.ProtoRevision != "" {
				q.Set("protoRevison", meta.ProtoRevision)
			}
			if len(q) > 0 {
				dst = fmt.Sprintf("%s?%s", dst, q.Encode())
			}
			return dst
		},
	}
	preferTarget, ok := fileProfile.Common.Annotation[AppBerrypostManagementInvokePreferTarget]
	if ok {
		page.PreferTarget = preferTarget
	}

	page.Services = make([]*Service, 0, len(fileProfile.ProtoPackage.FileDescriptor.GetServices()))
	for _, s := range fileProfile.ProtoPackage.FileDescriptor.GetServices() {
		ps := &Service{
			Name: s.GetName(),
		}
		ps.Methods = make([]*Method, 0, len(s.GetMethods()))
		for _, m := range s.GetMethods() {
			pm := &Method{
				Name:           m.GetName(),
				GRPCMethodName: fmt.Sprintf("/%s/%s", s.GetFullyQualifiedName(), m.GetName()),
				ServiceMethod:  fmt.Sprintf("%s.%s", s.GetName(), m.GetName()),
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

	return page, nil
}

func (m Management) firstServiceAlias(ctx context.Context) string {
	serviceAlias, err := m.resolveProtoManager(ctx).ListServiceAlias(ctx)
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
	meta, _ := metadata.FromContext(ctx)

	dst := path.Join("/management/invoke", serviceIdentifier)
	p := url.Values{}
	p.Set("protoRevision", meta.ProtoRevision)
	dst = fmt.Sprintf("%s?%s", dst, p.Encode())
	ctx.Redirect(http.StatusTemporaryRedirect, dst)
}

func (m Management) allProtoFiles(ctx context.Context) []*ProtoFileMeta {
	files, err := m.resolveProtoManager(ctx).ListProtoFiles(ctx)
	if err != nil {
		logrus.Error("Failed to list proto files: %+v", err)
		return nil
	}
	return files
}

func (m Management) invoke(ctx *gin.Context) {
	serviceIdentifier := ctx.Param("service-identifier")
	serviceIdentifier = strings.TrimPrefix(serviceIdentifier, "/")
	page, err := m.makeInvokePage(ctx, serviceIdentifier)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.HTML(http.StatusOK, "invoke.html", page)
}

func (m Management) emptyInvoke(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "invoke.html", &InvokePage{
		Meta:       m.server.Meta(),
		ProtoFiles: m.allProtoFiles(ctx),
	})
}

func (m Management) prepareMetadata(ctx *gin.Context) {
	meta := metadata.Metadata{
		ProtoRevision: ctx.Query("protoRevision"),
	}
	ctx.Set(metadata.ContextKey, meta)
}

func (m Management) Setup(s *server.Server) error {
	m.server = s

	r := s.Group("/management", m.prepareMetadata)
	r.GET("/rediect-to-example", m.redirectToFirstService)
	r.GET("/invoke", m.emptyInvoke)
	r.GET("/invoke/*service-identifier", m.invoke)

	rAPI := s.Group("/management/api", m.prepareMetadata)
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
