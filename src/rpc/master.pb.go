// Code generated by protoc-gen-go. DO NOT EDIT.
// source: master.proto

package rpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GetSlaveRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSlaveRequest) Reset()         { *m = GetSlaveRequest{} }
func (m *GetSlaveRequest) String() string { return proto.CompactTextString(m) }
func (*GetSlaveRequest) ProtoMessage()    {}
func (*GetSlaveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9c348dec43a6705, []int{0}
}

func (m *GetSlaveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSlaveRequest.Unmarshal(m, b)
}
func (m *GetSlaveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSlaveRequest.Marshal(b, m, deterministic)
}
func (m *GetSlaveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSlaveRequest.Merge(m, src)
}
func (m *GetSlaveRequest) XXX_Size() int {
	return xxx_messageInfo_GetSlaveRequest.Size(m)
}
func (m *GetSlaveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSlaveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSlaveRequest proto.InternalMessageInfo

func (m *GetSlaveRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type GetSlaveResponse struct {
	Addr                 []string `protobuf:"bytes,1,rep,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSlaveResponse) Reset()         { *m = GetSlaveResponse{} }
func (m *GetSlaveResponse) String() string { return proto.CompactTextString(m) }
func (*GetSlaveResponse) ProtoMessage()    {}
func (*GetSlaveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9c348dec43a6705, []int{1}
}

func (m *GetSlaveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSlaveResponse.Unmarshal(m, b)
}
func (m *GetSlaveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSlaveResponse.Marshal(b, m, deterministic)
}
func (m *GetSlaveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSlaveResponse.Merge(m, src)
}
func (m *GetSlaveResponse) XXX_Size() int {
	return xxx_messageInfo_GetSlaveResponse.Size(m)
}
func (m *GetSlaveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSlaveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSlaveResponse proto.InternalMessageInfo

func (m *GetSlaveResponse) GetAddr() []string {
	if m != nil {
		return m.Addr
	}
	return nil
}

func init() {
	proto.RegisterType((*GetSlaveRequest)(nil), "GetSlaveRequest")
	proto.RegisterType((*GetSlaveResponse)(nil), "GetSlaveResponse")
}

func init() { proto.RegisterFile("master.proto", fileDescriptor_f9c348dec43a6705) }

var fileDescriptor_f9c348dec43a6705 = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xc9, 0x4d, 0x2c, 0x2e,
	0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x52, 0xe6, 0xe2, 0x77, 0x4f, 0x2d, 0x09,
	0xce, 0x49, 0x2c, 0x4b, 0x0d, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x12, 0xe0, 0x62, 0xce,
	0x4e, 0xad, 0x94, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0x31, 0x95, 0xd4, 0xb8, 0x04, 0x10,
	0x8a, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85, 0x84, 0xb8, 0x58, 0x12, 0x53, 0x52, 0x8a, 0x24,
	0x18, 0x15, 0x98, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0x23, 0x4b, 0x2e, 0x36, 0x5f, 0xb0, 0xe1, 0x42,
	0xfa, 0x5c, 0x1c, 0x30, 0x1d, 0x42, 0x02, 0x7a, 0x68, 0x36, 0x48, 0x09, 0xea, 0xa1, 0x1b, 0xe7,
	0xc4, 0x1e, 0xc5, 0xaa, 0x67, 0x5d, 0x54, 0x90, 0x9c, 0xc4, 0x06, 0x76, 0x97, 0x31, 0x20, 0x00,
	0x00, 0xff, 0xff, 0x31, 0x23, 0x36, 0x69, 0xa7, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MasterClient is the client API for Master service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MasterClient interface {
	GetSlave(ctx context.Context, in *GetSlaveRequest, opts ...grpc.CallOption) (*GetSlaveResponse, error)
}

type masterClient struct {
	cc *grpc.ClientConn
}

func NewMasterClient(cc *grpc.ClientConn) MasterClient {
	return &masterClient{cc}
}

func (c *masterClient) GetSlave(ctx context.Context, in *GetSlaveRequest, opts ...grpc.CallOption) (*GetSlaveResponse, error) {
	out := new(GetSlaveResponse)
	err := c.cc.Invoke(ctx, "/Master/GetSlave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MasterServer is the server API for Master service.
type MasterServer interface {
	GetSlave(context.Context, *GetSlaveRequest) (*GetSlaveResponse, error)
}

func RegisterMasterServer(s *grpc.Server, srv MasterServer) {
	s.RegisterService(&_Master_serviceDesc, srv)
}

func _Master_GetSlave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSlaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).GetSlave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Master/GetSlave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).GetSlave(ctx, req.(*GetSlaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Master_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Master",
	HandlerType: (*MasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSlave",
			Handler:    _Master_GetSlave_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "master.proto",
}
