# 关键操作日志记录

## 1. 简介

ZZFrame 内置了基础的操作日志（Operation Log）功能，通过中间件自动记录所有写操作（POST/PUT/DELETE）的请求参数、响应结果和执行耗时。

然而，仅记录请求参数往往不足以满足审计需求。例如，当管理员修改系统配置时，我们不仅需要知道"修改了配置项"，还需要知道"修改前的值是 A，修改后的值是 B"。这就需要实现**数据变更审计**。

## 2. 数据库设计

为了支持高效的审计查询，我们需要对现有的 `sys_oper_log` 表进行扩展，增加业务关联字段和差异存储字段。

### 2.1 表结构变更


```sql
-- 新增审计相关字段
ALTER TABLE `sys_oper_log`
    ADD COLUMN `biz_id` varchar(64) DEFAULT '' COMMENT '业务ID' AFTER `title`,
    ADD COLUMN `biz_type` varchar(64) DEFAULT '' COMMENT '业务类型' AFTER `biz_id`,
    ADD COLUMN `diff` json DEFAULT NULL COMMENT '数据变更差异' AFTER `result`;

-- 添加联合索引，方便查询"某个配置项的所有变更记录"
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

### 2.3 权限控制

审计日志一旦写入就**不应被修改**，否则无法保证日志的真实性。建议按角色分配数据库权限：

| 角色 | 权限 | 说明 |
|------|------|------|
| **应用服务账号** | `INSERT` only | 只能写入，禁止修改删除 |
| **管理员查询** | `SELECT` | 只能查看日志 |
| **归档任务账号** | `SELECT, INSERT, DELETE` | 负责数据迁移 |

```sql
-- 1. 应用写入账号（只允许 INSERT）
CREATE USER 'oper_log_writer'@'%' IDENTIFIED BY 'strong_password';
GRANT INSERT ON your_database.sys_oper_log TO 'oper_log_writer'@'%';

-- 2. 只读查询账号
CREATE USER 'oper_log_reader'@'%' IDENTIFIED BY 'read_password';
GRANT SELECT ON your_database.sys_oper_log TO 'oper_log_reader'@'%';

-- 3. 归档任务账号
CREATE USER 'oper_log_archiver'@'localhost' IDENTIFIED BY 'archive_password';
GRANT SELECT, INSERT, DELETE ON your_database.sys_oper_log TO 'oper_log_archiver'@'localhost';
GRANT SELECT, INSERT, DELETE ON your_database.sys_oper_log_archive TO 'oper_log_archiver'@'localhost';
```

> **重要**：禁止授予 `UPDATE` 权限，防止日志被篡改。

### 2.4 更新代码模型

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

```
┌─────────────────────────────────────────────────────────────────────┐
│                        批量操作审计流程                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   Service 层                     Middleware                Database │
│   ┌─────────────┐               ┌─────────────┐          ┌────────┐ │
│   │ 1. 查旧数据  │               │ 4. 读取     │           │        │ │
│   │    oldList  │               │ audit_batch │          │        │ │
│   │             │               │             │          │        │ │
│   │ 2. 执行操作  │      ───►     │ 5. 遍历生成 │   ───►     │ INSERT │ │
│   │             │               │    多条日志 │            │ 多条   │ │
│   │ 3. SetBatch │               │             │          │        │ │
│   │   [item1,   │               │ [log1,log2, │          │        │ │
│   │    item2,..]│               │  log3,...]  │          │        │ │
│   └─────────────┘               └─────────────┘          └────────┘ │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

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

// SetAudit 记录单条审计数据
func SetAudit(ctx g.Ctx, bizType string, bizId string, oldData interface{}, newData interface{}) {
    data := AuditData{
        BizType: bizType,
        BizId:   bizId,
        Old:     oldData,
        New:     newData,
    }
    zcontext.SetData(ctx, "audit_data", data)
}

