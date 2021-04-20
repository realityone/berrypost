package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/realityone/berrypost/pkg/proxy"
	"github.com/realityone/berrypost/pkg/server/management"
	"github.com/sirupsen/logrus"
)

type Option func(*ServerConfig)
type ServerConfig struct {
	ProxyOptions []proxy.ServerOpt
}

func SetProxyOptions(in []proxy.ServerOpt) Option {
	return func(sc *ServerConfig) {
		sc.ProxyOptions = in
	}
}

type Server struct {
	*gin.Engine
	management *management.Management
	proxy      *proxy.ProxyServer

	component []string
}

func New(opts ...Option) *Server {
	cfg := &ServerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	engine := gin.Default()
	server := &Server{
		Engine:     engine,
		management: management.New(),
		proxy:      proxy.New(cfg.ProxyOptions...),
		component:  []string{},
	}
	server.setupRouter()
	return server
}

func (s *Server) Serve() {
	srv := &http.Server{
		Handler:      s.WrappedHandler(),
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}
	logrus.Infof("Starting server listen and serve at: %s...", srv.Addr)
	logrus.Fatal(srv.ListenAndServe())
}

func (s *Server) WrappedHandler() http.Handler {
	return s.Engine
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

func (s *Server) setupRouter() {
	s.GET("/berrypost/api/_intro", s.intro)
	// builtin components
	s.management.SetupRoute(s.Group("/berrypost/management"))
	s.proxy.SetupRoute(s.Group("/berrypost/invoke"))
}

func (s *Server) RegisterComponentAPI(component string, fn func(gin.IRouter)) {
	s.component = append(s.component, component)
	fn(s.Group(fmt.Sprintf("/berrypost/component/%s", component)))
}
