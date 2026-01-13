# 架构设计

ZZFrame 采用经典的分层架构设计，结合现代 Web 开发最佳实践。

## 整体架构

```
┌──────────────────────────────────────────────────────┐
│                   Presentation Layer                 │
│              (HTTP/RESTful API)                      │
└──────────────────────────────────────────────────────┘
                          │
                          ↓
┌──────────────────────────────────────────────────────┐
│                   Middleware Layer                   │
│  (Authentication | Authorization | Logging | CORS)   │
└──────────────────────────────────────────────────────┘
                          │
                          ↓
┌──────────────────────────────────────────────────────┐
│                   Controller Layer                   │
│            (Request/Response Handling)                │
└──────────────────────────────────────────────────────┘
                          │
                          ↓
┌──────────────────────────────────────────────────────┐
│                    Service Layer                     │
│               (Business Logic)                       │
└──────────────────────────────────────────────────────┘
                          │
                          ↓
┌──────────────────────────────────────────────────────┐
│                     DAO Layer                        │
│                (Data Access Object)                  │
└──────────────────────────────────────────────────────┘
                          │
                          ↓
┌──────────────────────────────────────────────────────┐
│                  Data Layer                          │
│  (Database | Cache | Queue | File System)            │
└──────────────────────────────────────────────────────┘
```

## 分层职责

### Presentation Layer (表示层)

- 提供 RESTful API 接口
- 处理 HTTP 请求和响应
- 数据格式转换
- 参数验证

**位置**: `zcontroller/`, `api/`

### Middleware Layer (中间件层)

- 请求认证
- 权限验证
- 日志记录
- 跨域处理
- 请求限流

**位置**: `zservice/logic/middleware/`

### Controller Layer (控制层)

- 接收 HTTP 请求
- 参数验证
- 调用 Service 层
- 返回响应

**位置**: `zcontroller/`

### Service Layer (服务层)

- 实现业务逻辑
- 事务管理
- 调用 DAO 层
- 数据转换

**位置**: `zservice/`

### DAO Layer (数据访问层)

- 数据库操作
- SQL 查询
- 数据映射

**位置**: `internal/dao/`

### Data Layer (数据层)

- 数据库存储
- 缓存存储
- 队列存储
- 文件系统

**位置**: MySQL, Redis, File System

## 核心模块设计

### 认证模块

```
┌─────────────┐
│   Login     │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│ Verify User │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│ Generate   │
│   Token    │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│   Return    │
│   Token    │
└─────────────┘
```

### 权限模块

```
┌─────────────┐
│   Request   │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│ Parse Token│
└──────┬──────┘
       │
       ↓
┌─────────────┐
│ Get User    │
│   Roles     │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│ Casbin     │
│  Enforce   │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│  Allow/     │
│  Deny       │
└─────────────┘
```

### 缓存模块

```
┌─────────────┐
│   Request   │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│   Check     │
│   Cache     │
└──────┬──────┘
       │
   ┌───┴───┐
   │       │
   ↓       ↓
  Hit    Miss
   │       │
   │       ↓
   │   ┌─────────────┐
   │   │ Query DB    │
   │   └──────┬──────┘
   │          │
   │          ↓
   │   ┌─────────────┐
   │   │ Set Cache   │
   │   └──────┬──────┘
   │          │
   └────┬─────┘
        │
        ↓
┌─────────────┐
│   Return    │
└─────────────┘
```

### 队列模块

```
┌─────────────┐
│   Producer  │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│   Push to   │
│   Queue     │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│   Consumer  │
│   Pull      │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│   Process   │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│   Complete  │
└─────────────┘
```

## 数据流设计

### 请求处理流程

```
1. Client 发送 HTTP 请求
   ↓
2. Middleware 进行认证和授权
   ↓
3. Controller 接收请求，验证参数
   ↓
4. Service 处理业务逻辑
   ↓
5. DAO 查询/更新数据库
   ↓
6. Service 返回结果
   ↓
7. Controller 格式化响应
   ↓
8. Client 接收响应
```

### 缓存命中流程

```
1. Service 查询缓存
   ↓
2. Cache 存在数据 → 返回
   ↓
3. Cache 不存在数据
   ↓
4. DAO 查询数据库
   ↓
5. 写入缓存
   ↓
6. 返回数据
```

### 队列处理流程

```
1. Service 产生消息
   ↓
2. Push 到队列
   ↓
3. Consumer 拉取消息
   ↓
4. 执行异步任务
   ↓
5. 标记消息已处理
```

## 设计模式

### 1. 单例模式

```go
var localAdminMember IAdminMember

func AdminMember() IAdminMember {
    if localAdminMember == nil {
        panic("AdminMember is not initialized")
    }
    return localAdminMember
}
```

### 2. 工厂模式

```go
func NewCache(adapter string) (cache.ICache, error) {
    switch adapter {
    case "redis":
        return redis.New()
    case "file":
        return file.New()
    default:
        return nil, errors.New("unsupported cache adapter")
    }
}
```

### 3. 策略模式

```go
type IQueue interface {
    Push(ctx context.Context, topic string, data interface{}) error
    Pull(ctx context.Context, topic string) ([]interface{}, error)
}

type QueueManager struct {
    queues map[string]IQueue
}

func (qm *QueueManager) GetQueue(name string) IQueue {
    return qm.queues[name]
}
```

### 4. 依赖注入

```go
// 接口定义
type IAdminMember interface {
    List(ctx g.Ctx, in *MemberListInput) (*MemberListOutput, int, error)
}

// 注册实现
func RegisterAdminMember(i IAdminMember) {
    localAdminMember = i
}

// 使用接口
func Controller() {
    service := AdminMember() // 通过接口访问
}
```

## 性能优化

### 1. 数据库优化

- 使用连接池
- 合理使用索引
- 避免 N+1 查询
- 使用批量操作

### 2. 缓存优化

- 合理设置缓存过期时间
- 使用缓存预热
- 缓存穿透防护
- 缓存雪崩防护

### 3. 代码优化

- 减少内存分配
- 使用对象池
- 避免不必要的类型转换
- 使用切片预分配

## 扩展性设计

### 1. 中间件扩展

```go
func Middleware(ctx ghttp.Request) {
    // 自定义逻辑
    ctx.Middleware.Next()
}
```

### 2. 缓存扩展

```go
type CustomCache struct{}

func (c *CustomCache) Set(ctx context.Context, key string, value interface{}, ttl int) error {
    // 自定义缓存实现
}
```

### 3. 队列扩展

```go
type CustomQueue struct{}

func (q *CustomQueue) Push(ctx context.Context, topic string, data interface{}) error {
    // 自定义队列实现
}
```

## 下一步

- 了解 [认证授权](./authentication)
- 学习 [权限管理](./authorization)
