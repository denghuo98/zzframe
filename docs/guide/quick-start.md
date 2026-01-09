# 5分钟快速搭建

本指南将帮助你快速搭建一个完整的后台管理系统，包含用户管理、角色管理、菜单管理和登录日志功能。

## 准备工作

确保你已经完成 [安装配置](./installation.md) 中的步骤。

## 步骤 1：创建项目

```bash
mkdir my-admin
cd my-admin
go mod init my-admin
```

## 步骤 2：添加依赖

```bash
go get github.com/denghuo98/zzframe
```

## 步骤 3：创建配置文件

在项目根目录创建 `config.yaml`。

### 使用 SQLite（最小配置）

```yaml
server:
  address: ":9090"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"

system:
  mode: "develop"
  superAdmin:
    password: "admin123"

database:
  default:
    link: "sqlite::@file(./data/my_admin.db)"
    Prefix: "zz_"
```

### 使用 MySQL（完整配置）

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
    password: "admin123"
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
    link: "mysql:root:your_password@tcp(127.0.0.1:3306)/my_admin?loc=Local&parseTime=true&charset=utf8mb4"
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

### 使用 MySQL

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
    password: "admin123"
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
    link: "mysql:root:your_password@tcp(127.0.0.1:3306)/my_admin?loc=Local&parseTime=true&charset=utf8mb4"
    Prefix: "zz_"
```

## 步骤 4：创建入口文件

在项目根目录创建 `main.go`：

```go
package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/denghuo98/zzframe/zcmd"
	"github.com/denghuo98/zzframe/zservice"

	_ "github.com/denghuo98/zzframe/zservice/logic"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	var ctx = gctx.GetInitCtx()

	// 初始化系统配置
	if err := zservice.SystemConfig().LoadConfig(ctx); err != nil {
		g.Log().Panicf(ctx, "初始化系统配置失败: %v", err)
	}

	zcmd.Main.Run(ctx)
}
```

## 步骤 5：创建数据库

```sql
CREATE DATABASE IF NOT EXISTS my_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 步骤 6：启动服务

```bash
go run main.go server
```

服务启动后，你将看到类似的输出：

```
PID: 12345
  ADDRESS :9090
  MODE    develop
  OPENAPI /api.json
  SWAGGER /swagger
```

## 步骤 7：验证功能

### 访问系统

打开浏览器访问：http://localhost:9090/swagger

### 登录系统

使用以下账号登录：
- 用户名：`superAdmin`
- 密码：`admin123`（配置文件中设置的密码）

## 可用功能

✅ **用户管理**
- 查看用户列表
- 新增用户
- 编辑用户信息
- 删除用户
- 修改密码

✅ **角色管理**
- 查看角色列表
- 创建角色
- 分配权限
- 删除角色

✅ **菜单管理**
- 查看菜单树
- 创建菜单
- 编辑菜单
- 删除菜单

✅ **登录日志**
- 查看登录记录
- 查看登录状态
- 追踪用户活动

✅ **权限控制**
- 基于 RBAC 的权限模型
- 灵活的权限配置
- 接口级别的权限控制

## API 文档

启动服务后，可以通过以下方式查看 API 文档：

- **Swagger UI**: http://localhost:9090/swagger
- **OpenAPI JSON**: http://localhost:9090/api.json

## 常用命令

### 启动服务

```bash
go run main.go server
```

### 停止服务

按 `Ctrl+C` 或使用命令：

```bash
go run main.go stop
```

### 查看帮助

```bash
go run main.go help
```

## 下一步

恭喜！你已经成功搭建了一个完整的后台管理系统。接下来你可以：

- 深入了解 [项目结构](./project-structure.md)
- 查看 [核心概念](../concepts/)
- 阅读 [开发指南](../development/) 进行自定义开发
- 查看 [API 参考](../api-reference/) 了解所有可用的接口

## 常见问题

### 数据库连接失败

请检查 `config.yaml` 中的数据库连接信息是否正确。

### 端口被占用

修改 `config.yaml` 中的 `server.address`，改为其他端口，如 `:9091`。

### 找不到超级管理员账号

确保配置文件中 `system.superAdmin.password` 已设置，并且服务启动时没有报错。
