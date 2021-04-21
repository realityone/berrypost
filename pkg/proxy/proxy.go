package proxy

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/realityone/berrypost/pkg/protohelper"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

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
	logrus.Debugf("Try to dial gRPC connection to service: %q", service)
	ps.clientLock.RLock()
	cli, ok := ps.clients[service]
	ps.clientLock.RUnlock()
	if ok {
		logrus.Debugf("Got %q gRPC connection from client store", service)
		return cli, nil
	}

	logrus.Debugf("Resolving service %q to dial gRPC connection", service)
	target, err := ps.resolver.ResolveOnce(ctx, service)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Dial gRPC connection to service: %q", service)
	newCC, err := grpc.DialContext(ctx, target, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	ps.clientLock.Lock()
	defer ps.clientLock.Unlock()
	cli, ok = ps.clients[service]
	if ok {
		logrus.Debugf("Already has established connection for service: %q", service)
		newCC.Close()
		return cli, nil
	}
	logrus.Debugf("Put new gRPC connection to client store: %q", service)
	newCliSet := &clientSet{
		cc: newCC,
	}
	ps.clients[service] = newCliSet
	return newCliSet, nil
}

func (ps *ProxyServer) ServeHTTP(ctx *gin.Context) {
	invokeCtx := &Context{
		Context: ctx,
		req:     ctx.Request,
		writer:  ctx.Writer,
	}

	service, method := ctx.Param("service"), ctx.Param("method")
	invokeCtx.serviceMethod = fmt.Sprintf("/%s/%s", service, method)
	logrus.Debugf("Received gRPC call from http: %q", invokeCtx.serviceMethod)

	reply, err := ps.Invoke(invokeCtx)
	if err != nil {
		logrus.Errorf("Failed to invoke backend on method: %q: %+v", invokeCtx.serviceMethod, err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	marshaler := &jsonpb.Marshaler{
		AnyResolver: AsContextedAnyResolver(ctx, ps.protoStore),
	}
	if err := marshaler.Marshal(ctx.Writer, reply); err != nil {
		logrus.Errorf("Failed to marshal reply on method: %q: %+v", invokeCtx.serviceMethod, err)
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
	logrus.DebugFn(func() []interface{} {
		return []interface{}{
			fmt.Sprintf(
				"Succeeded to get message type from proto store, request: %+v, reply: %+v",
				protohelper.AsEmptyMessageJSON(req),
				protohelper.AsEmptyMessageJSON(reply),
			),
		}
	})

	if err := jsonpb.Unmarshal(ctx.req.Body, req); err != nil {
		return nil, errors.Errorf("Failed to unmarshal json to request message: %+v", err)
	}

	if err := cli.cc.Invoke(ctx, ctx.serviceMethod, req, reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (p *ProxyServer) SetupRoute(in *gin.Engine) {
	in.POST("/invoke/:service/:method", p.ServeHTTP)
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
