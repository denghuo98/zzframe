# 项目配置

ZZFrame 提供了灵活的配置方式，支持 YAML、JSON、TOML 等多种配置格式。

## 配置文件位置

配置文件默认放置在项目根目录，命名规则如下：

- `config.yaml` - YAML 格式（推荐）
- `config.toml` - TOML 格式
- `config.json` - JSON 格式

## 最小配置

### SQLite 数据库（推荐用于快速开发）

```yaml
server:
  address: ":9090"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"

system:
  mode: "develop"
  superAdmin:
    password: "123456"

database:
  default:
    link: "sqlite::@file(./data/zzframe.db)"
    Prefix: "zz_"
```

### MySQL 数据库

```yaml
server:
  address: ":9090"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"

system:
  mode: "develop"
  superAdmin:
    password: "123456"

database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/zzframe?loc=Local&parseTime=true&charset=utf8mb4"
    Prefix: "zz_"
```

## 完整配置

```yaml
# 服务器配置
server:
  address: ":9090"           # 服务监听地址
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
  superAdmin:                # 超级管理员配置
    password: "123456"
  
  # 缓存配置
  cache:
    adapter: "file"          # 缓存适配器：file/redis
    fileDir: "tmp/cache"    # 文件缓存目录
  
  # Token 配置
  token:
    expires: 86400           # Token 过期时间（秒）
    refreshInterval: 3600    # Token 刷新间隔（秒）
    maxRefreshTimes: 10      # 最大刷新次数
    secretKey: "zzframe"     # Token 密钥
    multiLogin: true         # 是否允许多端登录

# 日志配置
logger:
  path: "logs"               # 日志目录
  level: "all"               # 日志级别：all/debug/info/warn/error/none

# 数据库配置
database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/zzframe?loc=Local&parseTime=true&charset=utf8mb4"
    Prefix: "zz_"

# 队列配置
queue:
  switch: true               # 是否启用队列
  driver: "disk"             # 队列驱动：disk/redis
  groupName: "default"       # 队列组名
  disk:
    path: "./tmp/diskqueue" # 磁盘队列目录
    batchSize: 100           # 批量处理数量
    batchTime: 1             # 批量处理时间（秒）
    segmentSize: 10485760    # 分段大小（字节）
    segmentLimit: 3000       # 分段数量限制
```

## 配置项说明

### Server 配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| address | 服务监听地址 | :9090 | 是 |
| openapiPath | OpenAPI 文档地址 | /api.json | 否 |
| swaggerPath | Swagger UI 地址 | /swagger | 否 |

### System 配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| mode | 运行模式 | develop | 是 |
| exceptAuth | 免认证接口列表 | [] | 否 |
| superAdmin.password | 超级管理员密码 | 123456 | 是 |
| cache.adapter | 缓存适配器 | file | 否 |
| cache.fileDir | 文件缓存目录 | tmp/cache | 否 |
| token.expires | Token 过期时间（秒） | 86400 | 否 |
| token.refreshInterval | Token 刷新间隔（秒） | 3600 | 否 |
| token.maxRefreshTimes | 最大刷新次数 | 10 | 否 |
| token.secretKey | Token 密钥 | zzframe | 否 |
| token.multiLogin | 是否允许多端登录 | true | 否 |

### Logger 配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| path | 日志目录 | logs | 否 |
| level | 日志级别 | all | 否 |

### Database 配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| link | 数据库连接字符串 | - | 是 |
| Prefix | 表名前缀 | zz_ | 否 |

#### SQLite 连接字符串

```
sqlite::@file(./data/zzframe.db)
```

#### MySQL 连接字符串

```
mysql:用户名:密码@tcp(地址:端口)/数据库?loc=Local&parseTime=true&charset=utf8mb4
```

### Queue 配置

| 配置项 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| switch | 是否启用队列 | true | 否 |
| driver | 队列驱动 | disk | 否 |
| groupName | 队列组名 | default | 否 |
| disk.path | 磁盘队列目录 | ./tmp/diskqueue | 否 |
| disk.batchSize | 批量处理数量 | 100 | 否 |
| disk.batchTime | 批量处理时间（秒） | 1 | 否 |
| disk.segmentSize | 分段大小（字节） | 10485760 | 否 |
| disk.segmentLimit | 分段数量限制 | 3000 | 否 |

## 环境变量

除了配置文件，你也可以通过环境变量覆盖配置：

```bash
# 设置服务器端口
export SERVER_PORT=9090

# 设置运行模式
export SYSTEM_MODE=production

# 设置数据库连接
export DB_LINK=mysql:root:password@tcp(127.0.0.1:3306)/zzframe
```

## 配置文件优先级

配置文件的加载优先级（从高到低）：

1. 命令行参数
2. 环境变量
3. 配置文件（config.yaml/config.toml/config.json）

## 开发环境 vs 生产环境

### 开发环境

```yaml
system:
  mode: "develop"

logger:
  level: "all"

database:
  default:
    link: "sqlite::@file(./data/zzframe.db)"
```

### 生产环境

```yaml
system:
  mode: "production"

logger:
  level: "error"

database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/zzframe?loc=Local&parseTime=true&charset=utf8mb4"
```

## 下一步

完成配置后，继续阅读 [第一个应用](./first-application)。
