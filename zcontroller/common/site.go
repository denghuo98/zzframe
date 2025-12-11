package common

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gmode"

	"github.com/denghuo98/zzframe/web/zcaptcha"
	commonApi "github.com/denghuo98/zzframe/zapi/common"
	"github.com/denghuo98/zzframe/zservice"
)

type cSite struct{}

var Site = cSite{}

func (c *cSite) Ping(ctx g.Ctx, req *commonApi.SitePingReq) (res *commonApi.SitePingRes, err error) {
	status := zservice.CommonSite().Ping(ctx)
	res = new(commonApi.SitePingRes)
	res.Status = status
	return
}

func (c *cSite) Captcha(ctx g.Ctx, req *commonApi.SiteCaptchaReq) (res *commonApi.SiteCaptchaRes, err error) {
	id, base64 := zcaptcha.Generate(ctx)
	res = new(commonApi.SiteCaptchaRes)
	res.Id = id
	res.Base64 = base64
	return
}

func (c *cSite) AccountLogin(ctx g.Ctx, req *commonApi.SiteAccountLoginReq) (res *commonApi.SiteAccountLoginRes, err error) {
	// 验证码(开发环境不验证)
	if !gmode.IsDevelop() && !zcaptcha.Verify(req.Cid, req.Code) {
		return nil, gerror.New("验证码错误")
	}

	out, err := zservice.CommonSite().AccountLogin(ctx, &req.SiteAccountLoginInput)
	if err != nil {
		return nil, err
	}
	res = new(commonApi.SiteAccountLoginRes)
	res.SiteLoginOutput = *out
	return
}