// SetAuditBatch 批量记录审计数据（用于批量操作场景）
func SetAuditBatch(ctx g.Ctx, items []AuditData) {
    zcontext.SetData(ctx, "audit_data_batch", items)
}
```

#### 批量操作使用示例

```go
// 批量删除用户
func (s *sUser) BatchDelete(ctx g.Ctx, ids []int) error {
    // 1. 查询所有旧数据
    oldList, _ := dao.User.Ctx(ctx).WhereIn("id", ids).All()
    
    // 2. 执行删除
    _, err := dao.User.Ctx(ctx).WhereIn("id", ids).Delete()
    if err != nil {
        return err
    }
    
    // 3. 构建批量审计记录
    var items []zutils.AuditData
    for _, old := range oldList {
        items = append(items, zutils.AuditData{
            BizType: "user",
            BizId:   gconv.String(old["id"]),
            Old:     old,
            New:     nil, // 删除操作 new 为 nil
        })
    }
    zutils.SetAuditBatch(ctx, items)
    return nil
}
```

> **注意**：Middleware 需要同时处理 `audit_data`（单条）和 `audit_data_batch`（批量）两种情况。

### 4.3 改造 Middleware

修改 `zservice/logic/middleware/oper_log.go`，在组装日志时读取 Context 中的审计数据。
Service 层虽然计算出了 Diff，但它不能直接写日志（否则代码会很乱，每个 Service 方法都要写一遍日志入库代码）。
所以我打算采用**"Service 生产，Middleware 收集"**的模式

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

## 6. 数据生命周期管理

操作日志作为辅助业务数据，不应永久存储于主数据库。合理的生命周期管理可以降低存储成本、提升查询性能。

### 6.1 归档策略

| 阶段 | 时间范围 | 存储位置 | 访问频率 | 说明 |
|------|----------|----------|----------|------|
| **热数据** | 0-3个月 | `sys_oper_log` 主表 | 高频 | 支持实时查询 |
| **温数据** | 3-12个月 | `sys_oper_log_archive` 归档表 | 低频 | 按需查询 |
| **冷数据** | 12个月+ | JSON/CSV 文件 | 极少 | 压缩存储，按需恢复 |

```
┌─────────────────────────────────────────────────────────────────┐
│                        操作日志生命周期                            │
├─────────────────────────────────────────────────────────────────┤
│   [热数据]          [温数据]                [冷数据]               │
│   0-3个月           3-12个月               12个月+                │
│   ┌─────────┐      ┌─────────────┐       ┌─────────────┐        │
│   │主数据库  │ ───► │ 归档表       │ ───►   │ 文件存储     │        │
│   │高频查询  │      │ 低频查询     │        │ 按需恢复    │         │
│   └─────────┘      └─────────────┘       └─────────────┘        │
└─────────────────────────────────────────────────────────────────┘
```

### 6.2 归档表结构

```sql
-- 创建归档表（结构与主表一致）
CREATE TABLE `sys_oper_log_archive` LIKE `sys_oper_log`;

-- 添加归档时间字段
ALTER TABLE `sys_oper_log_archive` 
    ADD COLUMN `archived_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '归档时间';
```

### 6.3 定时归档任务

#### Go 代码实现

```go
// internal/cron/archive_log.go
package cron

import (
    "context"
    "time"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gcron"
)

func init() {
    // 每月1号凌晨2点执行归档
    gcron.Add(ctx, "0 2 1 * *", func(ctx context.Context) {
        ArchiveOperLog(ctx)
    }, "archive_oper_log")
}

// ArchiveOperLog 归档3个月前的操作日志
func ArchiveOperLog(ctx context.Context) error {
    threshold := time.Now().AddDate(0, -3, 0).Format("2006-01-02")
    
    // 1. 迁移到归档表
    _, err := g.DB().Exec(ctx, `
        INSERT INTO sys_oper_log_archive 
        SELECT *, NOW() as archived_at FROM sys_oper_log 
        WHERE created_at < ?
    `, threshold)
    if err != nil {
        g.Log().Error(ctx, "归档失败:", err)
        return err
    }
    
    // 2. 清理主表
    result, _ := g.DB().Exec(ctx, `DELETE FROM sys_oper_log WHERE created_at < ?`, threshold)
    affected, _ := result.RowsAffected()
    g.Log().Infof(ctx, "归档完成，迁移 %d 条记录", affected)
    return nil
}
```

