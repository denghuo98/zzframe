package common

import (
	"github.com/gogf/gf/v2/frame/g"

	commonSchema "github.com/denghuo98/zzframe/zschema/common"
)

type SitePingReq struct {
	g.Meta `path:"/site/ping" method:"get" tags:"SYS-00-系统管理" summary:"心跳检测"`
}

type SitePingRes struct {
	Status string `json:"status" dc:"状态"`
}

type SiteCaptchaReq struct {
	g.Meta `path:"/site/captcha" method:"get" tags:"SYS-00-系统管理" summary:"获取验证码"`
}

type SiteCaptchaRes struct {
	Id     string `json:"id" dc:"验证码ID"`
	Base64 string `json:"base64" dc:"验证码图片"`
}

// SiteAccountLoginReq 账号登录
type SiteAccountLoginReq struct {
	g.Meta `path:"/site/account/login" method:"post" tags:"SYS-00-系统管理" summary:"账号登录"`
	commonSchema.SiteAccountLoginInput
}

type SiteAccountLoginRes struct {
	commonSchema.SiteLoginOutput
}
