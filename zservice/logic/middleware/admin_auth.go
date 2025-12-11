package middleware

import (
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/samber/lo"

	"github.com/denghuo98/zzframe/web/zcache"
	"github.com/denghuo98/zzframe/web/zresp"
	"github.com/denghuo98/zzframe/zservice"
)

// AdminAuth 后台登录认证
func (s *sMiddleware) AdminAuth(r *ghttp.Request) {
	var (
		ctx  = r.Context()
		path = r.URL.Path
	)

	if err := s.DeliverUserContext(r); err != nil {
		zresp.JsonExit(r, gcode.CodeNotAuthorized.Code(), err.Error())
		return
	}

	// 不需要验证权限
	if s.IsExceptAuth(ctx, path) {
		r.Middleware.Next()
		return
	}

	// 验证路由访问权限
	if !zservice.AdminRole().Verify(ctx, path, r.Method) {
		zresp.JsonExit(r, gcode.CodeNotAuthorized.Code(), "您没有权限访问该页面")
		return
	}
	r.Middleware.Next()
}

// IsExceptAuth 是否是不需要验证权限的路由地址
func (s *sMiddleware) IsExceptAuth(ctx g.Ctx, path string) bool {

	var exceptAuth []string
	// 从缓存中获取不需要验证权限的路由地址
	v, err := zcache.Instance().Get(ctx, "exceptAuth")
	if err != nil {
		g.Log().Error(ctx, "从缓存获取不需要验证权限的路由地址失败", err)
	}

	if v.IsEmpty() {
		exceptAuth = g.Cfg().MustGet(ctx, "system.exceptAuth").Strings()
		zcache.Instance().Set(ctx, "exceptAuth", exceptAuth, time.Hour)
	} else {
		exceptAuth = v.Strings()
	}

	return lo.Contains(exceptAuth, path)
}
