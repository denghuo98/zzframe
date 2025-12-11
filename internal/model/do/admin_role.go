// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRole is the golang structure of table zz_admin_role for DAO operations like Where/Data.
type AdminRole struct {
	g.Meta    `orm:"table:zz_admin_role, do:true"`
	Id        any         // 角色ID
	Name      any         // 角色名称
	Key       any         // 角色权限字符串
	Remark    any         // 备注
	Sort      any         // 排序
	Status    any         // 角色状态
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
