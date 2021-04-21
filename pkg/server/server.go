package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/realityone/berrypost"
	"github.com/realityone/berrypost/pkg/proxy"
	"github.com/realityone/berrypost/pkg/server/management"
	"github.com/sirupsen/logrus"
)

type Option func(*ServerConfig)
type ServerConfig struct {
	ProxyOptions []proxy.ServerOpt
	Meta         ServerMeta
}

func SetProxyOptions(in []proxy.ServerOpt) Option {
	return func(sc *ServerConfig) {
		sc.ProxyOptions = in
	}
}

func SetServerMeta(in ServerMeta) Option {
	return func(sc *ServerConfig) {
		sc.Meta = in
	}
}

type ServerMeta struct {
	Name        string
	Description string
}

type Server struct {
	*gin.Engine
	management *management.Management
	proxy      *proxy.ProxyServer

	component []string
	meta      ServerMeta
}

func New(opts ...Option) *Server {
	cfg := &ServerConfig{
		Meta: ServerMeta{
			Name:        "berrypost",
			Description: "Berrypost is a simple gRPC service debugging tool, built for human beings.",
		},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	engine := gin.Default()
	server := &Server{
		Engine:     engine,
		management: management.New(),
		proxy:      proxy.New(cfg.ProxyOptions...),
		component:  []string{},
		meta:       cfg.Meta,
	}

	templ := template.Must(template.ParseFS(berrypost.TemplateFS, "statics/templates/*.html"))
	engine.SetHTMLTemplate(templ)

	server.setupRouter()
	return server
}

func (s *Server) Serve() {
	srv := &http.Server{
		Handler:      s,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}
	logrus.Infof("Starting server listen and serve at: %s...", srv.Addr)
	logrus.Fatal(srv.ListenAndServe())
}

func (s *Server) intro(ctx *gin.Context) {
	introSchema := struct {
		Name      string   `json:"name"`
		Version   string   `json:"version"`
		Paths     []string `json:"paths"`
		Component []string `json:"component"`
	}{
		Name:      "berrypost-server",
		Version:   "0.0.1",
		Component: s.component,
	}
	for _, r := range s.Engine.Routes() {
		introSchema.Paths = append(introSchema.Paths, r.Path)
	}
	ctx.JSON(200, introSchema)
}

func (s *Server) index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", s.meta)
}

func (s *Server) setupRouter() {
	s.GET("/", s.index)
	s.GET("/api/_intro", s.intro)
	s.StaticFS("/assets", http.FS(berrypost.DistFS))

	// builtin components
	s.management.SetupRoute(s.Engine)
	s.proxy.SetupRoute(s.Engine)

}

func (s *Server) RegisterComponentAPI(component string, fn func(gin.IRouter)) {
	s.component = append(s.component, component)
	fn(s.Group(fmt.Sprintf("/component/%s", component)))
}
