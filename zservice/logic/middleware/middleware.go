package middleware

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/denghuo98/zzframe/web/zcontext"
	"github.com/denghuo98/zzframe/web/ztoken"
	webSchema "github.com/denghuo98/zzframe/zschema/zweb"
	"github.com/denghuo98/zzframe/zservice"
)

type sMiddleware struct {
	LoginUrl string // 登录路由地址
}

func init() {
	zservice.RegisterMiddleware(NewMiddleware())
}

func NewMiddleware() *sMiddleware {
	return &sMiddleware{
		LoginUrl: "/admin/login",
	}
}

// Ctx 初始化上下文
func (s *sMiddleware) Ctx(r *ghttp.Request) {

	data := make(g.Map)
	data["request.body"] = gjson.New(r.GetBodyString())

	zcontext.Init(r, &webSchema.Context{
		Data: data,
	})

	r.SetCtx(r.GetNeverDoneCtx())
	r.Middleware.Next()
}

// CORS 允许跨域请求
func (s *sMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// // DeliverUserContext 将用户信息传递到上下文中
func (s *sMiddleware) DeliverUserContext(r *ghttp.Request) (err error) {
	user, err := ztoken.ParseLoginUser(r)
	if err != nil {
		return
	}

	zservice.CommonSite().BindUserContext(r.Context(), user)
	return
}
