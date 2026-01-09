# 安装配置

## 环境要求

在使用 ZZFrame 之前，请确保你的开发环境满足以下要求：

- **Go**: 1.24 或更高版本
- **MySQL**: 5.7 或更高版本（可选，不配置则使用 SQLite）
- **Redis**: 5.0 或更高版本（可选，用于缓存）

> **提示**：如果不配置 MySQL，框架会自动创建 SQLite 文件作为数据库，适合快速开发和测试。

## 安装框架

### 方式一：克隆项目

```bash
git clone https://github.com/denghuo98/zzframe.git
cd zzframe
go mod tidy
```

### 方式二：作为依赖引入

在你的项目中通过 `go get` 引入：

```bash
go get github.com/denghuo98/zzframe
```

## 配置文件

ZZFrame 使用 YAML 格式的配置文件（也支持 JSON、TOML）。在项目根目录创建 `config.yaml`：

### 最小配置（SQLite）

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

### 完整配置（MySQL）

```yaml
server:
  address: ":9090"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"

system:
  mode: "develop"
  exceptAuth:
    - "/login"
    - "/logout"
    - "/admin/member/info"
    - "/admin/menu/dynamic"
  superAdmin:
    password: "123456"
  cache:
    adapter: "file"
    fileDir: "tmp/cache"
  token:
    expires: 86400
    refreshInterval: 3600
    maxRefreshTimes: 10
    secretKey: "zzframe"
    multiLogin: true

logger:
  path: "logs"
  level: "all"

database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/zzframe?loc=Local&parseTime=true&charset=utf8mb4"
    Prefix: "zz_"

queue:
  switch: true
  driver: "disk"
  groupName: "default"
  disk:
    path: "./tmp/diskqueue"
    batchSize: 100
    batchTime: 1
    segmentSize: 10485760
    segmentLimit: 3000
```

### 完整配置示例（使用 SQLite）

```yaml
server:
  address: ":9090"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"

system:
  mode: "develop"
  exceptAuth:
    - "/login"
    - "/logout"
    - "/admin/member/info"
    - "/admin/menu/dynamic"
  superAdmin:
    password: "123456"
  cache:
    adapter: "file"
    fileDir: "tmp/cache"
  token:
    expires: 86400
    refreshInterval: 3600
    maxRefreshTimes: 10
    secretKey: "zzframe"
    multiLogin: true

logger:
  path: "logs"
  level: "all"

database:
  default:
    link: "sqlite::@file(./data/zzframe.db)"
    Prefix: "zz_"
```

### 完整配置示例（使用 MySQL）

```yaml
server:
  address: ":9090"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"

system:
  mode: "develop"
  exceptAuth:
    - "/login"
    - "/logout"
    - "/admin/member/info"
    - "/admin/menu/dynamic"
  superAdmin:
    password: "123456"
  cache:
    adapter: "file"
    fileDir: "tmp/cache"
  token:
    expires: 86400
    refreshInterval: 3600
    maxRefreshTimes: 10
    secretKey: "zzframe"
    multiLogin: true

logger:
  path: "logs"
  level: "all"

database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/zzframe?loc=Local&parseTime=true&charset=utf8mb4"
    Prefix: "zz_"
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
| mode | 系统模式（develop/production） | develop |
| exceptAuth | 不需要认证的接口列表 | [] |
| superAdmin.password | 超级管理员密码 | 123456 |
| cache.adapter | 缓存适配器（file/redis） | file |
| cache.fileDir | 文件缓存目录 | tmp/cache |
| token.expires | Token 过期时间（秒） | 86400 |
| token.refreshInterval | Token 刷新间隔（秒） | 3600 |
| token.maxRefreshTimes | Token 最大刷新次数 | 10 |
| token.secretKey | Token 密钥 | zzframe |
| token.multiLogin | 是否允许多端登录 | true |

### Database 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| link | 数据库连接字符串 | - |
| Prefix | 表名前缀 | zz_ |

#### SQLite 配置

```yaml
database:
  default:
    link: "sqlite::@file(./data/zzframe.db)"
    Prefix: "zz_"
```

SQLite 会自动在项目根目录下创建 `data/zzframe.db` 文件，无需手动创建数据库。

#### MySQL 配置

```yaml
database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/zzframe?loc=Local&parseTime=true&charset=utf8mb4"
    Prefix: "zz_"
```

需要提前创建数据库，框架会自动创建表结构。

## 数据库初始化

### 使用 SQLite（自动创建）

无需手动创建数据库，框架会自动在项目根目录下创建 `data/zzframe.db` 文件，并自动初始化表结构。

### 使用 MySQL（手动创建）

如果你选择使用 MySQL，需要先创建数据库：

```sql
CREATE DATABASE IF NOT EXISTS zzframe CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

框架启动后会自动创建必要的表结构，你也可以手动导入 SQL 文件：

```sql
-- 用户表
CREATE TABLE IF NOT EXISTS `zz_admin_member` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
  `dept_id` bigint(20) DEFAULT '0' COMMENT '部门ID',
  `real_name` varchar(32) DEFAULT '' COMMENT '真实姓名',
  `username` varchar(20) NOT NULL DEFAULT '' COMMENT '帐号',
  `password_hash` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(16) NOT NULL COMMENT '密码盐',
  `password_reset_token` varchar(150) DEFAULT '' COMMENT '密码重置令牌',
  `avatar` char(150) DEFAULT '' COMMENT '头像',
  `sex` tinyint(1) DEFAULT '3' COMMENT '性别',
  `email` varchar(60) DEFAULT '' COMMENT '邮箱',
  `mobile` varchar(20) DEFAULT '' COMMENT '手机号码',
  `last_active_at` datetime DEFAULT NULL COMMENT '最后活跃时间',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_username` (`username`),
  KEY `idx_phone` (`mobile`),
  KEY `idx_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_用户表';
```

### 首次启动初始化

框架首次启动时，会自动创建：
- 超级管理员账号：`superAdmin`
- 超级管理员密码：配置文件中 `system.superAdmin.password` 指定的密码
- 默认角色和权限
- 所有必要的数据库表

> **提示**：如果使用 SQLite，数据库文件会自动创建在 `data/zzframe.db`，无需额外操作。

## 下一步

完成配置后，请继续阅读 [5分钟快速搭建](./quick-start.md)。
