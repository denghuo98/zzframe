# 日志方案

ZZFrame 提供了完善的日志系统，支持结构化日志、日志分级、日志轮转等功能。

## 日志级别

| 级别 | 说明 | 使用场景 |
|------|------|----------|
| all | 所有日志 | 调试环境 |
| debug | 调试日志 | 开发调试 |
| info | 信息日志 | 一般信息 |
| warn | 警告日志 | 警告信息 |
| error | 错误日志 | 错误信息 |
| none | 不记录日志 | 生产环境关闭 |

## 配置

```yaml
logger:
  path: "logs"           # 日志目录
  level: "all"           # 日志级别
```

### 开发环境配置

```yaml
logger:
  path: "logs"
  level: "all"
```

### 生产环境配置

```yaml
logger:
  path: "logs"
  level: "error"
```

## 基本用法

### 导入日志

```go
import "github.com/gogf/gf/v2/frame/g"
```

### 记录日志

```go
// Debug 日志
g.Log().Debugf(ctx, "调试信息: %v", data)

// Info 日志
g.Log().Infof(ctx, "用户登录: %s", username)

// Warn 日志
g.Log().Warningf(ctx, "缓存读取失败，使用默认值")

// Error 日志
g.Log().Errorf(ctx, "数据库查询失败: %v", err)
```

### 不同级别日志

```go
func ProcessUser(ctx g.Ctx, user *User) error {
    g.Log().Debugf(ctx, "开始处理用户: %+v", user)
    
    // 业务逻辑
    if user.Id == 0 {
        g.Log().Warningf(ctx, "用户 ID 为空: %+v", user)
        return errors.New("用户 ID 不能为空")
    }
    
    if err := saveUser(ctx, user); err != nil {
        g.Log().Errorf(ctx, "保存用户失败: %v, 用户: %+v", err, user)
        return err
    }
    
    g.Log().Infof(ctx, "用户处理成功: %d", user.Id)
    return nil
}
```

## 结构化日志

### 使用字段

```go
g.Log().Infof(ctx, `{"userId": %d, "action": "login", "ip": "%s"}`, userId, ip)
```

### 使用 JSON 格式

```go
type LogFields struct {
    UserId int64  `json:"userId"`
    Action string `json:"action"`
    IP     string `json:"ip"`
}

fields := LogFields{
    UserId: userId,
    Action: "login",
    IP:     ip,
}

g.Log().Infof(ctx, "%+v", fields)
```

## 操作日志

### 记录操作日志

```go
func (s *sAdminMember) Edit(ctx g.Ctx, in *admin.MemberEditInput) error {
    // 获取操作用户
    operatorId := ctx.GetVar("user_id").Int64()
    
    // 记录操作日志
    g.Log().Infof(ctx, "编辑用户: 操作人=%d, 操作对象=%d, 操作内容=%+v", 
        operatorId, in.Id, in)
    
    // 业务逻辑
    if err := dao.AdminMember.Update(ctx, &do.AdminMember{
        Id:       in.Id,
        Username: in.Username,
        // ...
    }); err != nil {
        g.Log().Errorf(ctx, "编辑用户失败: %v", err)
        return err
    }
    
    return nil
}
```

### 记录到数据库

```go
func SaveOperationLog(ctx g.Ctx, userId int64, action, content string) error {
    log := &do.SysLoginLog{
        UserId:  userId,
        Action:  action,
        Content: content,
        IP:      ctx.GetVar("ip").String(),
        Time:    time.Now(),
    }
    
    _, err := dao.SysLoginLog.Insert(ctx, log)
    return err
}

// 使用
SaveOperationLog(ctx, userId, "编辑用户", "编辑用户ID: 1")
```

## 登录日志

### 记录登录日志

```go
func (s *sAdmin) Login(ctx g.Ctx, req *LoginReq) (res *LoginRes, err error) {
    // 验证用户
    member, err := dao.AdminMember.GetByUsername(ctx, req.Username)
    if err != nil {
        // 记录登录失败
        g.Log().Warningf(ctx, "用户登录失败: 用户名不存在, IP=%s", getClientIP(ctx))
        return nil, errors.New("用户名或密码错误")
    }
    
    // 验证密码
    if !verifyPassword(member.Password, req.Password) {
        // 记录登录失败
        g.Log().Warningf(ctx, "用户登录失败: 密码错误, 用户=%d, IP=%s", member.Id, getClientIP(ctx))
        return nil, errors.New("用户名或密码错误")
    }
    
    // 生成 Token
    token, err := generateToken(member)
    if err != nil {
        g.Log().Errorf(ctx, "生成 Token 失败: %v", err)
        return nil, err
    }
    
    // 记录登录日志
    if err := recordLoginLog(ctx, member.Id, true); err != nil {
        g.Log().Errorf(ctx, "记录登录日志失败: %v", err)
    }
    
    return &LoginRes{Token: token}, nil
}

func recordLoginLog(ctx g.Ctx, userId int64, success bool) error {
    log := &do.SysLoginLog{
        UserId:   userId,
        Status:   1,
        IP:       getClientIP(ctx),
        Location: getLocation(ctx),
        Time:     time.Now(),
    }
    
    if !success {
        log.Status = 0
    }
    
    return dao.SysLoginLog.Insert(ctx, log)
}
```

