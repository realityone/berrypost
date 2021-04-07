package proxy

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type clientSet struct {
	cc  *grpc.ClientConn
	rcc *grpcreflect.Client
}

// Server is
type Server struct {
	*grpc.Server

	resolver   RuntimeServiceResolver
	protoStore RuntimeProtoStore

	clientLock sync.RWMutex
	clients    map[string]*clientSet
}

type dummyMessage struct {
	payload []byte
}

func (dm *dummyMessage) Reset()                   { dm.payload = dm.payload[:0] }
func (dm *dummyMessage) String() string           { return fmt.Sprintf("%q", dm.payload) }
func (dm *dummyMessage) ProtoMessage()            {}
func (dm *dummyMessage) Marshal() ([]byte, error) { return dm.payload, nil }
func (dm *dummyMessage) Unmarshal(in []byte) error {
	dm.payload = append(dm.payload[:0], in...)
	return nil
}

func resolveInOutMessage(rcc *grpcreflect.Client, serviceName string, methodName string) (proto.Message, proto.Message, error) {
	type resolver func() (proto.Message, proto.Message, error)

	byReflect := func() (proto.Message, proto.Message, error) {
		serviceDesc, err := rcc.ResolveService(serviceName)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		for _, method := range serviceDesc.GetMethods() {
			if method.GetName() != methodName {
				continue
			}
			tInput := method.GetInputType()
			tOutput := method.GetOutputType()
			return dynamic.NewMessage(tInput), dynamic.NewMessage(tOutput), nil
		}
		return nil, nil, errors.Errorf("service %q method %q not found by reflect server", serviceName, methodName)
	}
	byLocalProto := func() (proto.Message, proto.Message, error) {
		return nil, nil, errors.New("unimplement")
	}

	for _, resolv := range []resolver{byReflect, byLocalProto} {
		in, out, err := resolv()
		if err != nil {
			logrus.Debugf("Failed to resolve method in out message: %+v", err)
			continue
		}
		return in, out, nil
	}
	return &dummyMessage{}, &dummyMessage{}, nil
}

func (s *Server) client(ctx *Context, service string) (*clientSet, error) {
	s.clientLock.RLock()
	cli, ok := s.clients[service]
	s.clientLock.RUnlock()
	if ok {
		return cli, nil
	}

	target, err := s.resolver.ResolveOnce(ctx, service)
	if err != nil {
		return nil, err
	}

	newCC, err := grpc.DialContext(ctx, target, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	cli, ok = s.clients[service]
	if ok {
		logrus.Debugf("Already has established connection for %q", service)
		newCC.Close()
		return cli, nil
	}
	newCliSet := &clientSet{
		cc:  newCC,
		rcc: grpcreflect.NewClient(context.Background(), rpb.NewServerReflectionClient(newCC)),
	}
	s.clients[service] = newCliSet
	return newCliSet, nil
}

// Handler is
func (s *Server) Handler(ctx *Context) error {
	service, method, err := splitServiceMethod(ctx.ServiceMethod())
	if err != nil {
		return grpc.Errorf(codes.FailedPrecondition, err.Error())
	}

	logrus.Debugf("Handler: service: %+v: method: %+v", service, method)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		logrus.Debugf("In coming metadata: %+v", md)
	}

	cli, err := s.client(ctx, service)
	if err != nil {
		return err
	}

	req, reply, err := s.protoStore.GetMethodMessage(ctx, service, method)
	if err != nil {
		return err
	}

	stream := ctx.ServerStream()
	if err := stream.RecvMsg(req); err != nil {
		return err
	}
	if err := cli.cc.Invoke(ctx, ctx.ServiceMethod(), req, reply); err != nil {
		return err
	}
	if err := stream.SendHeader(md); err != nil {
		return err
	}
	if err := stream.SendMsg(reply); err != nil {
		return err
	}
	logrus.Debugf("Request: %+v, Reply: %+v", req, reply)
	return nil
}

func wrapped(handler func(*Context) error) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		serviceMethod, ok := grpc.MethodFromServerStream(stream)
		if !ok {
			return grpc.Errorf(codes.Internal, "failed to get method from stream")
		}
		ctx := &Context{
			Context:       stream.Context(),
			srv:           srv,
			serverStream:  stream,
			serviceMethod: serviceMethod,
		}
		if err := handler(ctx); err != nil {
			logrus.Errorf("Failed to handle request stream: method: %q: %+v", serviceMethod, err)
			return err
		}
		return nil
	}
}

// New is
func New() *Server {
	server := &Server{
		resolver:   &defaultRuntimeServiceResolver{},
		protoStore: &defaultRuntimeProtoStore{},
		clients:    map[string]*clientSet{},
	}
	server.Server = grpc.NewServer(
		grpc.UnknownServiceHandler(wrapped(server.Handler)),
	)
	return server
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
