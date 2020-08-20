// Code generated by protoc-gen-go. DO NOT EDIT.
// source: engine/driver/bigquery/proto/bigquery.proto

package proto

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

type Cursor struct {
	Dataset              string   `protobuf:"bytes,1,opt,name=dataset,proto3" json:"dataset,omitempty"`
	Table                string   `protobuf:"bytes,2,opt,name=table,proto3" json:"table,omitempty"`
	Token                string   `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	AvroSchema           string   `protobuf:"bytes,4,opt,name=avro_schema,json=avroSchema,proto3" json:"avro_schema,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Cursor) Reset()         { *m = Cursor{} }
func (m *Cursor) String() string { return proto.CompactTextString(m) }
func (*Cursor) ProtoMessage()    {}
func (*Cursor) Descriptor() ([]byte, []int) {
	return fileDescriptor_7a9e211ec3e97f0b, []int{0}
}

func (m *Cursor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Cursor.Unmarshal(m, b)
}
func (m *Cursor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Cursor.Marshal(b, m, deterministic)
}
func (m *Cursor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Cursor.Merge(m, src)
}
func (m *Cursor) XXX_Size() int {
	return xxx_messageInfo_Cursor.Size(m)
}
func (m *Cursor) XXX_DiscardUnknown() {
	xxx_messageInfo_Cursor.DiscardUnknown(m)
}

var xxx_messageInfo_Cursor proto.InternalMessageInfo

func (m *Cursor) GetDataset() string {
	if m != nil {
		return m.Dataset
	}
	return ""
}

func (m *Cursor) GetTable() string {
	if m != nil {
		return m.Table
	}
	return ""
}

func (m *Cursor) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *Cursor) GetAvroSchema() string {
	if m != nil {
		return m.AvroSchema
	}
	return ""
}

func init() {
	proto.RegisterType((*Cursor)(nil), "bigquery.Cursor")
}

func init() {
	proto.RegisterFile("engine/driver/bigquery/proto/bigquery.proto", fileDescriptor_7a9e211ec3e97f0b)
}

var fileDescriptor_7a9e211ec3e97f0b = []byte{
	// 154 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4e, 0xcd, 0x4b, 0xcf,
	0xcc, 0x4b, 0xd5, 0x4f, 0x29, 0xca, 0x2c, 0x4b, 0x2d, 0xd2, 0x4f, 0xca, 0x4c, 0x2f, 0x2c, 0x4d,
	0x2d, 0xaa, 0xd4, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x87, 0x73, 0xf5, 0xc0, 0x5c, 0x21, 0x0e, 0x18,
	0x5f, 0x29, 0x9f, 0x8b, 0xcd, 0xb9, 0xb4, 0xa8, 0x38, 0xbf, 0x48, 0x48, 0x82, 0x8b, 0x3d, 0x25,
	0xb1, 0x24, 0xb1, 0x38, 0xb5, 0x44, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xc6, 0x15, 0x12,
	0xe1, 0x62, 0x2d, 0x49, 0x4c, 0xca, 0x49, 0x95, 0x60, 0x02, 0x8b, 0x43, 0x38, 0x60, 0xd1, 0xfc,
	0xec, 0xd4, 0x3c, 0x09, 0x66, 0xa8, 0x28, 0x88, 0x23, 0x24, 0xcf, 0xc5, 0x9d, 0x58, 0x56, 0x94,
	0x1f, 0x5f, 0x9c, 0x9c, 0x91, 0x9a, 0x9b, 0x28, 0xc1, 0x02, 0x96, 0xe3, 0x02, 0x09, 0x05, 0x83,
	0x45, 0x9c, 0xd8, 0xa3, 0x58, 0xc1, 0x6e, 0x48, 0x62, 0x03, 0x53, 0xc6, 0x80, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xdf, 0xee, 0xe7, 0xbe, 0xb9, 0x00, 0x00, 0x00,
}