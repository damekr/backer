// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protoclnt/baclnt.proto

/*
Package protoclnt is a generated protocol buffer package.

It is generated from these files:
	protoclnt/baclnt.proto

It has these top-level messages:
	PingRequest
	PingResponse
	BackupRequest
	BackupResponse
	RestoreRequest
	RestoreResponse
*/
package protoclnt

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PingRequest struct {
	Ip string `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
}

func (m *PingRequest) Reset()                    { *m = PingRequest{} }
func (m *PingRequest) String() string            { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()               {}
func (*PingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PingRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

type PingResponse struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *PingResponse) Reset()                    { *m = PingResponse{} }
func (m *PingResponse) String() string            { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()               {}
func (*PingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PingResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type BackupRequest struct {
	Ip    string   `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
	Paths []string `protobuf:"bytes,2,rep,name=paths" json:"paths,omitempty"`
}

func (m *BackupRequest) Reset()                    { *m = BackupRequest{} }
func (m *BackupRequest) String() string            { return proto.CompactTextString(m) }
func (*BackupRequest) ProtoMessage()               {}
func (*BackupRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *BackupRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *BackupRequest) GetPaths() []string {
	if m != nil {
		return m.Paths
	}
	return nil
}

type BackupResponse struct {
	Validpaths []string `protobuf:"bytes,1,rep,name=validpaths" json:"validpaths,omitempty"`
}

func (m *BackupResponse) Reset()                    { *m = BackupResponse{} }
func (m *BackupResponse) String() string            { return proto.CompactTextString(m) }
func (*BackupResponse) ProtoMessage()               {}
func (*BackupResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *BackupResponse) GetValidpaths() []string {
	if m != nil {
		return m.Validpaths
	}
	return nil
}

