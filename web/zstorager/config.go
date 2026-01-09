package zstorager

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	webSchema "github.com/denghuo98/zzframe/zschema/zweb"
)

var config *webSchema.UploadConfig

func SetConfig(c *webSchema.UploadConfig) {
	config = c
}

func GetConfig() *webSchema.UploadConfig {
	return config
}

func GetModel(ctx g.Ctx) *gdb.Model {
	return g.Model("sys_attachment").Ctx(ctx)
}
