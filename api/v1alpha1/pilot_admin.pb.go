// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pilot_admin.proto

package v1alpha1

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// pilot管理接口配置
type PilotAdminSpec struct {
	// 通过pilot_selector选择目标pilot实例
	PilotSelector        map[string]string           `protobuf:"bytes,1,rep,name=pilot_selector,json=pilotSelector,proto3" json:"pilot_selector,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Limit                *PilotAdminSpec_Limit       `protobuf:"bytes,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Loadbalance          *PilotAdminSpec_LoadBalance `protobuf:"bytes,3,opt,name=loadbalance,proto3" json:"loadbalance,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *PilotAdminSpec) Reset()         { *m = PilotAdminSpec{} }
func (m *PilotAdminSpec) String() string { return proto.CompactTextString(m) }
func (*PilotAdminSpec) ProtoMessage()    {}
func (*PilotAdminSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_36e83212b50b1767, []int{0}
}
func (m *PilotAdminSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PilotAdminSpec.Unmarshal(m, b)
}
func (m *PilotAdminSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PilotAdminSpec.Marshal(b, m, deterministic)
}
func (m *PilotAdminSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PilotAdminSpec.Merge(m, src)
}
func (m *PilotAdminSpec) XXX_Size() int {
	return xxx_messageInfo_PilotAdminSpec.Size(m)
}
func (m *PilotAdminSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_PilotAdminSpec.DiscardUnknown(m)
}

var xxx_messageInfo_PilotAdminSpec proto.InternalMessageInfo

func (m *PilotAdminSpec) GetPilotSelector() map[string]string {
	if m != nil {
		return m.PilotSelector
	}
	return nil
}

func (m *PilotAdminSpec) GetLimit() *PilotAdminSpec_Limit {
	if m != nil {
		return m.Limit
	}
	return nil
}

func (m *PilotAdminSpec) GetLoadbalance() *PilotAdminSpec_LoadBalance {
	if m != nil {
		return m.Loadbalance
	}
	return nil
}

type PilotAdminSpec_Limit struct {
	Qps                  string   `protobuf:"bytes,1,opt,name=qps,proto3" json:"qps,omitempty"`
	Connections          string   `protobuf:"bytes,2,opt,name=connections,proto3" json:"connections,omitempty"`
	Clusters             string   `protobuf:"bytes,3,opt,name=clusters,proto3" json:"clusters,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PilotAdminSpec_Limit) Reset()         { *m = PilotAdminSpec_Limit{} }
func (m *PilotAdminSpec_Limit) String() string { return proto.CompactTextString(m) }
func (*PilotAdminSpec_Limit) ProtoMessage()    {}
func (*PilotAdminSpec_Limit) Descriptor() ([]byte, []int) {
	return fileDescriptor_36e83212b50b1767, []int{0, 1}
}
func (m *PilotAdminSpec_Limit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PilotAdminSpec_Limit.Unmarshal(m, b)
}
func (m *PilotAdminSpec_Limit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PilotAdminSpec_Limit.Marshal(b, m, deterministic)
}
func (m *PilotAdminSpec_Limit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PilotAdminSpec_Limit.Merge(m, src)
}
func (m *PilotAdminSpec_Limit) XXX_Size() int {
	return xxx_messageInfo_PilotAdminSpec_Limit.Size(m)
}
func (m *PilotAdminSpec_Limit) XXX_DiscardUnknown() {
	xxx_messageInfo_PilotAdminSpec_Limit.DiscardUnknown(m)
}

var xxx_messageInfo_PilotAdminSpec_Limit proto.InternalMessageInfo

func (m *PilotAdminSpec_Limit) GetQps() string {
	if m != nil {
		return m.Qps
	}
	return ""
}

func (m *PilotAdminSpec_Limit) GetConnections() string {
	if m != nil {
		return m.Connections
	}
	return ""
}

func (m *PilotAdminSpec_Limit) GetClusters() string {
	if m != nil {
		return m.Clusters
	}
	return ""
}

type PilotAdminSpec_LoadBalance struct {
	Weight               float32  `protobuf:"fixed32,1,opt,name=weight,proto3" json:"weight,omitempty"`
	Diff                 int32    `protobuf:"varint,2,opt,name=diff,proto3" json:"diff,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PilotAdminSpec_LoadBalance) Reset()         { *m = PilotAdminSpec_LoadBalance{} }
func (m *PilotAdminSpec_LoadBalance) String() string { return proto.CompactTextString(m) }
func (*PilotAdminSpec_LoadBalance) ProtoMessage()    {}
func (*PilotAdminSpec_LoadBalance) Descriptor() ([]byte, []int) {
	return fileDescriptor_36e83212b50b1767, []int{0, 2}
}
func (m *PilotAdminSpec_LoadBalance) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PilotAdminSpec_LoadBalance.Unmarshal(m, b)
}
func (m *PilotAdminSpec_LoadBalance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PilotAdminSpec_LoadBalance.Marshal(b, m, deterministic)
}
func (m *PilotAdminSpec_LoadBalance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PilotAdminSpec_LoadBalance.Merge(m, src)
}
func (m *PilotAdminSpec_LoadBalance) XXX_Size() int {
	return xxx_messageInfo_PilotAdminSpec_LoadBalance.Size(m)
}
func (m *PilotAdminSpec_LoadBalance) XXX_DiscardUnknown() {
	xxx_messageInfo_PilotAdminSpec_LoadBalance.DiscardUnknown(m)
}

