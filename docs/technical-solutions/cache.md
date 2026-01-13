# 缓存方案

ZZFrame 提供了灵活的缓存系统，支持文件缓存和 Redis 缓存。

## 缓存类型

### 文件缓存

适用于小型项目或开发环境，无需额外的缓存服务。

### Redis 缓存

适用于中大型项目，提供更好的性能和分布式缓存能力。

## 配置

### 文件缓存配置

```yaml
system:
  cache:
    adapter: "file"        # 使用文件缓存
    fileDir: "tmp/cache"   # 缓存文件目录
```

### Redis 缓存配置

```yaml
system:
  cache:
    adapter: "redis"       # 使用 Redis 缓存
    redis:
      host: "127.0.0.1:6379"
      password: ""
      db: 0
```

## 缓存接口

### 缓存接口定义

```go
type ICache interface {
    Set(ctx context.Context, key string, value interface{}, ttl int) error
    Get(ctx context.Context, key string, value interface{}) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    Clear(ctx context.Context) error
}
```

### 获取缓存实例

```go
config := zservice.SystemConfig()
cache, err := config.Cache()
if err != nil {
    return err
}
```

## 基本操作

### 设置缓存

```go
// 设置缓存，过期时间 300 秒
err := cache.Set(ctx, "user:1", user, 300)

// 设置永久缓存
err := cache.Set(ctx, "config:app", config, 0)
```

### 获取缓存

```go
var user User
err := cache.Get(ctx, "user:1", &user)
if err != nil {
    // 缓存不存在或出错
}

fmt.Println(user.Username)
```

### 删除缓存

```go
err := cache.Delete(ctx, "user:1")
```

### 检查缓存是否存在

```go
exists, err := cache.Exists(ctx, "user:1")
if exists {
    fmt.Println("缓存存在")
}
```

### 清空所有缓存

```go
err := cache.Clear(ctx)
```

## 高级用法

### 批量操作

```go
// 批量设置
err := cache.MSet(ctx, map[string]interface{}{
    "user:1": user1,
    "user:2": user2,
})

// 批量获取
keys := []string{"user:1", "user:2"}
values, err := cache.MGet(ctx, keys...)
```

### 过期时间操作

```go
// 获取过期时间
ttl, err := cache.TTL(ctx, "user:1")

// 设置过期时间
err := cache.Expire(ctx, "user:1", 600)
```

### 缓存标签

```go
// 设置带标签的缓存
err := cache.SetWithTags(ctx, "user:1", user, 300, []string{"user"})

// 按标签删除缓存
err := cache.DeleteByTags(ctx, []string{"user"})
```

## 缓存模式

### Cache Aside 模式

```go
func GetUser(ctx context.Context, id uint64) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    var user User
    
    // 1. 先查缓存
    if err := cache.Get(ctx, cacheKey, &user); err == nil && user.Id > 0 {
        return &user, nil
    }
    
    // 2. 缓存不存在，查数据库
    user, err := dao.User.GetById(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    _ = cache.Set(ctx, cacheKey, user, 300)
    
    return &user, nil
}
```

### Write Through 模式

```go
func UpdateUser(ctx context.Context, user *User) error {
    // 1. 更新数据库
    if err := dao.User.Update(ctx, user); err != nil {
        return err
    }
    
    // 2. 同步更新缓存
    cacheKey := fmt.Sprintf("user:%d", user.Id)
    _ = cache.Set(ctx, cacheKey, user, 300)
    
    return nil
}
```

### Write Behind 模式

```go
func UpdateUserAsync(ctx context.Context, user *User) error {
    // 1. 更新数据库
    if err := dao.User.Update(ctx, user); err != nil {
        return err
    }
    
    // 2. 异步更新缓存
    go func() {
        cacheKey := fmt.Sprintf("user:%d", user.Id)
        _ = cache.Set(context.Background(), cacheKey, user, 300)
    }()
    
    return nil
}
```

## 缓存预热

```go
func WarmUpCache(ctx context.Context) error {
    // 预热热门用户
    userIds := []uint64{1, 2, 3, 4, 5}
    
    for _, id := range userIds {
        user, err := dao.User.GetById(ctx, id)
        if err != nil {
            continue
        }
        
        cacheKey := fmt.Sprintf("user:%d", id)
        _ = cache.Set(ctx, cacheKey, user, 300)
    }
    
    return nil
}
```

## 缓存失效

### 主动失效

```go
func DeleteUser(ctx context.Context, id uint64) error {
    // 1. 删除数据库
    if err := dao.User.Delete(ctx, id); err != nil {
        return err
    }
    
    // 2. 删除缓存
    cacheKey := fmt.Sprintf("user:%d", id)
    _ = cache.Delete(ctx, cacheKey)
    
    return nil
}
```

### 按标签失效

```go
func UpdateRole(ctx context.Context, role *Role) error {
    // 更新角色后，删除所有用户缓存
    _ = cache.DeleteByTags(ctx, []string{"user"})
    return nil
}
```

## 缓存问题处理

### 缓存穿透

```go
func GetUser(ctx context.Context, id uint64) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    var user User
    
    // 查缓存
    if err := cache.Get(ctx, cacheKey, &user); err == nil && user.Id > 0 {
        return &user, nil
    }
    
    // 缓存不存在，查数据库
    user, err := dao.User.GetById(ctx, id)
    if err != nil {
        // 缓存空值，防止缓存穿透
        _ = cache.Set(ctx, cacheKey, nil, 60)
        return nil, err
    }
    
    // 写入缓存
    _ = cache.Set(ctx, cacheKey, user, 300)
    
    return &user, nil
}
```

### 缓存雪崩

```go
// 设置不同的过期时间
func GetUser(ctx context.Context, id uint64) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    var user User
    
    // 查缓存
    if err := cache.Get(ctx, cacheKey, &user); err == nil && user.Id > 0 {
        return &user, nil
    }
    
    // 查数据库
    user, err := dao.User.GetById(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 随机过期时间，避免雪崩
    ttl := 300 + rand.Intn(60) // 300-360 秒
    _ = cache.Set(ctx, cacheKey, user, ttl)
    
    return &user, nil
}
```

### 缓存击穿

```go
var userMutex sync.Map

func GetUser(ctx context.Context, id uint64) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    var user User
    
    // 查缓存
    if err := cache.Get(ctx, cacheKey, &user); err == nil && user.Id > 0 {
        return &user, nil
    }
    
    // 使用互斥锁，只让一个请求查数据库
    mu, _ := userMutex.LoadOrStore(id, &sync.Mutex{})
    mutex := mu.(*sync.Mutex)
    
    mutex.Lock()
    defer mutex.Unlock()
    
    // 再次检查缓存
    if err := cache.Get(ctx, cacheKey, &user); err == nil && user.Id > 0 {
        return &user, nil
    }
    
    // 查数据库
    user, err := dao.User.GetById(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 写入缓存
    _ = cache.Set(ctx, cacheKey, user, 300)
    
    return &user, nil
}
```

## 最佳实践

1. **合理设置过期时间**: 根据数据更新频率设置合理的 TTL
2. **缓存预热**: 应用启动时预热热点数据
3. **缓存穿透防护**: 缓存空值或使用布隆过滤器
4. **缓存雪崩防护**: 使用随机过期时间
5. **缓存击穿防护**: 使用互斥锁
6. **监控缓存**: 监控缓存命中率和内存使用
7. **定期清理**: 定期清理过期缓存

## 下一步

- 学习 [队列方案](./queue)
- 了解 [日志方案](./logging)
