# 系统配置

ZZFrame 使用 YAML 格式的配置文件，支持灵活的系统配置。

## 配置文件位置

配置文件默认位于项目根目录：

```
my-admin/
├── config.yaml          # 主配置文件
├── config.dev.yaml      # 开发环境配置
└── config.prod.yaml     # 生产环境配置
```

## 完整配置示例

```yaml
# 服务器配置
server:
  address: ":9090"           # 监听地址
  openapiPath: "/api.json"   # OpenAPI 文档地址
  swaggerPath: "/swagger"    # Swagger UI 地址

# 系统配置
system:
  mode: "develop"            # 运行模式：develop/production
  exceptAuth:                # 免认证接口
    - "/login"
    - "/logout"
    - "/admin/member/info"
    - "/admin/menu/dynamic"

  # 超级管理员配置
  superAdmin:
    password: "123456"       # 超级管理员密码

  # 缓存配置
  cache:
    adapter: "file"          # 缓存适配器：file/redis
    fileDir: "tmp/cache"     # 文件缓存目录

  # Token 配置
  token:
    expires: 86400           # Token 过期时间（秒）
    refreshInterval: 3600    # Token 刷新间隔（秒）
    maxRefreshTimes: 10      # Token 最大刷新次数
    secretKey: "zzframe"     # Token 密钥
    multiLogin: true         # 是否允许多端登录

# 日志配置
logger:
  path: "logs"               # 日志目录
  level: "all"               # 日志级别：all/debug/info/warning/error
  stdoutColorDisabled: false # 是否关闭终端颜色

# 数据库配置
database:
  logger:
    path: "logs/database"
    level: "all"
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/zzframe?loc=Local&parseTime=true&charset=utf8mb4"
    Prefix: "zz_"            # 表名前缀
    debug: true              # 是否打印 SQL

# 队列配置
queue:
  switch: true               # 是否启用队列
  driver: "disk"             # 队列驱动：disk/redis
  groupName: "default"
  disk:
    path: "./tmp/diskqueue"
    batchSize: 100
    batchTime: 1
    segmentSize: 10485760    # 10MB
    segmentLimit: 3000

# Redis 配置（可选）
redis:
  address: "127.0.0.1:6379"
  db: 0
  pass: ""
```

## 配置说明

### Server 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| address | 服务监听地址 | :9090 |
| openapiPath | OpenAPI 文档地址 | /api.json |
| swaggerPath | Swagger UI 地址 | /swagger |

### System 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| mode | 运行模式（develop/production） | develop |
| exceptAuth | 免认证接口列表 | [] |
| superAdmin.password | 超级管理员密码 | 123456 |
| cache.adapter | 缓存适配器（file/redis） | file |
| cache.fileDir | 文件缓存目录 | tmp/cache |
| token.expires | Token 过期时间（秒） | 86400 |
| token.refreshInterval | Token 刷新间隔（秒） | 3600 |
| token.maxRefreshTimes | Token 最大刷新次数 | 10 |
| token.secretKey | Token 密钥 | zzframe |
| token.multiLogin | 是否允许多端登录 | true |

### Logger 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| path | 日志目录 | logs |
| level | 日志级别 | all |
| stdoutColorDisabled | 是否关闭终端颜色 | false |

### Database 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| link | 数据库连接字符串 | - |
| Prefix | 表名前缀 | zz_ |
| debug | 是否打印 SQL | true |

## 多环境配置

支持为不同环境创建不同配置文件：

```yaml
# config.dev.yaml - 开发环境
system:
  mode: "develop"
database:
  default:
    link: "mysql:root:dev_password@tcp(127.0.0.1:3306)/zzframe_dev"
```

```yaml
# config.prod.yaml - 生产环境
system:
  mode: "production"
database:
  default:
    link: "mysql:root:prod_password@tcp(prod-db:3306)/zzframe_prod"
```

启动时指定配置文件：

```bash
# 开发环境
go run main.go server -c config.dev.yaml

# 生产环境
./server -c config.prod.yaml
```

## 环境变量

支持通过环境变量覆盖配置：

```bash
export SERVER_ADDRESS=:8080
export DATABASE_LINK="mysql:root:password@tcp(127.0.0.1:3306)/zzframe"
```

## 配置加载

框架启动时自动加载配置：

```go
func main() {
    var ctx = gctx.GetInitCtx()

    // 加载配置
    if err := zservice.SystemConfig().LoadConfig(ctx); err != nil {
        g.Log().Panicf(ctx, "初始化系统配置失败: %v", err)
    }

    zcmd.Main.Run(ctx)
}
```

## 配置验证

框架会在启动时验证配置的合法性：

- 数据库连接是否成功
- 必填配置项是否存在
- 配置格式是否正确

如果验证失败，框架会记录日志并退出。

## 配置规范

1. **敏感信息**：不要在配置文件中硬编码密码，使用环境变量
2. **多环境**：为不同环境创建独立的配置文件
3. **版本控制**：配置文件可以纳入版本控制，但要移除敏感信息
4. **文档化**：为自定义配置项添加注释说明
5. **热更新**：生产环境避免频繁修改配置

## 下一步

- [Controller 开发](../controller/)
- [Service 开发](../service/)
