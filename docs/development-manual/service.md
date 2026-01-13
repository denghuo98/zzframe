# Service 开发

Service 层包含核心业务逻辑，负责处理业务规则、事务管理、数据转换等。

## 基本结构

```go
package admin

import (
    "github.com/gogf/gf/v2/frame/g"

    "github.com/denghuo98/zzframe/internal/dao"
    "github.com/denghuo98/zzframe/zschema/admin"
)

type sMember struct{}

var Member = sMember{}
```

## 接口定义

为了更好的代码解耦，建议先定义接口：

```go
type IAdminMember interface {
    // VerifyUnique 验证成员唯一性
    VerifyUnique(ctx g.Ctx, in *admin.VerifyUniqueInput) error
    
    // Info 获取当前登录用户信息
    Info(ctx g.Ctx) (*admin.LoginMemberInfoOutput, error)
    
    // List 获取用户列表
    List(ctx g.Ctx, in *admin.MemberListInput) (*admin.MemberListOutput, int, error)
    
    // Edit 编辑成员（新增或更新）
    Edit(ctx g.Ctx, in *admin.MemberEditInput) error
    
    // Delete 删除成员
    Delete(ctx g.Ctx, in *admin.MemberDeleteInput) error
}
```

## 注册接口

```go
var (
    localAdminMember IAdminMember
)

func AdminMember() IAdminMember {
    if localAdminMember == nil {
        panic("AdminMember is not initialized, please register it first")
    }
    return localAdminMember
}

func RegisterAdminMember(i IAdminMember) {
    localAdminMember = i
}
```

## 常见方法实现

### 列表查询

```go
func (s *sMember) List(ctx g.Ctx, in *admin.MemberListInput) (*admin.MemberListOutput, int, error) {
    condition := "1=1"
    args := []interface{}{}
    
    // 构建查询条件
    if in.Username != "" {
        condition += " AND username LIKE ?"
        args = append(args, "%"+in.Username+"%")
    }
    
    if in.Status > 0 {
        condition += " AND status = ?"
        args = append(args, in.Status)
    }
    
    // 查询列表
    list, err := dao.AdminMember.List(ctx, condition, args...)
    if err != nil {
        return nil, 0, err
    }
    
    // 统计总数
    total, err := dao.AdminMember.Count(ctx, condition, args...)
    if err != nil {
        return nil, 0, err
    }
    
    // 转换为输出格式
    output := new(admin.MemberListOutput)
    output.List = make([]*admin.MemberItem, 0, len(list))
    
    for _, item := range list {
        output.List = append(output.List, &admin.MemberItem{
            Id:       item.Id,
            Username: item.Username,
            RealName: item.RealName,
            Status:   item.Status,
        })
    }
    
    return output, total, nil
}
```

### 新增/编辑

```go
func (s *sMember) Edit(ctx g.Ctx, in *admin.MemberEditInput) error {
    // 验证唯一性
    if err := s.VerifyUnique(ctx, &admin.VerifyUniqueInput{
        Id:       in.Id,
        Username: in.Username,
        Email:    in.Email,
        Phone:    in.Phone,
    }); err != nil {
        return err
    }
    
    // 根据是否有 ID 判断是新增还是更新
    if in.Id > 0 {
        // 更新
        _, err := dao.AdminMember.Update(ctx, &do.AdminMember{
            Id:       in.Id,
            Username: in.Username,
            RealName: in.RealName,
            Email:    in.Email,
            Phone:    in.Phone,
        })
        return err
    } else {
        // 新增
        id, err := dao.AdminMember.Insert(ctx, &do.AdminMember{
            Username: in.Username,
            RealName: in.RealName,
            Email:    in.Email,
            Phone:    in.Phone,
            Status:   1,
        })
        if err != nil {
            return err
        }
        
        // 分配默认角色
        return s.AssignDefaultRole(ctx, id)
    }
}
```

### 删除

```go
func (s *sMember) Delete(ctx g.Ctx, in *admin.MemberDeleteInput) error {
    // 验证是否是超级管理员
    if s.VerifySuperAdmin(ctx, in.Id) {
        return errors.New("超级管理员不能删除")
    }
    
    // 软删除：设置状态为禁用
    _, err := dao.AdminMember.Update(ctx, &do.AdminMember{
        Id:     in.Id,
        Status: 0,
    })
    return err
}
```

### 验证逻辑

```go
func (s *sMember) VerifyUnique(ctx g.Ctx, in *admin.VerifyUniqueInput) error {
    condition := "id != ? AND (username = ? OR email = ? OR phone = ?)"
    args := []interface{}{in.Id, in.Username, in.Email, in.Phone}
    
    count, err := dao.AdminMember.Count(ctx, condition, args...)
    if err != nil {
        return err
    }
    
    if count > 0 {
        return errors.New("账号、邮箱或手机号已存在")
    }
    
    return nil
}
```

## 事务处理

```go
func (s *sMember) EditWithTransaction(ctx g.Ctx, in *admin.MemberEditInput) error {
    err := dao.AdminMember.DB.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
        // 使用事务进行数据库操作
        _, err := tx.Model(do.AdminMember{}).Where("id", in.Id).Update(do.AdminMember{
            Username: in.Username,
            // ...
        })
        if err != nil {
            return err
        }
        
        // 其他操作...
        
        return nil
    })
    
    return err
}
```

## 缓存使用

```go
func (s *sMember) GetCachedInfo(ctx g.Ctx, id uint64) (*do.AdminMember, error) {
    // 尝试从缓存获取
    cacheKey := fmt.Sprintf("member:info:%d", id)
    cache, err := zservice.SystemConfig().Cache()
    if err != nil {
        return nil, err
    }
    
    // 从缓存获取
    var member do.AdminMember
    if err := cache.Get(ctx, cacheKey, &member); err == nil && member.Id > 0 {
        return &member, nil
    }
    
    // 从数据库获取
    info, err := dao.AdminMember.GetById(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 设置缓存
    _ = cache.Set(ctx, cacheKey, info, 300) // 缓存 5 分钟
    
    return info, nil
}
```

## 错误处理

```go
func (s *sMember) Edit(ctx g.Ctx, in *admin.MemberEditInput) error {
    // 业务验证
    if in.Username == "" {
        return errors.New("用户名不能为空")
    }
    
    // 唯一性验证
    if err := s.VerifyUnique(ctx, ...); err != nil {
        return err
    }
    
    // 数据库操作
    if err := dao.AdminMember.Update(ctx, ...); err != nil {
        return err
    }
    
    return nil
}
```

## 日志记录

```go
func (s *sMember) Edit(ctx g.Ctx, in *admin.MemberEditInput) error {
    // 记录操作日志
    g.Log().Infof(ctx, "编辑用户: %+v", in)
    
    // 业务逻辑...
    
    // 记录成功日志
    g.Log().Infof(ctx, "编辑用户成功: %d", in.Id)
    
    return nil
}
```

## 最佳实践

1. **接口定义**: 先定义接口，再实现，便于测试和扩展
2. **事务处理**: 需要原子操作时使用事务
3. **缓存使用**: 合理使用缓存提升性能
4. **错误处理**: 提供清晰的错误信息
5. **日志记录**: 记录关键操作和错误
6. **职责单一**: 每个方法只做一件事

## 下一步

- 了解 [数据库操作](./database)
- 查看 [Controller 开发](./controller)
