# 开发规范

遵循良好的开发规范可以提高代码质量和团队协作效率。

## 代码风格

### 命名规范

#### 包名

- 使用小写字母
- 简洁明了，不要过长
- 避免使用下划线或混合大小写

```go
// 好的包名
package admin
package service
package controller

// 不好的包名
package admin_service
package MyService
package adminService
```

#### 变量名

- 使用驼峰命名法
- 私有变量小写开头
- 公开变量大写开头

```go
// 好的变量名
var userList []User
var userName string
var memberService *MemberService

// 不好的变量名
var user_list []User
var user_name string
var memberservice *MemberService
```

#### 函数名

- 使用驼峰命名法
- 动词开头，描述函数行为

```go
// 好的函数名
func GetUserList() {}
func AddUser() {}
func UpdateUserInfo() {}

// 不好的函数名
func user_list() {}
func User() {}
func info() {}
```

#### 接口名

- 接口名以 `I` 开头
- 或者在 `er` 结尾

```go
// 好的接口名
type IUserService interface {}
type UserRepository interface {}

// 不好的接口名
type UserService interface {}
type IUser interface {}
```

### 注释规范

#### 包注释

```go
// Package admin 提供后台管理的控制器和服务
package admin
```

#### 函数注释

```go
// GetUserList 获取用户列表
// 支持按用户名、状态过滤，支持分页
//
// 参数:
//   ctx: 上下文
//   req: 查询请求参数
//
// 返回:
//   *UserListRes: 用户列表响应
//   error: 错误信息
func GetUserList(ctx g.Ctx, req *UserListReq) (*UserListRes, error) {
    // ...
}
```

#### 结构体注释

```go
// UserListReq 用户列表查询请求
type UserListReq struct {
    Page   int    `json:"page" v:"required#页码不能为空"`   // 页码
    Size   int    `json:"size" v:"required#每页数量不能为空"` // 每页数量
    Status int    `json:"status"`                             // 状态
    Name   string `json:"name"`                               // 用户名
}
```

## 项目结构规范

### 目录组织

```
my-admin/
├── api/              # API 定义
│   └── admin/
│       └── member.go
├── controller/       # 控制器层
│   └── admin/
│       └── member.go
├── service/          # 服务层
│   └── admin/
│       └── member.go
├── internal/         # 私有代码
│   ├── model/
│   │   ├── entity/
│   │   └── do/
│   └── dao/
└── main.go          # 入口文件
```

### 分层规范

| 层级 | 职责 | 不应包含 |
|------|------|---------|
| API | 数据结构定义 | 业务逻辑 |
| Controller | 请求处理 | 复杂业务逻辑、直接数据库操作 |
| Service | 业务逻辑 | HTTP 请求处理、参数验证 |
| DAO | 数据库操作 | 业务逻辑 |

## 编码规范

### 错误处理

```go
// 好的错误处理
func GetUser(ctx g.Ctx, id uint64) (*User, error) {
    user, err := dao.User.GetById(ctx, id)
    if err != nil {
        g.Log().Errorf(ctx, "获取用户失败: %v", err)
        return nil, err
    }
    return user, nil
}

// 不好的错误处理
func GetUser(ctx g.Ctx, id uint64) (*User, error) {
    user, _ := dao.User.GetById(ctx, id)
    return user, nil
}
```

### 日志记录

```go
// 好的日志记录
func EditUser(ctx g.Ctx, req *EditUserReq) error {
    g.Log().Infof(ctx, "开始编辑用户: %+v", req)
    
    if err := dao.User.Update(ctx, req); err != nil {
        g.Log().Errorf(ctx, "编辑用户失败: %v", err)
        return err
    }
    
    g.Log().Infof(ctx, "编辑用户成功: %d", req.Id)
    return nil
}

// 不好的日志记录
func EditUser(ctx g.Ctx, req *EditUserReq) error {
    if err := dao.User.Update(ctx, req); err != nil {
        return err
    }
    return nil
}
```

### 参数验证

```go
// 好的参数验证
type EditUserReq struct {
    Id       uint64 `json:"id" v:"required#用户ID不能为空"`
    Username string `json:"username" v:"required|length:3,20#用户名不能为空|用户名长度为3-20"`
    Email    string `json:"email" v:"email#邮箱格式不正确"`
}

// 不好的参数验证
type EditUserReq struct {
    Id       uint64 `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
```

## 性能规范

### 数据库查询

```go
// 好的做法 - 只查询需要的字段
func GetUserSimpleList(ctx g.Ctx) ([]*User, error) {
    var list []*User
    err := dao.User.Model(ctx).
        Fields("id", "username", "real_name").
        Scan(&list)
    return list, err
}

// 不好的做法 - 查询所有字段
func GetUserSimpleList(ctx g.Ctx) ([]*User, error) {
    var list []*User
    err := dao.User.Model(ctx).Scan(&list)
    return list, err
}
```

### 缓存使用

```go
// 好的做法 - 使用缓存
func GetUser(ctx g.Ctx, id uint64) (*User, error) {
    // 先查缓存
    var user User
    cacheKey := fmt.Sprintf("user:%d", id)
    
    if err := cache.Get(ctx, cacheKey, &user); err == nil && user.Id > 0 {
        return &user, nil
    }
    
    // 查数据库
    user, err := dao.User.GetById(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 设置缓存
    cache.Set(ctx, cacheKey, user, 300)
    
    return user, nil
}

// 不好的做法 - 不使用缓存
func GetUser(ctx g.Ctx, id uint64) (*User, error) {
    return dao.User.GetById(ctx, id)
}
```

## 测试规范

### 单元测试

```go
// 测试文件命名: xxx_test.go
func TestUserEdit(t *testing.T) {
    ctx := context.Background()
    
    req := &EditUserReq{
        Id:       1,
        Username: "test",
    }
    
    err := service.User().Edit(ctx, req)
    
    assert.NoError(t, err)
}
```

## Git 提交规范

### 提交信息格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

### 示例

```
feat(admin): 添加用户管理功能

- 添加用户列表查询
- 添加用户新增/编辑
- 添加用户删除

Closes #123
```

## 最佳实践

1. **保持简洁**: 代码应该简洁明了，避免过度设计
2. **注释适度**: 代码应该自解释，只对复杂逻辑添加注释
3. **错误处理**: 妥善处理所有可能的错误
4. **日志记录**: 记录关键操作和错误
5. **单元测试**: 为核心业务逻辑编写单元测试
6. **代码审查**: 提交代码前进行自我审查

## 下一步

- 学习 [Controller 开发](./controller)
- 了解 [Service 开发](./service)
