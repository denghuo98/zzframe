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


## 数据库初始化

### 使用 SQLite（自动创建）

无需手动创建数据库，框架会自动在项目根目录下创建 `data/zzframe.db` 文件，并自动初始化表结构。

### 使用 MySQL（手动创建）

如果你选择使用 MySQL，框架会自动访问数据库判断所依赖的数据库表是否存在，不存在则自动创建：


### 首次启动初始化

框架首次启动时，会自动创建：
- 超级管理员账号：`superAdmin`
- 超级管理员密码：配置文件中 `system.superAdmin.password` 指定的密码
- 默认角色和权限
- 所有必要的数据库表

> **提示**：如果使用 SQLite，数据库文件会自动创建在 `data/zzframe.db`，无需额外操作。

## 下一步

想要了解更多配置细节，请继续阅读 [项目配置](./configuration.md)。
