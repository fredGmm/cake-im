// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chat/grpc/chat.proto

// cmd 进入rpc目录  protoc chat/grpc/chat.proto --go_out=plugins=grpc:./   未生成client时，加上plugins

package grpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PingReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingReq) Reset()         { *m = PingReq{} }
func (m *PingReq) String() string { return proto.CompactTextString(m) }
func (*PingReq) ProtoMessage()    {}
func (*PingReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_332787bef78344c2, []int{0}
}

func (m *PingReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingReq.Unmarshal(m, b)
}
func (m *PingReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingReq.Marshal(b, m, deterministic)
}
func (m *PingReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingReq.Merge(m, src)
}
func (m *PingReq) XXX_Size() int {
	return xxx_messageInfo_PingReq.Size(m)
}
func (m *PingReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PingReq.DiscardUnknown(m)
}

var xxx_messageInfo_PingReq proto.InternalMessageInfo

type PingReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingReply) Reset()         { *m = PingReply{} }
func (m *PingReply) String() string { return proto.CompactTextString(m) }
func (*PingReply) ProtoMessage()    {}
func (*PingReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_332787bef78344c2, []int{1}
}

func (m *PingReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingReply.Unmarshal(m, b)
}
func (m *PingReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingReply.Marshal(b, m, deterministic)
}
func (m *PingReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingReply.Merge(m, src)
}
func (m *PingReply) XXX_Size() int {
	return xxx_messageInfo_PingReply.Size(m)
}
func (m *PingReply) XXX_DiscardUnknown() {
	xxx_messageInfo_PingReply.DiscardUnknown(m)
}

var xxx_messageInfo_PingReply proto.InternalMessageInfo

type PushMsgReq struct {
	MsgId                string      `protobuf:"bytes,1,opt,name=msgId,proto3" json:"msgId,omitempty"`
	Content              string      `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	UserId               string      `protobuf:"bytes,3,opt,name=userId,proto3" json:"userId,omitempty"`
	GroupId              string      `protobuf:"bytes,4,opt,name=groupId,proto3" json:"groupId,omitempty"`
	ToUserId             string      `protobuf:"bytes,5,opt,name=toUserId,proto3" json:"toUserId,omitempty"`
	Type                 string      `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
	Redis                *RedisParam `protobuf:"bytes,7,opt,name=redis,proto3" json:"redis,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *PushMsgReq) Reset()         { *m = PushMsgReq{} }
func (m *PushMsgReq) String() string { return proto.CompactTextString(m) }
func (*PushMsgReq) ProtoMessage()    {}
func (*PushMsgReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_332787bef78344c2, []int{2}
}

func (m *PushMsgReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushMsgReq.Unmarshal(m, b)
}
func (m *PushMsgReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushMsgReq.Marshal(b, m, deterministic)
}
func (m *PushMsgReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushMsgReq.Merge(m, src)
}
func (m *PushMsgReq) XXX_Size() int {
	return xxx_messageInfo_PushMsgReq.Size(m)
}
func (m *PushMsgReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PushMsgReq.DiscardUnknown(m)
}

var xxx_messageInfo_PushMsgReq proto.InternalMessageInfo

func (m *PushMsgReq) GetMsgId() string {
	if m != nil {
		return m.MsgId
	}
	return ""
}

func (m *PushMsgReq) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *PushMsgReq) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *PushMsgReq) GetGroupId() string {
	if m != nil {
		return m.GroupId
	}
	return ""
}

func (m *PushMsgReq) GetToUserId() string {
	if m != nil {
		return m.ToUserId
	}
	return ""
}

func (m *PushMsgReq) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *PushMsgReq) GetRedis() *RedisParam {
	if m != nil {
		return m.Redis
	}
	return nil
}

type RedisParam struct {
	Addr                 string   `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RedisParam) Reset()         { *m = RedisParam{} }
func (m *RedisParam) String() string { return proto.CompactTextString(m) }
func (*RedisParam) ProtoMessage()    {}
func (*RedisParam) Descriptor() ([]byte, []int) {
	return fileDescriptor_332787bef78344c2, []int{3}
}

func (m *RedisParam) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RedisParam.Unmarshal(m, b)
}
func (m *RedisParam) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RedisParam.Marshal(b, m, deterministic)
}
func (m *RedisParam) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RedisParam.Merge(m, src)
}
func (m *RedisParam) XXX_Size() int {
	return xxx_messageInfo_RedisParam.Size(m)
}
func (m *RedisParam) XXX_DiscardUnknown() {
	xxx_messageInfo_RedisParam.DiscardUnknown(m)
}

