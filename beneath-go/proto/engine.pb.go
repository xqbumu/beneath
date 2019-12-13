// Code generated by protoc-gen-go. DO NOT EDIT.
// source: engine.proto

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

type Record struct {
	AvroData             []byte   `protobuf:"bytes,1,opt,name=avro_data,json=avroData,proto3" json:"avro_data,omitempty"`
	Timestamp            int64    `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Record) Reset()         { *m = Record{} }
func (m *Record) String() string { return proto.CompactTextString(m) }
func (*Record) ProtoMessage()    {}
func (*Record) Descriptor() ([]byte, []int) {
	return fileDescriptor_770b178c3aab763f, []int{0}
}

func (m *Record) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Record.Unmarshal(m, b)
}
func (m *Record) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Record.Marshal(b, m, deterministic)
}
func (m *Record) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Record.Merge(m, src)
}
func (m *Record) XXX_Size() int {
	return xxx_messageInfo_Record.Size(m)
}
func (m *Record) XXX_DiscardUnknown() {
	xxx_messageInfo_Record.DiscardUnknown(m)
}

var xxx_messageInfo_Record proto.InternalMessageInfo

func (m *Record) GetAvroData() []byte {
	if m != nil {
		return m.AvroData
	}
	return nil
}

func (m *Record) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type WriteRecordsReport struct {
	InstanceId           []byte   `protobuf:"bytes,1,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	BytesWritten         int64    `protobuf:"varint,2,opt,name=bytes_written,json=bytesWritten,proto3" json:"bytes_written,omitempty"`
	OffsetStart          int64    `protobuf:"varint,3,opt,name=offset_start,json=offsetStart,proto3" json:"offset_start,omitempty"`
	OffsetEnd            int32    `protobuf:"varint,4,opt,name=offset_end,json=offsetEnd,proto3" json:"offset_end,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WriteRecordsReport) Reset()         { *m = WriteRecordsReport{} }
func (m *WriteRecordsReport) String() string { return proto.CompactTextString(m) }
func (*WriteRecordsReport) ProtoMessage()    {}
func (*WriteRecordsReport) Descriptor() ([]byte, []int) {
	return fileDescriptor_770b178c3aab763f, []int{1}
}

func (m *WriteRecordsReport) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WriteRecordsReport.Unmarshal(m, b)
}
func (m *WriteRecordsReport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WriteRecordsReport.Marshal(b, m, deterministic)
}
func (m *WriteRecordsReport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WriteRecordsReport.Merge(m, src)
}
func (m *WriteRecordsReport) XXX_Size() int {
	return xxx_messageInfo_WriteRecordsReport.Size(m)
}
func (m *WriteRecordsReport) XXX_DiscardUnknown() {
	xxx_messageInfo_WriteRecordsReport.DiscardUnknown(m)
}

var xxx_messageInfo_WriteRecordsReport proto.InternalMessageInfo

func (m *WriteRecordsReport) GetInstanceId() []byte {
	if m != nil {
		return m.InstanceId
	}
	return nil
}

func (m *WriteRecordsReport) GetBytesWritten() int64 {
	if m != nil {
		return m.BytesWritten
	}
	return 0
}

func (m *WriteRecordsReport) GetOffsetStart() int64 {
	if m != nil {
		return m.OffsetStart
	}
	return 0
}

func (m *WriteRecordsReport) GetOffsetEnd() int32 {
	if m != nil {
		return m.OffsetEnd
	}
	return 0
}

func init() {
	proto.RegisterType((*Record)(nil), "proto.Record")
	proto.RegisterType((*WriteRecordsReport)(nil), "proto.WriteRecordsReport")
}

func init() { proto.RegisterFile("engine.proto", fileDescriptor_770b178c3aab763f) }

var fileDescriptor_770b178c3aab763f = []byte{
	// 234 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0xd0, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x06, 0x60, 0xd6, 0xda, 0x62, 0xa6, 0xeb, 0x65, 0x41, 0x08, 0xa8, 0x18, 0xeb, 0xc1, 0x9c,
	0x7a, 0xf1, 0x0d, 0xa2, 0x1e, 0xbc, 0x95, 0xf5, 0xd0, 0x63, 0x98, 0x74, 0xa7, 0xba, 0x48, 0x67,
	0xc3, 0xee, 0x60, 0xf1, 0x5d, 0x7c, 0x58, 0x49, 0x36, 0xd2, 0xd3, 0xf0, 0x7f, 0xfc, 0xfc, 0x87,
	0x01, 0x4d, 0xfc, 0xe1, 0x99, 0xd6, 0x7d, 0x0c, 0x12, 0xcc, 0x7c, 0x3c, 0xab, 0x67, 0x58, 0x58,
	0xda, 0x85, 0xe8, 0xcc, 0x35, 0x14, 0xf8, 0x1d, 0x43, 0xeb, 0x50, 0xb0, 0x54, 0x95, 0xaa, 0xb5,
	0xbd, 0x18, 0xe0, 0x05, 0x05, 0xcd, 0x0d, 0x14, 0xe2, 0x0f, 0x94, 0x04, 0x0f, 0x7d, 0x79, 0x56,
	0xa9, 0x7a, 0x66, 0x4f, 0xb0, 0xfa, 0x55, 0x60, 0xb6, 0xd1, 0x0b, 0xe5, 0xa9, 0x64, 0xa9, 0x0f,
	0x51, 0xcc, 0x1d, 0x2c, 0x3d, 0x27, 0x41, 0xde, 0x51, 0xeb, 0xdd, 0xb4, 0x09, 0xff, 0xf4, 0xe6,
	0xcc, 0x03, 0x5c, 0x76, 0x3f, 0x42, 0xa9, 0x3d, 0x46, 0x2f, 0x42, 0x3c, 0x2d, 0xeb, 0x11, 0xb7,
	0xd9, 0xcc, 0x3d, 0xe8, 0xb0, 0xdf, 0x27, 0x92, 0x36, 0x09, 0x46, 0x29, 0x67, 0x63, 0x67, 0x99,
	0xed, 0x7d, 0x20, 0x73, 0x0b, 0x30, 0x55, 0x88, 0x5d, 0x79, 0x5e, 0xa9, 0x7a, 0x6e, 0x8b, 0x2c,
	0xaf, 0xec, 0x9a, 0x47, 0xb8, 0x62, 0x92, 0x63, 0x88, 0x5f, 0xeb, 0x8e, 0x98, 0x50, 0x3e, 0xf3,
	0x0f, 0x1a, 0xdd, 0xe4, 0xb8, 0x19, 0xd2, 0x46, 0x75, 0x8b, 0x91, 0x9f, 0xfe, 0x02, 0x00, 0x00,
	0xff, 0xff, 0xaf, 0x08, 0xda, 0xfd, 0x2a, 0x01, 0x00, 0x00,
}
