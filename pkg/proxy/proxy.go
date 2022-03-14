package proxy

import (
	"encoding/base64"
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

var (
	_headerPrefix  = http.CanonicalHeaderKey("X-Berrypost-Md-")
	_trailerPrefix = http.CanonicalHeaderKey("X-Berrypost-Md-Trailer-")
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

type metadataSet struct {
	header  metadata.MD
	trailer metadata.MD
}

func (ps *ProxyServer) client(ctx *Context, service string) (*clientSet, error) {
	userDefinedTarget, _ := GetUserDefinedTarget(ctx)
	clientKey := clientID{service, userDefinedTarget}

	toForward, err := extractIncommingGRPCMetadata(ctx.req.Header)
	if err != nil {
		return nil, errors.Wrap(err, "extract metadata")
	}
	dialCtx := metadata.NewOutgoingContext(ctx, toForward)

	logrus.Debugf("Resolving service %+v to dial gRPC connection", clientKey)
	target, err := ps.resolver.ResolveOnce(dialCtx, &ResolveOnceRequest{
		ServiceFullyQualifiedName: service,
		UserDefinedTarget:         userDefinedTarget,
	})
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Dial gRPC connection to service: %+v with target: %q", clientKey, target)
	newCC, err := grpc.DialContext(dialCtx, target, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	newCliSet := &clientSet{
		cc: newCC,
	}
	return newCliSet, nil
}

func asBerrypostHeader(in string) string {
	if strings.HasPrefix(in, _headerPrefix) {
		return in
	}
	return http.CanonicalHeaderKey(fmt.Sprintf("%s%s", _headerPrefix, in))
}

func asBerrypostTrailer(in string) string {
	if strings.HasPrefix(in, _trailerPrefix) {
		return in
	}
	return http.CanonicalHeaderKey(fmt.Sprintf("%s%s", _trailerPrefix, in))
}

func writeMetadataAlways(mdSet *metadataSet, dst http.Header) {
	for k, v := range mdSet.header {
		dst[asBerrypostHeader(k)] = v
	}
	for k, v := range mdSet.trailer {
		dst[asBerrypostTrailer(k)] = v
	}
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

	reply, mdSet, err := ps.Invoke(invokeCtx)
	if mdSet != nil {
		writeMetadataAlways(mdSet, ctx.Writer.Header())
	}
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

func decodeBinHeader(v string) ([]byte, error) {
	if len(v)%4 == 0 {
		// Input was padded, or padding was not necessary.
		return base64.StdEncoding.DecodeString(v)
	}
	return base64.RawStdEncoding.DecodeString(v)
}

func decodeMetadataHeader(k, v string) (string, error) {
	const binHdrSuffix = "-Bin"
	if strings.HasSuffix(k, binHdrSuffix) {
		b, err := decodeBinHeader(v)
		return string(b), err
	}
	return v, nil
}

func extractIncommingGRPCMetadata(header http.Header) (metadata.MD, error) {
	const base64Prefix = "base64://"
	out := metadata.MD{}
	for k, vs := range header {
		if !strings.HasPrefix(k, _headerPrefix) {
			continue
		}
		name := strings.TrimPrefix(k, _headerPrefix)
		if name == "" {
			continue
		}
		parsed := make([]string, 0, len(vs))
		for _, v := range vs {
			if strings.HasPrefix(v, base64Prefix) {
				v = strings.TrimPrefix(v, base64Prefix)
				decoded, err := decodeMetadataHeader(k, v)
				if err != nil {
					return nil, errors.Wrapf(err, "Invalid base64 string: %s:%+v", k, v)
				}
				parsed = append(parsed, string(decoded))
				continue
			}
			parsed = append(parsed, v)
		}
		out.Append(name, parsed...)
	}
	return out, nil
}

func (ps *ProxyServer) Invoke(ctx *Context) (proto.Message, *metadataSet, error) {
	service, method, err := splitServiceMethod(ctx.serviceMethod)
	if err != nil {
		return nil, nil, err
	}

	cli, err := ps.client(ctx, service)
	if err != nil {
		return nil, nil, err
	}
	defer cli.Close()

	toForward, err := extractIncommingGRPCMetadata(ctx.req.Header)
	if err != nil {
		return nil, nil, errors.Wrap(err, "extract metadata")
	}
	invokeCtx := metadata.NewOutgoingContext(ctx, toForward)

	req, reply, err := ps.protoStore.GetMethodMessage(invokeCtx, service, method)
	if err != nil {
		return nil, nil, err
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
			AnyResolver: AsContextedAnyResolver(invokeCtx, ps.protoStore),
		},
	}
	if err := unmarshaler.Unmarshal(ctx.req.Body, req); err != nil {
		return nil, nil, errors.Errorf("Failed to unmarshal json to request message: %+v", err)
	}

	mdSet := &metadataSet{
		header:  metadata.MD{},
		trailer: metadata.MD{},
	}
	if err := cli.cc.Invoke(invokeCtx, ctx.serviceMethod, req, reply, grpc.Header(&mdSet.header), grpc.Trailer(&mdSet.trailer)); err != nil {
		return nil, mdSet, err
	}
	return reply, mdSet, nil
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
