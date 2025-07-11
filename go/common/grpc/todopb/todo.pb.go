// Code generated by protoc-gen-go. DO NOT EDIT.
// source: todo.proto

package todopb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type User struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_0e4b95d0c4e09639, []int{0}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Todo struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Text                 string   `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	Done                 bool     `protobuf:"varint,3,opt,name=done,proto3" json:"done,omitempty"`
	User                 *User    `protobuf:"bytes,4,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Todo) Reset()         { *m = Todo{} }
func (m *Todo) String() string { return proto.CompactTextString(m) }
func (*Todo) ProtoMessage()    {}
func (*Todo) Descriptor() ([]byte, []int) {
	return fileDescriptor_0e4b95d0c4e09639, []int{1}
}

func (m *Todo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Todo.Unmarshal(m, b)
}
func (m *Todo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Todo.Marshal(b, m, deterministic)
}
func (m *Todo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Todo.Merge(m, src)
}
func (m *Todo) XXX_Size() int {
	return xxx_messageInfo_Todo.Size(m)
}
func (m *Todo) XXX_DiscardUnknown() {
	xxx_messageInfo_Todo.DiscardUnknown(m)
}

var xxx_messageInfo_Todo proto.InternalMessageInfo

func (m *Todo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Todo) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *Todo) GetDone() bool {
	if m != nil {
		return m.Done
	}
	return false
}

func (m *Todo) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type CreateTodoRequest struct {
	Todo                 *Todo    `protobuf:"bytes,1,opt,name=todo,proto3" json:"todo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateTodoRequest) Reset()         { *m = CreateTodoRequest{} }
func (m *CreateTodoRequest) String() string { return proto.CompactTextString(m) }
func (*CreateTodoRequest) ProtoMessage()    {}
func (*CreateTodoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0e4b95d0c4e09639, []int{2}
}

func (m *CreateTodoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateTodoRequest.Unmarshal(m, b)
}
func (m *CreateTodoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateTodoRequest.Marshal(b, m, deterministic)
}
func (m *CreateTodoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateTodoRequest.Merge(m, src)
}
func (m *CreateTodoRequest) XXX_Size() int {
	return xxx_messageInfo_CreateTodoRequest.Size(m)
}
func (m *CreateTodoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateTodoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateTodoRequest proto.InternalMessageInfo

func (m *CreateTodoRequest) GetTodo() *Todo {
	if m != nil {
		return m.Todo
	}
	return nil
}

type CreateTodoResponse struct {
	Response             string   `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateTodoResponse) Reset()         { *m = CreateTodoResponse{} }
func (m *CreateTodoResponse) String() string { return proto.CompactTextString(m) }
func (*CreateTodoResponse) ProtoMessage()    {}
func (*CreateTodoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0e4b95d0c4e09639, []int{3}
}

func (m *CreateTodoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateTodoResponse.Unmarshal(m, b)
}
func (m *CreateTodoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateTodoResponse.Marshal(b, m, deterministic)
}
func (m *CreateTodoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateTodoResponse.Merge(m, src)
}
func (m *CreateTodoResponse) XXX_Size() int {
	return xxx_messageInfo_CreateTodoResponse.Size(m)
}
func (m *CreateTodoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateTodoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateTodoResponse proto.InternalMessageInfo

func (m *CreateTodoResponse) GetResponse() string {
	if m != nil {
		return m.Response
	}
	return ""
}

func init() {
	proto.RegisterType((*User)(nil), "todopb.User")
	proto.RegisterType((*Todo)(nil), "todopb.Todo")
	proto.RegisterType((*CreateTodoRequest)(nil), "todopb.CreateTodoRequest")
	proto.RegisterType((*CreateTodoResponse)(nil), "todopb.CreateTodoResponse")
}

func init() {
	proto.RegisterFile("todo.proto", fileDescriptor_0e4b95d0c4e09639)
}

var fileDescriptor_0e4b95d0c4e09639 = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0x95, 0x60, 0x95, 0xf6, 0x8a, 0x90, 0xb8, 0xc9, 0x64, 0xb2, 0x3c, 0x45, 0x0c, 0x01,
	0x15, 0xf1, 0x07, 0xe8, 0x3f, 0x30, 0xb0, 0x30, 0xd1, 0xe2, 0x1b, 0x32, 0x90, 0x0b, 0xb6, 0x8b,
	0xf8, 0xf9, 0xe8, 0x2e, 0x2d, 0x54, 0xa2, 0xdb, 0xf3, 0xbd, 0xef, 0xee, 0x3d, 0x19, 0xa0, 0x70,
	0xe4, 0x6e, 0x4c, 0x5c, 0x18, 0x67, 0xa2, 0xc7, 0xad, 0xbf, 0x01, 0xf3, 0x92, 0x29, 0xe1, 0x25,
	0xd4, 0x7d, 0xb4, 0x95, 0xab, 0xda, 0x45, 0xa8, 0xfb, 0x88, 0x08, 0x66, 0xd8, 0x7c, 0x90, 0xad,
	0x75, 0xa2, 0xda, 0xbf, 0x81, 0x79, 0xe6, 0xc8, 0xa7, 0xd8, 0x42, 0xdf, 0xe5, 0xc0, 0x8a, 0x96,
	0x59, 0xe4, 0x81, 0xec, 0x99, 0xab, 0xda, 0x79, 0x50, 0x8d, 0x0e, 0xcc, 0x2e, 0x53, 0xb2, 0xc6,
	0x55, 0xed, 0x72, 0x75, 0xd1, 0x4d, 0x15, 0x3a, 0xc9, 0x0f, 0xea, 0xf8, 0x07, 0xb8, 0x5a, 0x27,
	0xda, 0x14, 0x92, 0x9c, 0x40, 0x9f, 0x3b, 0xca, 0x45, 0xd6, 0x84, 0xd4, 0xc0, 0xa3, 0x35, 0x45,
	0xd4, 0xf1, 0x77, 0x80, 0xc7, 0x6b, 0x79, 0xe4, 0x21, 0x13, 0x36, 0x30, 0x4f, 0x7b, 0xbd, 0x2f,
	0xfb, 0xfb, 0x5e, 0x05, 0x58, 0x0a, 0xfb, 0x44, 0xe9, 0xab, 0x7f, 0x27, 0x5c, 0x03, 0xfc, 0x1d,
	0xc0, 0xeb, 0x43, 0xc4, 0xbf, 0x2e, 0x4d, 0x73, 0xca, 0x9a, 0x6e, 0x3e, 0x2e, 0x5e, 0xcf, 0x6f,
	0x27, 0x77, 0x3b, 0xd3, 0x4f, 0xbe, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0x76, 0x58, 0x54, 0x5b,
	0x72, 0x01, 0x00, 0x00,
}
