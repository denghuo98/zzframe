package zservice

import (
	"github.com/gogf/gf/v2/frame/g"

	"github.com/denghuo98/zzframe/internal/model/entity"
	systemSchema "github.com/denghuo98/zzframe/zschema/system"
	webSchema "github.com/denghuo98/zzframe/zschema/zweb"
)

type ISystemConfig interface {
	LoadConfig(ctx g.Ctx) (err error)
	GetSuperAdmin(ctx g.Ctx) (conf *webSchema.SuperAdminConfig, err error)
}

type ISysLoginLog interface {
	Push(ctx g.Ctx, in *systemSchema.SysLoginLogPushInput) (err error)
	RealWrite(ctx g.Ctx, data entity.SysLoginLog) (err error)
	Delete(ctx g.Ctx, in *systemSchema.SysLoginLogDeleteInput) (err error)
	List(ctx g.Ctx, in *systemSchema.SysLoginLogListInput) (out *systemSchema.SysLoginLogListOutput, totalCount int, err error)
}

var (
	localSystemConfig ISystemConfig
	localSysLoginLog  ISysLoginLog
)

func SystemConfig() ISystemConfig {
	if localSystemConfig == nil {
		panic("SystemConfig is not initialized, please register it first")
	}
	return localSystemConfig
}

func RegisterSystemConfig(i ISystemConfig) {
	localSystemConfig = i
}

func SysLoginLog() ISysLoginLog {
	if localSysLoginLog == nil {
		panic("SysLoginLog is not initialized, please register it first")
	}
	return localSysLoginLog
}

func RegisterSysLoginLog(i ISysLoginLog) {
	localSysLoginLog = i
}
