# 队列方案

ZZFrame 提供了可靠的消息队列系统，支持磁盘队列和 Redis 队列，用于异步任务处理。

## 队列类型

### 磁盘队列

持久化到磁盘，适合单机应用，消息不会丢失。

### Redis 队列

基于 Redis 实现，支持分布式场景。

## 配置

### 磁盘队列配置

```yaml
queue:
  switch: true              # 是否启用队列
  driver: "disk"           # 队列驱动
  groupName: "default"     # 队列组名
  disk:
    path: "./tmp/diskqueue"     # 磁盘队列目录
    batchSize: 100             # 批量处理数量
    batchTime: 1               # 批量处理时间（秒）
    segmentSize: 10485760      # 分段大小（字节）
    segmentLimit: 3000         # 分段数量限制
```

### Redis 队列配置

```yaml
queue:
  switch: true
  driver: "redis"
  redis:
    host: "127.0.0.1:6379"
    password: ""
    db: 0
```

## 队列接口

### 队列接口定义

```go
type IQueue interface {
    Push(ctx context.Context, topic string, data interface{}) error
    Pull(ctx context.Context, topic string) ([]interface{}, error)
    PushDelay(ctx context.Context, topic string, data interface{}, delay int64) error
}
```

### 获取队列实例

```go
config := zservice.SystemConfig()
queue, err := config.Queue()
if err != nil {
    return err
}
```

## 基本操作

### 生产者推送消息

```go
// 推送消息
err := queue.Push(ctx, "email", map[string]interface{}{
    "to":      "user@example.com",
    "subject": "欢迎注册",
    "body":    "欢迎注册我们的平台",
})
```

### 消费者拉取消息

```go
// 拉取消息
messages, err := queue.Pull(ctx, "email")
if err != nil {
    return err
}

// 处理消息
for _, msg := range messages {
    emailMsg := msg.(map[string]interface{})
    to := emailMsg["to"].(string)
    subject := emailMsg["subject"].(string)
    body := emailMsg["body"].(string)
    
    // 发送邮件
    sendEmail(to, subject, body)
}
```

### 延迟消息

```go
// 延迟 60 秒后执行
err := queue.PushDelay(ctx, "email", emailData, 60)
```

## 使用场景

### 1. 异步发送邮件

```go
// 注册成功后异步发送欢迎邮件
func (s *sUser) Register(ctx g.Ctx, req *RegisterReq) error {
    // 注册用户
    userId, err := dao.User.Insert(ctx, &do.User{
        Username: req.Username,
        Email:    req.Email,
    })
    if err != nil {
        return err
    }
    
    // 异步发送欢迎邮件
    queue.Push(ctx, "email", map[string]interface{}{
        "type":    "welcome",
        "userId":  userId,
        "to":      req.Email,
    })
    
    return nil
}
```

### 2. 异步处理图片

```go
// 上传图片后异步处理缩略图
func (s *sUpload) UploadImage(ctx g.Ctx, file *ghttp.UploadFile) (string, error) {
    // 保存原图
    url, err := saveImage(file)
    if err != nil {
        return "", err
    }
    
    // 异步生成缩略图
    queue.Push(ctx, "image", map[string]interface{}{
        "type":  "thumbnail",
        "url":   url,
        "sizes": []int{100, 200, 400},
    })
    
    return url, nil
}
```

### 3. 异步生成报表

```go
// 定时生成日报
func GenerateDailyReport(ctx g.Ctx, date string) {
    queue.PushDelay(ctx, "report", map[string]interface{}{
        "type": "daily",
        "date": date,
    }, 3600) // 延迟 1 小时
}
```

### 4. 异步清理数据

```go
// 清理过期数据
func CleanExpiredData(ctx g.Ctx) {
    queue.Push(ctx, "cleanup", map[string]interface{}{
        "type": "expired",
        "days": 30,
    })
}
```

## 消费者实现

### 消费者示例

