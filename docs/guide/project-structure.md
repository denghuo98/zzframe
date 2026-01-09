# 项目结构

理解 ZZFrame 的项目结构对于高效开发非常重要。

## 典型项目结构

```
my-admin/
├── main.go                    # 程序入口
├── config.yaml                # 配置文件
├── go.mod                     # Go 模块定义
├── go.sum                     # Go 模块依赖锁定
│
├── internal/                  # 私有应用代码（不会被外部引用）
│   ├── model/                # 数据模型
│   ├── logic/                # 业务逻辑
│   └── consts/               # 常量定义
│
├── api/                      # API 接口定义
│   ├── admin/               # 后台管理 API
│   │   ├── member.go
│   │   ├── role.go
│   │   └── menu.go
│   └── system/              # 系统 API
│
├── controller/              # 控制器层
│   ├── admin/
│   │   ├── member.go
│   │   ├── role.go
│   │   └── menu.go
│   └── system/
│
├── service/                 # 服务层
│   ├── admin/
│   │   ├── member.go
│   │   ├── role.go
│   │   └── menu.go
│   └── system/
│
├── router/                  # 路由配置
├── dao/                     # 数据访问对象
├── do/                      # 数据对象
│
├── logs/                    # 日志目录
├── tmp/                     # 临时文件
│   ├── cache/              # 缓存文件
│   └── diskqueue/          # 队列文件
│
└── docs/                    # 文档目录
```

## 目录说明

### 主入口

- **main.go**: 程序的入口文件，负责初始化系统配置和启动服务
- **config.yaml**: 系统配置文件，包含数据库、缓存、日志等配置

### API 层（api/）

定义请求和响应的数据结构：

```go
type MemberListReq struct {
    g.Meta `path:"/admin/member/list" method:"get"`
    Page    int    `json:"page" v:"required#页码不能为空"`
    Size    int    `json:"size" v:"required#每页数量不能为空"`
}

type MemberListRes struct {
    MemberListOutput
    PageRes
}
```

### Controller 层（controller/）

控制器层负责处理 HTTP 请求，调用 Service 层：

```go
type cMember struct{}

func (c *cMember) List(ctx g.Ctx, req *api.MemberListReq) (res *api.MemberListRes, err error) {
    list, totalCount, err := service.AdminMember().List(ctx, &req.MemberListInput)
    if err != nil {
        return nil, err
    }
    res = new(api.MemberListRes)
    res.MemberListOutput = *list
    res.PageRes.Pack(req, int(totalCount))
    return res, nil
}
```

### Service 层（service/）

服务层包含核心业务逻辑：

```go
type sAdminMember struct{}

func (s *sAdminMember) List(ctx g.Ctx, input *service.MemberListInput) (output *service.MemberListOutput, total int64, err error) {
    // 业务逻辑处理
    // ...
}
```

### Model 层（internal/model/）

定义数据模型和常量：

```go
const (
    TableNameAdminMember = "zz_admin_member"
)

type AdminMember struct {
    Id       uint64    `json:"id"`
    Username string    `json:"username"`
    // ...
}
```

### DAO 层（dao/）

数据访问对象，封装数据库操作：

```go
var AdminMember = adminMember{}

type adminMember struct{}

func (d *adminMember) List(ctx g.Context, condition string, args ...interface{}) ([]*model.AdminMember, error) {
    // 数据库查询
}
```

### DO 层（do/）

数据对象，用于与数据库交互：

```go
type AdminMember struct {
    g.Meta `orm:"table:zz_admin_member"`
    Id       uint64    `json:"id"`
    Username string    `json:"username"`
    // ...
}
```

## 命名规范

### 目录命名

- 使用小写字母
- 多个单词使用下划线分隔（如：admin_member）
- 复数形式表示集合（如：controllers、services）

### 文件命名

- 使用小写字母
- 多个单词使用下划线分隔（如：member_service.go）
- 测试文件以 `_test.go` 结尾

### 变量命名

- 驼峰命名法（如：memberList、userName）
- 私有变量小写开头（如：memberList）
- 公开变量大写开头（如：MemberList）

## 分层架构原则

ZZFrame 采用经典的分层架构：

```
请求 → Controller → Service → DAO → Database
    ↓            ↓
  参数验证    业务逻辑
    ↓            ↓
  响应处理    数据处理
```

### 各层职责

| 层级 | 职责 | 不应包含 |
|------|------|---------|
| Controller | 处理 HTTP 请求、参数验证、调用 Service | 复杂业务逻辑、直接数据库操作 |
| Service | 业务逻辑、事务管理、调用 DAO | HTTP 请求处理、参数验证 |
| DAO | 数据库操作、SQL 查询 | 业务逻辑、HTTP 响应 |
| Model | 数据模型定义、常量定义 | 业务逻辑 |

## 扩展项目

### 添加新功能模块

1. 在 `api/` 下定义接口
2. 在 `controller/` 下创建控制器
3. 在 `service/` 下实现业务逻辑
4. 在 `model/` 下定义模型
5. 在 `dao/` 下实现数据访问
6. 在 `router/` 下配置路由

### 大型项目组织方式

大型项目可按业务模块组织：

```
internal/
├── user/
│   ├── controller/
│   ├── service/
│   ├── model/
│   └── api/
├── order/
│   ├── controller/
│   ├── service/
│   ├── model/
│   └── api/
└── product/
    ├── controller/
    ├── service/
    ├── model/
    └── api/
```

## 下一步

了解项目结构后，继续阅读 [第一个应用](./first-application.md)。
