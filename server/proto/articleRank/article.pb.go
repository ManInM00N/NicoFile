// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v3.12.4
// source: server/proto/articleRank/article.proto

package articleRank

import (
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

type GetArticleRankRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ArticleNum    int64                  `protobuf:"varint,1,opt,name=articleNum,proto3" json:"articleNum,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetArticleRankRequest) Reset() {
	*x = GetArticleRankRequest{}
	mi := &file_server_proto_articleRank_article_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetArticleRankRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetArticleRankRequest) ProtoMessage() {}

func (x *GetArticleRankRequest) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_articleRank_article_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetArticleRankRequest.ProtoReflect.Descriptor instead.
func (*GetArticleRankRequest) Descriptor() ([]byte, []int) {
	return file_server_proto_articleRank_article_proto_rawDescGZIP(), []int{0}
}

func (x *GetArticleRankRequest) GetArticleNum() int64 {
	if x != nil {
		return x.ArticleNum
	}
	return 0
}

type Article struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ArticleId     uint32                 `protobuf:"varint,1,opt,name=articleId,proto3" json:"articleId,omitempty"`
	ArticleTitle  string                 `protobuf:"bytes,2,opt,name=articleTitle,proto3" json:"articleTitle,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Article) Reset() {
	*x = Article{}
	mi := &file_server_proto_articleRank_article_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Article) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Article) ProtoMessage() {}

func (x *Article) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_articleRank_article_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Article.ProtoReflect.Descriptor instead.
func (*Article) Descriptor() ([]byte, []int) {
	return file_server_proto_articleRank_article_proto_rawDescGZIP(), []int{1}
}

func (x *Article) GetArticleId() uint32 {
	if x != nil {
		return x.ArticleId
	}
	return 0
}

func (x *Article) GetArticleTitle() string {
	if x != nil {
		return x.ArticleTitle
	}
	return ""
}

type GetArticleRankResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ArticleList   []*Article             `protobuf:"bytes,2,rep,name=articleList,proto3" json:"articleList,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetArticleRankResponse) Reset() {
	*x = GetArticleRankResponse{}
	mi := &file_server_proto_articleRank_article_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetArticleRankResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetArticleRankResponse) ProtoMessage() {}

func (x *GetArticleRankResponse) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_articleRank_article_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetArticleRankResponse.ProtoReflect.Descriptor instead.
func (*GetArticleRankResponse) Descriptor() ([]byte, []int) {
	return file_server_proto_articleRank_article_proto_rawDescGZIP(), []int{2}
}

func (x *GetArticleRankResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *GetArticleRankResponse) GetArticleList() []*Article {
	if x != nil {
		return x.ArticleList
	}
	return nil
}

var File_server_proto_articleRank_article_proto protoreflect.FileDescriptor

var file_server_proto_articleRank_article_proto_rawDesc = []byte{
	0x0a, 0x26, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x2f, 0x61, 0x72, 0x74, 0x69, 0x63,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c,
	0x65, 0x52, 0x61, 0x6e, 0x6b, 0x22, 0x37, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x41, 0x72, 0x74, 0x69,
	0x63, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e,
	0x0a, 0x0a, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x4e, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x4e, 0x75, 0x6d, 0x22, 0x4b,
	0x0a, 0x07, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x72, 0x74,
	0x69, 0x63, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x61, 0x72,
	0x74, 0x69, 0x63, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x61, 0x72, 0x74, 0x69, 0x63,
	0x6c, 0x65, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x22, 0x6a, 0x0a, 0x16, 0x47,
	0x65, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x36, 0x0a, 0x0b, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x61,
	0x6e, 0x6b, 0x2e, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x0b, 0x61, 0x72, 0x74, 0x69,
	0x63, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x32, 0x6f, 0x0a, 0x12, 0x41, 0x72, 0x74, 0x69, 0x63,
	0x6c, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x59, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x12,
	0x22, 0x2e, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x2e, 0x47, 0x65,
	0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x61, 0x6e,
	0x6b, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x61, 0x6e, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1c, 0x5a, 0x1a, 0x2e, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x72, 0x74, 0x69, 0x63,
	0x6c, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_server_proto_articleRank_article_proto_rawDescOnce sync.Once
	file_server_proto_articleRank_article_proto_rawDescData = file_server_proto_articleRank_article_proto_rawDesc
)

func file_server_proto_articleRank_article_proto_rawDescGZIP() []byte {
	file_server_proto_articleRank_article_proto_rawDescOnce.Do(func() {
		file_server_proto_articleRank_article_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_proto_articleRank_article_proto_rawDescData)
	})
	return file_server_proto_articleRank_article_proto_rawDescData
}

var file_server_proto_articleRank_article_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_server_proto_articleRank_article_proto_goTypes = []any{
	(*GetArticleRankRequest)(nil),  // 0: articleRank.GetArticleRankRequest
	(*Article)(nil),                // 1: articleRank.Article
	(*GetArticleRankResponse)(nil), // 2: articleRank.GetArticleRankResponse
}
var file_server_proto_articleRank_article_proto_depIdxs = []int32{
	1, // 0: articleRank.GetArticleRankResponse.articleList:type_name -> articleRank.Article
	0, // 1: articleRank.ArticleRankService.GetArticleRank:input_type -> articleRank.GetArticleRankRequest
	2, // 2: articleRank.ArticleRankService.GetArticleRank:output_type -> articleRank.GetArticleRankResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_server_proto_articleRank_article_proto_init() }
func file_server_proto_articleRank_article_proto_init() {
	if File_server_proto_articleRank_article_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_server_proto_articleRank_article_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_server_proto_articleRank_article_proto_goTypes,
		DependencyIndexes: file_server_proto_articleRank_article_proto_depIdxs,
		MessageInfos:      file_server_proto_articleRank_article_proto_msgTypes,
	}.Build()
	File_server_proto_articleRank_article_proto = out.File
	file_server_proto_articleRank_article_proto_rawDesc = nil
	file_server_proto_articleRank_article_proto_goTypes = nil
	file_server_proto_articleRank_article_proto_depIdxs = nil
}
