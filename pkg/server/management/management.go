package management

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/pkg/errors"
	"github.com/realityone/berrypost/pkg/etcd"
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

func (m Management) addressUpdate(ctx *gin.Context) {
	type KDRespBody struct {
		TargetAddr string `json:"targetAddrInput"`
		ProtoName  string `json:"serviceInput"`
	}
	var reqInfo KDRespBody
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.Error(err)
		return
	}
	reqInfo.ProtoName = strings.Replace(reqInfo.ProtoName, " ", "", -1)
	reqInfo.TargetAddr = strings.Replace(reqInfo.TargetAddr, " ", "", -1)
	err := etcd.Dao.Put(reqInfo.ProtoName, reqInfo.TargetAddr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) getMethods(ctx *gin.Context) {
	type RespBody struct {
		FileName string `json:"fileName"`
	}
	var reqInfo RespBody
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.Error(err)
		return
	}
	methods, err := m.getMethodsByService(ctx, reqInfo.FileName)
	if err != nil {
		ctx.Error(err)
	}
	ctx.JSON(http.StatusOK, methods)
}

// find proto file by service identifier in this order:
// - proto file identifier
// - proto file name
// - proto package name
// - service alias
func (m Management) findProtoFileByServiceIdentifier(ctx context.Context, serviceIdentifier string) (*ProtoFileProfile, bool) {
	files, err := m.protoManager.ListProtoFiles(ctx)
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
		packages, err := m.protoManager.ListPackages(ctx)
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
		alias, err := m.protoManager.ListServiceAlias(ctx)
		if err != nil {
			return "", false, nil
		}
		packages, err := m.protoManager.ListPackages(ctx)
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

	stayStill := func() (string, bool, error) {
		return serviceIdentifier, true, nil
	}

	for _, fn := range []func() (string, bool, error){
		fromProtoImportPath,
		fromProtoFilename,
		fromProtoPackageName,
		fromServiceAlias,
		stayStill,
	} {
		importPath, ok, err := fn()
		if err != nil {
			logrus.Error("Failed to find proto import: %+v", err)
			continue
		}
		if !ok {
			logrus.Error("Failed to find proto file desctrption with given service identifier: %q", serviceIdentifier)
			continue
		}
		profile, err := m.protoManager.GetProtoFile(ctx, &GetProtoFileRequest{
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

func (m Management) makeInvokePage(ctx context.Context, serviceIdentifier string) (*InvokePage, error) {
	// todo: userid
	userid := "test_user"
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
		Blueprints:        m.allUserBlueprints(ctx, userid),
	}
	preferTarget, ok := fileProfile.Common.Annotation[AppBerrypostManagementInvokePreferTarget]
	if ok {
		page.PreferTarget = preferTarget
	}
	defaultTarget, ok := fileProfile.Common.Annotation[AppBerrypostManagementInvokeDefaultTarget]
	if ok {
		page.DefaultTarget = defaultTarget
	}

	page.Services = make([]*Service, 0, len(fileProfile.ProtoPackage.FileDescriptor.GetServices()))
	for _, s := range fileProfile.ProtoPackage.FileDescriptor.GetServices() {
		ps := &Service{
			Name:     s.GetName(),
			FileName: serviceIdentifier,
		}
		ps.Methods = make([]*Method, 0, len(s.GetMethods()))
		for _, m := range s.GetMethods() {
			pm := &Method{
				Name:           m.GetName(),
				GRPCMethodName: fmt.Sprintf("/%s/%s", s.GetFullyQualifiedName(), m.GetName()),
				ServiceMethod:  fmt.Sprintf("%s.%s", s.GetName(), m.GetName()),
				PreferTarget:   preferTarget,
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

func (m Management) makeBlueprintPage(ctx context.Context, blueprintIdentifier string) (*BlueprintPage, error) {
	//todo: userid
	userid := "test_user"
	info, err := m.blueprintMethods(ctx, userid, blueprintIdentifier)
	if err != nil {
		return nil, err
	}
	meta := &BlueprintMeta{
		blueprintIdentifier: blueprintIdentifier,
		Methods:             info,
	}
	page := &BlueprintPage{
		Meta:                m.server.Meta(),
		BlueprintIdentifier: blueprintIdentifier,
		//PackageName:       fileProfile.ProtoPackage.FileDescriptor.GetFullyQualifiedName(),
		PreferTarget: blueprintIdentifier,
		ProtoFiles:   m.allProtoFiles(ctx),
		Blueprints:   m.allUserBlueprints(ctx, userid),
	}
	page.Services = make([]*Service, 0, len(meta.Methods))
	for _, info := range meta.Methods {
		serviceIdentifier := info.Filename
		fileProfile, ok := m.findProtoFileByServiceIdentifier(ctx, serviceIdentifier)

		if !ok {
			return nil, errors.Errorf("Failed to find package profile from service identifier: %q", serviceIdentifier)
		}
		preferTarget, ok := fileProfile.Common.Annotation[AppBerrypostManagementInvokePreferTarget]
		if ok {
			page.PreferTarget = preferTarget
		}
		defaultTarget, ok := fileProfile.Common.Annotation[AppBerrypostManagementInvokeDefaultTarget]
		if ok {
			page.DefaultTarget = defaultTarget
		}
		for _, s := range fileProfile.ProtoPackage.FileDescriptor.GetServices() {
			ps := &Service{
				Name:     s.GetName(),
				FileName: serviceIdentifier,
			}
			ps.Methods = make([]*Method, 0, len(s.GetMethods()))
			for _, m := range s.GetMethods() {
				if m.GetName() != info.MethodName {
					continue
				}
				pm := &Method{
					Name:           m.GetName(),
					GRPCMethodName: fmt.Sprintf("/%s/%s", s.GetFullyQualifiedName(), m.GetName()),
					ServiceMethod:  fmt.Sprintf("%s.%s", s.GetName(), m.GetName()),
					PreferTarget:   preferTarget,
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
				break
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

func (m Management) allProtoFiles(ctx context.Context) []*ProtoFileMeta {
	files, err := m.protoManager.ListProtoFiles(ctx)
	if err != nil {
		logrus.Error("Failed to list proto files: %+v", err)
		return nil
	}
	return files
}

func (m Management) login(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", &LoginPage{
		Meta: m.server.Meta(),
	})
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

func (m Management) emptyBlueprint(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "blueprint.html", &BlueprintPage{
		Meta:       m.server.Meta(),
		Blueprints: m.allUserBlueprintsMeta(ctx, "test_user"),
	})
}

func (m Management) blueprint(ctx *gin.Context) {
	blueprintIdentifier := ctx.Param("blueprint-identifier")
	blueprintIdentifier = strings.TrimPrefix(blueprintIdentifier, "/")
	page, err := m.makeBlueprintPage(ctx, blueprintIdentifier)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.HTML(http.StatusOK, "blueprint.html", page)
}

func (m Management) emptyDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "dashboard.html", &InvokePage{
		Meta:       m.server.Meta(),
		ProtoFiles: m.allProtoFiles(ctx),
	})
}

func (m Management) dashboard(ctx *gin.Context) {
	serviceIdentifier := ctx.Param("service-identifier")
	serviceIdentifier = strings.TrimPrefix(serviceIdentifier, "/")
	page, err := m.makeInvokePage(ctx, serviceIdentifier)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.HTML(http.StatusOK, "dashboard.html", page)
}

func (m Management) setting(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "setting.html", nil)
}

func (m Management) Setup(s *server.Server) error {
	m.server = s

	r := s.Group("/management")
	r.GET("/rediect-to-example", m.redirectToFirstService)
	r.GET("/invoke", m.emptyInvoke)
	r.GET("/invoke/*service-identifier", m.invoke)
	r.GET("/blueprint", m.emptyBlueprint)
	r.GET("/blueprint/*blueprint-identifier", m.blueprint)
	r.GET("/login", m.login)

	rAPI := s.Group("/management/api")
	rAPI.GET("/_intro", m.intro)
	rAPI.GET("/packages", m.listPackages)
	rAPI.GET("/packages/:package_name", m.getPackage)
	rAPI.GET("/service-alias", m.listServiceAlias)
	rAPI.POST("/update", m.addressUpdate)
	rAPI.POST("/getMethods", m.getMethods)
	b := rAPI.Group("/blueprint")
	b.POST("/new", m.newBlueprint)
	b.POST("/delete", m.delBlueprint)
	b.POST("/copy", m.copyBlueprint)
	b.POST("/append", m.savetoBlueprint)
	b.POST("/appendList", m.appendBlueprint)
	b.POST("/reduce", m.deleteBlueprintMethod)

	s.GET("/dashboard", m.emptyDashboard)
	s.GET("/dashboard/*service-identifier", m.dashboard)
	s.GET("/setting", m.setting)
	return nil
}

func (m Management) Name() string {
	return "management-server"
}

func (m Management) Meta() map[string]string {
	return nil
}
