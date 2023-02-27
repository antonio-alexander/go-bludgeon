//
//go_bludgeon_healthcheck defines a set of types for use with the healthcheck service

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: healthcheck.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// HealthCheckResponse
type HealthCheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// healthcheck
	Healthcheck *HealthCheck `protobuf:"bytes,1,opt,name=healthcheck,proto3" json:"healthcheck,omitempty"`
}

func (x *HealthCheckResponse) Reset() {
	*x = HealthCheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_healthcheck_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthCheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckResponse) ProtoMessage() {}

func (x *HealthCheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_healthcheck_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheckResponse.ProtoReflect.Descriptor instead.
func (*HealthCheckResponse) Descriptor() ([]byte, []int) {
	return file_healthcheck_proto_rawDescGZIP(), []int{0}
}

func (x *HealthCheckResponse) GetHealthcheck() *HealthCheck {
	if x != nil {
		return x.Healthcheck
	}
	return nil
}

// HealthCheck
type HealthCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// time
	Time int64 `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *HealthCheck) Reset() {
	*x = HealthCheck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_healthcheck_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheck) ProtoMessage() {}

func (x *HealthCheck) ProtoReflect() protoreflect.Message {
	mi := &file_healthcheck_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheck.ProtoReflect.Descriptor instead.
func (*HealthCheck) Descriptor() ([]byte, []int) {
	return file_healthcheck_proto_rawDescGZIP(), []int{1}
}

func (x *HealthCheck) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

// Wrapper describes a basic data type for conversion of any
// other data type
type Wrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// type is a string identifying the type of payload
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// payload describes the payload as protobuf.any
	Payload *anypb.Any `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Wrapper) Reset() {
	*x = Wrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_healthcheck_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Wrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Wrapper) ProtoMessage() {}

func (x *Wrapper) ProtoReflect() protoreflect.Message {
	mi := &file_healthcheck_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Wrapper.ProtoReflect.Descriptor instead.
func (*Wrapper) Descriptor() ([]byte, []int) {
	return file_healthcheck_proto_rawDescGZIP(), []int{2}
}

func (x *Wrapper) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Wrapper) GetPayload() *anypb.Any {
	if x != nil {
		return x.Payload
	}
	return nil
}

// Bytes makes it easier to provide
type Bytes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// bytes describes the payload
	Bytes []byte `protobuf:"bytes,1,opt,name=bytes,proto3" json:"bytes,omitempty"`
}

func (x *Bytes) Reset() {
	*x = Bytes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_healthcheck_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bytes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bytes) ProtoMessage() {}

func (x *Bytes) ProtoReflect() protoreflect.Message {
	mi := &file_healthcheck_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bytes.ProtoReflect.Descriptor instead.
func (*Bytes) Descriptor() ([]byte, []int) {
	return file_healthcheck_proto_rawDescGZIP(), []int{3}
}

func (x *Bytes) GetBytes() []byte {
	if x != nil {
		return x.Bytes
	}
	return nil
}

// Error
type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// error
	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_healthcheck_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_healthcheck_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_healthcheck_proto_rawDescGZIP(), []int{4}
}

func (x *Error) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// Empty
type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_healthcheck_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_healthcheck_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_healthcheck_proto_rawDescGZIP(), []int{5}
}

var File_healthcheck_proto protoreflect.FileDescriptor

var file_healthcheck_proto_rawDesc = []byte{
	0x0a, 0x11, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x17, 0x67, 0x6f, 0x5f, 0x62, 0x6c, 0x75, 0x64, 0x67, 0x65, 0x6f, 0x6e,
	0x5f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x1a, 0x19, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5d, 0x0a, 0x13, 0x48, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46,
	0x0a, 0x0b, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x67, 0x6f, 0x5f, 0x62, 0x6c, 0x75, 0x64, 0x67, 0x65, 0x6f,
	0x6e, 0x5f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x0b, 0x68, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x22, 0x21, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x4d, 0x0a, 0x07, 0x57, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x1d, 0x0a, 0x05, 0x42, 0x79, 0x74, 0x65,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x22, 0x1d, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32,
	0x6d, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x12,
	0x5d, 0x0a, 0x0b, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x1e,
	0x2e, 0x67, 0x6f, 0x5f, 0x62, 0x6c, 0x75, 0x64, 0x67, 0x65, 0x6f, 0x6e, 0x5f, 0x68, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2c,
	0x2e, 0x67, 0x6f, 0x5f, 0x62, 0x6c, 0x75, 0x64, 0x67, 0x65, 0x6f, 0x6e, 0x5f, 0x68, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x3e,
	0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6e, 0x74,
	0x6f, 0x6e, 0x69, 0x6f, 0x2d, 0x61, 0x6c, 0x65, 0x78, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x67,
	0x6f, 0x2d, 0x62, 0x6c, 0x75, 0x64, 0x67, 0x65, 0x6f, 0x6e, 0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_healthcheck_proto_rawDescOnce sync.Once
	file_healthcheck_proto_rawDescData = file_healthcheck_proto_rawDesc
)

func file_healthcheck_proto_rawDescGZIP() []byte {
	file_healthcheck_proto_rawDescOnce.Do(func() {
		file_healthcheck_proto_rawDescData = protoimpl.X.CompressGZIP(file_healthcheck_proto_rawDescData)
	})
	return file_healthcheck_proto_rawDescData
}

var file_healthcheck_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_healthcheck_proto_goTypes = []interface{}{
	(*HealthCheckResponse)(nil), // 0: go_bludgeon_healthcheck.HealthCheckResponse
	(*HealthCheck)(nil),         // 1: go_bludgeon_healthcheck.HealthCheck
	(*Wrapper)(nil),             // 2: go_bludgeon_healthcheck.Wrapper
	(*Bytes)(nil),               // 3: go_bludgeon_healthcheck.Bytes
	(*Error)(nil),               // 4: go_bludgeon_healthcheck.Error
	(*Empty)(nil),               // 5: go_bludgeon_healthcheck.Empty
	(*anypb.Any)(nil),           // 6: google.protobuf.Any
}
var file_healthcheck_proto_depIdxs = []int32{
	1, // 0: go_bludgeon_healthcheck.HealthCheckResponse.healthcheck:type_name -> go_bludgeon_healthcheck.HealthCheck
	6, // 1: go_bludgeon_healthcheck.Wrapper.payload:type_name -> google.protobuf.Any
	5, // 2: go_bludgeon_healthcheck.HealthChecks.healthcheck:input_type -> go_bludgeon_healthcheck.Empty
	0, // 3: go_bludgeon_healthcheck.HealthChecks.healthcheck:output_type -> go_bludgeon_healthcheck.HealthCheckResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_healthcheck_proto_init() }
func file_healthcheck_proto_init() {
	if File_healthcheck_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_healthcheck_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthCheckResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_healthcheck_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthCheck); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_healthcheck_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Wrapper); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_healthcheck_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bytes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_healthcheck_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_healthcheck_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_healthcheck_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_healthcheck_proto_goTypes,
		DependencyIndexes: file_healthcheck_proto_depIdxs,
		MessageInfos:      file_healthcheck_proto_msgTypes,
	}.Build()
	File_healthcheck_proto = out.File
	file_healthcheck_proto_rawDesc = nil
	file_healthcheck_proto_goTypes = nil
	file_healthcheck_proto_depIdxs = nil
}
