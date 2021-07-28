package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/realityone/berrypost/pkg/protohelper"
	"github.com/realityone/berrypost/pkg/server"
	"github.com/realityone/berrypost/pkg/server/contrib/errorhandler"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
}

type clientID struct {
	service string
	target  string
}

func (cs *clientSet) Close() error {
	return cs.cc.Close()
}

func (ps *ProxyServer) client(ctx *Context, service string) (*clientSet, error) {
	userDefinedTarget, _ := GetUserDefinedTarget(ctx)
	clientKey := clientID{service, userDefinedTarget}

	logrus.Debugf("Resolving service %+v to dial gRPC connection", clientKey)
	target, err := ps.resolver.ResolveOnce(ctx, &ResolveOnceRequest{
		ServiceFullyQualifiedName: service,
		UserDefinedTarget:         userDefinedTarget,
	})
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Dial gRPC connection to service: %+v with target: %q", clientKey, target)
	newCC, err := grpc.DialContext(ctx, target, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	newCliSet := &clientSet{
		cc: newCC,
	}
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
		AnyResolver: protohelper.WrappedAnyResolver{
			AnyResolver: AsContextedAnyResolver(ctx, ps.protoStore),
		},
		Indent: "    ",
	}
	if err := marshaler.Marshal(ctx.Writer, reply); err != nil {
		logrus.Errorf("Failed to marshal reply on method: %q: %+v", invokeCtx.serviceMethod, err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
}

func extractIncommingGRPCMetadata(header http.Header) metadata.MD {
	prefix := http.CanonicalHeaderKey("X-Berrypost-Md-")
	out := metadata.MD{}
	for k, vs := range header {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		name := strings.TrimPrefix(k, prefix)
		if name == "" {
			continue
		}
		out.Append(name, vs...)
	}
	return out
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
	defer cli.Close()

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

	unmarshaler := jsonpb.Unmarshaler{
		AnyResolver: protohelper.WrappedAnyResolver{
			AnyResolver: AsContextedAnyResolver(ctx, ps.protoStore),
		},
	}
	if err := unmarshaler.Unmarshal(ctx.req.Body, req); err != nil {
		return nil, errors.Errorf("Failed to unmarshal json to request message: %+v", err)
	}

	toForward := extractIncommingGRPCMetadata(ctx.req.Header)
	invokeCtx := metadata.NewOutgoingContext(ctx, toForward)
	replyMD := metadata.MD{}
	if err := cli.cc.Invoke(invokeCtx, ctx.serviceMethod, req, reply, grpc.Header(&replyMD)); err != nil {
		return nil, err
	}
	return reply, nil
}

func (p *ProxyServer) Name() string {
	return "proxy-server"
}

func (p *ProxyServer) Meta() map[string]string {
	return nil
}

func (p *ProxyServer) Setup(s *server.Server) error {
	s.POST("/invoke/:service/:method", errorhandler.JSONErrorHandler(), p.ServeHTTP)
	return nil
}

// New is
func New(opts ...ServerOpt) *ProxyServer {
	ps := &ProxyServer{
		resolver:   &defaultRuntimeServiceResolver{},
		protoStore: &defaultRuntimeProtoStore{},
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
