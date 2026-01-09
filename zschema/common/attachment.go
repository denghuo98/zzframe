package common

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/samber/lo"

	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/web/zstorager"
	"github.com/denghuo98/zzframe/zschema/zform"
)

// AttachmentDeleteInp 删除附件
type AttachmentDeleteInp struct {
	Id interface{} `json:"id" v:"required#附件ID不能为空" dc:"附件ID"`
}

type AttachmentDeleteModel struct{}

// AttachmentViewInp 获取附件信息
type AttachmentViewInp struct {
	Id int64 `json:"id" v:"required#附件ID不能为空" dc:"附件ID"`
}

type AttachmentViewModel struct {
	entity.SysAttachment
}

// AttachmentClearKindInp 清空上传类型
type AttachmentClearKindInp struct {
	Kind string `json:"kind" v:"required#上传类型不能为空" dc:"上传类型"`
}

func (in *AttachmentClearKindInp) Filter(ctx context.Context) (err error) {
	if !lo.Contains(zstorager.KindSlice, in.Kind) {
		err = gerror.New("上传类型是无效的")
		return
	}
	return
}

// AttachmentListInp 获取附件列表
type AttachmentListInp struct {
	zform.PageReq
	MemberId  int64         `json:"memberId"  dc:"用户ID"`
	Name      string        `json:"name"       dc:"文件名称"`
	Drive     string        `json:"drive"      dc:"驱动"`
	Kind      string        `json:"kind"       dc:"上传类型"`
	UpdatedAt []*gtime.Time `json:"updatedAt"  dc:"更新时间"`
}

type AttachmentListModel struct {
	entity.SysAttachment
	SizeFormat string `json:"sizeFormat"      dc:"大小"`
}

// AttachmentChooserListInp 获取附件列表
type AttachmentChooserListInp struct {
	zform.PageReq
	Drive     string  `json:"drive"      dc:"驱动"`
	Kind      string  `json:"kind"       dc:"上传类型"`
	UpdatedAt []int64 `json:"updatedAt"  dc:"更新时间"`
}

type AttachmentChooserListModel struct {
	entity.SysAttachment
	SizeFormat string `json:"sizeFormat"      dc:"大小"`
}

// AttachmentClearInp 清空分类
type AttachmentClearInp struct {
	Kind string `json:"kind"       dc:"上传类型"`
}

// CheckMultipartInp 检查文件分片
type CheckMultipartInp struct {
	*zstorager.CheckMultipartParams
}

type CheckMultipartModel struct {
	*zstorager.CheckMultipartModel
}

// UploadPartInp 上传分片
type UploadPartInp struct {
	*zstorager.UploadPartParams
}

type UploadPartModel struct {
	*zstorager.UploadPartModel
}
