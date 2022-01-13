package management

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	reqInfo.ProtoName = m.trim(reqInfo.ProtoName)
	reqInfo.TargetAddr = m.trim(reqInfo.TargetAddr)
	err := etcd.Dao.Put(reqInfo.ProtoName, reqInfo.TargetAddr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
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

func (m Management) makeInvokePage(ctx context.Context, serviceIdentifier string, userid string) (*InvokePage, error) {
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
	if userid != "" {
		page.UserId = userid
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
		for _, mt := range s.GetMethods() {
			pm := &Method{
				Name:           mt.GetName(),
				GRPCMethodName: fmt.Sprintf("/%s/%s", s.GetFullyQualifiedName(), mt.GetName()),
				ServiceMethod:  fmt.Sprintf("%s.%s", s.GetName(), mt.GetName()),
				PreferTarget:   preferTarget,
			}
			descMarshaler := jsonpb.Marshaler{
				EmitDefaults: true,
				Indent:       "    ",
			}
			inputSchema, err := descMarshaler.MarshalToString(dynamic.NewMessage(mt.GetInputType()))
			if err != nil {
				logrus.Warn("Failed to marshal method: %q input type as string: %+v", mt.GetFullyQualifiedName(), err)
			}
			if userid != "" {
				hisKey := m.historyKey(userid, "", pm.GRPCMethodName)
				his, ok, err := etcd.Dao.Get(hisKey)
				if err != nil {
					logrus.Warn(err)
				}
				if ok {
					inputSchema = his
				}
			}
			pm.InputSchema = inputSchema
			ps.Methods = append(ps.Methods, pm)
		}
		page.Services = append(page.Services, ps)
	}
	return page, nil
}

func (m Management) makeBlueprintPage(ctx context.Context, blueprintIdentifier string, userid string) (*BlueprintPage, error) {
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
		PreferTarget:        blueprintIdentifier,
		ProtoFiles:          m.allProtoFiles(ctx),
		Blueprints:          m.allUserBlueprints(ctx, userid),
		UserId:              userid,
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
			for _, mt := range s.GetMethods() {
				if mt.GetName() != info.MethodName {
					continue
				}
				pm := &Method{
					Name:           mt.GetName(),
					GRPCMethodName: fmt.Sprintf("/%s/%s", s.GetFullyQualifiedName(), mt.GetName()),
					ServiceMethod:  fmt.Sprintf("%s.%s", s.GetName(), mt.GetName()),
					PreferTarget:   preferTarget,
				}
				descMarshaler := jsonpb.Marshaler{
					EmitDefaults: true,
					Indent:       "    ",
				}
				inputSchema, err := descMarshaler.MarshalToString(dynamic.NewMessage(mt.GetInputType()))
				if err != nil {
					logrus.Warn("Failed to marshal method: %q input type as string: %+v", mt.GetFullyQualifiedName(), err)
				}
				hisKey := m.historyKey(userid, blueprintIdentifier, pm.GRPCMethodName)
				his, ok, err := etcd.Dao.Get(hisKey)
				if err != nil {
					logrus.Warn(err)
				}
				if ok {
					inputSchema = his
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
func (m Management) register(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", &LoginPage{
		Meta: m.server.Meta(),
	})
}

func (m Management) signIn(ctx *gin.Context) {
	type RespBody struct {
		Userid   string `json:"userid"`
		Password string `json:"password"`
	}
	var reqInfo RespBody
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.Error(err)
		return
	}
	ok, err := m.userSignIn(ctx, reqInfo.Userid, reqInfo.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, ok)
}

func (m Management) userSignIn(ctx *gin.Context, userid string, password string) (bool, error) {
	ok, err := m.userVerify(ctx, userid, password)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	var hmacSampleSecret []byte
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userid,
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return false, err
	}
	ctx.SetCookie("session", tokenString, 3600, "", "", true, true)
	ctx.SetCookie("userid", userid, 3600, "", "", true, true)
	return true, nil
}