var xxx_messageInfo_RedisParam proto.InternalMessageInfo

func (m *RedisParam) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *RedisParam) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type PushMsgReply struct {
	MsgId                string   `protobuf:"bytes,1,opt,name=msgId,proto3" json:"msgId,omitempty"`
	ServerId             string   `protobuf:"bytes,2,opt,name=serverId,proto3" json:"serverId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushMsgReply) Reset()         { *m = PushMsgReply{} }
func (m *PushMsgReply) String() string { return proto.CompactTextString(m) }
func (*PushMsgReply) ProtoMessage()    {}
func (*PushMsgReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_332787bef78344c2, []int{4}
}

func (m *PushMsgReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushMsgReply.Unmarshal(m, b)
}
func (m *PushMsgReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushMsgReply.Marshal(b, m, deterministic)
}
func (m *PushMsgReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushMsgReply.Merge(m, src)
}
func (m *PushMsgReply) XXX_Size() int {
	return xxx_messageInfo_PushMsgReply.Size(m)
}
func (m *PushMsgReply) XXX_DiscardUnknown() {
	xxx_messageInfo_PushMsgReply.DiscardUnknown(m)
}

var xxx_messageInfo_PushMsgReply proto.InternalMessageInfo

func (m *PushMsgReply) GetMsgId() string {
	if m != nil {
		return m.MsgId
	}
	return ""
}

func (m *PushMsgReply) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

// The request message containing the user's name.
type HelloRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloRequest) Reset()         { *m = HelloRequest{} }
func (m *HelloRequest) String() string { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()    {}
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_332787bef78344c2, []int{5}
}

func (m *HelloRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloRequest.Unmarshal(m, b)
}
func (m *HelloRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloRequest.Marshal(b, m, deterministic)
}
func (m *HelloRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloRequest.Merge(m, src)
}
func (m *HelloRequest) XXX_Size() int {
	return xxx_messageInfo_HelloRequest.Size(m)
}
func (m *HelloRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HelloRequest proto.InternalMessageInfo

func (m *HelloRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// The response message containing the greetings
type HelloReply struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloReply) Reset()         { *m = HelloReply{} }
func (m *HelloReply) String() string { return proto.CompactTextString(m) }
func (*HelloReply) ProtoMessage()    {}
func (*HelloReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_332787bef78344c2, []int{6}
}

func (m *HelloReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloReply.Unmarshal(m, b)
}
func (m *HelloReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloReply.Marshal(b, m, deterministic)
}
func (m *HelloReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloReply.Merge(m, src)
}
func (m *HelloReply) XXX_Size() int {
	return xxx_messageInfo_HelloReply.Size(m)
}
func (m *HelloReply) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloReply.DiscardUnknown(m)
}

var xxx_messageInfo_HelloReply proto.InternalMessageInfo

func (m *HelloReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type RoomRequest struct {
	RoomId               string   `protobuf:"bytes,2,opt,name=roomId,proto3" json:"roomId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomRequest) Reset()         { *m = RoomRequest{} }
func (m *RoomRequest) String() string { return proto.CompactTextString(m) }
func (*RoomRequest) ProtoMessage()    {}
func (*RoomRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_332787bef78344c2, []int{7}
}

func (m *RoomRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomRequest.Unmarshal(m, b)
}
func (m *RoomRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomRequest.Marshal(b, m, deterministic)
}
func (m *RoomRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomRequest.Merge(m, src)
}
func (m *RoomRequest) XXX_Size() int {
	return xxx_messageInfo_RoomRequest.Size(m)
}
func (m *RoomRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RoomRequest proto.InternalMessageInfo

func (m *RoomRequest) GetRoomId() string {
	if m != nil {
		return m.RoomId
	}
	return ""
}

type RoomReply struct {
	Count                string   `protobuf:"bytes,2,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomReply) Reset()         { *m = RoomReply{} }
func (m *RoomReply) String() string { return proto.CompactTextString(m) }
func (*RoomReply) ProtoMessage()    {}
func (*RoomReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_332787bef78344c2, []int{8}
}

func (m *RoomReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomReply.Unmarshal(m, b)
}
func (m *RoomReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomReply.Marshal(b, m, deterministic)
}
func (m *RoomReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomReply.Merge(m, src)
}
func (m *RoomReply) XXX_Size() int {
	return xxx_messageInfo_RoomReply.Size(m)
}
func (m *RoomReply) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomReply.DiscardUnknown(m)
}

var xxx_messageInfo_RoomReply proto.InternalMessageInfo

func (m *RoomReply) GetCount() string {
	if m != nil {
		return m.Count
	}
	return ""
}

func init() {
	proto.RegisterType((*PingReq)(nil), "chat.PingReq")
	proto.RegisterType((*PingReply)(nil), "chat.PingReply")
	proto.RegisterType((*PushMsgReq)(nil), "chat.PushMsgReq")
	proto.RegisterType((*RedisParam)(nil), "chat.RedisParam")
	proto.RegisterType((*PushMsgReply)(nil), "chat.PushMsgReply")
	proto.RegisterType((*HelloRequest)(nil), "chat.HelloRequest")
	proto.RegisterType((*HelloReply)(nil), "chat.HelloReply")
	proto.RegisterType((*RoomRequest)(nil), "chat.RoomRequest")
	proto.RegisterType((*RoomReply)(nil), "chat.RoomReply")
}

func init() { proto.RegisterFile("chat/grpc/chat.proto", fileDescriptor_332787bef78344c2) }

var fileDescriptor_332787bef78344c2 = []byte{
	// 405 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0x4d, 0x8f, 0xd3, 0x30,
	0x10, 0x25, 0x90, 0xe6, 0x63, 0xba, 0x88, 0x65, 0xb4, 0x5a, 0x59, 0x39, 0x2d, 0x96, 0xa8, 0x7a,
	0x6a, 0x51, 0xb9, 0x72, 0x40, 0x70, 0x80, 0x1e, 0x90, 0xaa, 0x00, 0x17, 0x6e, 0x26, 0xb1, 0xd2,
	0x4a, 0x49, 0xec, 0xda, 0x0e, 0x28, 0x3f, 0x10, 0x7e, 0x17, 0xf2, 0x47, 0xd2, 0x82, 0xc4, 0x6d,
	0xde, 0x9b, 0x19, 0xcf, 0x9b, 0xf1, 0x83, 0xbb, 0xea, 0xc8, 0xcc, 0xb6, 0x51, 0xb2, 0xda, 0xda,
	0x68, 0x23, 0x95, 0x30, 0x02, 0x63, 0x1b, 0xd3, 0x1c, 0xd2, 0xc3, 0xa9, 0x6f, 0x4a, 0x7e, 0xa6,
	0x4b, 0xc8, 0x7d, 0x28, 0xdb, 0x91, 0xfe, 0x8e, 0x00, 0x0e, 0x83, 0x3e, 0x7e, 0xd2, 0x36, 0x87,
	0x77, 0xb0, 0xe8, 0x74, 0xb3, 0xaf, 0x49, 0xf4, 0x10, 0xad, 0xf3, 0xd2, 0x03, 0x24, 0x90, 0x56,
	0xa2, 0x37, 0xbc, 0x37, 0xe4, 0xb1, 0xe3, 0x27, 0x88, 0xf7, 0x90, 0x0c, 0x9a, 0xab, 0x7d, 0x4d,
	0x9e, 0xb8, 0x44, 0x40, 0xb6, 0xa3, 0x51, 0x62, 0x90, 0xfb, 0x9a, 0xc4, 0xbe, 0x23, 0x40, 0x2c,
	0x20, 0x33, 0xe2, 0xab, 0xef, 0x59, 0xb8, 0xd4, 0x8c, 0x11, 0x21, 0x36, 0xa3, 0xe4, 0x24, 0x71,
	0xbc, 0x8b, 0x71, 0x05, 0x0b, 0xc5, 0xeb, 0x93, 0x26, 0xe9, 0x43, 0xb4, 0x5e, 0xee, 0x6e, 0x37,
	0x6e, 0xb5, 0xd2, 0x52, 0x07, 0xa6, 0x58, 0x57, 0xfa, 0x34, 0x7d, 0x03, 0x70, 0x21, 0xed, 0x4b,
	0xac, 0xae, 0x55, 0x58, 0xc3, 0xc5, 0x76, 0xb2, 0x64, 0x5a, 0xff, 0x14, 0xaa, 0x0e, 0x6b, 0xcc,
	0x98, 0xbe, 0x85, 0x9b, 0xf9, 0x0a, 0xb2, 0x1d, 0xff, 0x73, 0x87, 0x02, 0x32, 0xcd, 0xd5, 0x0f,
	0xa7, 0x3d, 0xbc, 0x30, 0x61, 0x4a, 0xe1, 0xe6, 0x23, 0x6f, 0x5b, 0x51, 0xf2, 0xf3, 0xc0, 0xb5,
	0xb1, 0x0a, 0x7a, 0xd6, 0xf1, 0x49, 0x81, 0x8d, 0xe9, 0x0a, 0x20, 0xd4, 0xd8, 0x19, 0x04, 0xd2,
	0x8e, 0x6b, 0xcd, 0x9a, 0xa9, 0x68, 0x82, 0xf4, 0x25, 0x2c, 0x4b, 0x21, 0xba, 0xe9, 0xa9, 0x7b,
	0x48, 0x94, 0x10, 0xdd, 0x3c, 0x34, 0x20, 0xfa, 0x02, 0x72, 0x5f, 0x16, 0x14, 0x57, 0x62, 0x98,
	0x7f, 0xc8, 0x83, 0xdd, 0xaf, 0x08, 0xe2, 0xf7, 0x47, 0x66, 0x70, 0x07, 0xd9, 0x67, 0x36, 0xba,
	0xe9, 0x88, 0xfe, 0x86, 0xd7, 0x72, 0x8b, 0xdb, 0xbf, 0x38, 0xeb, 0x8c, 0x47, 0xf8, 0x0a, 0xb2,
	0x0f, 0xdc, 0x7c, 0x11, 0x86, 0xb5, 0xf8, 0x3c, 0xdc, 0xfd, 0x22, 0xab, 0x78, 0x76, 0x4d, 0xf9,
	0x8e, 0x15, 0xc4, 0xd6, 0x5a, 0xf8, 0xd4, 0xa7, 0x82, 0xe3, 0xa6, 0xca, 0xd9, 0x75, 0xb8, 0x85,
	0x34, 0x9c, 0x1b, 0xc3, 0xe0, 0x8b, 0x07, 0x0b, 0xfc, 0x87, 0x91, 0xed, 0xf8, 0x2e, 0xf9, 0x16,
	0x5b, 0x5f, 0x7f, 0x4f, 0x9c, 0xa7, 0x5f, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0xeb, 0xbb, 0xc7,
	0x89, 0xeb, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChatClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
	GetTotal(ctx context.Context, in *RoomRequest, opts ...grpc.CallOption) (*RoomReply, error)
	Ping(ctx context.Context, in *PingReq, opts ...grpc.CallOption) (*PingReply, error)
	PushMsg(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error)
}

type chatClient struct {
	cc *grpc.ClientConn
}

func NewChatClient(cc *grpc.ClientConn) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/chat.Chat/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) GetTotal(ctx context.Context, in *RoomRequest, opts ...grpc.CallOption) (*RoomReply, error) {
	out := new(RoomReply)
	err := c.cc.Invoke(ctx, "/chat.Chat/GetTotal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) Ping(ctx context.Context, in *PingReq, opts ...grpc.CallOption) (*PingReply, error) {
	out := new(PingReply)
	err := c.cc.Invoke(ctx, "/chat.Chat/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) PushMsg(ctx context.Context, in *PushMsgReq, opts ...grpc.CallOption) (*PushMsgReply, error) {
	out := new(PushMsgReply)
	err := c.cc.Invoke(ctx, "/chat.Chat/PushMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
type ChatServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	GetTotal(context.Context, *RoomRequest) (*RoomReply, error)
	Ping(context.Context, *PingReq) (*PingReply, error)
	PushMsg(context.Context, *PushMsgReq) (*PushMsgReply, error)
}

// UnimplementedChatServer can be embedded to have forward compatible implementations.
type UnimplementedChatServer struct {
}

func (*UnimplementedChatServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (*UnimplementedChatServer) GetTotal(ctx context.Context, req *RoomRequest) (*RoomReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTotal not implemented")
}
func (*UnimplementedChatServer) Ping(ctx context.Context, req *PingReq) (*PingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedChatServer) PushMsg(ctx context.Context, req *PushMsgReq) (*PushMsgReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMsg not implemented")
}

func RegisterChatServer(s *grpc.Server, srv ChatServer) {
	s.RegisterService(&_Chat_serviceDesc, srv)
}

func _Chat_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Chat/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_GetTotal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).GetTotal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Chat/GetTotal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).GetTotal(ctx, req.(*RoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Chat/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).Ping(ctx, req.(*PingReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_PushMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).PushMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Chat/PushMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).PushMsg(ctx, req.(*PushMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Chat_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chat.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Chat_SayHello_Handler,
		},
		{
			MethodName: "GetTotal",
			Handler:    _Chat_GetTotal_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Chat_Ping_Handler,
		},
		{
			MethodName: "PushMsg",
			Handler:    _Chat_PushMsg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat/grpc/chat.proto",
}