# 关键操作日志记录

## 1. 简介

ZZFrame 内置了基础的操作日志（Operation Log）功能，通过中间件自动记录所有写操作（POST/PUT/DELETE）的请求参数、响应结果和执行耗时。

然而，仅记录请求参数往往不足以满足审计需求。例如，当管理员修改系统配置时，我们不仅需要知道“修改了配置项”，还需要知道“修改前的值是 A，修改后的值是 B”。这就需要实现**数据变更审计**。

**适用场景**：
*   管理员账号管理（新增、修改权限）
*   系统参数配置变更
*   敏感数据修改（如费率、黑白名单）

## 2. 数据库设计

为了支持高效的审计查询，我们需要对现有的 `sys_oper_log` 表进行扩展，增加业务关联字段和差异存储字段。

### 2.1 表结构变更


```sql
-- 新增审计相关字段
ALTER TABLE `sys_oper_log`
    ADD COLUMN `biz_id` varchar(64) DEFAULT '' COMMENT '业务ID' AFTER `title`,
    ADD COLUMN `biz_type` varchar(64) DEFAULT '' COMMENT '业务类型' AFTER `biz_id`,
    ADD COLUMN `diff` json DEFAULT NULL COMMENT '数据变更差异' AFTER `result`;

-- 添加联合索引，方便查询“某个配置项的所有变更记录”
ALTER TABLE `sys_oper_log` ADD INDEX `idx_biz` (`biz_type`, `biz_id`);
```

### 2.2 字段说明

| 字段名 | 类型 | 说明 | 示例 |
|--------|------|------|------|
| `biz_id` | varchar | 业务对象的唯一标识 | `1001` (用户ID) 或 `sys_config_01` (配置键) |
| `biz_type` | varchar | 业务对象的类型标识 | `admin_member` (管理员) 或 `sys_config` (系统配置) |
| `diff` | json | 数据变更前后的差异 | `{"old": {"status": 1}, "new": {"status": 2}}` |

> **注意**：
> 1. `biz_id` 使用 `varchar` 是为了兼容非数字 ID（如 UUID 或 字符串主键）。
> 2. `biz_type` 建议使用表名或统一的业务常量。

### 2.3 更新代码模型

执行完 SQL 变更后，请运行 GF 工具重新生成 DAO 和 Entity：

```bash
# 在项目根目录下运行
gf gen dao
```

或者手动更新 `internal/model/entity/sys_oper_log.go`：

```go
type SysOperLog struct {
    // ... 原有字段
    BizId     string      `json:"bizId"     orm:"biz_id"     description:"业务ID"`
    BizType   string      `json:"bizType"   orm:"biz_type"   description:"业务类型"`
    Diff      *gjson.Json `json:"diff"      orm:"diff"       description:"数据变更差异"`
    // ...
}
```

## 3. 核心思路

利用框架提供的 Context 上下文传递机制，在 Service 业务层捕获数据变更（Diff），并将其传递给 Middleware 中间件，最终由中间件统一入库。

**流程图：**

1.  **Service 层**：查询旧数据 -> 执行更新 -> 记录 Diff + BizInfo -> 写入 Context。
2.  **Middleware 层**：请求结束 -> 从 Context 提取 Diff/BizInfo -> 填充到日志对象。
3.  **Queue/Database**：异步写入数据库。

## 4. 实现步骤

### 4.1 定义 Context Key

在 `zconsts/context.go` 中定义用于存储审计数据的 Key。

```go
// zconsts/context.go

const (
    // ... 其他常量
    ContextKeyAudit = "audit_data" // 审计日志数据 Key
)
```

### 4.2 封装辅助工具