func (m Management) userSignUp(ctx *gin.Context, userid string, password string) (bool, error) {
	userKey := m.userKey(userid)
	_, ok, err := etcd.Dao.Get(userKey)
	if err != nil {
		return false, err
	}
	if ok {
		return false, nil
	}
	if err = etcd.Dao.Put(userKey, password); err != nil {
		return false, err
	}
	return true, nil
}

func (m Management) signUp(ctx *gin.Context) {
	type RespBody struct {
		Userid   string `json:"userid"`
		Password string `json:"password"`
	}
	var reqInfo RespBody
	if err := ctx.BindJSON(&reqInfo); err != nil {
		ctx.Error(err)
		return
	}
	ok, err := m.userSignUp(ctx, reqInfo.Userid, reqInfo.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, ok)
}

func (m Management) signOut(ctx *gin.Context) {
	ctx.SetCookie("session", "", -1, "", "", true, true)
	ctx.SetCookie("userid", "", -1, "", "", true, true)
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) saveHistory(ctx *gin.Context) {
	userid, err := ctx.Cookie("userid")
	if err != nil {
		ctx.JSON(http.StatusOK, nil)
		return
	}
	type RespBody struct {
		Blueprint string `json:"blueprint"`
		Service   string `json:"service"`
		Method    string `json:"method"`
		ReqBody   string `json:"reqBody"`
	}
	var reqInfo RespBody
	if err := ctx.BindJSON(&reqInfo); err != nil {
		logrus.Error(err)
		return
	}
	reqInfo.Blueprint = m.trim(reqInfo.Blueprint)
	reqInfo.Service = m.trim(reqInfo.Service)
	reqInfo.Method = m.trim(reqInfo.Method)
	if reqInfo.Blueprint == reqInfo.Service {
		reqInfo.Blueprint = ""
	}
	key := m.historyKey(userid, reqInfo.Blueprint, reqInfo.Method)
	if err := m.saveReq(key, reqInfo.ReqBody); err != nil {
		logrus.Error("Failed to save req history %+v", err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (m Management) historyKey(userid string, blueprint string, method string) string {
	return "/history/" + userid + "/" + blueprint + method
}

func (m Management) saveReq(key string, req string) error {
	if err := etcd.Dao.Put(key, req); err != nil {
		return err
	}
	return nil
}

func (m Management) userVerify(ctx context.Context, userid string, password string) (bool, error) {
	userKey := m.userKey(userid)
	pw, ok, err := etcd.Dao.Get(userKey)
	if err != nil {
		return false, err
	}
	if ok && pw != password {
		return false, nil
	}
	return true, nil
}

func (m Management) invoke(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	serviceIdentifier := ctx.Param("service-identifier")
	serviceIdentifier = strings.TrimPrefix(serviceIdentifier, "/")
	page, err := m.makeInvokePage(ctx, serviceIdentifier, userid)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.HTML(http.StatusOK, "invoke.html", page)
}

func (m Management) emptyInvoke(ctx *gin.Context) {
	page := &InvokePage{
		Meta:       m.server.Meta(),
		ProtoFiles: m.allProtoFiles(ctx),
	}
	userid, err := ctx.Cookie("userid")
	if err == nil {
		page.UserId = userid
	}
	ctx.HTML(http.StatusOK, "invoke.html", page)
}

func (m Management) emptyBlueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	page := &BlueprintPage{
		Meta:       m.server.Meta(),
		Blueprints: m.allUserBlueprintsMeta(ctx, userid),
	}
	page.UserId = userid
	ctx.HTML(http.StatusOK, "blueprint.html", page)
}

func (m Management) blueprint(ctx *gin.Context) {
	userid, _ := ctx.Cookie("userid")
	blueprintIdentifier := ctx.Param("blueprint-identifier")
	blueprintIdentifier = strings.TrimPrefix(blueprintIdentifier, "/")
	page, err := m.makeBlueprintPage(ctx, blueprintIdentifier, userid)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.HTML(http.StatusOK, "blueprint.html", page)
}

func (m Management) publicBlueprint(ctx *gin.Context) {
	req := new(struct {
		Token string `form:"token"`
	})
	if err := ctx.Bind(req); err != nil {
		return
	}
	claims, err := m.JwtDecode(ctx, req.Token)
	if err != nil {
		ctx.Error(err)
		return
	}
	userid := claims["userid"].(string)
	blueprintIdentifier := claims["blueprintName"].(string)
	page, err := m.makeBlueprintPage(ctx, blueprintIdentifier, userid)
	if err != nil {
		ctx.Error(err)
		return
	}
	userid, _ = ctx.Cookie("userid")
	page.UserId = userid
	ctx.HTML(http.StatusOK, "public.html", page)
}

func (m Management) emptyDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "dashboard.html", &InvokePage{
		Meta:       m.server.Meta(),
		ProtoFiles: m.allProtoFiles(ctx),
	})
}

