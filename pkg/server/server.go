package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/realityone/berrypost/pkg/proxy"
	"github.com/sirupsen/logrus"
)

type Server struct {
	*mux.Router
	management *Management
	proxy      *proxy.Server
}

func New() *Server {
	router := mux.NewRouter()
	server := &Server{
		Router:     router,
		management: &Management{},
		proxy:      proxy.New(),
	}
	server.setupRouter()
	return server
}

func (s *Server) Serve() {
	srv := &http.Server{
		Handler:      s.Router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}
	logrus.Infof("Starting server listen and serve at: %s...", srv.Addr)
	logrus.Fatal(srv.ListenAndServe())
}

func intro(w http.ResponseWriter, req *http.Request) {
	introSchema := struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}{
		Name:    "berrypost-server",
		Version: "0.0.1",
	}
	json.NewEncoder(w).Encode(introSchema)
}

func (s *Server) setupRouter() {
	s.PathPrefix("/berrypost").Methods("GET").Path("/api/_intro").HandlerFunc(intro)
	s.management.SetupRoute(s.NewRoute().PathPrefix("/berrypost/management"))
}
