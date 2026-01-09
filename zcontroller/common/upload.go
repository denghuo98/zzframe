package common

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"

	"github.com/denghuo98/zzframe/web/zstorager"
	commonApi "github.com/denghuo98/zzframe/zapi/common"
	"github.com/denghuo98/zzframe/zservice"
)

var Upload = new(cUpload)

type cUpload struct{}

// UploadFile 上传文件
func (c *cUpload) UploadFile(ctx context.Context, _ *commonApi.UploadFileReq) (res commonApi.UploadFileRes, err error) {
	r := g.RequestFromCtx(ctx)
	uploadType := r.Header.Get("uploadType")
	if uploadType != "default" && !lo.Contains(zstorager.KindSlice, uploadType) {
		err = gerror.New("上传类型是无效的")
		return
	}

	file := r.GetUploadFile("file")
	if file == nil {
		err = gerror.New("没有找到上传的文件")
		return
	}
	return zservice.CommonUpload().UploadFile(ctx, uploadType, file)
}

// CheckMultipart 检查文件分片
func (c *cUpload) CheckMultipart(ctx context.Context, req *commonApi.CheckMultipartReq) (res *commonApi.CheckMultipartRes, err error) {
	data, err := zservice.CommonUpload().CheckMultipart(ctx, &req.CheckMultipartInp)
	if err != nil {
		return nil, err
	}
	res = new(commonApi.CheckMultipartRes)
	res.CheckMultipartModel = data
	return
}

// UploadPart 上传分片
func (c *cUpload) UploadPart(ctx context.Context, req *commonApi.UploadPartReq) (res *commonApi.UploadPartRes, err error) {
	data, err := zservice.CommonUpload().UploadPart(ctx, &req.UploadPartInp)
	if err != nil {
		return nil, err
	}
	res = new(commonApi.UploadPartRes)
	res.UploadPartModel = data
	return
}
