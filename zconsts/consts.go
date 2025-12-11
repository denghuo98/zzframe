package zconsts

const (
	DemoTips                  = "演示系统已隐藏" // 演示系统敏感数据打码
	NilJsonToString           = "{}"      // 空json初始化值
	RegionSpilt               = " / "     // 地区分隔符
	Unknown                   = "Unknown" // Unknown
	SuperRoleKey              = "super"   // 超管角色唯一标识符，通过角色验证超管
	DefaultSuperAdminUsername = "superAdmin"
	DefaultSuperAdminPassword = "zzframe@admin123"
	MaxServeLogContentLen     = 2048 // 最大保留服务日志内容大小
)

// RequestEncryptKey
// 请求加密密钥用于敏感数据加密，16位字符，前后端需保持一致
// 安全起见，生产环境运行时请注意使用配置文件修改
var RequestEncryptKey = []byte("zzframe@admin123")
