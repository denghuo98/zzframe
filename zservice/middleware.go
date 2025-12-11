package zservice

import "github.com/gogf/gf/v2/net/ghttp"

type IMiddleware interface {
	Ctx(r *ghttp.Request)
	CORS(r *ghttp.Request)
	ResponseHandler(r *ghttp.Request)
	AdminAuth(r *ghttp.Request)
}

var (
	localMiddleware IMiddleware
)

func Middleware() IMiddleware {
	if localMiddleware == nil {
		panic("Middleware is not initialized, please register it first")
	}
	return localMiddleware
}

func RegisterMiddleware(i IMiddleware) {
	localMiddleware = i
}
