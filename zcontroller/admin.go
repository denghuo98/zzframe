package zcontroller

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	adminController "github.com/denghuo98/zzframe/zcontroller/admin"
	commonController "github.com/denghuo98/zzframe/zcontroller/common"
	"github.com/denghuo98/zzframe/zservice"
)

func Admin(ctx g.Ctx, group *ghttp.RouterGroup) {

	// 兼容后台登录入口
	group.ALL("/login", func(r *ghttp.Request) {
		r.Response.RedirectTo("/admin/login")
	})

	group.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Bind(
			commonController.Site, // 站点公共接口
		)
		group.Middleware(zservice.Middleware().AdminAuth)
		group.Bind(
			adminController.Menu,
			adminController.Role,
			adminController.Member,
		)
	})
}
