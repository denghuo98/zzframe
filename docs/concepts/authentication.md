# 认证机制

ZZFrame 提供了完善的用户认证机制，支持基于 Token 的认证方式。

## 认证流程

```
┌─────────┐       ┌─────────┐       ┌─────────┐
│  Client │       │  Server │       │  Cache  │
└────┬────┘       └────┬────┘       └────┬────┘
     │                  │                  │
     │  1. Login        │                  │
     │─────────────────>│                  │
     │                  │                  │
     │  2. Validate     │                  │
     │<─────────────────│                  │
     │                  │                  │
     │  3. Generate     │                  │
     │                  │  4. Save Token  │
     │                  │─────────────────>│
     │  5. Return Token │                  │
     │<─────────────────│                  │
     │                  │                  │
     │  6. API Request   │                  │
     │─────────────────>│                  │
     │  with Token      │                  │
     │                  │  7. Verify Token │
     │                  │─────────────────>│
     │                  │                  │
     │  8. Return Data   │                  │
     │<─────────────────│                  │
```

## 登录流程

### 1. 用户登录

用户提交用户名和密码：

```go
// 登录请求
type LoginReq struct {
    g.Meta `path:"/login" method:"post" tags:"认证"`
    Username string `json:"username" v:"required#用户名不能为空"`
    Password string `json:"password" v:"required#密码不能为空"`
    Captcha  string `json:"captcha" v:"required#验证码不能为空"`
}

// 登录处理
func (c *cAuth) Login(ctx g.Ctx, req *api.LoginReq) (res *api.LoginRes, err error) {
    // 1. 验证验证码
    if err := zcaptcha.Verify(req.Captcha); err != nil {
        return nil, err
    }

    // 2. 验证用户名密码
    member, err := dao.AdminMember.GetByUsername(ctx, req.Username)
    if err != nil {
        return nil, err
    }
    if member == nil {
        return nil, gerror.New("用户名或密码错误")
    }

    // 3. 生成 Token
    token, err := ztoken.Generate(ctx, member.Id, member.Username)
    if err != nil {
        return nil, err
    }

    // 4. 保存登录日志
    zqueue.PushLoginLog(ctx, member.Id, req.Username, "登录成功")

    return &api.LoginRes{Token: token}, nil
}
```

### 2. Token 生成

Token 包含用户基本信息和签名：

```go
type TokenClaims struct {
    UserId   uint64 `json:"userId"`
    Username string `json:"username"`
    jwt.StandardClaims
}

// 生成 Token
func Generate(ctx g.Context, userId uint64, username string) (string, error) {
    claims := &TokenClaims{
        UserId:   userId,
        Username: username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(Expires).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(SecretKey))
}
```

### 3. Token 验证

中间件自动验证 Token：

```go
func AuthMiddleware(r *ghttp.Request) {
    // 1. 跳过白名单
    if isExceptAuth(r.URL.Path) {
        r.Middleware.Next()
        return
    }

    // 2. 获取 Token
    token := r.Header.Get("Authorization")
    if token == "" {
        r.Response.WriteJsonExit(ghttp.Code{
            Code:    401,
            Message: "未授权",
        })
        return
    }

    // 3. 验证 Token
    claims, err := ztoken.Parse(token)
    if err != nil {
        r.Response.WriteJsonExit(ghttp.Code{
            Code:    401,
            Message: "Token 无效",
        })
        return
    }

    // 4. 设置用户上下文
    r.SetCtxVar("userId", claims.UserId)
    r.SetCtxVar("username", claims.Username)

    r.Middleware.Next()
}
```

## Token 配置

在 `config.yaml` 中配置 Token：

```yaml
system:
  token:
    # Token 过期时间（秒）
    expires: 86400
    # Token 刷新间隔（秒）
    refreshInterval: 3600
    # Token 最大刷新次数
    maxRefreshTimes: 10
    # Token 密钥
    secretKey: "zzframe"
    # 是否允许多端登录
    multiLogin: true
```

### 配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| expires | Token 过期时间（秒） | 86400（24小时） |
| refreshInterval | Token 刷新间隔（秒） | 3600（1小时） |
| maxRefreshTimes | Token 最大刷新次数 | 10 |
| secretKey | Token 签名密钥 | zzframe |
| multiLogin | 是否允许多端登录 | true |

## Token 刷新

Token 支持刷新机制：

```go
// 刷新 Token
func (c *cAuth) RefreshToken(ctx g.Ctx, req *api.RefreshTokenReq) (*api.RefreshTokenRes, error) {
    // 1. 验证旧 Token
    claims, err := ztoken.Parse(req.OldToken)
    if err != nil {
        return nil, err
    }

    // 2. 检查刷新次数
    if claims.RefreshCount >= MaxRefreshTimes {
        return nil, gerror.New("Token 刷新次数已达上限")
    }

    // 3. 生成新 Token
    newToken, err := ztoken.GenerateWithCount(ctx, claims.UserId, claims.Username, claims.RefreshCount+1)
    if err != nil {
        return nil, err
    }

    return &api.RefreshTokenRes{Token: newToken}, nil
}
```

## 多端登录

配置 `multiLogin: true` 时，用户可以在多个设备同时登录：

```yaml
system:
  token:
    multiLogin: true
```

如果需要限制单点登录，设置为 `false`，用户登录时会清除之前的 Token。

## 登录日志

系统自动记录用户登录日志：

```go
// 记录登录日志
type LoginLog struct {
    UserId      uint64    `json:"userId"`
    Username    string    `json:"username"`
    Ip          string    `json:"ip"`
    Location    string    `json:"location"`
    UserAgent   string    `json:"userAgent"`
    Status      string    `json:"status"`
    LoginTime   time.Time `json:"loginTime"`
}

// 异步写入日志
zqueue.PushLoginLog(ctx, userId, username, "登录成功")
```

## 登出

用户登出时清除 Token：

```go
func (c *cAuth) Logout(ctx g.Ctx, req *api.LogoutReq) (*api.LogoutRes, error) {
    // 1. 获取当前用户
    userId := ctx.Value("userId").(uint64)

    // 2. 清除 Token
    zcache.Delete(ctx, fmt.Sprintf("token:%d", userId))

    // 3. 记录登出日志
    zqueue.PushLoginLog(ctx, userId, username, "登出成功")

    return &api.LogoutRes{}, nil
}
```

## 白名单配置

有些接口不需要认证，可以配置到白名单：

```yaml
system:
  exceptAuth:
    - "/login"
    - "/logout"
    - "/admin/member/info"
    - "/admin/menu/dynamic"
```

## 验证码

登录时可以要求输入验证码：

```go
// 生成验证码
func (c *cAuth) Captcha(ctx g.Ctx, req *api.CaptchaReq) (*api.CaptchaRes, error) {
    id, image, err := zcaptcha.Generate()
    if err != nil {
        return nil, err
    }
    return &api.CaptchaRes{
        CaptchaId: id,
        Image:     image,
    }, nil
}
```

## 安全要求

1. **Token 过期时间**：根据业务需求设置
2. **密钥安全**：生产环境使用强密钥
3. **HTTPS**：生产环境使用 HTTPS 传输 Token
4. **Token 存储**：客户端使用 HttpOnly Cookie 或 LocalStorage
5. **刷新机制**：设置合理的刷新间隔
6. **登录日志**：记录登录日志

## 下一步

了解 [授权机制](./authorization.md)。
