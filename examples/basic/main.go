package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/denghuo98/zzframe/zcmd"
	"github.com/denghuo98/zzframe/zservice"

	_ "github.com/denghuo98/zzframe/zdb/zmigrate"
	_ "github.com/denghuo98/zzframe/zservice/logic"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	var ctx = gctx.GetInitCtx()

	// 初始化系统配置
	if err := zservice.SystemConfig().LoadConfig(ctx); err != nil {
		g.Log().Panicf(ctx, "初始化系统配置失败: %v", err)
	}

	zcmd.Main.Run(ctx)
}