func (m Management) dashboard(ctx *gin.Context) {
	userid, err := ctx.Cookie("userid")
	if err != nil {
		//todo admin鉴权
		ctx.Error(err)
		return
	}
	serviceIdentifier := ctx.Param("service-identifier")
	serviceIdentifier = strings.TrimPrefix(serviceIdentifier, "/")
	page, err := m.makeInvokePage(ctx, serviceIdentifier, userid)
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
	r.GET("/blueprint", m.Authorize(), m.emptyBlueprint)
	r.GET("/blueprint/*blueprint-identifier", m.Authorize(), m.blueprint)
	r.GET("/public", m.publicBlueprint)
	r.GET("/login", m.login)
	r.GET("/register", m.register)

	rAPI := s.Group("/management/api")
	rAPI.GET("/_intro", m.intro)
	rAPI.GET("/packages", m.listPackages)
	rAPI.GET("/packages/:package_name", m.getPackage)
	rAPI.GET("/service-alias", m.listServiceAlias)

	rAPI.POST("/addressUpdate", m.addressUpdate)
	rAPI.POST("/getMethods", m.getMethods)
	rAPI.POST("/signIn", m.signIn)
	rAPI.POST("/signUp", m.signUp)
	rAPI.POST("/signOut", m.signOut)
	rAPI.POST("/saveHistory", m.saveHistory)

	b := rAPI.Group("/blueprint", m.Authorize())
	b.POST("/new", m.newBlueprint)
	b.POST("/delete", m.delBlueprint)
	b.POST("/copyFromFile", m.copyBlueprintFromFile)
	b.POST("/copy", m.copyBlueprint)
	b.POST("/append", m.savetoBlueprint)
	b.POST("/appendList", m.appendBlueprint)
	b.POST("/reduce", m.deleteBlueprintMethod)
	b.POST("/share", m.shareBlueprint)

	s.GET("/dashboard", m.emptyDashboard)
	s.GET("/dashboard/*service-identifier", m.dashboard)
	s.GET("/setting", m.setting)
	return nil
}

func (m Management) Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cSession, err1 := ctx.Request.Cookie("session")
		cUserid, err2 := ctx.Request.Cookie("userid")
		if err1 == nil && err2 == nil {
			session := cSession.Value
			userid := cUserid.Value
			claims, err := m.JwtDecode(ctx, session)
			if err == nil && claims["userid"].(string) == userid {
				ctx.Next()
				return
			}
		}
		ctx.Abort()
		ctx.HTML(http.StatusOK, "login.html", &LoginPage{
			Meta: m.server.Meta(),
		})
		return
	}
}

func (m Management) JwtDecode(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	var hmacSampleSecret []byte
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func (m Management) Name() string {
	return "management-server"
}

func (m Management) Meta() map[string]string {
	return nil
}

func (m Management) userKey(userid string) string {
	//todo: 键设计
	return "/users/" + userid
}
