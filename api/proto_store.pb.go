// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto_store.proto

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

type GetProtoRequest struct {
	Service              string   `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Method               string   `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetProtoRequest) Reset()         { *m = GetProtoRequest{} }
func (m *GetProtoRequest) String() string { return proto.CompactTextString(m) }
func (*GetProtoRequest) ProtoMessage()    {}
func (*GetProtoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_063e5d34d8874bae, []int{0}
}
func (m *GetProtoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetProtoRequest.Unmarshal(m, b)
}
func (m *GetProtoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetProtoRequest.Marshal(b, m, deterministic)
}
func (m *GetProtoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetProtoRequest.Merge(m, src)
}
func (m *GetProtoRequest) XXX_Size() int {
	return xxx_messageInfo_GetProtoRequest.Size(m)
}
func (m *GetProtoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetProtoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetProtoRequest proto.InternalMessageInfo

func (m *GetProtoRequest) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *GetProtoRequest) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

type GetProtoResponse struct {
	Files                []*ProtoFile `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *GetProtoResponse) Reset()         { *m = GetProtoResponse{} }
func (m *GetProtoResponse) String() string { return proto.CompactTextString(m) }
func (*GetProtoResponse) ProtoMessage()    {}
func (*GetProtoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_063e5d34d8874bae, []int{1}
}
func (m *GetProtoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetProtoResponse.Unmarshal(m, b)
}
func (m *GetProtoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetProtoResponse.Marshal(b, m, deterministic)
}
func (m *GetProtoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetProtoResponse.Merge(m, src)
}
func (m *GetProtoResponse) XXX_Size() int {
	return xxx_messageInfo_GetProtoResponse.Size(m)
}
func (m *GetProtoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetProtoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetProtoResponse proto.InternalMessageInfo

func (m *GetProtoResponse) GetFiles() []*ProtoFile {
	if m != nil {
		return m.Files
	}
	return nil
}

type ProtoFile struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Content              []byte   `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProtoFile) Reset()         { *m = ProtoFile{} }
func (m *ProtoFile) String() string { return proto.CompactTextString(m) }
func (*ProtoFile) ProtoMessage()    {}
func (*ProtoFile) Descriptor() ([]byte, []int) {
	return fileDescriptor_063e5d34d8874bae, []int{2}
}
func (m *ProtoFile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProtoFile.Unmarshal(m, b)
}
func (m *ProtoFile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProtoFile.Marshal(b, m, deterministic)
}
func (m *ProtoFile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtoFile.Merge(m, src)
}
func (m *ProtoFile) XXX_Size() int {
	return xxx_messageInfo_ProtoFile.Size(m)
}
func (m *ProtoFile) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtoFile.DiscardUnknown(m)
}

var xxx_messageInfo_ProtoFile proto.InternalMessageInfo

func (m *ProtoFile) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProtoFile) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterType((*GetProtoRequest)(nil), "berrypost.v1.GetProtoRequest")
	proto.RegisterType((*GetProtoResponse)(nil), "berrypost.v1.GetProtoResponse")
	proto.RegisterType((*ProtoFile)(nil), "berrypost.v1.ProtoFile")
}

func init() { proto.RegisterFile("proto_store.proto", fileDescriptor_063e5d34d8874bae) }

