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

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
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

func (x *Container) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Started struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Started bool `protobuf:"varint,1,opt,name=started,proto3" json:"started,omitempty"`
}

func (x *Started) Reset() {
	*x = Started{}
	if protoimpl.UnsafeEnabled {
		mi := &file_caked_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Started) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Started) ProtoMessage() {}

func (x *Started) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Started.ProtoReflect.Descriptor instead.
func (*Started) Descriptor() ([]byte, []int) {
	return file_caked_proto_rawDescGZIP(), []int{1}
}

func (x *Started) GetStarted() bool {
	if x != nil {
		return x.Started
	}
	return false
}

var File_caked_proto protoreflect.FileDescriptor

var file_caked_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x61, 0x6b, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63,
	0x61, 0x6b, 0x65, 0x22, 0x1b, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x23, 0x0a, 0x07, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x65, 0x64, 0x32, 0x3b, 0x0a, 0x05, 0x43, 0x61, 0x6b, 0x65, 0x64, 0x12, 0x32,
	0x0a, 0x0e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x12, 0x0f, 0x2e, 0x63, 0x61, 0x6b, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x1a, 0x0d, 0x2e, 0x63, 0x61, 0x6b, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64,
	0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x3b, 0x63, 0x61, 0x6b, 0x65, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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
	(*Container)(nil), // 0: cake.Container
	(*Started)(nil),   // 1: cake.Started
}
var file_caked_proto_depIdxs = []int32{
	0, // 0: cake.Caked.StartContainer:input_type -> cake.Container
	1, // 1: cake.Caked.StartContainer:output_type -> cake.Started
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
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
			switch v := v.(*Started); i {
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