```go

package zutils

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/denghuo98/zzframe/web/zcontext"
)

// AuditData 审计数据结构
type AuditData struct {
    BizId   string      `json:"biz_id"`
    BizType string      `json:"biz_type"`
    Old     interface{} `json:"old"`
    New     interface{} `json:"new"`
}

// SetAudit 记录审计数据
// bizType: 业务类型 (如 "admin_member", "sys_config")
// bizId: 业务ID (如 "1", "config_key")
// oldData: 修改前的数据
// newData: 修改后的数据（或变更的字段）
func SetAudit(ctx g.Ctx, bizType string, bizId string, oldData interface{}, newData interface{}) {
    data := AuditData{
        BizType: bizType,
        BizId:   bizId,
        Old:     oldData,
        New:     newData,
    }
    // 挂载到 Context
    zcontext.SetData(ctx, "audit_data", data)
}
```

### 4.3 改造 Middleware

修改 `zservice/logic/middleware/oper_log.go`，在组装日志时读取 Context 中的审计数据。
Service 层虽然计算出了 Diff，但它不能直接写日志（否则代码会很乱，每个 Service 方法都要写一遍日志入库代码）。
所以我打算采用**“Service 生产，Middleware 收集”**的模式

```go
// ...
    // 获取审计数据
    auditData := zcontext.GetData(ctx, "audit_data")
    
    // 初始化字段
    var (
        bizId   = ""
        bizType = ""
        diff    *gjson.Json
    )

    if auditData != nil {
        // 解析数据 (假设 auditData 是我们在 zutils 中定义的结构体或 Map)
        // 这里简化处理，转为 Json 再取值，或者使用类型断言
        j := gjson.New(auditData)
        bizId = j.Get("biz_id").String()
        bizType = j.Get("biz_type").String()
        
        // 构建 Diff 结构
        diff = gjson.New(g.Map{
            "old": j.Get("old"),
            "new": j.Get("new"),
        })
    }

    data := entity.SysOperLog{
        // ... 原有字段
        BizId:     bizId,
        BizType:   bizType,
        Diff:      diff,
        // ...
    }
// ...
```

### 4.4 业务层调用 (Service)

在具体的业务逻辑中（例如 `zservice/logic/admin/member.go`），记录变更。

```go
func (s *sAdminMember) Edit(ctx g.Ctx, in *adminSchema.MemberEditInput) (err error) {
    // 1. 查询旧数据 (用于审计)
    // 注意：尽量只查询需要的字段，避免性能浪费
    oldData, _ := dao.AdminMember.Ctx(ctx).WherePri(in.Id).One()

    // 2. 执行更新
    // ... update logic ...

    // 3. 记录变更
    if oldData != nil {
        // 这里可以对 oldData 做一些处理，比如去除敏感字段（密码等）
        oldData.Password = "***"
        
        // 记录审计信息: 业务类型=admin_member, ID=用户ID
        zutils.SetAudit(ctx, "admin_member", gconv.String(in.Id), oldData, in)
    }
    
    return nil
}
```

## 5. 设计考量：为什么选择异步？

在实现日志记录时，选择了**异步入库**（先推入队列，后台消费入库），而非同步写入。主要基于以下考量：

1.  **低延迟（Performance）**：日志记录属于辅助功能，不应阻塞主业务流程。同步写入数据库会增加接口响应时间（RT），特别是在高并发或数据库负载较高时。
2.  **高可用与解耦（Decoupling）**：如果日志数据库发生抖动或宕机，同步写入会导致主业务接口报错。异步模式下，队列可以作为缓冲区，即使数据库暂时不可用，日志也不会立即丢失，且不会拖垮主业务。


> **注意**：异步模式存在极小概率的日志丢失风险。

## 6. 总结

通过增加 `biz_id`, `biz_type` 和 `diff` 字段，我们将原本扁平的操作日志升级为了具备业务语义的审计日志。

*   **BizType/BizId**: 解决了“查找某个对象变更历史”的问题（例如：谁修改了这个系统参数？）。
*   **Diff**: 解决了“发生了什么具体变化”的问题（例如：将重试次数从 3 改为了 5）。
