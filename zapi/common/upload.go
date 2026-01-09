package common

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	commonSchema "github.com/denghuo98/zzframe/zschema/common"
)

// UploadFileReq 上传文件
type UploadFileReq struct {
	g.Meta     `path:"/upload/file" tags:"SYS-11-附件" method:"post" summary:"上传附件" mime:"multipart/form-data"`
	UploadType string            `json:"uploadType" in:"header" dc:"上传类型" v:"required|in:default,image,doc,audio,video,zip,other"`
	File       *ghttp.UploadFile `json:"file" mime:"file" dc:"文件"`
}

type UploadFileRes *commonSchema.AttachmentListModel

// CheckMultipartReq 检查文件分片
type CheckMultipartReq struct {
	g.Meta `path:"/upload/checkMultipart" tags:"SYS-11-附件" method:"post" summary:"检查文件分片"`
	commonSchema.CheckMultipartInp
}

type CheckMultipartRes struct {
	*commonSchema.CheckMultipartModel
}

// UploadPartReq 分片上传
type UploadPartReq struct {
	g.Meta `path:"/upload/uploadPart" tags:"SYS-11-附件" method:"post" summary:"分片上传"`
	commonSchema.UploadPartInp
}

type UploadPartRes struct {
	*commonSchema.UploadPartModel
}
