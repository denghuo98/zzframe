package system

import (
	"path/filepath"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"

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

// Export 导出登录日志
func (c *cSysLoginLog) Export(ctx g.Ctx, req *systemApi.SysLoginLogExportReq) (res *systemApi.SysLoginLogExportRes, err error) {
	filePath, err := zservice.SysLoginLog().Export(ctx, &req.SysLoginLogExportInput)
	if err != nil {
		return nil, err
	}

	// 获取请求对象
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.New("无法获取请求对象")
	}

	// 检查文件是否存在
	if !gfile.Exists(filePath) {
		return nil, gerror.New("导出文件不存在")
	}

	// 设置响应头，支持文件下载
	fileName := filepath.Base(filePath)
	r.Response.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	r.Response.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	r.Response.Header().Set("Content-Transfer-Encoding", "binary")

	// 直接返回文件内容
	r.Response.ServeFile(filePath)

	// 退出当前请求处理，避免后续的 JSON 响应处理
	r.Exit()

	return nil, nil
}
