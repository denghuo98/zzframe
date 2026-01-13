# 数据库操作

ZZFrame 基于 GoFrame 的 ORM 框架，提供强大的数据库操作能力。

## DAO 层结构

```go
package dao

import (
    "context"
    
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
)

var AdminMember = adminMember{}

type adminMember struct{}

func (d *adminMember) Model(ctx context.Context, args ...gdb.ModelFunc) *gdb.Model {
    return g.Model("zz_admin_member").Ctx(ctx).Where("deleted_at IS NULL").Safe(args...)
}
```

## 基本操作

### 查询单条记录

```go
func (d *adminMember) GetById(ctx context.Context, id uint64) (*do.AdminMember, error) {
    var info *do.AdminMember
    err := d.Model(ctx).Where("id", id).Scan(&info)
    return info, err
}
```

### 查询多条记录

```go
func (d *adminMember) List(ctx context.Context, condition string, args ...interface{}) ([]*do.AdminMember, error) {
    var list []*do.AdminMember
    err := d.Model(ctx).Where(condition, args...).OrderDesc("id").Scan(&list)
    return list, err
}
```

### 统计记录数

```go
func (d *adminMember) Count(ctx context.Context, condition string, args ...interface{}) (int64, error) {
    return d.Model(ctx).Where(condition, args...).Count()
}
```

### 新增记录

```go
func (d *adminMember) Insert(ctx context.Context, data *do.AdminMember) (int64, error) {
    result, err := d.Model(ctx).Data(data).Insert()
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}
```

### 更新记录

```go
func (d *adminMember) Update(ctx context.Context, data *do.AdminMember) (int64, error) {
    result, err := d.Model(ctx).Where("id", data.Id).Data(data).Update()
    if err != nil {
        return 0, err
    }
    return result.RowsAffected()
}
```

### 删除记录

```go
func (d *adminMember) Delete(ctx context.Context, id uint64) (int64, error) {
    result, err := d.Model(ctx).Where("id", id).Delete()
    if err != nil {
        return 0, err
    }
    return result.RowsAffected()
}
```

## 高级查询

### 分页查询

```go
func (d *adminMember) Page(ctx context.Context, page, size int, condition string, args ...interface{}) ([]*do.AdminMember, int64, error) {
    total, err := d.Model(ctx).Where(condition, args...).Count()
    if err != nil {
        return nil, 0, err
    }
    
    var list []*do.AdminMember
    err = d.Model(ctx).
        Where(condition, args...).
        Page(page, size).
        OrderDesc("id").
        Scan(&list)
    
    return list, total, err
}
```

### 联表查询

```go
func (d *adminMember) ListWithRole(ctx context.Context) ([]*do.AdminMember, error) {
    var list []*do.AdminMember
    err := d.Model(ctx).
        LeftJoin("zz_admin_member_role", "zz_admin_member.id = zz_admin_member_role.member_id").
        LeftJoin("zz_admin_role", "zz_admin_role.id = zz_admin_member_role.role_id").
        Fields("zz_admin_member.*, zz_admin_role.name as role_name").
        Scan(&list)
    return list, err
}
```

### 条件构建

```go
func (d *adminMember) Search(ctx context.Context, keyword string, status int) ([]*do.AdminMember, error) {
    model := d.Model(ctx)
    
    // 动态条件
    if keyword != "" {
        model = model.WhereLike("username", "%"+keyword+"%").
                     WhereOrLike("real_name", "%"+keyword+"%").
                     WhereOrLike("email", "%"+keyword+"%")
    }
    
    if status > 0 {
        model = model.Where("status", status)
    }
    
    var list []*do.AdminMember
    err := model.OrderDesc("id").Scan(&list)
    return list, err
}
```

### 字段选择

```go
func (d *adminMember) GetSimpleList(ctx context.Context) ([]*do.AdminMember, error) {
    var list []*do.AdminMember
    err := d.Model(ctx).
        Fields("id", "username", "real_name").
        Scan(&list)
    return list, err
}
```

## 软删除

```go
func (d *adminMember) SoftDelete(ctx context.Context, id uint64) error {
    _, err := d.Model(ctx).
        Where("id", id).
        Update(g.Map{
            "deleted_at": gtime.Now(),
            "status":     0,
        })
    return err
}
```

## 事务处理

```go
func (d *adminMember) Transfer(ctx context.Context, fromId, toId uint64, amount float64) error {
    return dao.AdminMember.DB.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
        // 扣减
        _, err := tx.Model(do.AdminMember{}).Where("id", fromId).
                 Decrement("balance", amount)
        if err != nil {
            return err
        }
        
        // 增加
        _, err = tx.Model(do.AdminMember{}).Where("id", toId).
                 Increment("balance", amount)
        if err != nil {
            return err
        }
        
        return nil
    })
}
```

## 批量操作

### 批量插入

```go
func (d *adminMember) BatchInsert(ctx context.Context, list []*do.AdminMember) error {
    _, err := d.Model(ctx).Data(list).Insert()
    return err
}
```

### 批量更新

```go
func (d *adminMember) BatchUpdate(ctx context.Context, ids []uint64, status int) error {
    _, err := d.Model(ctx).WhereIn("id", ids).Update(g.Map{
        "status": status,
    })
    return err
}
```

### 批量删除

```go
func (d *adminMember) BatchDelete(ctx context.Context, ids []uint64) error {
    _, err := d.Model(ctx).WhereIn("id", ids).Delete()
    return err
}
```

## 复杂查询示例

### 带子查询的查询

```go
func (d *adminMember) GetUsersWithRoles(ctx context.Context) ([]*do.AdminMember, error) {
    var list []*do.AdminMember
    err := d.Model(ctx).
        Where("id IN (?)", g.Model("zz_admin_member_role").Fields("member_id")).
        Scan(&list)
    return list, err
}
```

### 分组统计

```go
func (d *adminMember) GetStatusStats(ctx context.Context) ([]map[string]interface{}, error) {
    var result []map[string]interface{}
    err := d.Model(ctx).
        Fields("status, COUNT(*) as count").
        Group("status").
        Scan(&result)
    return result, err
}
```

## 性能优化

### 使用索引

```go
// 在 model 层定义索引
func init() {
    dao.AdminMember.DB.Model(do.AdminMember{}).
        Index("idx_username", "username").
        Index("idx_phone", "phone").
        Create()
}
```

### 批量查询优化

```go
func (d *adminMember) GetByIds(ctx context.Context, ids []uint64) ([]*do.AdminMember, error) {
    var list []*do.AdminMember
    err := d.Model(ctx).WhereIn("id", ids).Scan(&list)
    return list, err
}
```

### 只读查询

```go
func (d *adminMember) ListReadOnly(ctx context.Context) ([]*do.AdminMember, error) {
    var list []*do.AdminMember
    err := d.Model(ctx).Cache(gdb.CacheOption{
        Duration: time.Minute * 5,
        Name:     "member_list",
    }).Scan(&list)
    return list, err
}
```

## 最佳实践

1. **使用 DO 层**: 使用 do 包中的结构体进行数据操作
2. **软删除**: 推荐使用软删除而不是物理删除
3. **事务处理**: 需要保证数据一致性时使用事务
4. **批量操作**: 大量数据时使用批量操作提升性能
5. **索引优化**: 为常用查询字段添加索引
6. **错误处理**: 妥善处理数据库错误

## 下一步

- 学习 [Service 开发](./service)
- 了解 [配置管理](./configuration)
