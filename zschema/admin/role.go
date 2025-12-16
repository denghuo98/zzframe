package admin

import (
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/zschema/zform"
)

// RoleEditInput 服务层-角色编辑入参
type RoleEditInput struct {
	entity.AdminRole
	MenuIds []int64 `json:"menuIds" dc:"菜单ID列表"`
}

// RoleDeleteInput 服务层-角色删除入参
type RoleDeleteInput struct {
	Id int64 `json:"id" dc:"角色ID" v:"required|min:1#角色ID不能为空|角色ID不能小于1"`
}

// RoleListInput 服务层-角色列表入参
type RoleListInput struct {
	zform.PageReq
	Name   string `json:"name" dc:"角色名称"`
	Key    string `json:"key" dc:"角色权限字符串"`
	Status int    `json:"status" dc:"角色状态"`
}

// RoleListOutputItem 服务层-角色列表输出项
type RoleListOutputItem struct {
	*entity.AdminRole
	MenuIds []int64 `json:"menuIds" dc:"绑定的菜单ID列表"`
}

// RoleListOutput 服务层-角色列表出参
type RoleListOutput struct {
	List []*RoleListOutputItem `json:"list" dc:"角色列表"`
}

// RoleUpdateFields 角色允许更新的数据字段
type RoleUpdateFields struct {
	Id     int64  `json:"id"           dc:"角色ID"`
	Name   string `json:"name"         dc:"角色名称"`
	Key    string `json:"key"          dc:"角色权限字符串"`
	Remark string `json:"remark"       dc:"备注"`
	Sort   int    `json:"sort"         dc:"排序"`
	Status int    `json:"status"       dc:"角色状态"`
}

// RoleUpdateMenuInput 服务层-角色更新菜单入参
type RoleUpdateMenuInput struct {
	Id      int64   `json:"id" dc:"角色ID"`
	MenuIds []int64 `json:"menuIds" dc:"菜单ID列表"`
}
