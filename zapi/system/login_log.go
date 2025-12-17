package system

import (
	"github.com/gogf/gf/v2/frame/g"

	systemSchema "github.com/denghuo98/zzframe/zschema/system"
	"github.com/denghuo98/zzframe/zschema/zform"
)

// SysLoginLogDeleteReq 删除登录日志
type SysLoginLogDeleteReq struct {
	g.Meta `path:"/system/login-log/delete" method:"delete" tags:"SYS-10-登录日志" summary:"删除登录日志"`
	systemSchema.SysLoginLogDeleteInput
}

type SysLoginLogDeleteRes struct{}

// SysLoginLogListReq 获取登录日志列表
type SysLoginLogListReq struct {
	g.Meta `path:"/system/login-log/list" method:"get" tags:"SYS-10-登录日志" summary:"获取登录日志列表"`
	systemSchema.SysLoginLogListInput
}

type SysLoginLogListRes struct {
	zform.PageRes
	systemSchema.SysLoginLogListOutput
}
