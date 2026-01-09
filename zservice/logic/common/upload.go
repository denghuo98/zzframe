package common

import (
	"context"

	"github.com/dustin/go-humanize"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/denghuo98/zzframe/web/zstorager"
	commonSchema "github.com/denghuo98/zzframe/zschema/common"
	"github.com/denghuo98/zzframe/zservice"
)

type sCommonUpload struct{}

func NewCommonUpload() *sCommonUpload {
	return &sCommonUpload{}
}

func init() {
	zservice.RegisterCommonUpload(NewCommonUpload())
}

// UploadFile 上传文件
func (s *sCommonUpload) UploadFile(ctx context.Context, uploadType string, file *ghttp.UploadFile) (res *commonSchema.AttachmentListModel, err error) {
	attachment, err := zstorager.DoUpload(ctx, uploadType, file)
	if err != nil {
		return
	}

	attachment.FileUrl = zstorager.LastUrl(ctx, attachment.FileUrl, attachment.Drive)
	res = &commonSchema.AttachmentListModel{
		SysAttachment: *attachment,
		SizeFormat:    humanize.IBytes(gconv.Uint64(attachment.Size)),
	}
	return
}

// CheckMultipart 检查文件分片
func (s *sCommonUpload) CheckMultipart(ctx context.Context, in *commonSchema.CheckMultipartInp) (res *commonSchema.CheckMultipartModel, err error) {
	data, err := zstorager.CheckMultipart(ctx, in.CheckMultipartParams)
	if err != nil {
		return nil, err
	}
	res = new(commonSchema.CheckMultipartModel)
	res.CheckMultipartModel = data
	return
}

// UploadPart 上传分片
func (s *sCommonUpload) UploadPart(ctx context.Context, in *commonSchema.UploadPartInp) (res *commonSchema.UploadPartModel, err error) {
	data, err := zstorager.UploadPart(ctx, in.UploadPartParams)
	if err != nil {
		return nil, err
	}
	res = new(commonSchema.UploadPartModel)
	res.UploadPartModel = data
	return
}
