package zconsts

type CtxKey string

// ContextKey 上下文键
const (
	ContextHTTPKey        CtxKey = "httpContext" // http上下文
	ContextKeyCronArgsKey CtxKey = "cronArgs"    // 定时任务参数
	ContextKeyCronSn      CtxKey = "cronSn"      // 定时任务序列号
)
