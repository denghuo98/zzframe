# 架构设计

ZZFrame 采用了清晰的分层架构和模块化设计，便于理解、维护和扩展。

## 分层架构

```
┌─────────────────────────────────────────────┐
│              Presentation Layer             │
│              (Controller/API)                │
│  ┌─────────────────────────────────────────┐ │
│  │   Router   →   Controller   →   API    │ │
│  └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────┐
│              Business Layer                 │
│              (Service)                       │
│  ┌─────────────────────────────────────────┐ │
│  │         Business Logic                 │ │
│  │         Transaction Management          │ │
│  │         Data Validation                │ │
│  └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────┐
│            Data Access Layer                │
│              (DAO/Model)                    │
│  ┌─────────────────────────────────────────┐ │
│  │   SQL   ←   ORM   ←   Data Mapping     │ │
│  └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────┐
│              Data Storage                   │
│         (MySQL/Redis/Cache)                 │
└─────────────────────────────────────────────┘
```

## 各层职责

### Presentation Layer（表现层）

**Controller 层**：
- 处理 HTTP 请求和响应
- 参数验证和绑定
- 调用 Service 层处理业务
- 返回统一格式的响应

```go
func (c *cMember) List(ctx g.Ctx, req *api.MemberListReq) (res *api.MemberListRes, err error) {
    // 调用 Service 层
    list, totalCount, err := service.AdminMember().List(ctx, &req.MemberListInput)
    if err != nil {
        return nil, err
    }
    // 组装响应
    res = new(api.MemberListRes)
    res.MemberListOutput = *list
    res.PageRes.Pack(req, int(totalCount))
    return res, nil
}
```

**API 层**：
- 定义请求和响应的数据结构
- 提供接口文档
- 参数验证规则

```go
type MemberListReq struct {
    g.Meta `path:"/admin/member/list" method:"get" tags:"用户管理"`
    Page   int    `json:"page" v:"required#页码不能为空"`
    Size   int    `json:"size" v:"required#每页数量不能为空"`
}
```

### Business Layer（业务层）

**Service 层**：
- 实现核心业务逻辑
- 管理事务
- 数据验证和处理
- 调用 DAO 层访问数据

```go
func (s *sAdminMember) List(ctx context.Context, input *service.MemberListInput) (*service.MemberListOutput, int64, error) {
    // 业务逻辑
    condition := s.buildCondition(input)
    list, err := dao.AdminMember.List(ctx, condition)
    if err != nil {
        return nil, 0, err
    }
    // 返回结果
    return &service.MemberListOutput{List: list}, count, nil
}
```

### Data Access Layer（数据访问层）

**DAO 层**：
- 封装数据库操作
- 执行 SQL 查询
- 数据持久化

```go
func (d *adminMemberDao) List(ctx context.Context, condition string, args ...interface{}) ([]*model.AdminMember, error) {
    return AdminMember.Model(ctx).Where(condition, args...).All()
}
```

**Model 层**：
- 定义数据模型
- 数据结构映射

## 模块化设计

框架按功能模块组织代码：

```
zzframe/
├── zapi/                 # API 定义模块
│   ├── admin/           # 后台管理 API
│   ├── system/          # 系统 API
│   └── common/          # 通用 API
│
├── zcontroller/          # 控制器模块
│   ├── admin/           # 后台管理控制器
│   ├── system/          # 系统控制器
│   └── common/          # 通用控制器
│
├── zservice/             # 服务层模块
│   ├── admin/           # 后台管理服务
│   ├── system/          # 系统服务
│   └── common/          # 通用服务
│
├── web/                  # Web 功能模块
│   ├── ztoken/          # Token 管理
│   ├── zcache/          # 缓存管理
│   ├── zcaptcha/        # 验证码
│   └── zqueue/          # 队列
│
└── zschema/             # 数据模型模块
    ├── admin/           # 后台管理模型
    ├── system/          # 系统模型
    └── common/          # 通用模型
```

## 核心组件

### 1. 路由系统

- 自动路由注册
- 支持分组路由
- 中间件支持
- 参数绑定

### 2. 中间件系统

- 认证中间件
- 权限中间件
- 日志中间件
- 自定义中间件

### 3. 认证授权

- 基于 Token 的认证
- Casbin RBAC 权限控制
- 动态权限配置

### 4. 缓存系统

- 文件缓存
- Redis 缓存
- 统一缓存接口

### 5. 队列系统

- 磁盘队列
- 异步任务处理
- 消息持久化

### 6. 日志系统

- 结构化日志
- 日志分级
- 文件切割
- 异步写入

## 依赖注入

框架使用依赖注入的方式管理组件：

```go
// Controller 注入 Service
type cMember struct{}

func (c *cMember) List(ctx g.Ctx, req *api.MemberListReq) (*api.MemberListRes, error) {
    // 通过 service 包获取单例
    return service.AdminMember().List(ctx, &req.MemberListInput)
}
```

## 错误处理

框架提供统一的错误处理机制：

```go
// 统一错误响应
if err != nil {
    return nil, gerror.New("操作失败")
}

// 错误码
const (
    CodeSuccess       = 0
    CodeInvalidParam  = 1
    CodeUnauthorized  = 2
    CodeForbidden     = 3
    CodeNotFound      = 4
    CodeInternalError = 5
)
```

## 事务管理

框架支持声明式事务：

```go
// 使用事务
return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
    // 业务操作
    if err := dao.Member.Update(ctx, member); err != nil {
        return err
    }
    // 返回 nil 提交事务，返回 error 回滚事务
    return nil
})
```

## 开发规范

1. **分层清晰**：各层职责明确，不要越层调用
2. **依赖注入**：通过接口解耦
3. **统一错误处理**：使用框架提供的错误处理机制
4. **缓存使用**：热点数据需要缓存
5. **事务管理**：事务在 Service 层管理
6. **接口文档**：使用 Swagger 自动生成文档

## 下一步

了解 [认证机制](./authentication.md) 和 [授权机制](./authorization.md)。
