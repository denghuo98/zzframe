# Controller 开发

Controller 层负责处理 HTTP 请求，包括参数验证、调用 Service 层、返回响应。

## 基本结构

```go
package admin

import (
    "github.com/gogf/gf/v2/frame/g"

    "github.com/denghuo98/zzframe/zschema/admin"
    "github.com/denghuo98/zzframe/zservice"
)

type cMember struct{}

var Member = cMember{}
```

## 常见方法模式

### 列表查询

```go
func (c *cMember) List(ctx g.Ctx, req *admin.MemberListReq) (res *admin.MemberListRes, err error) {
    output, total, err := zservice.AdminMember().List(ctx, req)
    if err != nil {
        return nil, err
    }
    
    res = new(admin.MemberListRes)
    res.List = output.List
    res.Page = req.Page
    res.Size = req.Size
    res.Total = total
    return res, nil
}
```

### 新增/编辑

```go
func (c *cMember) Edit(ctx g.Ctx, req *admin.MemberEditReq) (res *admin.MemberEditRes, err error) {
    err = zservice.AdminMember().Edit(ctx, req)
    if err != nil {
        return nil, err
    }
    return new(admin.MemberEditRes), nil
}
```

### 删除

```go
func (c *cMember) Delete(ctx g.Ctx, req *admin.MemberDeleteReq) (res *admin.MemberDeleteRes, err error) {
    err = zservice.AdminMember().Delete(ctx, req)
    if err != nil {
        return nil, err
    }
    return new(admin.MemberDeleteRes), nil
}
```

### 获取信息

```go
func (c *cMember) Info(ctx g.Ctx, req *admin.MemberInfoReq) (res *admin.MemberInfoRes, err error) {
    info, err := zservice.AdminMember().Info(ctx)
    if err != nil {
        return nil, err
    }
    res = new(admin.MemberInfoRes)
    res.MemberInfo = *info
    return res, nil
}
```

## 参数验证

### 在 API 定义中添加验证规则

```go
type UserAddReq struct {
    g.Meta      `path:"/user/add" method:"post" tags:"用户管理"`
    Username    string `json:"username" v:"required|length:3,20#用户名不能为空|用户名长度为3-20"`
    Password    string `json:"password" v:"required|min-length:6#密码不能为空|密码长度至少6位"`
    Email       string `json:"email" v:"email#邮箱格式不正确"`
    Phone       string `json:"phone" v:"phone#手机号格式不正确"`
}
```

### 常用验证规则

| 规则 | 说明 | 示例 |
|------|------|------|
| required | 必填 | `v:"required#字段不能为空"` |
| length | 长度范围 | `v:"length:3,20#长度为3-20"` |
| min-length | 最小长度 | `v:"min-length:6#最小长度为6"` |
| max-length | 最大长度 | `v:"max-length:50#最大长度为50"` |
| email | 邮箱格式 | `v:"email#邮箱格式不正确"` |
| phone | 手机号格式 | `v:"phone#手机号格式不正确"` |
| numeric | 数字 | `v:"numeric#必须是数字"` |
| gt | 大于 | `v:"gt:0#必须大于0"` |
| gte | 大于等于 | `v:"gte:0#必须大于等于0"` |

## 错误处理

```go
func (c *cMember) Edit(ctx g.Ctx, req *admin.MemberEditReq) (res *admin.MemberEditRes, err error) {
    // Service 层会处理业务逻辑错误
    err = zservice.AdminMember().Edit(ctx, req)
    if err != nil {
        // 返回错误，框架会自动处理响应
        return nil, err
    }
    return new(admin.MemberEditRes), nil
}
```

## 获取当前用户信息

```go
func (c *cMember) UpdateProfile(ctx g.Ctx, req *admin.MemberUpdateProfileReq) (res *admin.MemberUpdateProfileRes, err error) {
    // 从上下文中获取当前用户 ID
    userId := ctx.GetVar("user_id").Int64()
    req.Id = userId
    
    err = zservice.AdminMember().UpdateProfile(ctx, req)
    if err != nil {
        return nil, err
    }
    return new(admin.MemberUpdateProfileRes), nil
}
```

## 文件上传

```go
func (c *cUpload) Image(ctx g.Ctx, req *admin.UploadImageReq) (res *admin.UploadImageRes, err error) {
    // 从请求中获取文件
    file := req.File
    if file == nil {
        return nil, errors.New("请选择文件")
    }
    
    // 调用 Service 层处理上传
    url, err := zservice.Common().UploadImage(ctx, file)
    if err != nil {
        return nil, err
    }
    
    res = new(admin.UploadImageRes)
    res.Url = url
    return res, nil
}
```

## 最佳实践

1. **保持简洁**: Controller 层应该尽可能简洁，只负责接收请求和返回响应
2. **参数验证**: 在 API 定义中添加验证规则，让框架自动处理验证
3. **错误处理**: 直接返回 Service 层的错误，不要在 Controller 层捕获
4. **类型安全**: 使用 API 定义的类型，不要使用 map 或 interface
5. **命名规范**: Controller 方法名应该与 API 方法名对应

## 完整示例

```go
package admin

import (
    "github.com/gogf/gf/v2/frame/g"

    "github.com/denghuo98/zzframe/zschema/admin"
    "github.com/denghuo98/zzframe/zservice"
)

type cMember struct{}

var Member = cMember{}

// List 获取用户列表
func (c *cMember) List(ctx g.Ctx, req *admin.MemberListReq) (res *admin.MemberListRes, err error) {
    output, total, err := zservice.AdminMember().List(ctx, req)
    if err != nil {
        return nil, err
    }
    
    res = new(admin.MemberListRes)
    res.List = output.List
    res.Page = req.Page
    res.Size = req.Size
    res.Total = total
    return res, nil
}

// Edit 编辑用户
func (c *cMember) Edit(ctx g.Ctx, req *admin.MemberEditReq) (res *admin.MemberEditRes, err error) {
    err = zservice.AdminMember().Edit(ctx, req)
    if err != nil {
        return nil, err
    }
    return new(admin.MemberEditRes), nil
}

// Delete 删除用户
func (c *cMember) Delete(ctx g.Ctx, req *admin.MemberDeleteReq) (res *admin.MemberDeleteRes, err error) {
    err = zservice.AdminMember().Delete(ctx, req)
    if err != nil {
        return nil, err
    }
    return new(admin.MemberDeleteRes), nil
}

// Info 获取当前用户信息
func (c *cMember) Info(ctx g.Ctx, req *admin.MemberInfoReq) (res *admin.MemberInfoRes, err error) {
    info, err := zservice.AdminMember().Info(ctx)
    if err != nil {
        return nil, err
    }
    res = new(admin.MemberInfoRes)
    res.MemberInfo = *info
    return res, nil
}
```

## 下一步

- 学习 [Service 开发](./service)
- 了解 [数据库操作](./database)