```go
package worker

import (
    "context"
    "github.com/gogf/gf/v2/frame/g"
)

func StartConsumer(ctx context.Context) {
    config := zservice.SystemConfig()
    queue, _ := config.Queue()
    
    // 启动邮件消费者
    go EmailConsumer(ctx, queue)
    
    // 启动图片处理消费者
    go ImageConsumer(ctx, queue)
    
    // 启动报表生成消费者
    go ReportConsumer(ctx, queue)
}

func EmailConsumer(ctx context.Context, queue IQueue) {
    for {
        // 拉取消息
        messages, err := queue.Pull(ctx, "email")
        if err != nil {
            g.Log().Errorf(ctx, "拉取邮件消息失败: %v", err)
            continue
        }
        
        // 处理消息
        for _, msg := range messages {
            emailMsg := msg.(map[string]interface{})
            if err := processEmail(ctx, emailMsg); err != nil {
                g.Log().Errorf(ctx, "处理邮件失败: %v", err)
                // 失败重新入队
                queue.Push(ctx, "email", emailMsg)
            }
        }
    }
}

func processEmail(ctx context.Context, msg map[string]interface{}) error {
    typ := msg["type"].(string)
    
    switch typ {
    case "welcome":
        return sendWelcomeEmail(ctx, msg)
    case "reset":
        return sendResetEmail(ctx, msg)
    default:
        return errors.New("未知的邮件类型")
    }
}
```

## 消息可靠性

### 消息确认机制

```go
func Consumer(ctx context.Context) {
    for {
        messages, err := queue.Pull(ctx, "topic")
        if err != nil {
            continue
        }
        
        for _, msg := range messages {
            // 处理消息
            if err := processMessage(ctx, msg); err != nil {
                // 处理失败，重新入队
                queue.Push(ctx, "topic", msg)
                continue
            }
            
            // 处理成功，确认消息
            // queue.Ack(ctx, msg)
        }
    }
}
```

### 死信队列

```go
func ConsumerWithRetry(ctx context.Context) {
    retryCount := make(map[interface{}]int)
    
    for {
        messages, err := queue.Pull(ctx, "topic")
        if err != nil {
            continue
        }
        
        for _, msg := range messages {
            // 检查重试次数
            if retryCount[msg] >= 3 {
                // 超过重试次数，放入死信队列
                queue.Push(ctx, "topic:dead", msg)
                delete(retryCount, msg)
                continue
            }
            
            // 处理消息
            if err := processMessage(ctx, msg); err != nil {
                // 处理失败，重新入队
                retryCount[msg]++
                queue.Push(ctx, "topic", msg)
            } else {
                // 处理成功
                delete(retryCount, msg)
            }
        }
    }
}
```

## 性能优化

### 批量处理

```go
func BatchConsumer(ctx context.Context) {
    batchSize := 100
    
    for {
        // 批量拉取消息
        messages, err := queue.PullBatch(ctx, "topic", batchSize)
        if err != nil {
            continue
        }
        
        // 批量处理
        if err := processBatch(ctx, messages); err != nil {
            // 失败重新入队
            for _, msg := range messages {
                queue.Push(ctx, "topic", msg)
            }
        }
    }
}
```

### 并发处理

```go
func ConcurrentConsumer(ctx context.Context) {
    workerCount := 10
    workers := make(chan struct{}, workerCount)
    
    for {
        messages, err := queue.Pull(ctx, "topic")
        if err != nil {
            continue
        }
        
        for _, msg := range messages {
            workers <- struct{}{}
            go func(m interface{}) {
                defer func() { <-workers }()
                
                if err := processMessage(ctx, m); err != nil {
                    queue.Push(ctx, "topic", m)
                }
            }(msg)
        }
    }
}
```

## 监控和日志

### 消息处理日志

```go
func processMessage(ctx context.Context, msg interface{}) error {
    start := time.Now()
    
    g.Log().Infof(ctx, "开始处理消息: %+v", msg)
    
    // 处理消息
    if err := doProcess(ctx, msg); err != nil {
        g.Log().Errorf(ctx, "处理消息失败: %v, 耗时: %v", err, time.Since(start))
        return err
    }
    
    g.Log().Infof(ctx, "处理消息成功, 耗时: %v", time.Since(start))
    return nil
}
```

### 队列监控

```go
func MonitorQueue(ctx context.Context) {
    for {
        // 监控队列长度
        len, err := queue.Len(ctx, "topic")
        if err != nil {
            g.Log().Errorf(ctx, "获取队列长度失败: %v", err)
        } else {
            g.Log().Infof(ctx, "队列长度: %d", len)
        }
        
        // 监控消费者状态
        // ...
        
        time.Sleep(60 * time.Second)
    }
}
```

## 最佳实践

1. **消息幂等**: 确保消息处理幂等，避免重复消费
2. **错误处理**: 妥善处理失败消息，考虑重试机制
3. **批量处理**: 使用批量处理提升性能
4. **并发消费**: 使用并发消费者提升吞吐量
5. **监控告警**: 监控队列长度和处理速度
6. **资源隔离**: 不同类型的消息使用不同的队列
7. **超时控制**: 设置消息处理超时时间

## 下一步

- 学习 [日志方案](./logging)
- 了解 [架构设计](./architecture)
