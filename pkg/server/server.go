package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	*mux.Router
	management *management.Management
	proxy      *proxy.ProxyServer
}

func New(opts ...Option) *Server {
	cfg := &ServerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	router := mux.NewRouter()
	server := &Server{
		Router:     router,
		management: &management.Management{},
		proxy:      proxy.New(cfg.ProxyOptions...),
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
	out := handlers.CombinedLoggingHandler(logrus.StandardLogger().Out, s.Router)
	out = handlers.CORS()(out)
	out = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(out)
	return out
}

func (s *Server) intro(w http.ResponseWriter, req *http.Request) {
	introSchema := struct {
		Name    string   `json:"name"`
		Version string   `json:"version"`
		Paths   []string `json:"paths"`
	}{
		Name:    "berrypost-server",
		Version: "0.0.1",
	}
	s.Router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		introSchema.Paths = append(introSchema.Paths, pathTemplate)
		return nil
	})
	json.NewEncoder(w).Encode(introSchema)
}

func (s *Server) setupRouter() {
	s.PathPrefix("/berrypost").Methods("GET").Path("/api/_intro").HandlerFunc(s.intro)
	// builtin components
	s.management.SetupRoute(s.PathPrefix("/berrypost/management").Subrouter())
	s.proxy.SetupRoute(s.PathPrefix("/berrypost/invoke").Subrouter())
}

func (s *Server) RegisterComponentAPI(component string, fn func(*mux.Router)) {
	fn(s.PathPrefix(fmt.Sprintf("/berrypost/component/%s", component)).Subrouter())
}