var fileDescriptor_063e5d34d8874bae = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x28, 0xca, 0x2f,
	0xc9, 0x8f, 0x2f, 0x2e, 0xc9, 0x2f, 0x4a, 0xd5, 0x03, 0xb3, 0x85, 0x78, 0x92, 0x52, 0x8b, 0x8a,
	0x2a, 0x0b, 0xf2, 0x8b, 0x4b, 0xf4, 0xca, 0x0c, 0x95, 0x9c, 0xb9, 0xf8, 0xdd, 0x53, 0x4b, 0x02,
	0x40, 0x32, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x12, 0x5c, 0xec, 0xc5, 0xa9, 0x45,
	0x65, 0x99, 0xc9, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x30, 0xae, 0x90, 0x18, 0x17,
	0x5b, 0x6e, 0x6a, 0x49, 0x46, 0x7e, 0x8a, 0x04, 0x13, 0x58, 0x02, 0xca, 0x53, 0x72, 0xe4, 0x12,
	0x40, 0x18, 0x52, 0x5c, 0x90, 0x9f, 0x57, 0x9c, 0x2a, 0xa4, 0xcb, 0xc5, 0x9a, 0x96, 0x99, 0x93,
	0x5a, 0x2c, 0xc1, 0xa8, 0xc0, 0xac, 0xc1, 0x6d, 0x24, 0xae, 0x87, 0x6c, 0xad, 0x1e, 0x58, 0xad,
	0x5b, 0x66, 0x4e, 0x6a, 0x10, 0x44, 0x95, 0x92, 0x25, 0x17, 0x27, 0x5c, 0x4c, 0x48, 0x88, 0x8b,
	0x25, 0x2f, 0x31, 0x17, 0x66, 0x3d, 0x98, 0x0d, 0x72, 0x55, 0x72, 0x7e, 0x5e, 0x49, 0x6a, 0x5e,
	0x09, 0xd8, 0x72, 0x9e, 0x20, 0x18, 0xd7, 0x28, 0x81, 0x4b, 0xd8, 0x09, 0x64, 0x76, 0x40, 0x7e,
	0x31, 0xc4, 0x0d, 0xc1, 0x20, 0xdf, 0x0a, 0x79, 0x72, 0x71, 0xc0, 0x1c, 0x25, 0x24, 0x8b, 0x6a,
	0x3b, 0x9a, 0x8f, 0xa5, 0xe4, 0x70, 0x49, 0x43, 0xfc, 0xe2, 0xc4, 0x1a, 0xc5, 0x9c, 0x58, 0x90,
	0x99, 0xc4, 0x06, 0x0e, 0x40, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb1, 0xbd, 0x48, 0x96,
	0x55, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BerryPostProtoStoreClient is the client API for BerryPostProtoStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BerryPostProtoStoreClient interface {
	GetProto(ctx context.Context, in *GetProtoRequest, opts ...grpc.CallOption) (*GetProtoResponse, error)
}

type berryPostProtoStoreClient struct {
	cc *grpc.ClientConn
}

func NewBerryPostProtoStoreClient(cc *grpc.ClientConn) BerryPostProtoStoreClient {
	return &berryPostProtoStoreClient{cc}
}

func (c *berryPostProtoStoreClient) GetProto(ctx context.Context, in *GetProtoRequest, opts ...grpc.CallOption) (*GetProtoResponse, error) {
	out := new(GetProtoResponse)
	err := c.cc.Invoke(ctx, "/berrypost.v1.BerryPostProtoStore/GetProto", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BerryPostProtoStoreServer is the server API for BerryPostProtoStore service.
type BerryPostProtoStoreServer interface {
	GetProto(context.Context, *GetProtoRequest) (*GetProtoResponse, error)
}

// UnimplementedBerryPostProtoStoreServer can be embedded to have forward compatible implementations.
type UnimplementedBerryPostProtoStoreServer struct {
}

func (*UnimplementedBerryPostProtoStoreServer) GetProto(ctx context.Context, req *GetProtoRequest) (*GetProtoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProto not implemented")
}

func RegisterBerryPostProtoStoreServer(s *grpc.Server, srv BerryPostProtoStoreServer) {
	s.RegisterService(&_BerryPostProtoStore_serviceDesc, srv)
}

func _BerryPostProtoStore_GetProto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProtoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BerryPostProtoStoreServer).GetProto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/berrypost.v1.BerryPostProtoStore/GetProto",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BerryPostProtoStoreServer).GetProto(ctx, req.(*GetProtoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BerryPostProtoStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "berrypost.v1.BerryPostProtoStore",
	HandlerType: (*BerryPostProtoStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProto",
			Handler:    _BerryPostProtoStore_GetProto_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto_store.proto",
}