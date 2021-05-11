package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/realityone/berrypost"
	"github.com/realityone/berrypost/pkg/server/contrib/cacheablefs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Option func(*ServerConfig)
type ServerConfig struct {
	Components     []Component
	Meta           ServerMeta
	GinMiddlewares []gin.HandlerFunc
}

type Component interface {
	Name() string
	Meta() map[string]string
	Setup(*Server) error
}

func SetComponents(in []Component) Option {
	return func(sc *ServerConfig) {
		sc.Components = in
	}
}

func SetServerMeta(in ServerMeta) Option {
	return func(sc *ServerConfig) {
		sc.Meta = in
	}
}

func SetGinMiddlewares(in []gin.HandlerFunc) Option {
	return func(sc *ServerConfig) {
		sc.GinMiddlewares = in
	}
}

type ServerMeta struct {
	Name        string
	Description string
	GitHubLink  bool
}

type Server struct {
	*gin.Engine

	components []Component
	meta       ServerMeta
}

func New(opts ...Option) *Server {
	cfg := &ServerConfig{
		Meta: ServerMeta{
			Name:        "berrypost",
			Description: "Berrypost is a simple gRPC service debugging tool, built for human beings.",
			GitHubLink:  true,
		},
		GinMiddlewares: []gin.HandlerFunc{gin.Logger(), gin.Recovery()},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	engine := gin.New()
	engine.Use(cfg.GinMiddlewares...)
	server := &Server{
		Engine:     engine,
		components: cfg.Components,
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

func (s *Server) Meta() ServerMeta {
	return s.meta
}

func (s *Server) componentsMeta() map[string]map[string]string {
	out := map[string]map[string]string{}
	for _, c := range s.components {
		out[c.Name()] = c.Meta()
	}
	return out
}

func (s *Server) intro(ctx *gin.Context) {
	introSchema := struct {
		Name      string                       `json:"name"`
		Version   string                       `json:"version"`
		Paths     []string                     `json:"paths"`
		Component map[string]map[string]string `json:"component"`
	}{
		Name:      "berrypost-server",
		Version:   "0.0.1",
		Component: s.componentsMeta(),
	}
	for _, r := range s.Engine.Routes() {
		introSchema.Paths = append(introSchema.Paths, r.Path)
	}
	ctx.JSON(http.StatusOK, introSchema)
}

func (s *Server) index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", s.meta)
}

func (s *Server) favicon(ctx *gin.Context) {
	ctx.Data(http.StatusOK, http.DetectContentType(berrypost.Icon), berrypost.Icon)
}

func (s *Server) setupRouter() {
	s.GET("/", s.index)
	s.GET("/favicon.ico", s.favicon)
	s.GET("/api/_intro", s.intro)
	s.StaticFS("/assets", http.FS(cacheablefs.Wrap(berrypost.DistFS)))

	for _, c := range s.components {
		if err := s.SetComponent(c); err != nil {
			logrus.Error("Failed to setup component: %+v: %+v", c.Name(), err)
			continue
		}
	}
}

func (s *Server) SetComponent(c Component) error {
	return c.Setup(s)
}
