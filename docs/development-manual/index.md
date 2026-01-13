# 开发手册

本章节将详细介绍如何使用 ZZFrame 进行开发，包括开发规范、控制器开发、服务层开发、数据库操作、配置管理等。

## 开发流程

ZZFrame 采用经典的 MVC 分层架构，开发流程如下：

```
定义 API → 实现 DAO 层 → 实现 Service 层 → 实现 Controller 层 → 配置路由
```

## 各层职责

| 层级 | 职责 | 包位置 |
|------|------|--------|
| API | 定义请求和响应数据结构 | `api/` |
| Controller | 处理 HTTP 请求、参数验证、调用 Service | `zcontroller/` |
| Service | 业务逻辑、事务管理、调用 DAO | `zservice/` |
| DAO | 数据库操作、SQL 查询 | `internal/dao/` |
| Model | 数据模型定义、常量定义 | `internal/model/` |

## 快速链接

- [开发规范](./development-standards)
- [Controller 开发](./controller)
- [Service 开发](./service)
- [数据库操作](./database)
- [配置管理](./configuration)
- [API 参考](./api)

## 代码示例

### API 定义

```go
type UserListReq struct {
    g.Meta `path:"/user/list" method:"get" tags:"用户管理"`
    Page   int    `json:"page" v:"required#页码不能为空"`
    Size   int    `json:"size" v:"required#每页数量不能为空"`
}
```

### Controller 实现

```go
type cUser struct{}

func (c *cUser) List(ctx g.Ctx, req *api.UserListReq) (res *api.UserListRes, err error) {
    return service.User().List(ctx, req)
}
```

### Service 实现

```go
type sUser struct{}

func (s *sUser) List(ctx g.Ctx, req *api.UserListReq) (res *api.UserListRes, err error) {
    // 业务逻辑
    list, total, err := dao.User.List(ctx, condition, args...)
    // ...
}
```

## 下一步

选择一个主题开始学习：

1. 了解 [开发规范](./development-standards) 掌握代码规范
2. 学习 [Controller 开发](./controller) 掌握控制器开发
3. 学习 [Service 开发](./service) 掌握服务层开发
4. 学习 [数据库操作](./database) 掌握数据库访问
