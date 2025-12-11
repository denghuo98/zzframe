package zweb

type SystemConfig struct {
	Mode string `json:"mode" dc:"系统模式"` // dev | test | prod
}

type SuperAdminConfig struct {
	Password string `json:"password"`
}

type CacheConfig struct {
	Adapter string `json:"adapter" dc:"缓存适配器"`  // redis | file
	FileDir string `json:"fileDir" dc:"文件缓存目录"` // 文件缓存目录，当adapter为file时必填
}

// TokenConfig 登录令牌配置
type TokenConfig struct {
	SecretKey       string `json:"secretKey"`
	Expires         int64  `json:"expires"`
	AutoRefresh     bool   `json:"autoRefresh"`
	RefreshInterval int64  `json:"refreshInterval"`
	MaxRefreshTimes int64  `json:"maxRefreshTimes"`
	MultiLogin      bool   `json:"multiLogin"`
}
