# 权限管理

ZZFrame 基于 Casbin 实现了强大的 RBAC (基于角色的访问控制) 权限管理系统。

## 权限模型

### RBAC 模型

```
用户 (User) → 角色 (Role) → 权限 (Permission)
```

### 数据模型

```go
// 用户-角色关联
type AdminMemberRole struct {
    Id       uint64 `json:"id"`
    MemberId uint64 `json:"memberId"`  // 用户 ID
    RoleId   uint64 `json:"roleId"`    // 角色 ID
}

// 角色-菜单关联
type AdminRoleMenu struct {
    Id     uint64 `json:"id"`
    RoleId uint64 `json:"roleId"`    // 角色 ID
    MenuId uint64 `json:"menuId"`    // 菜单 ID
}

// 角色
type AdminRole struct {
    Id          uint64   `json:"id"`
    Name        string   `json:"name"`       // 角色名称
    Code        string   `json:"code"`       // 角色标识
    Status      int      `json:"status"`     // 状态
    Description string   `json:"description"` // 描述
}

// 菜单
type AdminMenu struct {
    Id       uint64 `json:"id"`
    ParentId uint64 `json:"parentId"` // 父级 ID
    Name     string `json:"name"`     // 菜单名称
    Path     string `json:"path"`     // 路径
    Method   string `json:"method"`   // 请求方法
    Status   int    `json:"status"`   // 状态
}
```

## Casbin 配置

### Model 配置

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

### 策略规则

Casbin 的策略规则格式：

```
p, 角色, 资源, 操作
g, 用户, 角色
```

示例：

```
p, super_admin, /admin/member/list, GET
p, super_admin, /admin/member/edit, POST
p, admin, /admin/member/list, GET
g, 1, super_admin
g, 2, admin
```

## 权限验证

### 中间件验证

```go
package middleware

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/casbin/casbin/v2"
)

func AdminAuth(r *ghttp.Request) {
    // 获取用户信息
    userId := r.GetCtxVar("user_id").Int64()
    
    // 获取请求路径和方法
    path := r.URL.Path
    method := r.Method
    
    // 权限验证
    allowed, err := verifyPermission(userId, path, method)
    if err != nil || !allowed {
        r.Response.WriteJson(g.Map{
            "code":    403,
            "message": "没有权限",
        })
        r.Exit()
    }
    
    r.Middleware.Next()
}

func verifyPermission(userId int64, path, method string) (bool, error) {
    enforcer := getCasbinEnforcer()
    
    // 获取用户角色
    roles, err := getUserRoles(userId)
    if err != nil {
        return false, err
    }
    
    // 验证权限
    for _, role := range roles {
        allowed, _ := enforcer.Enforce(role, path, method)
        if allowed {
            return true, nil
        }
    }
    
    return false, nil
}
```

### Service 层验证

```go
func (s *sAdminRole) Verify(ctx g.Ctx, path, method string) bool {
    userId := ctx.GetVar("user_id").Int64()
    
    // 超级管理员拥有所有权限
    if s.VerifySuperAdmin(ctx, userId) {
        return true
    }
    
    // 获取用户角色
    roles, err := s.GetMemberRoles(ctx, userId)
    if err != nil {
        return false
    }
    
    // 验证权限
    enforcer := getCasbinEnforcer()
    for _, role := range roles {
        allowed, _ := enforcer.Enforce(role.Code, path, method)
        if allowed {
            return true, nil
        }
    }
    
    return false
}
```

## 权限管理

### 分配角色

