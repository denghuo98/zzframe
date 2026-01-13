# 认证授权

ZZFrame 采用 JWT (JSON Web Token) 进行用户认证，提供安全、高效的身份验证机制。

## 认证流程

### 登录流程

```
1. 用户提交账号密码
   ↓
2. 验证用户信息
   ↓
3. 生成 JWT Token
   ↓
4. 返回 Token 给客户端
   ↓
5. 客户端存储 Token
```

### 请求认证流程

```
1. 客户端发送请求携带 Token
   ↓
2. 中间件拦截请求
   ↓
3. 验证 Token 有效性
   ↓
4. 解析 Token 获取用户信息
   ↓
5. 将用户信息放入上下文
   ↓
6. 继续处理请求
```

## Token 结构

### Token Payload

```json
{
  "user_id": 1,
  "username": "admin",
  "roles": ["super_admin"],
  "exp": 1640995200,
  "iat": 1640908800
}
```

### Token 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | int | 用户 ID |
| username | string | 用户名 |
| roles | []string | 角色列表 |
| exp | int | 过期时间 |
| iat | int | 签发时间 |

## 配置说明

```yaml
system:
  token:
    expires: 86400           # Token 过期时间（秒），默认 24 小时
    refreshInterval: 3600    # Token 刷新间隔（秒），默认 1 小时
    maxRefreshTimes: 10      # 最大刷新次数
    secretKey: "zzframe"     # Token 密钥
    multiLogin: true         # 是否允许多端登录
```

## 认证中间件

### 中间件实现

```go
package middleware

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/golang-jwt/jwt/v5"
)

func AdminAuth(r *ghttp.Request) {
    // 获取 Token
    token := r.Header.Get("Authorization")
    if token == "" {
        r.Response.WriteJson(g.Map{
            "code":    401,
            "message": "未登录",
        })
        r.Exit()
    }
    
    // 去掉 Bearer 前缀
    token = strings.Replace(token, "Bearer ", "", 1)
    
    // 验证 Token
    claims, err := verifyToken(token)
    if err != nil {
        r.Response.WriteJson(g.Map{
            "code":    401,
            "message": "Token 无效",
        })
        r.Exit()
    }
    
    // 将用户信息放入上下文
    r.SetCtxVar("user_id", claims.UserId)
    r.SetCtxVar("username", claims.Username)
    r.SetCtxVar("roles", claims.Roles)
    
    // 继续处理请求
    r.Middleware.Next()
}
```

### Token 验证

```go
func verifyToken(tokenString string) (*Claims, error) {
    secretKey := zservice.SystemConfig().System().Token.SecretKey
    
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("token invalid")
}
```

### Claims 定义

```go
type Claims struct {
    UserId   uint64   `json:"user_id"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}
```

## Token 生成

### 登录生成 Token

```go
func (s *sAdmin) Login(ctx g.Ctx, req *LoginReq) (res *LoginRes, err error) {
    // 验证用户
    member, err := dao.AdminMember.GetByUsername(ctx, req.Username)
    if err != nil {
        return nil, errors.New("用户不存在")
    }
    
    // 验证密码
    if !verifyPassword(member.Password, req.Password) {
        return nil, errors.New("密码错误")
    }
    
    // 生成 Token
    token, err := generateToken(member)
    if err != nil {
        return nil, err
    }
    
    // 返回 Token
    res = &LoginRes{
        Token: token,
        User: member,
    }
    
    return res, nil
}

func generateToken(member *AdminMember) (string, error) {
    config := zservice.SystemConfig()
    expires := config.System().Token.Expires
    secretKey := config.System().Token.SecretKey
    
    claims := &Claims{
        UserId:   member.Id,
        Username: member.Username,
        Roles:    getMemberRoles(member.Id),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expires) * time.Second)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}
```

## Token 刷新

### 刷新 Token

```go
func RefreshToken(oldToken string) (string, error) {
    // 验证旧 Token
    claims, err := verifyToken(oldToken)
    if err != nil {
        return "", err
    }
    
    // 检查是否可以刷新
    if !canRefreshToken(claims) {
        return "", errors.New("token 已过期，请重新登录")
    }
    
    // 生成新 Token
    config := zservice.SystemConfig()
    expires := config.System().Token.Expires
    secretKey := config.System().Token.SecretKey
    
    newClaims := &Claims{
        UserId:   claims.UserId,
        Username: claims.Username,
        Roles:    claims.Roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expires) * time.Second)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
    return token.SignedString([]byte(secretKey))
}

func canRefreshToken(claims *Claims) bool {
    // 检查是否在刷新时间内
    now := time.Now()
    refreshInterval := time.Duration(zservice.SystemConfig().System().Token.RefreshInterval) * time.Second
    
    return now.Sub(claims.IssuedAt.Time) < refreshInterval
}
```

## 多端登录

### 支持多端登录

```yaml
system:
  token:
    multiLogin: true  # 允许多端登录
```

### 禁用多端登录（单点登录）

```yaml
system:
  token:
    multiLogin: false  # 禁止多端登录
```

当 `multiLogin: false` 时，每次登录会生成新的 Token，旧的 Token 会失效。

## 安全性

### Token 安全建议

1. **使用 HTTPS**: 在生产环境中使用 HTTPS 传输 Token
2. **设置合理的过期时间**: Token 过期时间不宜过长
3. **定期更换密钥**: 定期更换 Token 密钥
4. **存储安全**: 客户端安全存储 Token（如 HttpOnly Cookie）
5. **验证 Token 每次请求**: 每次请求都验证 Token 有效性

### 密码安全

```go
// 密码加密
func hashPassword(password, salt string) string {
    return gmd5.MustEncryptString(password + salt)
}

// 验证密码
func verifyPassword(hashedPassword, password, salt string) bool {
    return hashedPassword == hashPassword(password, salt)
}
```

## 使用示例

### 登录

```bash
curl -X POST http://localhost:9090/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "real_name": "超级管理员"
    }
  }
}
```

### 带Token请求

```bash
curl http://localhost:9090/admin/member/list \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## 最佳实践

1. **Token 存储**: 推荐使用 HttpOnly Cookie 存储 Token，防止 XSS 攻击
2. **刷新机制**: 合理设置刷新间隔，避免频繁刷新
3. **错误处理**: Token 失效时引导用户重新登录
4. **单点登录**: 根据业务需求选择是否支持多端登录
5. **日志记录**: 记录登录和登出日志

## 下一步

- 学习 [权限管理](./authorization)
- 了解 [架构设计](./architecture)