var xxx_messageInfo_PilotAdminSpec_LoadBalance proto.InternalMessageInfo

func (m *PilotAdminSpec_LoadBalance) GetWeight() float32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *PilotAdminSpec_LoadBalance) GetDiff() int32 {
	if m != nil {
		return m.Diff
	}
	return 0
}

// 记录pilot的负载状态
type PilotAdminStatus struct {
	// pilot副本数量
	Replicas int64 `protobuf:"varint,1,opt,name=replicas,proto3" json:"replicas,omitempty"`
	// value： {pod_name} key: ep状态
	Endpoints            map[string]*PilotAdminStatus_EndpointStatus `protobuf:"bytes,2,rep,name=endpoints,proto3" json:"endpoints,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *PilotAdminStatus) Reset()         { *m = PilotAdminStatus{} }
func (m *PilotAdminStatus) String() string { return proto.CompactTextString(m) }
func (*PilotAdminStatus) ProtoMessage()    {}
func (*PilotAdminStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_36e83212b50b1767, []int{1}
}
func (m *PilotAdminStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PilotAdminStatus.Unmarshal(m, b)
}
func (m *PilotAdminStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PilotAdminStatus.Marshal(b, m, deterministic)
}
func (m *PilotAdminStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PilotAdminStatus.Merge(m, src)
}
func (m *PilotAdminStatus) XXX_Size() int {
	return xxx_messageInfo_PilotAdminStatus.Size(m)
}
func (m *PilotAdminStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_PilotAdminStatus.DiscardUnknown(m)
}

var xxx_messageInfo_PilotAdminStatus proto.InternalMessageInfo

func (m *PilotAdminStatus) GetReplicas() int64 {
	if m != nil {
		return m.Replicas
	}
	return 0
}

func (m *PilotAdminStatus) GetEndpoints() map[string]*PilotAdminStatus_EndpointStatus {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

// pilot单个endpoint的状态记录
type PilotAdminStatus_EndpointStatus struct {
	Status               map[string]string `protobuf:"bytes,1,rep,name=status,proto3" json:"status,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *PilotAdminStatus_EndpointStatus) Reset()         { *m = PilotAdminStatus_EndpointStatus{} }
func (m *PilotAdminStatus_EndpointStatus) String() string { return proto.CompactTextString(m) }
func (*PilotAdminStatus_EndpointStatus) ProtoMessage()    {}
func (*PilotAdminStatus_EndpointStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_36e83212b50b1767, []int{1, 0}
}
func (m *PilotAdminStatus_EndpointStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PilotAdminStatus_EndpointStatus.Unmarshal(m, b)
}
func (m *PilotAdminStatus_EndpointStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PilotAdminStatus_EndpointStatus.Marshal(b, m, deterministic)
}
func (m *PilotAdminStatus_EndpointStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PilotAdminStatus_EndpointStatus.Merge(m, src)
}
func (m *PilotAdminStatus_EndpointStatus) XXX_Size() int {
	return xxx_messageInfo_PilotAdminStatus_EndpointStatus.Size(m)
}
func (m *PilotAdminStatus_EndpointStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_PilotAdminStatus_EndpointStatus.DiscardUnknown(m)
}

var xxx_messageInfo_PilotAdminStatus_EndpointStatus proto.InternalMessageInfo

func (m *PilotAdminStatus_EndpointStatus) GetStatus() map[string]string {
	if m != nil {
		return m.Status
	}
	return nil
}

func init() {
	proto.RegisterType((*PilotAdminSpec)(nil), "slime.microservice.pilotadmin.v1alpha1.PilotAdminSpec")
	proto.RegisterMapType((map[string]string)(nil), "slime.microservice.pilotadmin.v1alpha1.PilotAdminSpec.PilotSelectorEntry")
	proto.RegisterType((*PilotAdminSpec_Limit)(nil), "slime.microservice.pilotadmin.v1alpha1.PilotAdminSpec.Limit")
	proto.RegisterType((*PilotAdminSpec_LoadBalance)(nil), "slime.microservice.pilotadmin.v1alpha1.PilotAdminSpec.LoadBalance")
	proto.RegisterType((*PilotAdminStatus)(nil), "slime.microservice.pilotadmin.v1alpha1.PilotAdminStatus")
	proto.RegisterMapType((map[string]*PilotAdminStatus_EndpointStatus)(nil), "slime.microservice.pilotadmin.v1alpha1.PilotAdminStatus.EndpointsEntry")
	proto.RegisterType((*PilotAdminStatus_EndpointStatus)(nil), "slime.microservice.pilotadmin.v1alpha1.PilotAdminStatus.EndpointStatus")
	proto.RegisterMapType((map[string]string)(nil), "slime.microservice.pilotadmin.v1alpha1.PilotAdminStatus.EndpointStatus.StatusEntry")
}

func init() { proto.RegisterFile("pilot_admin.proto", fileDescriptor_36e83212b50b1767) }

var fileDescriptor_36e83212b50b1767 = []byte{
	// 447 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xdf, 0x8a, 0x13, 0x31,
	0x14, 0xc6, 0x99, 0xce, 0xb6, 0xd8, 0x33, 0x58, 0xd6, 0x20, 0x52, 0xe6, 0xaa, 0xec, 0x85, 0xf4,
	0x6a, 0x86, 0x5d, 0x11, 0xff, 0xe0, 0x85, 0x16, 0x16, 0x11, 0xbc, 0x90, 0xf4, 0x42, 0x10, 0x44,
	0xb2, 0x99, 0xd3, 0xdd, 0xb0, 0x99, 0x24, 0x3b, 0x49, 0x2b, 0xfb, 0x00, 0xbe, 0x80, 0x4f, 0xe3,
	0xa3, 0xf8, 0x38, 0x92, 0x64, 0x66, 0x37, 0x22, 0x82, 0x56, 0xaf, 0x7a, 0xbe, 0xd3, 0x9c, 0xdf,
	0xf7, 0xe5, 0x34, 0x14, 0xee, 0x19, 0x21, 0xb5, 0xfb, 0xc4, 0x9a, 0x56, 0xa8, 0xca, 0x74, 0xda,
	0x69, 0xf2, 0xd0, 0x4a, 0xd1, 0x62, 0xd5, 0x0a, 0xde, 0x69, 0x8b, 0xdd, 0x4e, 0x70, 0xac, 0xc2,
	0xa9, 0x78, 0x68, 0x77, 0xcc, 0xa4, 0xb9, 0x60, 0xc7, 0x47, 0x5f, 0x0f, 0x60, 0xf6, 0xce, 0xf7,
	0x5f, 0xf9, 0xfe, 0xda, 0x20, 0x27, 0x06, 0x66, 0x91, 0x67, 0x51, 0x22, 0x77, 0xba, 0x9b, 0x67,
	0x8b, 0x7c, 0x59, 0x9c, 0xbc, 0xa9, 0xfe, 0x8c, 0x59, 0xfd, 0xcc, 0x8b, 0x72, 0xdd, 0xb3, 0x4e,
	0x95, 0xeb, 0xae, 0xe9, 0x5d, 0x93, 0xf6, 0x08, 0x85, 0xb1, 0x14, 0xad, 0x70, 0xf3, 0xd1, 0x22,
	0x5b, 0x16, 0x27, 0x2f, 0xf6, 0x34, 0x7a, 0xeb, 0x19, 0x34, 0xa2, 0x48, 0x03, 0x85, 0xd4, 0xac,
	0x39, 0x63, 0x92, 0x29, 0x8e, 0xf3, 0x3c, 0x90, 0x57, 0xfb, 0x92, 0x35, 0x6b, 0x56, 0x91, 0x44,
	0x53, 0x6c, 0xf9, 0x12, 0xc8, 0xaf, 0xd7, 0x23, 0x87, 0x90, 0x5f, 0xe2, 0xf5, 0x3c, 0x5b, 0x64,
	0xcb, 0x29, 0xf5, 0x25, 0xb9, 0x0f, 0xe3, 0x1d, 0x93, 0x5b, 0x0c, 0x37, 0x9c, 0xd2, 0x28, 0x9e,
	0x8f, 0x9e, 0x66, 0xe5, 0x7b, 0x18, 0x87, 0xdc, 0x7e, 0xe8, 0xca, 0xd8, 0x61, 0xe8, 0xca, 0x58,
	0xb2, 0x80, 0x82, 0x6b, 0xa5, 0x90, 0x3b, 0xa1, 0x95, 0xed, 0x47, 0xd3, 0x16, 0x29, 0xe1, 0x0e,
	0x97, 0x5b, 0xeb, 0xb0, 0xb3, 0xe1, 0x86, 0x53, 0x7a, 0xa3, 0xcb, 0x67, 0x50, 0x24, 0xb1, 0xc9,
	0x03, 0x98, 0x7c, 0x46, 0x71, 0x7e, 0xe1, 0x82, 0xc3, 0x88, 0xf6, 0x8a, 0x10, 0x38, 0x68, 0xc4,
	0x66, 0x13, 0xe8, 0x63, 0x1a, 0xea, 0xa3, 0xef, 0x39, 0x1c, 0x26, 0x1b, 0x70, 0xcc, 0x6d, 0x83,
	0x57, 0x87, 0x46, 0x0a, 0xce, 0x62, 0xc8, 0x9c, 0xde, 0x68, 0x82, 0x30, 0x45, 0xd5, 0x18, 0x2d,
	0x94, 0xf3, 0x39, 0xfd, 0x6b, 0x79, 0xbd, 0xc7, 0xaa, 0x83, 0x51, 0x75, 0x3a, 0x90, 0xe2, 0x5b,
	0xb9, 0x25, 0x97, 0xdf, 0x32, 0x98, 0x0d, 0xdf, 0xf6, 0xa9, 0x2e, 0x61, 0x62, 0x43, 0xd5, 0x3f,
	0xd2, 0xf5, 0x3f, 0xdb, 0xf6, 0x32, 0x7e, 0xc4, 0x08, 0xbd, 0x85, 0x5f, 0x69, 0xd2, 0xfe, 0xab,
	0x9f, 0xf9, 0x4b, 0x12, 0xfd, 0xb7, 0xe3, 0x1f, 0xd3, 0xf1, 0xff, 0xb1, 0xc2, 0x28, 0x93, 0x1c,
	0xab, 0x27, 0x1f, 0x1e, 0x47, 0xa8, 0xd0, 0x75, 0x28, 0xea, 0x56, 0x37, 0x5b, 0x89, 0xb6, 0xbe,
	0x05, 0xd7, 0xcc, 0x88, 0x9a, 0x6b, 0xb5, 0x11, 0xe7, 0xf5, 0xe0, 0x71, 0x36, 0x09, 0xff, 0x2b,
	0x8f, 0x7e, 0x04, 0x00, 0x00, 0xff, 0xff, 0x77, 0x7b, 0x80, 0x58, 0x6c, 0x04, 0x00, 0x00,
}
