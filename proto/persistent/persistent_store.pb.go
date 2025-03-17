// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v6.30.1
// source: proto/persistent/persistent_store.proto

package persistent

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PersistentStore struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CreatedAtUnix int64                  `protobuf:"varint,1,opt,name=created_at_unix,json=createdAtUnix,proto3" json:"created_at_unix,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Kv            *KeyValueStore         `protobuf:"bytes,3,opt,name=kv,proto3" json:"kv,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PersistentStore) Reset() {
	*x = PersistentStore{}
	mi := &file_proto_persistent_persistent_store_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PersistentStore) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistentStore) ProtoMessage() {}

func (x *PersistentStore) ProtoReflect() protoreflect.Message {
	mi := &file_proto_persistent_persistent_store_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersistentStore.ProtoReflect.Descriptor instead.
func (*PersistentStore) Descriptor() ([]byte, []int) {
	return file_proto_persistent_persistent_store_proto_rawDescGZIP(), []int{0}
}

func (x *PersistentStore) GetCreatedAtUnix() int64 {
	if x != nil {
		return x.CreatedAtUnix
	}
	return 0
}

func (x *PersistentStore) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *PersistentStore) GetKv() *KeyValueStore {
	if x != nil {
		return x.Kv
	}
	return nil
}

type KeyValueStore struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Store         map[string]*Value      `protobuf:"bytes,1,rep,name=store,proto3" json:"store,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KeyValueStore) Reset() {
	*x = KeyValueStore{}
	mi := &file_proto_persistent_persistent_store_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KeyValueStore) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValueStore) ProtoMessage() {}

func (x *KeyValueStore) ProtoReflect() protoreflect.Message {
	mi := &file_proto_persistent_persistent_store_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValueStore.ProtoReflect.Descriptor instead.
func (*KeyValueStore) Descriptor() ([]byte, []int) {
	return file_proto_persistent_persistent_store_proto_rawDescGZIP(), []int{1}
}

func (x *KeyValueStore) GetStore() map[string]*Value {
	if x != nil {
		return x.Store
	}
	return nil
}

type Value struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Value         string                 `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Values        []string               `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
	Datatype      string                 `protobuf:"bytes,3,opt,name=datatype,proto3" json:"datatype,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Value) Reset() {
	*x = Value{}
	mi := &file_proto_persistent_persistent_store_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_proto_persistent_persistent_store_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_proto_persistent_persistent_store_proto_rawDescGZIP(), []int{2}
}

func (x *Value) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Value) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *Value) GetDatatype() string {
	if x != nil {
		return x.Datatype
	}
	return ""
}

var File_proto_persistent_persistent_store_proto protoreflect.FileDescriptor

var file_proto_persistent_persistent_store_proto_rawDesc = string([]byte{
	0x0a, 0x27, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x74, 0x2f, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x70, 0x65, 0x72, 0x73, 0x69,
	0x73, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x83, 0x01, 0x0a, 0x0f, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73,
	0x74, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x5f, 0x75, 0x6e, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x55, 0x6e, 0x69,
	0x78, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x29, 0x0a, 0x02, 0x6b, 0x76, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70,
	0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x02, 0x6b, 0x76, 0x22, 0x98, 0x01, 0x0a, 0x0d,
	0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x3a, 0x0a,
	0x05, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x70,
	0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x05, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x4b, 0x0a, 0x0a, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x27, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x69,
	0x73, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x51, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x1a, 0x0a,
	0x08, 0x64, 0x61, 0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x64, 0x61, 0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x42, 0x14, 0x5a, 0x12, 0x67, 0x6f, 0x72,
	0x65, 0x64, 0x69, 0x73, 0x2f, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_persistent_persistent_store_proto_rawDescOnce sync.Once
	file_proto_persistent_persistent_store_proto_rawDescData []byte
)

func file_proto_persistent_persistent_store_proto_rawDescGZIP() []byte {
	file_proto_persistent_persistent_store_proto_rawDescOnce.Do(func() {
		file_proto_persistent_persistent_store_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_persistent_persistent_store_proto_rawDesc), len(file_proto_persistent_persistent_store_proto_rawDesc)))
	})
	return file_proto_persistent_persistent_store_proto_rawDescData
}

var file_proto_persistent_persistent_store_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_persistent_persistent_store_proto_goTypes = []any{
	(*PersistentStore)(nil), // 0: persistent.PersistentStore
	(*KeyValueStore)(nil),   // 1: persistent.KeyValueStore
	(*Value)(nil),           // 2: persistent.Value
	nil,                     // 3: persistent.KeyValueStore.StoreEntry
}
var file_proto_persistent_persistent_store_proto_depIdxs = []int32{
	1, // 0: persistent.PersistentStore.kv:type_name -> persistent.KeyValueStore
	3, // 1: persistent.KeyValueStore.store:type_name -> persistent.KeyValueStore.StoreEntry
	2, // 2: persistent.KeyValueStore.StoreEntry.value:type_name -> persistent.Value
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_persistent_persistent_store_proto_init() }
func file_proto_persistent_persistent_store_proto_init() {
	if File_proto_persistent_persistent_store_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_persistent_persistent_store_proto_rawDesc), len(file_proto_persistent_persistent_store_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_persistent_persistent_store_proto_goTypes,
		DependencyIndexes: file_proto_persistent_persistent_store_proto_depIdxs,
		MessageInfos:      file_proto_persistent_persistent_store_proto_msgTypes,
	}.Build()
	File_proto_persistent_persistent_store_proto = out.File
	file_proto_persistent_persistent_store_proto_goTypes = nil
	file_proto_persistent_persistent_store_proto_depIdxs = nil
}
