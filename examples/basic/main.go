package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/denghuo98/zzframe/web/zcasbin"
	"github.com/denghuo98/zzframe/zcontroller"
	"github.com/denghuo98/zzframe/zservice"

	_ "github.com/denghuo98/zzframe/zservice/logic"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	var ctx = gctx.GetInitCtx()
	g.Log().Info(ctx, "Hello World")

	s := g.Server()

	// 注册全局中间件
	s.BindMiddleware("/*any", []ghttp.HandlerFunc{
		zservice.Middleware().Ctx,
		zservice.Middleware().CORS,
		zservice.Middleware().ResponseHandler,
	}...)

	// 注册路由，后台管理
	s.Group("/", func(group *ghttp.RouterGroup) {
		zcontroller.Admin(ctx, group)
	})
	s.SetPort(9090)
	s.SetOpenApiPath("/api.json")
	s.SetSwaggerPath("/swagger")

	// 初始化 casbin
	zcasbin.InitEnforcer(ctx)

	// 初始化系统配置
	if err := zservice.SystemConfig().LoadConfig(ctx); err != nil {
		g.Log().Panicf(ctx, "初始化系统配置失败: %v", err)
		return
	}

	s.Run()
}
