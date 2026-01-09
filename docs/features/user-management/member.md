# 用户管理

ZZFrame 提供了完整的用户管理功能，包括用户列表、新增、编辑、删除、密码修改等。

## 功能特性

- 用户列表（分页、搜索）
- 新增用户
- 编辑用户信息
- 删除用户
- 修改密码
- 重置密码
- 用户状态管理
- 用户角色分配

## 数据模型

```sql
CREATE TABLE IF NOT EXISTS `zz_admin_member` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
  `dept_id` bigint(20) DEFAULT '0' COMMENT '部门ID',
  `real_name` varchar(32) DEFAULT '' COMMENT '真实姓名',
  `username` varchar(20) NOT NULL DEFAULT '' COMMENT '帐号',
  `password_hash` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(16) NOT NULL COMMENT '密码盐',
  `password_reset_token` varchar(150) DEFAULT '' COMMENT '密码重置令牌',
  `avatar` char(150) DEFAULT '' COMMENT '头像',
  `sex` tinyint(1) DEFAULT '3' COMMENT '性别',
  `email` varchar(60) DEFAULT '' COMMENT '邮箱',
  `mobile` varchar(20) DEFAULT '' COMMENT '手机号码',
  `last_active_at` datetime DEFAULT NULL COMMENT '最后活跃时间',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_username` (`username`),
  KEY `idx_phone` (`mobile`),
  KEY `idx_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_用户表';
```

## API 接口

### 获取用户信息

```http
GET /admin/member/info
Authorization: Bearer {token}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "superAdmin",
    "realName": "超级管理员",
    "avatar": "",
    "roles": ["admin"],
    "permissions": ["*:*:*"]
  }
}
```

### 获取用户列表

```http
GET /admin/member/list?page=1&size=10&username=admin
Authorization: Bearer {token}
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| size | int | 是 | 每页数量 |
| username | string | 否 | 用户名（模糊搜索） |
| realName | string | 否 | 真实姓名（模糊搜索） |
| mobile | string | 否 | 手机号 |
| status | int | 否 | 状态：1-正常，0-禁用 |

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "realName": "管理员",
        "mobile": "13800138000",
        "email": "admin@example.com",
        "status": 1,
        "createdAt": "2024-01-01 10:00:00"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  }
}
```

### 新增用户

```http
POST /admin/member/add
Authorization: Bearer {token}
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456",
  "realName": "测试用户",
  "mobile": "13800138000",
  "email": "test@example.com",
  "roleIds": [1, 2]
}
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（唯一） |
| password | string | 是 | 密码 |
| realName | string | 是 | 真实姓名 |
| mobile | string | 否 | 手机号 |
| email | string | 否 | 邮箱 |
| roleIds | array | 否 | 角色ID列表 |

### 编辑用户

```http
POST /admin/member/edit
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": 1,
  "realName": "管理员",
  "mobile": "13800138000",
  "email": "admin@example.com",
  "roleIds": [1, 2]
}
```

### 删除用户

```http
POST /admin/member/delete
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": 1
}
```

### 修改密码

```http
POST /admin/member/updatePwd
Authorization: Bearer {token}
Content-Type: application/json

{
  "oldPassword": "123456",
  "newPassword": "654321"
}
```

### 重置密码

```http
POST /admin/member/resetPwd
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": 1,
  "newPassword": "123456"
}
```

## 使用示例

### 获取当前登录用户信息

```go
func (c *cMember) Info(ctx g.Ctx, req *api.MemberInfoReq) (res *api.MemberInfoRes, err error) {
    out, err := service.AdminMember().Info(ctx)
    if err != nil {
        return nil, err
    }
    res = new(api.MemberInfoRes)
    res.LoginMemberInfoOutput = out
    return res, nil
}
```

### 获取用户列表

```go
func (c *cMember) List(ctx g.Ctx, req *api.MemberListReq) (res *api.MemberListRes, err error) {
    list, totalCount, err := service.AdminMember().List(ctx, &req.MemberListInput)
    if err != nil {
        return nil, err
    }
    res = new(api.MemberListRes)
    res.MemberListOutput = *list
    res.PageRes.Pack(req, int(totalCount))
    return res, nil
}
```

### 新增用户

```go
func (c *cMember) Add(ctx g.Ctx, req *api.MemberAddReq) (res *api.MemberAddRes, err error) {
    if err = service.AdminMember().Add(ctx, &req.MemberAddInput); err != nil {
        return nil, err
    }
    return &api.MemberAddRes{}, nil
}
```

## 权限要求

| 接口 | 权限标识 |
|------|---------|
| GET /admin/member/info | - |
| GET /admin/member/list | admin:member:list |
| POST /admin/member/add | admin:member:add |
| POST /admin/member/edit | admin:member:edit |
| POST /admin/member/delete | admin:member:delete |
| POST /admin/member/updatePwd | - |
| POST /admin/member/resetPwd | admin:member:resetPwd |

## 功能说明

1. 用户名必须唯一
2. 手机号和邮箱不能重复
3. 删除用户会同时删除用户的角色关联
4. 修改密码需要验证旧密码
5. 重置密码由管理员操作

## 扩展功能

### 自定义用户字段

在模型中添加自定义字段：

```go
type AdminMember struct {
    Id          uint64    `json:"id"`
    Username    string    `json:"username"`
    // 自定义字段
    CustomField string    `json:"customField"`
}
```

### 用户验证规则

添加自定义验证规则：

```go
type MemberAddReq struct {
    Username string `json:"username" v:"required|length:4,20#用户名不能为空|用户名长度4-20位"`
    Password string `json:"password" v:"required|length:6,32#密码不能为空|密码长度6-32位"`
}
```

## 下一步

- [角色管理](./role.md)
- [菜单管理](./menu.md)
