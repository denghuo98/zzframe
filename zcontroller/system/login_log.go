package system

import (
	"github.com/gogf/gf/v2/frame/g"

	systemApi "github.com/denghuo98/zzframe/zapi/system"
	"github.com/denghuo98/zzframe/zservice"
)

type cSysLoginLog struct{}

var SysLoginLog = cSysLoginLog{}

// Delete 删除登录日志
func (c *cSysLoginLog) Delete(ctx g.Ctx, req *systemApi.SysLoginLogDeleteReq) (res *systemApi.SysLoginLogDeleteRes, err error) {
	err = zservice.SysLoginLog().Delete(ctx, &req.SysLoginLogDeleteInput)
	return
}

// List 获取登录日志列表
func (c *cSysLoginLog) List(ctx g.Ctx, req *systemApi.SysLoginLogListReq) (res *systemApi.SysLoginLogListRes, err error) {
	out, totalCount, err := zservice.SysLoginLog().List(ctx, &req.SysLoginLogListInput)
	if err != nil {
		return nil, err
	}
	res = new(systemApi.SysLoginLogListRes)
	res.SysLoginLogListOutput = *out
	res.PageRes.Pack(req, int(totalCount))
	return
}
