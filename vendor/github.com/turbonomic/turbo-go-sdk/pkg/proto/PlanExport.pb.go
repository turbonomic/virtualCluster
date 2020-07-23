// Code generated by protoc-gen-go. DO NOT EDIT.
// source: PlanExport.proto

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

// Indicates the current export status of an export operation.
type PlanExportResponse_PlanExportResponseState int32

const (
	// No export had been attempted yet.
	PlanExportResponse_NONE PlanExportResponse_PlanExportResponseState = 0
	// Export request was rejected.
	PlanExportResponse_REJECTED PlanExportResponse_PlanExportResponseState = 1
	// Export is in progress.
	PlanExportResponse_IN_PROGRESS PlanExportResponse_PlanExportResponseState = 2
	// Export finished successfully.
	PlanExportResponse_SUCCEEDED PlanExportResponse_PlanExportResponseState = 3
	// Export terminated due to error.
	PlanExportResponse_FAILED PlanExportResponse_PlanExportResponseState = 4
)

var PlanExportResponse_PlanExportResponseState_name = map[int32]string{
	0: "NONE",
	1: "REJECTED",
	2: "IN_PROGRESS",
	3: "SUCCEEDED",
	4: "FAILED",
}

var PlanExportResponse_PlanExportResponseState_value = map[string]int32{
	"NONE":        0,
	"REJECTED":    1,
	"IN_PROGRESS": 2,
	"SUCCEEDED":   3,
	"FAILED":      4,
}

func (x PlanExportResponse_PlanExportResponseState) Enum() *PlanExportResponse_PlanExportResponseState {
	p := new(PlanExportResponse_PlanExportResponseState)
	*p = x
	return p
}

func (x PlanExportResponse_PlanExportResponseState) String() string {
	return proto.EnumName(PlanExportResponse_PlanExportResponseState_name, int32(x))
}

func (x *PlanExportResponse_PlanExportResponseState) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(PlanExportResponse_PlanExportResponseState_value, data, "PlanExportResponse_PlanExportResponseState")
	if err != nil {
		return err
	}
	*x = PlanExportResponse_PlanExportResponseState(value)
	return nil
}

func (PlanExportResponse_PlanExportResponseState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b584bfb073803d20, []int{1, 0}
}

type PlanExportDTO struct {
	// The ID of the market from which the plan export was created.
	MarketId *string `protobuf:"bytes,1,req,name=marketId" json:"marketId,omitempty"`
	// The human readable name for this plan as seen in the UI (eg "Migrate to Public Cloud 8")
	PlanName *string `protobuf:"bytes,2,req,name=planName" json:"planName,omitempty"`
	// Actions generated by the plan
	Actions              []*ActionExecutionDTO `protobuf:"bytes,3,rep,name=actions" json:"actions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *PlanExportDTO) Reset()         { *m = PlanExportDTO{} }
func (m *PlanExportDTO) String() string { return proto.CompactTextString(m) }
func (*PlanExportDTO) ProtoMessage()    {}
func (*PlanExportDTO) Descriptor() ([]byte, []int) {
	return fileDescriptor_b584bfb073803d20, []int{0}
}

func (m *PlanExportDTO) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlanExportDTO.Unmarshal(m, b)
}
func (m *PlanExportDTO) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlanExportDTO.Marshal(b, m, deterministic)
}
func (m *PlanExportDTO) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlanExportDTO.Merge(m, src)
}
func (m *PlanExportDTO) XXX_Size() int {
	return xxx_messageInfo_PlanExportDTO.Size(m)
}
func (m *PlanExportDTO) XXX_DiscardUnknown() {
	xxx_messageInfo_PlanExportDTO.DiscardUnknown(m)
}

var xxx_messageInfo_PlanExportDTO proto.InternalMessageInfo

func (m *PlanExportDTO) GetMarketId() string {
	if m != nil && m.MarketId != nil {
		return *m.MarketId
	}
	return ""
}

func (m *PlanExportDTO) GetPlanName() string {
	if m != nil && m.PlanName != nil {
		return *m.PlanName
	}
	return ""
}

func (m *PlanExportDTO) GetActions() []*ActionExecutionDTO {
	if m != nil {
		return m.Actions
	}
	return nil
}

// This class holds response information about an plan export operation. It contains: state: the
// PlanExportState code representing the state of the export, progress: a percentage complete
// indicator, description: a message notifying detailed information about the current status.
type PlanExportResponse struct {
	// current action state
	State *PlanExportResponse_PlanExportResponseState `protobuf:"varint,1,req,name=state,enum=common_dto.PlanExportResponse_PlanExportResponseState" json:"state,omitempty"`
	// current action progress (0..100)
	Progress *int32 `protobuf:"varint,2,req,name=progress" json:"progress,omitempty"`
	// action state description, for example ("Moving VM...")
	Description          *string  `protobuf:"bytes,3,req,name=description" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlanExportResponse) Reset()         { *m = PlanExportResponse{} }
func (m *PlanExportResponse) String() string { return proto.CompactTextString(m) }
func (*PlanExportResponse) ProtoMessage()    {}
func (*PlanExportResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b584bfb073803d20, []int{1}
}

