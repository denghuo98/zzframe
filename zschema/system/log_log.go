package system

import (
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/denghuo98/zzframe/internal/model/entity"
	commonSchema "github.com/denghuo98/zzframe/zschema/common"
	"github.com/denghuo98/zzframe/zschema/zform"
)

type SysLoginLogPushInput struct {
	Response *commonSchema.SiteLoginOutput
	Error    error
}

type SysLoginLogDeleteInput struct {
	Id int64
}

type SysLoginLogListInput struct {
	zform.PageReq
	Username string        `json:"username" dc:"账号"`
	Status   int           `json:"status" dc:"状态"`
	LoginAt  []*gtime.Time `json:"loginAt" dc:"登录时间"`
	LoginIp  string        `json:"loginIp" dc:"登录IP"`
}

type SysLoginLogListOuputItem struct {
	entity.SysLoginLog
	Os      string `json:"os" dc:"操作系统"`
	Browser string `json:"browser" dc:"浏览器"`
}

type SysLoginLogListOutput struct {
	List []*SysLoginLogListOuputItem `json:"list" dc:"列表"`
}
