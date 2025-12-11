package zcontext

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/denghuo98/zzframe/zconsts"
	webSchema "github.com/denghuo98/zzframe/zschema/zweb"
)

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程可以直接修改
func Init(r *ghttp.Request, customCtx *webSchema.Context) {
	r.SetCtxVar(zconsts.ContextHTTPKey, customCtx)
}

// Get 获得上下文变量，如果没有设置，那么返回 nil
func Get(ctx g.Ctx) *webSchema.Context {
	value := ctx.Value(zconsts.ContextHTTPKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*webSchema.Context); ok {
		return localCtx
	}
	return nil
}

// SetUser 将上下文信息设置到上下文请求中，注意是完整覆盖
func SetUser(ctx g.Ctx, user *webSchema.Identity) {
	c := Get(ctx)
	if c == nil {
		g.Log().Warning(ctx, "contexts.SetUser, c == nil ")
		return
	}
	c.User = user
}

// SetResponse 设置组件响应 用于访问日志使用
func SetResponse(ctx g.Ctx, response *webSchema.Response) {
	c := Get(ctx)
	if c == nil {
		g.Log().Warning(ctx, "上下文对象不存在")
		return
	}
	c.Response = response
}

// GetUser 获取用户信息
func GetUser(ctx g.Ctx) *webSchema.Identity {
	c := Get(ctx)
	if c == nil {
		return nil
	}
	return c.User
}

// GetUserId 获取用户ID
func GetUserId(ctx g.Ctx) int64 {
	user := GetUser(ctx)
	if user == nil {
		return 0
	}
	return user.Id
}

// SetData 设置额外数据
func SetData(ctx g.Ctx, key string, value interface{}) {
	c := Get(ctx)
	if c == nil {
		g.Log().Warning(ctx, "上下文对象不存在")
		return
	}
	c.Data[key] = value
}

// GetData 获得额外数据
func GetData(ctx g.Ctx, key string) g.Map {
	c := Get(ctx)
	if c == nil {
		return nil
	}
	return c.Data
}

// SetDataMap 设置额外数据
func SetDataMap(ctx g.Ctx, data g.Map) {
	c := Get(ctx)
	if c == nil {
		g.Log().Warning(ctx, "上下文对象不存在")
		return
	}
	c.Data = data
}
