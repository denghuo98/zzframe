package zservice

import (
	"github.com/gogf/gf/v2/frame/g"

	webSchema "github.com/denghuo98/zzframe/zschema/zweb"
)

type ISystemConfig interface {
	LoadConfig(ctx g.Ctx) (err error)
	GetSuperAdmin(ctx g.Ctx) (conf *webSchema.SuperAdminConfig, err error)
}

var (
	localSystemConfig ISystemConfig
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
