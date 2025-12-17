package system

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/web/zlocation"
	"github.com/denghuo98/zzframe/web/zqueue"
	"github.com/denghuo98/zzframe/web/zutils"
	"github.com/denghuo98/zzframe/zconsts"
	commonSchema "github.com/denghuo98/zzframe/zschema/common"
	systemSchema "github.com/denghuo98/zzframe/zschema/system"
	"github.com/denghuo98/zzframe/zservice"
)

func init() {
	zservice.RegisterSysLoginLog(NewSysLoginLog())
}

type sSysLoginLog struct{}

func NewSysLoginLog() *sSysLoginLog {
	return &sSysLoginLog{}
}

// Push 推送登录日志（通常异步推送）
func (s *sSysLoginLog) Push(ctx g.Ctx, in *systemSchema.SysLoginLogPushInput) (err error) {
	if in.Response == nil {
		in.Response = new(commonSchema.SiteLoginOutput)
	}

	r := g.RequestFromCtx(ctx)
	if r == nil {
		g.Log().Warningf(ctx, "request from ctx is nil")
		return
	}

	clientIp := zlocation.GetClientIp(r)
	ipData, err := zlocation.GetLocation(ctx, clientIp)
	if err != nil {
		g.Log().Warningf(ctx, "get location fail, ip:%v, err:%v", clientIp, err)
		return
	}

	if ipData == nil {
		ipData = new(zlocation.IpLocationData)
	}

	var models entity.SysLoginLog
	models.ReqId = gctx.CtxId(ctx)
	models.MemberId = in.Response.Id
	models.Username = in.Response.Username
	models.LoginAt = gtime.Now()
	models.LoginIp = clientIp
	models.UserAgent = r.UserAgent()
	models.Province = ipData.Province
	models.City = ipData.City
	models.Status = zconsts.StatusEnabled

	if in.Error != nil {
		models.Status = zconsts.StatusDisable
		models.ErrMsg = in.Error.Error()
	}

	models.Response = gjson.New(zconsts.NilJsonToString)
	if in.Response != nil {
		models.Response = gjson.New(in.Response)
	}

	if err = zqueue.Push(zconsts.QueueLoginLogTopic, models); err != nil {
		g.Log().Warningf(ctx, "push login log to queue failed, err:%+v, data:%+v", err, models)
	}
	return
}

func (s *sSysLoginLog) RealWrite(ctx g.Ctx, data entity.SysLoginLog) (err error) {
	_, err = dao.SysLoginLog.Ctx(ctx).Data(data).OmitEmpty().Insert()
	return
}

func (s *sSysLoginLog) Delete(ctx g.Ctx, in *systemSchema.SysLoginLogDeleteInput) (err error) {
	if in.Id <= 0 {
		return gerror.New("ID不能为空")
	}
	_, err = dao.SysLoginLog.Ctx(ctx).WherePri(in.Id).Delete()
	return
}

func (s *sSysLoginLog) List(ctx g.Ctx, in *systemSchema.SysLoginLogListInput) (out *systemSchema.SysLoginLogListOutput, totalCount int, err error) {
	var (
		m    = dao.SysLoginLog.Ctx(ctx)
		cols = dao.SysLoginLog.Columns()
	)

	// 筛选状态
	if in.Status > 0 {
		m = m.Where(cols.Status, in.Status)
	}

	// 筛选登录时间
	if len(in.LoginAt) == 2 {
		m = m.WhereBetween(cols.LoginAt, in.LoginAt[0], in.LoginAt[1])
	}

	// 筛选登录IP
	if in.LoginIp != "" {
		m = m.Where(cols.LoginIp, in.LoginIp)
	}

	// 筛选账号
	if in.Username != "" {
		m = m.Where(cols.Username, in.Username)
	}

	totalCount, err = m.Count()
	if err != nil || totalCount == 0 {
		return
	}

	var list []*systemSchema.SysLoginLogListOuputItem
	if err = m.Page(in.Page, in.PerPage).OrderDesc(cols.Id).Scan(&list); err != nil {
		return nil, 0, gerror.Wrap(err, zconsts.ErrorORM)
	}

	for _, item := range list {
		item.Os = zutils.GetOs(item.UserAgent)
		item.Browser = zutils.GetBrowser(item.UserAgent)
	}
	out = new(systemSchema.SysLoginLogListOutput)
	out.List = list
	return
}
