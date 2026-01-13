# 配置管理

ZZFrame 提供了灵活的配置管理机制，支持多种配置方式和环境。

## 配置获取

```go
import "github.com/denghuo98/zzframe/zservice"

// 获取系统配置实例
config := zservice.SystemConfig()
```

## 常用配置项

### 服务器配置

```go
address := config.Server().Address
openapiPath := config.Server().OpenapiPath
swaggerPath := config.Server().SwaggerPath
```

### 系统配置

```go
mode := config.System().Mode
exceptAuth := config.System().ExceptAuth
superAdminPassword := config.System().SuperAdmin.Password
```

### Token 配置

```go
expires := config.System().Token.Expires
refreshInterval := config.System().Token.RefreshInterval
secretKey := config.System().Token.SecretKey
multiLogin := config.System().Token.MultiLogin
```

### 缓存配置

```go
adapter := config.System().Cache.Adapter
fileDir := config.System().Cache.FileDir
```

### 日志配置

```go
path := config.Logger().Path
level := config.Logger().Level
```

### 数据库配置

```go
link := config.Database().Default.Link
prefix := config.Database().Default.Prefix
```

## 动态配置

### 获取缓存实例

```go
cache, err := config.Cache()
if err != nil {
    return err
}

// 使用缓存
cache.Set(ctx, "key", "value", 300)
cache.Get(ctx, "key", &result)
```

### 获取数据库实例

```go
db := config.DB()
db.Model(...).Where(...).Scan(...)
```

## 配置文件读取

### 读取自定义配置

```yaml
# config.yaml
custom:
  feature:
    enabled: true
    settings:
      timeout: 30
```

```go
import "github.com/gogf/gf/v2/frame/g"

// 读取自定义配置
enabled := g.Cfg().MustGet(ctx, "custom.feature.enabled").Bool()
timeout := g.Cfg().MustGet(ctx, "custom.feature.settings.timeout").Int()
```

## 环境变量

### 读取环境变量

```go
import (
    "os"
    "github.com/gogf/gf/v2/os/genv"
)

// 方式1：使用 os 包
port := os.Getenv("SERVER_PORT")

// 方式2：使用 GoFrame
port := genv.Get("SERVER_PORT")
```

### 设置默认值

```go
import "github.com/gogf/gf/v2/os/genv"

// 如果环境变量不存在，使用默认值
port := genv.Get("SERVER_PORT", "9090")
mode := genv.Get("SYSTEM_MODE", "develop")
```

## 配置覆盖

ZZFrame 的配置加载优先级：

1. 命令行参数
2. 环境变量
3. 配置文件

### 命令行参数覆盖

```bash
# 启动时指定配置
./server --server.address=9090 --system.mode=production
```

### 环境变量覆盖

```bash
# 设置环境变量后启动
export SERVER_PORT=9090
export SYSTEM_MODE=production
./server
```

## 配置热更新

```go
import "github.com/gogf/gf/v2/frame/g"

// 监听配置文件变化
g.Cfg().AddChangeListener(func(event gcfg.Event) {
    fmt.Println("配置文件已更新:", event)
})
```

## 最佳实践

1. **统一管理**: 所有配置通过 `zservice.SystemConfig()` 获取
2. **环境区分**: 使用环境变量区分开发/生产环境
3. **敏感信息**: 不要将密码等敏感信息硬编码到代码中
4. **默认值**: 为配置设置合理的默认值
5. **类型安全**: 使用类型化的配置获取方法

## 下一步

- 了解 [API 参考](./api)
- 查看 [技术方案](../technical-solutions/)
