// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.17.3
// source: caked.proto

package cakepb

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Container struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImageName          string `protobuf:"bytes,1,opt,name=ImageName,proto3" json:"ImageName,omitempty"`
	Tag                string `protobuf:"bytes,2,opt,name=Tag,proto3" json:"Tag,omitempty"`
	Registry           string `protobuf:"bytes,3,opt,name=Registry,proto3" json:"Registry,omitempty"`
	PreviousDigest     string `protobuf:"bytes,4,opt,name=PreviousDigest,proto3" json:"PreviousDigest,omitempty"`
	PreviousDigestTime int64  `protobuf:"varint,5,opt,name=PreviousDigestTime,proto3" json:"PreviousDigestTime,omitempty"`
	LatestDigest       string `protobuf:"bytes,6,opt,name=LatestDigest,proto3" json:"LatestDigest,omitempty"`
	LatestDigestTime   int64  `protobuf:"varint,7,opt,name=LatestDigestTime,proto3" json:"LatestDigestTime,omitempty"`
	LastChecked        int64  `protobuf:"varint,8,opt,name=LastChecked,proto3" json:"LastChecked,omitempty"`
	LastUpdated        int64  `protobuf:"varint,9,opt,name=LastUpdated,proto3" json:"LastUpdated,omitempty"`
	Architecture       string `protobuf:"bytes,10,opt,name=Architecture,proto3" json:"Architecture,omitempty"`
	OS                 string `protobuf:"bytes,11,opt,name=OS,proto3" json:"OS,omitempty"`
}

func (x *Container) Reset() {
	*x = Container{}
	if protoimpl.UnsafeEnabled {
		mi := &file_caked_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Container) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Container) ProtoMessage() {}

func (x *Container) ProtoReflect() protoreflect.Message {
	mi := &file_caked_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Container.ProtoReflect.Descriptor instead.
func (*Container) Descriptor() ([]byte, []int) {
	return file_caked_proto_rawDescGZIP(), []int{0}
}

func (x *Container) GetImageName() string {
	if x != nil {
		return x.ImageName
	}
	return ""
}

func (x *Container) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *Container) GetRegistry() string {
	if x != nil {
		return x.Registry
	}
	return ""
}

func (x *Container) GetPreviousDigest() string {
	if x != nil {
		return x.PreviousDigest
	}
	return ""
}

func (x *Container) GetPreviousDigestTime() int64 {
	if x != nil {
		return x.PreviousDigestTime
	}
	return 0
}

func (x *Container) GetLatestDigest() string {
	if x != nil {
		return x.LatestDigest
	}
	return ""
}

func (x *Container) GetLatestDigestTime() int64 {
	if x != nil {
		return x.LatestDigestTime
	}
	return 0
}

func (x *Container) GetLastChecked() int64 {
	if x != nil {
		return x.LastChecked
	}
	return 0
}

func (x *Container) GetLastUpdated() int64 {
	if x != nil {
		return x.LastUpdated
	}
	return 0
}

func (x *Container) GetArchitecture() string {
	if x != nil {
		return x.Architecture
	}
	return ""
}

func (x *Container) GetOS() string {
	if x != nil {
		return x.OS
	}
	return ""
}

type ContainerStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status      int32  `protobuf:"varint,1,opt,name=Status,proto3" json:"Status,omitempty"`
	ContainerId string `protobuf:"bytes,2,opt,name=ContainerId,proto3" json:"ContainerId,omitempty"`
	Message     string `protobuf:"bytes,3,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *ContainerStatus) Reset() {
	*x = ContainerStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_caked_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContainerStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerStatus) ProtoMessage() {}

func (x *ContainerStatus) ProtoReflect() protoreflect.Message {
	mi := &file_caked_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerStatus.ProtoReflect.Descriptor instead.
func (*ContainerStatus) Descriptor() ([]byte, []int) {
	return file_caked_proto_rawDescGZIP(), []int{1}
}

func (x *ContainerStatus) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *ContainerStatus) GetContainerId() string {
	if x != nil {
		return x.ContainerId
	}
	return ""
}

func (x *ContainerStatus) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_caked_proto protoreflect.FileDescriptor

var file_caked_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x61, 0x6b, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63,
	0x61, 0x6b, 0x65, 0x22, 0xf7, 0x02, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x12, 0x1c, 0x0a, 0x09, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x54, 0x61,
	0x67, 0x12, 0x1a, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x12, 0x26, 0x0a,
	0x0e, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x44,
	0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x12, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75,
	0x73, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x12, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x44, 0x69, 0x67, 0x65, 0x73,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x44,
	0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x4c, 0x61, 0x74,
	0x65, 0x73, 0x74, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x10, 0x4c, 0x61, 0x74,
	0x65, 0x73, 0x74, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x10, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x44, 0x69, 0x67, 0x65, 0x73,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x4c, 0x61, 0x73, 0x74, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x65, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x4c, 0x61, 0x73, 0x74,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x4c, 0x61, 0x73, 0x74, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x4c, 0x61,
	0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x41, 0x72, 0x63,
	0x68, 0x69, 0x74, 0x65, 0x63, 0x74, 0x75, 0x72, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x41, 0x72, 0x63, 0x68, 0x69, 0x74, 0x65, 0x63, 0x74, 0x75, 0x72, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x4f, 0x53, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x4f, 0x53, 0x22, 0x65, 0x0a,
	0x0f, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x43,
	0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x32, 0x7e, 0x0a, 0x05, 0x43, 0x61, 0x6b, 0x65, 0x64, 0x12, 0x3a, 0x0a,
	0x0e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12,
	0x0f, 0x2e, 0x63, 0x61, 0x6b, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x1a, 0x15, 0x2e, 0x63, 0x61, 0x6b, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x0d, 0x53, 0x74, 0x6f,
	0x70, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x0f, 0x2e, 0x63, 0x61, 0x6b,
	0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x1a, 0x15, 0x2e, 0x63, 0x61,
	0x6b, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x3b, 0x63, 0x61, 0x6b, 0x65, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_caked_proto_rawDescOnce sync.Once
	file_caked_proto_rawDescData = file_caked_proto_rawDesc
)

func file_caked_proto_rawDescGZIP() []byte {
	file_caked_proto_rawDescOnce.Do(func() {
		file_caked_proto_rawDescData = protoimpl.X.CompressGZIP(file_caked_proto_rawDescData)
	})
	return file_caked_proto_rawDescData
}

var file_caked_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_caked_proto_goTypes = []interface{}{
	(*Container)(nil),       // 0: cake.Container
	(*ContainerStatus)(nil), // 1: cake.ContainerStatus
}
var file_caked_proto_depIdxs = []int32{
	0, // 0: cake.Caked.StartContainer:input_type -> cake.Container
	0, // 1: cake.Caked.StopContainer:input_type -> cake.Container
	1, // 2: cake.Caked.StartContainer:output_type -> cake.ContainerStatus
	1, // 3: cake.Caked.StopContainer:output_type -> cake.ContainerStatus
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_caked_proto_init() }
func file_caked_proto_init() {
	if File_caked_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_caked_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Container); i {
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
		file_caked_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContainerStatus); i {
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
			RawDescriptor: file_caked_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_caked_proto_goTypes,
		DependencyIndexes: file_caked_proto_depIdxs,
		MessageInfos:      file_caked_proto_msgTypes,
	}.Build()
	File_caked_proto = out.File
	file_caked_proto_rawDesc = nil
	file_caked_proto_goTypes = nil
	file_caked_proto_depIdxs = nil
}