func (m *PlanExportResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlanExportResponse.Unmarshal(m, b)
}
func (m *PlanExportResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlanExportResponse.Marshal(b, m, deterministic)
}
func (m *PlanExportResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlanExportResponse.Merge(m, src)
}
func (m *PlanExportResponse) XXX_Size() int {
	return xxx_messageInfo_PlanExportResponse.Size(m)
}
func (m *PlanExportResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PlanExportResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PlanExportResponse proto.InternalMessageInfo

func (m *PlanExportResponse) GetState() PlanExportResponse_PlanExportResponseState {
	if m != nil && m.State != nil {
		return *m.State
	}
	return PlanExportResponse_NONE
}

func (m *PlanExportResponse) GetProgress() int32 {
	if m != nil && m.Progress != nil {
		return *m.Progress
	}
	return 0
}

func (m *PlanExportResponse) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func init() {
	proto.RegisterEnum("common_dto.PlanExportResponse_PlanExportResponseState", PlanExportResponse_PlanExportResponseState_name, PlanExportResponse_PlanExportResponseState_value)
	proto.RegisterType((*PlanExportDTO)(nil), "common_dto.PlanExportDTO")
	proto.RegisterType((*PlanExportResponse)(nil), "common_dto.PlanExportResponse")
}

func init() {
	proto.RegisterFile("PlanExport.proto", fileDescriptor_b584bfb073803d20)
}

var fileDescriptor_b584bfb073803d20 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x8f, 0x4f, 0x4f, 0xc2, 0x40,
	0x10, 0xc5, 0x2d, 0x05, 0x85, 0x41, 0xa4, 0xd9, 0xc4, 0xd8, 0x70, 0x50, 0xc2, 0x89, 0x8b, 0x3d,
	0x70, 0x30, 0x5e, 0x81, 0xae, 0x06, 0x43, 0x5a, 0xb2, 0xc5, 0xa3, 0x21, 0xb5, 0x5d, 0x0d, 0x91,
	0xed, 0x34, 0xbb, 0x0b, 0xc1, 0xa3, 0x9f, 0x5c, 0xd3, 0x56, 0x28, 0xf1, 0xcf, 0x6d, 0xdf, 0x9b,
	0xb7, 0x6f, 0x7e, 0x03, 0xd6, 0x6c, 0x15, 0x26, 0x74, 0x9b, 0xa2, 0xd4, 0x4e, 0x2a, 0x51, 0x23,
	0x81, 0x08, 0x85, 0xc0, 0x64, 0x11, 0x6b, 0xec, 0x9c, 0x0f, 0x23, 0xbd, 0xc4, 0x84, 0x6e, 0x79,
	0xb4, 0xce, 0x1e, 0x45, 0xa4, 0xd3, 0x76, 0x97, 0x2a, 0xc2, 0x0d, 0x97, 0xef, 0x85, 0xd1, 0xfb,
	0x30, 0xa0, 0x55, 0x16, 0xb9, 0x73, 0x9f, 0x74, 0xa0, 0x2e, 0x42, 0xf9, 0xc6, 0xf5, 0x24, 0xb6,
	0x8d, 0x6e, 0xa5, 0xdf, 0x60, 0x7b, 0x9d, 0xcd, 0xd2, 0x55, 0x98, 0x78, 0xa1, 0xe0, 0x76, 0xa5,
	0x98, 0xed, 0x34, 0xb9, 0x85, 0x93, 0x30, 0xdf, 0xa9, 0x6c, 0xb3, 0x6b, 0xf6, 0x9b, 0x83, 0x4b,
	0xa7, 0xe4, 0x71, 0x7e, 0xe0, 0xb8, 0x73, 0x9f, 0xed, 0xe2, 0xbd, 0x4f, 0x03, 0x48, 0xc9, 0xc0,
	0xb8, 0x4a, 0x31, 0x51, 0x9c, 0x4c, 0xa1, 0xa6, 0x74, 0xa8, 0x79, 0x4e, 0x71, 0x36, 0xb8, 0x39,
	0xac, 0xfb, 0x1d, 0xff, 0xc3, 0x0a, 0xb2, 0xdf, 0xac, 0x28, 0xc9, 0xd1, 0x25, 0xbe, 0x4a, 0xae,
	0x54, 0x8e, 0x5e, 0x63, 0x7b, 0x4d, 0xba, 0xd0, 0x8c, 0xb9, 0x8a, 0xe4, 0x32, 0xcd, 0x80, 0x6c,
	0x33, 0xbf, 0xec, 0xd0, 0xea, 0x3d, 0xc1, 0xc5, 0x3f, 0xfd, 0xa4, 0x0e, 0x55, 0xcf, 0xf7, 0xa8,
	0x75, 0x44, 0x4e, 0xa1, 0xce, 0xe8, 0x03, 0x1d, 0xcf, 0xa9, 0x6b, 0x19, 0xa4, 0x0d, 0xcd, 0x89,
	0xb7, 0x98, 0x31, 0xff, 0x9e, 0xd1, 0x20, 0xb0, 0x2a, 0xa4, 0x05, 0x8d, 0xe0, 0x71, 0x3c, 0xa6,
	0xd4, 0xa5, 0xae, 0x65, 0x12, 0x80, 0xe3, 0xbb, 0xe1, 0x64, 0x4a, 0x5d, 0xab, 0x3a, 0xba, 0x86,
	0xab, 0x08, 0x85, 0xb3, 0x11, 0x7a, 0x2d, 0x9f, 0xd1, 0x49, 0x57, 0xa1, 0x7e, 0x41, 0x29, 0xbe,
	0x2f, 0x76, 0x62, 0x8d, 0x23, 0x28, 0xf7, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0x24, 0xf7, 0x04,
	0xbd, 0xfb, 0x01, 0x00, 0x00,
}