type RestoreRequest struct {
	Ip             string   `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
	AssetID        int32    `protobuf:"varint,2,opt,name=assetID" json:"assetID,omitempty"`
	WholeBackup    bool     `protobuf:"varint,3,opt,name=wholeBackup" json:"wholeBackup,omitempty"`
	RestoreObjects []string `protobuf:"bytes,4,rep,name=restoreObjects" json:"restoreObjects,omitempty"`
	BasePath       string   `protobuf:"bytes,5,opt,name=basePath" json:"basePath,omitempty"`
}

func (m *RestoreRequest) Reset()                    { *m = RestoreRequest{} }
func (m *RestoreRequest) String() string            { return proto.CompactTextString(m) }
func (*RestoreRequest) ProtoMessage()               {}
func (*RestoreRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *RestoreRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *RestoreRequest) GetAssetID() int32 {
	if m != nil {
		return m.AssetID
	}
	return 0
}

func (m *RestoreRequest) GetWholeBackup() bool {
	if m != nil {
		return m.WholeBackup
	}
	return false
}

func (m *RestoreRequest) GetRestoreObjects() []string {
	if m != nil {
		return m.RestoreObjects
	}
	return nil
}

func (m *RestoreRequest) GetBasePath() string {
	if m != nil {
		return m.BasePath
	}
	return ""
}

type RestoreResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *RestoreResponse) Reset()                    { *m = RestoreResponse{} }
func (m *RestoreResponse) String() string            { return proto.CompactTextString(m) }
func (*RestoreResponse) ProtoMessage()               {}
func (*RestoreResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *RestoreResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*PingRequest)(nil), "protoclnt.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "protoclnt.PingResponse")
	proto.RegisterType((*BackupRequest)(nil), "protoclnt.BackupRequest")
	proto.RegisterType((*BackupResponse)(nil), "protoclnt.BackupResponse")
	proto.RegisterType((*RestoreRequest)(nil), "protoclnt.RestoreRequest")
	proto.RegisterType((*RestoreResponse)(nil), "protoclnt.RestoreResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Baclnt service

type BaclntClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Backup(ctx context.Context, in *BackupRequest, opts ...grpc.CallOption) (*BackupResponse, error)
	Restore(ctx context.Context, in *RestoreRequest, opts ...grpc.CallOption) (*RestoreResponse, error)
}

type baclntClient struct {
	cc *grpc.ClientConn
}

func NewBaclntClient(cc *grpc.ClientConn) BaclntClient {
	return &baclntClient{cc}
}

func (c *baclntClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := grpc.Invoke(ctx, "/protoclnt.Baclnt/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *baclntClient) Backup(ctx context.Context, in *BackupRequest, opts ...grpc.CallOption) (*BackupResponse, error) {
	out := new(BackupResponse)
	err := grpc.Invoke(ctx, "/protoclnt.Baclnt/Backup", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *baclntClient) Restore(ctx context.Context, in *RestoreRequest, opts ...grpc.CallOption) (*RestoreResponse, error) {
	out := new(RestoreResponse)
	err := grpc.Invoke(ctx, "/protoclnt.Baclnt/Restore", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Baclnt service

type BaclntServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Backup(context.Context, *BackupRequest) (*BackupResponse, error)
	Restore(context.Context, *RestoreRequest) (*RestoreResponse, error)
}

func RegisterBaclntServer(s *grpc.Server, srv BaclntServer) {
	s.RegisterService(&_Baclnt_serviceDesc, srv)
}

func _Baclnt_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BaclntServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protoclnt.Baclnt/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BaclntServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Baclnt_Backup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BackupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BaclntServer).Backup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protoclnt.Baclnt/Backup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BaclntServer).Backup(ctx, req.(*BackupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Baclnt_Restore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BaclntServer).Restore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protoclnt.Baclnt/Restore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BaclntServer).Restore(ctx, req.(*RestoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Baclnt_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protoclnt.Baclnt",
	HandlerType: (*BaclntServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Baclnt_Ping_Handler,
		},
		{
			MethodName: "Backup",
			Handler:    _Baclnt_Backup_Handler,
		},
		{
			MethodName: "Restore",
			Handler:    _Baclnt_Restore_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protoclnt/baclnt.proto",
}

func init() { proto.RegisterFile("protoclnt/baclnt.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x51, 0xcd, 0x4e, 0xf3, 0x30,
	0x10, 0xfc, 0x9c, 0xfe, 0x6f, 0x3f, 0x82, 0x64, 0xa1, 0x62, 0x22, 0x81, 0x22, 0x1f, 0x50, 0xb8,
	0x14, 0x04, 0xe2, 0xc0, 0x09, 0xa9, 0xe2, 0xc2, 0x89, 0xca, 0x6f, 0xe0, 0x14, 0xab, 0x0d, 0x94,
	0xc6, 0x64, 0x1d, 0x78, 0x1d, 0x1e, 0x87, 0xc7, 0x42, 0xb1, 0x9d, 0x28, 0x85, 0xf6, 0x14, 0xcd,
	0xec, 0x64, 0x76, 0xbc, 0x03, 0x13, 0x5d, 0xe4, 0x26, 0x5f, 0xac, 0x37, 0xe6, 0x32, 0x95, 0xd5,
	0x67, 0x6a, 0x09, 0x3a, 0x6a, 0x78, 0x7e, 0x0a, 0xe3, 0x79, 0xb6, 0x59, 0x0a, 0xf5, 0x5e, 0x2a,
	0x34, 0x34, 0x84, 0x20, 0xd3, 0x8c, 0xc4, 0x24, 0x19, 0x89, 0x20, 0xd3, 0x3c, 0x81, 0xff, 0x6e,
	0x8c, 0x3a, 0xdf, 0xa0, 0xa2, 0x0c, 0x06, 0x6f, 0x0a, 0x51, 0x2e, 0x95, 0x17, 0xd5, 0x90, 0xdf,
	0xc2, 0xc1, 0x4c, 0x2e, 0x5e, 0x4b, 0xbd, 0xc7, 0x8a, 0x1e, 0x41, 0x4f, 0x4b, 0xb3, 0x42, 0x16,
	0xc4, 0x9d, 0x64, 0x24, 0x1c, 0xe0, 0x57, 0x10, 0xd6, 0xbf, 0xf9, 0x15, 0x67, 0x00, 0x1f, 0x72,
	0x9d, 0x3d, 0x3b, 0x31, 0xb1, 0xe2, 0x16, 0xc3, 0xbf, 0x08, 0x84, 0x42, 0xa1, 0xc9, 0x0b, 0xb5,
	0x6f, 0x15, 0x83, 0x81, 0x44, 0x54, 0xe6, 0xf1, 0x81, 0x05, 0x31, 0x49, 0x7a, 0xa2, 0x86, 0x34,
	0x86, 0xf1, 0xe7, 0x2a, 0x5f, 0x2b, 0xb7, 0x93, 0x75, 0x62, 0x92, 0x0c, 0x45, 0x9b, 0xa2, 0xe7,
	0x10, 0x16, 0xce, 0xfd, 0x29, 0x7d, 0x51, 0x0b, 0x83, 0xac, 0x6b, 0x23, 0xfc, 0x62, 0x69, 0x04,
	0xc3, 0x54, 0xa2, 0x9a, 0x4b, 0xb3, 0x62, 0x3d, 0xbb, 0xb9, 0xc1, 0xfc, 0x02, 0x0e, 0x9b, 0x84,
	0xfe, 0x55, 0x13, 0xe8, 0xa3, 0x91, 0xa6, 0x44, 0x1f, 0xd3, 0xa3, 0xeb, 0x6f, 0x02, 0xfd, 0x99,
	0xed, 0x86, 0xde, 0x41, 0xb7, 0xba, 0x35, 0x9d, 0x4c, 0x9b, 0x7a, 0xa6, 0xad, 0x6e, 0xa2, 0xe3,
	0x3f, 0xbc, 0xf3, 0xe6, 0xff, 0xe8, 0xbd, 0x35, 0xa9, 0xe2, 0xb3, 0x96, 0x68, 0xab, 0x8f, 0xe8,
	0x64, 0xc7, 0xa4, 0x31, 0x98, 0xc1, 0xc0, 0x27, 0xa6, 0x6d, 0xdd, 0xf6, 0x9d, 0xa3, 0x68, 0xd7,
	0xa8, 0xf6, 0x48, 0xfb, 0x76, 0x78, 0xf3, 0x13, 0x00, 0x00, 0xff, 0xff, 0x1e, 0x75, 0x0c, 0xe9,
	0x76, 0x02, 0x00, 0x00,
}