// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: resolver.proto

package api

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ResolveOnceRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResolveOnceRequest) Reset()         { *m = ResolveOnceRequest{} }
func (m *ResolveOnceRequest) String() string { return proto.CompactTextString(m) }
func (*ResolveOnceRequest) ProtoMessage()    {}
func (*ResolveOnceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5838971722c666f, []int{0}
}
func (m *ResolveOnceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResolveOnceRequest.Unmarshal(m, b)
}
func (m *ResolveOnceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResolveOnceRequest.Marshal(b, m, deterministic)
}
func (m *ResolveOnceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResolveOnceRequest.Merge(m, src)
}
func (m *ResolveOnceRequest) XXX_Size() int {
	return xxx_messageInfo_ResolveOnceRequest.Size(m)
}
func (m *ResolveOnceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ResolveOnceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ResolveOnceRequest proto.InternalMessageInfo

func (m *ResolveOnceRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type ResolveOnceResponse struct {
	Addrs                []string `protobuf:"bytes,1,rep,name=addrs,proto3" json:"addrs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResolveOnceResponse) Reset()         { *m = ResolveOnceResponse{} }
func (m *ResolveOnceResponse) String() string { return proto.CompactTextString(m) }
func (*ResolveOnceResponse) ProtoMessage()    {}
func (*ResolveOnceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5838971722c666f, []int{1}
}
func (m *ResolveOnceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResolveOnceResponse.Unmarshal(m, b)
}
func (m *ResolveOnceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResolveOnceResponse.Marshal(b, m, deterministic)
}
func (m *ResolveOnceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResolveOnceResponse.Merge(m, src)
}
func (m *ResolveOnceResponse) XXX_Size() int {
	return xxx_messageInfo_ResolveOnceResponse.Size(m)
}
func (m *ResolveOnceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ResolveOnceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ResolveOnceResponse proto.InternalMessageInfo

func (m *ResolveOnceResponse) GetAddrs() []string {
	if m != nil {
		return m.Addrs
	}
	return nil
}

func init() {
	proto.RegisterType((*ResolveOnceRequest)(nil), "berrypost.v1.ResolveOnceRequest")
	proto.RegisterType((*ResolveOnceResponse)(nil), "berrypost.v1.ResolveOnceResponse")
}

func init() { proto.RegisterFile("resolver.proto", fileDescriptor_f5838971722c666f) }

var fileDescriptor_f5838971722c666f = []byte{
	// 166 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x4a, 0x2d, 0xce,
	0xcf, 0x29, 0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x49, 0x4a, 0x2d, 0x2a,
	0xaa, 0x2c, 0xc8, 0x2f, 0x2e, 0xd1, 0x2b, 0x33, 0x54, 0xd2, 0xe0, 0x12, 0x0a, 0x82, 0xc8, 0xfb,
	0xe7, 0x25, 0xa7, 0x06, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x09, 0x71, 0xb1, 0xe4, 0x25,
	0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x4a, 0xda, 0x5c, 0xc2, 0x28,
	0x2a, 0x8b, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85, 0x44, 0xb8, 0x58, 0x13, 0x53, 0x52, 0x8a, 0x8a,
	0x25, 0x18, 0x15, 0x98, 0x35, 0x38, 0x83, 0x20, 0x1c, 0xa3, 0x74, 0x2e, 0x41, 0x27, 0x90, 0x35,
	0x01, 0xf9, 0xc5, 0x25, 0x50, 0x5d, 0x45, 0x42, 0x41, 0x5c, 0xdc, 0x48, 0x26, 0x08, 0x29, 0xe8,
	0x21, 0xbb, 0x44, 0x0f, 0xd3, 0x19, 0x52, 0x8a, 0x78, 0x54, 0x40, 0xac, 0x77, 0x62, 0x8d, 0x62,
	0x4e, 0x2c, 0xc8, 0x4c, 0x62, 0x03, 0xfb, 0xcd, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xce, 0xa3,
	0x5e, 0x07, 0xed, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BerryPostResolverClient is the client API for BerryPostResolver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BerryPostResolverClient interface {
	ResolveOnce(ctx context.Context, in *ResolveOnceRequest, opts ...grpc.CallOption) (*ResolveOnceResponse, error)
}

type berryPostResolverClient struct {
	cc *grpc.ClientConn
}

func NewBerryPostResolverClient(cc *grpc.ClientConn) BerryPostResolverClient {
	return &berryPostResolverClient{cc}
}

func (c *berryPostResolverClient) ResolveOnce(ctx context.Context, in *ResolveOnceRequest, opts ...grpc.CallOption) (*ResolveOnceResponse, error) {
	out := new(ResolveOnceResponse)
	err := c.cc.Invoke(ctx, "/berrypost.v1.BerryPostResolver/ResolveOnce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BerryPostResolverServer is the server API for BerryPostResolver service.
type BerryPostResolverServer interface {
	ResolveOnce(context.Context, *ResolveOnceRequest) (*ResolveOnceResponse, error)
}

// UnimplementedBerryPostResolverServer can be embedded to have forward compatible implementations.
type UnimplementedBerryPostResolverServer struct {
}

func (*UnimplementedBerryPostResolverServer) ResolveOnce(ctx context.Context, req *ResolveOnceRequest) (*ResolveOnceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveOnce not implemented")
}

func RegisterBerryPostResolverServer(s *grpc.Server, srv BerryPostResolverServer) {
	s.RegisterService(&_BerryPostResolver_serviceDesc, srv)
}

func _BerryPostResolver_ResolveOnce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResolveOnceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BerryPostResolverServer).ResolveOnce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/berrypost.v1.BerryPostResolver/ResolveOnce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BerryPostResolverServer).ResolveOnce(ctx, req.(*ResolveOnceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BerryPostResolver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "berrypost.v1.BerryPostResolver",
	HandlerType: (*BerryPostResolverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResolveOnce",
			Handler:    _BerryPostResolver_ResolveOnce_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "resolver.proto",
}