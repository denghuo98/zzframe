package zservice

import (
	"github.com/gogf/gf/v2/frame/g"

	commonSchema "github.com/denghuo98/zzframe/zschema/common"
	"github.com/denghuo98/zzframe/zschema/zweb"
)

type ICommonSite interface {
	Ping(ctx g.Ctx) string
	InitSuperAdmin(ctx g.Ctx) error
	AccountLogin(ctx g.Ctx, in *commonSchema.SiteAccountLoginInput) (out *commonSchema.SiteLoginOutput, err error)
	BindUserContext(ctx g.Ctx, claims *zweb.Identity) error
}

var (
	localCommonSite ICommonSite
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