```go
func (s *sAdminMember) AssignRoles(ctx g.Ctx, userId uint64, roleIds []uint64) error {
    // 删除原有角色
    _, err := dao.AdminMemberRole.DeleteByMemberId(ctx, userId)
    if err != nil {
        return err
    }
    
    // 分配新角色
    for _, roleId := range roleIds {
        _, err := dao.AdminMemberRole.Insert(ctx, &do.AdminMemberRole{
            MemberId: userId,
            RoleId:   roleId,
        })
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### 分配权限

```go
func (s *sAdminRole) AssignMenus(ctx g.Ctx, roleId uint64, menuIds []uint64) error {
    // 删除原有权限
    _, err := dao.AdminRoleMenu.DeleteByRoleId(ctx, roleId)
    if err != nil {
        return err
    }
    
    // 分配新权限
    for _, menuId := range menuIds {
        _, err := dao.AdminRoleMenu.Insert(ctx, &do.AdminRoleMenu{
            RoleId: roleId,
            MenuId: menuId,
        })
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### 同步到 Casbin

```go
func SyncToCasbin(ctx g.Ctx) error {
    enforcer := getCasbinEnforcer()
    
    // 清空策略
    enforcer.ClearPolicy()
    
    // 同步所有角色权限
    roles, err := dao.AdminRole.List(ctx, "status = 1")
    if err != nil {
        return err
    }
    
    for _, role := range roles {
        // 获取角色的所有菜单
        menus, err := dao.AdminMenu.GetByRoleId(ctx, role.Id)
        if err != nil {
            return err
        }
        
        // 添加策略
        for _, menu := range menus {
            enforcer.AddPolicy(role.Code, menu.Path, menu.Method)
        }
    }
    
    return enforcer.SavePolicy()
}
```

## 获取用户权限

### 获取用户菜单

```go
func (s *sAdminMenu) GetDynamicMenus(ctx g.Ctx) ([]*admin.MenuDynamicItem, error) {
    userId := ctx.GetVar("user_id").Int64()
    
    // 获取用户角色
    roles, err := s.GetMemberRoles(ctx, userId)
    if err != nil {
        return nil, err
    }
    
    // 获取角色对应的菜单
    menuIds := make([]uint64, 0)
    for _, role := range roles {
        ids, err := dao.AdminMenu.GetMenuIdsByRoleId(ctx, role.Id)
        if err != nil {
            return nil, err
        }
        menuIds = append(menuIds, ids...)
    }
    
    // 去重
    menuIds = unique(menuIds)
    
    // 查询菜单详情
    menus, err := dao.AdminMenu.GetByIds(ctx, menuIds)
    if err != nil {
        return nil, err
    }
    
    // 构建树形结构
    return buildMenuTree(menus, 0), nil
}
```

### 获取用户权限列表

```go
func (s *sAdminMember) GetPermissions(ctx g.Ctx) ([]string, error) {
    userId := ctx.GetVar("user_id").Int64()
    
    // 获取用户角色
    roles, err := s.GetMemberRoles(ctx, userId)
    if err != nil {
        return nil, err
    }
    
    // 获取所有权限
    permissions := make([]string, 0)
    for _, role := range roles {
        menus, err := dao.AdminMenu.GetByRoleId(ctx, role.Id)
        if err != nil {
            return nil, err
        }
        
        for _, menu := range menus {
            permission := fmt.Sprintf("%s:%s", menu.Path, menu.Method)
            permissions = append(permissions, permission)
        }
    }
    
    // 去重
    return unique(permissions), nil
}
```

## 权限控制粒度

### 菜单级权限

通过菜单 ID 控制用户可以看到哪些菜单。

### 接口级权限

通过接口路径和方法控制用户可以访问哪些接口。

### 数据级权限

可以根据业务需求实现更细粒度的数据权限控制。

```go
func (s *sAdminMember) List(ctx g.Ctx, req *MemberListReq) (*MemberListOutput, int, error) {
    userId := ctx.GetVar("user_id").Int64()
    
    // 检查是否有查看所有用户的权限
    if !hasPermission(ctx, userId, "/admin/member/list", "GET") {
        // 只能查看自己
        condition += " AND id = ?"
        args = append(args, userId)
    }
    
    // 查询数据...
}
```

## 超级管理员

### 超级管理员特权

```go
func (s *sAdminMember) VerifySuperAdmin(ctx g.Ctx, userId int64) bool {
    // 检查是否是超级管理员
    superAdminUsername := zservice.SystemConfig().System().SuperAdmin.Username
    if superAdminUsername == "" {
        return false
    }
    
    member, err := dao.AdminMember.GetById(ctx, uint64(userId))
    if err != nil {
        return false
    }
    
    return member.Username == superAdminUsername
}
```

### 超级管理员权限

超级管理员拥有所有权限，不需要进行权限验证。

## 最佳实践

1. **最小权限原则**: 只分配用户所需的最小权限
2. **角色设计**: 合理设计角色，避免角色过多
3. **权限审计**: 定期审计用户权限
4. **权限缓存**: 缓存用户权限，提升性能
5. **动态权限**: 根据业务需求动态调整权限

## 下一步

- 学习 [缓存方案](./cache)
- 了解 [队列方案](./queue)
