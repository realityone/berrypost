package proxy

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var pbMarshaler = &jsonpb.Marshaler{}

type clientSet struct {
	cc *grpc.ClientConn
}

type ServerOpt func(*ProxyServer)

// ProxyServer is
// TODO: treat as a real gRPC server.
type ProxyServer struct {
	resolver   RuntimeServiceResolver
	protoStore RuntimeProtoStore

	clientLock sync.RWMutex
	clients    map[string]*clientSet
}

func (ps *ProxyServer) client(ctx *Context, service string) (*clientSet, error) {
	ps.clientLock.RLock()
	cli, ok := ps.clients[service]
	ps.clientLock.RUnlock()
	if ok {
		return cli, nil
	}

	target, err := ps.resolver.ResolveOnce(ctx, service)
	if err != nil {
		return nil, err
	}

	newCC, err := grpc.DialContext(ctx, target, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	ps.clientLock.Lock()
	defer ps.clientLock.Unlock()
	cli, ok = ps.clients[service]
	if ok {
		logrus.Debugf("Already has established connection for %q", service)
		newCC.Close()
		return cli, nil
	}
	newCliSet := &clientSet{
		cc: newCC,
	}
	ps.clients[service] = newCliSet
	return newCliSet, nil
}

func (ps *ProxyServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Context: context.Background(),
		req:     req,
		writer:  w,
	}

	vars := mux.Vars(req)
	service, method := vars["service"], vars["method"]
	ctx.serviceMethod = fmt.Sprintf("/%s/%s", service, method)

	reply, err := ps.Invoke(ctx)
	if err != nil {
		logrus.Errorf("Failed to invoke backend on method: %q: %+v", ctx.serviceMethod, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := pbMarshaler.Marshal(w, reply); err != nil {
		logrus.Errorf("Failed to marshal reply on method: %q: %+v", ctx.serviceMethod, err)
		return
	}
}

func (ps *ProxyServer) Invoke(ctx *Context) (proto.Message, error) {
	service, method, err := splitServiceMethod(ctx.serviceMethod)
	if err != nil {
		return nil, err
	}

	cli, err := ps.client(ctx, service)
	if err != nil {
		return nil, err
	}

	req, reply, err := ps.protoStore.GetMethodMessage(ctx, service, method)
	if err != nil {
		return nil, err
	}

	if err := jsonpb.Unmarshal(ctx.req.Body, req); err != nil {
		return nil, errors.Errorf("Failed to unmarshal json to request message: %+v", err)
	}

	if err := cli.cc.Invoke(ctx, ctx.serviceMethod, req, reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (p *ProxyServer) SetupRoute(in *mux.Router) {
	in.Path("/{service}/{method}").Methods("POST").Handler(p)
}

// New is
func New(opts ...ServerOpt) *ProxyServer {
	ps := &ProxyServer{
		resolver:   &defaultRuntimeServiceResolver{},
		protoStore: &defaultRuntimeProtoStore{},
		clients:    map[string]*clientSet{},
	}
	for _, opt := range opts {
		opt(ps)
	}
	return ps
}

func splitServiceMethod(serviceMethod string) (string, string, error) {
	if serviceMethod != "" && serviceMethod[0] == '/' {
		serviceMethod = serviceMethod[1:]
	}
	pos := strings.LastIndex(serviceMethod, "/")
	if pos == -1 {
		return "", "", errors.Errorf("malformed method name: %q", serviceMethod)
	}
	service := serviceMethod[:pos]
	method := serviceMethod[pos+1:]
	return service, method, nil
}

func SetResolver(in RuntimeServiceResolver) ServerOpt {
	return func(s *ProxyServer) {
		s.resolver = in
	}
}

func SetProtoStore(in RuntimeProtoStore) ServerOpt {
	return func(s *ProxyServer) {
		s.protoStore = in
	}
}