### 6.4 冷数据导出

超过12个月的归档数据，导出为压缩文件存储：

```go
func ExportToFile(ctx context.Context) error {
    threshold := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
    result, _ := g.DB().GetAll(ctx, `SELECT * FROM sys_oper_log_archive WHERE created_at < ?`, threshold)
    
    filename := fmt.Sprintf("oper_log_%s.json.gz", time.Now().Format("200601"))
    file, _ := os.Create(filename)
    defer file.Close()
    
    gzWriter := gzip.NewWriter(file)
    defer gzWriter.Close()
    json.NewEncoder(gzWriter).Encode(result)
    
    g.DB().Exec(ctx, `DELETE FROM sys_oper_log_archive WHERE created_at < ?`, threshold)
    return nil
}
```

## 7. 查询接口设计

### 7.1 查询参数

```go
type OperLogListReq struct {
    g.Meta    `path:"/oper-log" method:"get" tags:"操作日志"`
    BizType   string `json:"bizType"   dc:"业务类型"`
    BizId     string `json:"bizId"     dc:"业务ID"`
    Operator  string `json:"operator"  dc:"操作人"`
    StartTime string `json:"startTime" dc:"开始时间"`
    EndTime   string `json:"endTime"   dc:"结束时间"`
    Page      int    `json:"page" d:"1"`
    PageSize  int    `json:"pageSize" d:"20"`
}
```

### 7.2 查询服务

```go
func (s *sOperLog) List(ctx g.Ctx, in *OperLogListReq) (list []entity.SysOperLog, total int, err error) {
    m := dao.SysOperLog.Ctx(ctx)
    if in.BizType != "" { m = m.Where("biz_type", in.BizType) }
    if in.BizId != "" { m = m.Where("biz_id", in.BizId) }
    if in.StartTime != "" && in.EndTime != "" {
        m = m.WhereBetween("created_at", in.StartTime, in.EndTime)
    }
    total, _ = m.Count()
    err = m.Page(in.Page, in.PageSize).OrderDesc("id").Scan(&list)
    return
}

// GetHistoryWithArchive 跨表查询（含归档数据）
func (s *sOperLog) GetHistoryWithArchive(ctx g.Ctx, bizType, bizId string) (list []entity.SysOperLog, err error) {
    sql := `(SELECT * FROM sys_oper_log WHERE biz_type=? AND biz_id=?)
            UNION ALL
            (SELECT * FROM sys_oper_log_archive WHERE biz_type=? AND biz_id=?)
            ORDER BY id DESC`
    err = g.DB().GetScan(ctx, &list, sql, bizType, bizId, bizType, bizId)
    return
}
```

## 8. 容量规划

### 8.1 索引优化

```sql
ALTER TABLE `sys_oper_log` 
    ADD INDEX `idx_created_at` (`created_at`),
    ADD INDEX `idx_oper_id` (`oper_id`);
```

### 8.2 监控告警

```go
func CheckTableSize(ctx context.Context) {
    count, _ := dao.SysOperLog.Ctx(ctx).Count()
    if count > 1000000 {
        g.Log().Warning(ctx, "操作日志主表数据量超过100万，请检查归档任务")
    }
}
```

## 9. 总结

*   **BizType/BizId**: 解决"查找某个对象变更历史"的问题
*   **Diff**: 解决"发生了什么具体变化"的问题
*   **数据归档**: 热/温/冷三级存储策略，平衡查询性能与存储成本
*   **查询接口**: 支持多维度筛选和跨归档表查询
