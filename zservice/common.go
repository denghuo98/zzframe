package zservice

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	commonSchema "github.com/denghuo98/zzframe/zschema/common"
	"github.com/denghuo98/zzframe/zschema/zweb"
)

type ICommonSite interface {
	Ping(ctx g.Ctx) string
	InitSuperAdmin(ctx g.Ctx) error
	AccountLogin(ctx g.Ctx, in *commonSchema.SiteAccountLoginInput) (out *commonSchema.SiteLoginOutput, err error)
	BindUserContext(ctx g.Ctx, claims *zweb.Identity) error
}

type ICommonUpload interface {
	UploadFile(ctx g.Ctx, uploadType string, file *ghttp.UploadFile) (res *commonSchema.AttachmentListModel, err error)
	CheckMultipart(ctx g.Ctx, in *commonSchema.CheckMultipartInp) (res *commonSchema.CheckMultipartModel, err error)
	UploadPart(ctx g.Ctx, in *commonSchema.UploadPartInp) (res *commonSchema.UploadPartModel, err error)
}

var (
	localCommonSite   ICommonSite
	localCommonUpload ICommonUpload
)

func CommonSite() ICommonSite {
	if localCommonSite == nil {
		panic("CommonSite is not initialized, please register it first")
	}
	return localCommonSite
}

func RegisterCommonSite(i ICommonSite) {
	localCommonSite = i
}

func CommonUpload() ICommonUpload {
	if localCommonUpload == nil {
		panic("CommonUpload is not initialized, please register it first")
	}
	return localCommonUpload
}

func RegisterCommonUpload(i ICommonUpload) {
	localCommonUpload = i
}
