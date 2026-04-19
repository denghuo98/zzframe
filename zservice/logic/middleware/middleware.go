package middleware

import (
	"strings"

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

	// 只对非 multipart 请求解析 body，避免文件上传时切片越界
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/") {
		data["request.body"] = gjson.New(r.GetBodyString())
	}

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

// DeliverUserContext 将用户信息传递到上下文中
func (s *sMiddleware) DeliverUserContext(r *ghttp.Request) (err error) {
	user, err := ztoken.ParseLoginUser(r)
	if err != nil {
		// 若解析失败，检查是否启用了匿名身份
		anonymousCfg, cfgErr := zservice.SystemConfig().GetAnonymousConfig(r.Context())
		if cfgErr == nil && anonymousCfg.Enabled {
			zcontext.SetUser(r.Context(), &anonymousCfg.Identity)
			return nil
		}
		return
	}

	zservice.CommonSite().BindUserContext(r.Context(), user)
	return
}