## 错误日志

### 记录错误日志

```go
func ProcessData(ctx g.Ctx) error {
    data, err := fetchFromAPI(ctx)
    if err != nil {
        // 记录错误详情
        g.Log().Errorf(ctx, "获取 API 数据失败: %v", err)
        g.Log().Stack(ctx) // 打印堆栈
        return err
    }
    
    // 处理数据...
    return nil
}
```

### 错误跟踪

```go
func SafeProcess(ctx g.Ctx) {
    defer func() {
        if r := recover(); r != nil {
            g.Log().Errorf(ctx, "发生 panic: %v", r)
            g.Log().Stack(ctx)
        }
    }()
    
    // 可能发生 panic 的代码
    doSomething()
}
```

## 性能日志

### 记录执行时间

```go
func ProcessWithTiming(ctx g.Ctx) error {
    start := time.Now()
    defer func() {
        g.Log().Infof(ctx, "处理完成, 耗时: %v", time.Since(start))
    }()
    
    // 业务逻辑
    // ...
    
    return nil
}
```

### 分段计时

```go
func ProcessWithSegments(ctx g.Ctx) error {
    segments := make(map[string]time.Duration)
    
    // 段 1
    start := time.Now()
    if err := fetchFromDB(ctx); err != nil {
        return err
    }
    segments["fetchFromDB"] = time.Since(start)
    
    // 段 2
    start = time.Now()
    if err := processData(ctx); err != nil {
        return err
    }
    segments["processData"] = time.Since(start)
    
    // 段 3
    start = time.Now()
    if err := saveToDB(ctx); err != nil {
        return err
    }
    segments["saveToDB"] = time.Since(start)
    
    // 记录分段耗时
    for name, duration := range segments {
        g.Log().Infof(ctx, "分段 %s: 耗时 %v", name, duration)
    }
    
    return nil
}
```

## 日志轮转

### 按时间轮转

```yaml
logger:
  path: "logs"
  level: "all"
  stdout: true
  rotateExpire: "1d"    # 每天轮转
  rotateBackupExpire: "7d" # 保留 7 天
  rotateBackupLimit: 30 # 最多保留 30 个文件
  rotateSize: "100M"    # 文件大小超过 100M 时轮转
```

### 按大小轮转

```yaml
logger:
  path: "logs"
  level: "all"
  rotateSize: "100M"    # 文件大小超过 100M 时轮转
  rotateBackupLimit: 10 # 最多保留 10 个文件
```

## 日志过滤

### 按模块过滤

```go
// 使用不同的 logger
userLogger := g.Log("user")
orderLogger := g.Log("order")

// 配置不同级别
userLogger.SetLevel(glog.LEVEL_ALL)
orderLogger.SetLevel(glog.LEVEL_WARN)

// 记录日志
userLogger.Infof(ctx, "用户登录: %s", username)
orderLogger.Warningf(ctx, "订单处理异常: %v", err)
```

### 条件过滤

```go
func ProcessData(ctx g.Ctx) error {
    data, err := fetchData(ctx)
    if err != nil {
        // 只在开发环境记录详细错误
        if config.System().Mode == "develop" {
            g.Log().Errorf(ctx, "获取数据失败: %v, 数据: %+v", err, data)
        } else {
            g.Log().Errorf(ctx, "获取数据失败: %v", err)
        }
        return err
    }
    return nil
}
```

## 最佳实践

1. **合理使用级别**: 根据日志重要性选择合适的级别
2. **结构化日志**: 使用结构化格式便于日志分析
3. **避免过度日志**: 不要记录过多无关的日志
4. **记录关键操作**: 记录关键业务操作
5. **错误详情**: 错误日志包含足够的上下文信息
6. **性能日志**: 记录关键操作的执行时间
7. **日志轮转**: 配置日志轮转避免日志文件过大
8. **敏感信息**: 不要记录密码等敏感信息

## 下一步

- 了解 [架构设计](./architecture)
- 查看 [认证授权](./authentication)
