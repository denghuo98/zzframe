package middleware

import (
	"github.com/denghuo98/zzframe/web/zresp"
	"github.com/denghuo98/zzframe/zconsts"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gmeta"
)

// ResponseHandler HTTP 响应预处理
func (s *sMiddleware) ResponseHandler(r *ghttp.Request) {

	r.Middleware.Next()

	// 错误状态码接管
	switch r.Response.Status {
	case 403:
		r.Response.Writeln("403 - 网站拒绝显示此网页")
		return
	case 404:
		r.Response.Writeln("404 - 网站不存在")
		return
	}

	// 已存在响应
	contentType := getContentType(r)
	if r.Response.BufferLength() > 0 && contentType != "text/event-stream" {
		return
	}

	s.responseJson(r)

}

func getContentType(r *ghttp.Request) (contentType string) {
	contentType = r.Response.Header().Get("Content-Type")
	if contentType != "" {
		return
	}

	mime := gmeta.Get(r.GetHandlerResponse(), "mime").String()
	if mime == "" {
		contentType = "application/json"
	} else {
		contentType = mime
	}
	return
}

// responseJson json响应
func (s *sMiddleware) responseJson(r *ghttp.Request) {
	code, message, data := parseResponse(r)
	zresp.RJson(r, code, message, data)
}

// parseResponse 解析响应数据
func parseResponse(r *ghttp.Request) (code int, message string, resp interface{}) {
	ctx := r.Context()
	err := r.GetError()
	if err == nil {
		return gcode.CodeOK.Code(), "操作成功", r.GetHandlerResponse()
	}

	message = zconsts.ErrorMessage(gerror.Current(err))

	code = gerror.Code(err).Code()

	// 记录异常日志
	// 如果你想对错误做不同的处理，可以通过定义不同的错误码来区分
	// 默认-1为安全可控错误码只记录文件日志，非-1为不可控错误，记录文件日志+服务日志并打印堆栈
	if code == gcode.CodeNil.Code() {
		g.Log().Stdout(false).Infof(ctx, "exception:%v", err)
	} else {
		g.Log().Errorf(ctx, "exception:%v", err)
	}
	return
}
