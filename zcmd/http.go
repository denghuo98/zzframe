package zcmd

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/denghuo98/zzframe/web/zcasbin"
	"github.com/denghuo98/zzframe/zcontroller"
	"github.com/denghuo98/zzframe/zservice"
)

var (
	Http = &gcmd.Command{
		Name:  "http",
		Usage: "http",
		Brief: "HTTP服务，也可以称之为主服务",
		Func: func(ctx g.Ctx, parser *gcmd.Parser) (err error) {
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

			// 初始化 casbin
			zcasbin.InitEnforcer(ctx)

			// 信号监听
			SignalListen(ctx, SignalHandlerForOverall)

			go func() {
				<-serverCloseSignal
				_ = s.Shutdown() // 关闭http服务，主服务建议放在最后一个关闭
				g.Log().Debug(ctx, "http successfully closed ..")
				serverWg.Done()
			}()
			s.Run()
			return
		},
	}
)
