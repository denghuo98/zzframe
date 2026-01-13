# API 参考

本章节提供 ZZFrame 核心模块的 API 参考。

## 后台管理 API

### 用户管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/member/list` | GET | 获取用户列表 |
| `/admin/member/info` | GET | 获取当前用户信息 |
| `/admin/member/edit` | POST | 新增/编辑用户 |
| `/admin/member/delete` | POST | 删除用户 |
| `/admin/member/update-password` | POST | 修改密码 |
| `/admin/member/update-profile` | POST | 更新个人资料 |

### 角色管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/role/list` | GET | 获取角色列表 |
| `/admin/role/edit` | POST | 新增/编辑角色 |
| `/admin/role/delete` | POST | 删除角色 |
| `/admin/role/info` | GET | 获取角色详情 |

### 菜单管理

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/menu/list` | GET | 获取菜单列表 |
| `/admin/menu/edit` | POST | 新增/编辑菜单 |
| `/admin/menu/delete` | POST | 删除菜单 |
| `/admin/menu/dynamic` | GET | 获取动态菜单 |

### 登录认证

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/login` | POST | 用户登录 |
| `/admin/logout` | POST | 用户登出 |
| `/admin/refresh-token` | POST | 刷新 Token |

## 系统模块 API

### 文件上传

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/upload/image` | POST | 上传图片 |
| `/admin/upload/file` | POST | 上传文件 |

### 操作日志

| 接口 | 方法 | 说明 |
|------|------|------|
| `/admin/sys-login-log/list` | GET | 获取登录日志 |

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin"
  }
}
```

### 失败响应

```json
{
  "code": 400,
  "message": "参数验证失败",
  "data": null
}
```

### 列表响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin"
      }
    ],
    "page": 1,
    "size": 10,
    "total": 100
  }
}
```

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 使用示例

### Go 调用示例

```go
import "github.com/gogf/gf/v2/frame/g"

// 调用用户列表接口
var result struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    struct {
        List  []interface{} `json:"list"`
        Page  int           `json:"page"`
        Size  int           `json:"size"`
        Total int64         `json:"total"`
    } `json:"data"`
}

resp, err := g.Client().
    SetHeader("Authorization", "Bearer "+token).
    Get(ctx, "http://localhost:9090/admin/member/list?page=1&size=10")

if err != nil {
    return err
}

err = resp.Scan(&result)
if err != nil {
    return err
}

if result.Code != 0 {
    return errors.New(result.Message)
}
```

### cURL 调用示例

```bash
# 用户登录
curl -X POST http://localhost:9090/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'

# 获取用户列表
curl http://localhost:9090/admin/member/list?page=1&size=10 \
  -H "Authorization: Bearer YOUR_TOKEN"

# 新增用户
curl -X POST http://localhost:9090/admin/member/edit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "username": "testuser",
    "real_name": "测试用户",
    "email": "test@example.com",
    "phone": "13800138000"
  }'
```

## Swagger 文档

项目启动后，可以通过以下地址访问 Swagger 文档：

```
http://localhost:9090/swagger
```

Swagger 文档提供了完整的 API 接口说明，可以直接在线测试。

## OpenAPI 文档

也可以直接下载 OpenAPI JSON 格式的文档：

```
http://localhost:9090/api.json
```

## 下一步

- 学习 [Controller 开发](./controller)
- 查看 [技术方案](../technical-solutions/)
