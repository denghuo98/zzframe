// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminMenu is the golang structure of table zz_admin_menu for DAO operations like Where/Data.
type AdminMenu struct {
	g.Meta      `orm:"table:zz_admin_menu, do:true"`
	Id          any         // 菜单ID
	Pid         any         // 父菜单ID
	Title       any         // 菜单名称
	Name        any         // 名称编码
	Path        any         // 路由地址
	Icon        any         // 菜单图标
	Type        any         // 菜单类型（1目录 2菜单 3按钮）
	Redirect    any         // 重定向地址
	Component   any         // 组件路径
	IsFrame     any         // 是否内嵌
	FrameSrc    any         // 内联外部地址
	Hidden      any         // 是否隐藏
	Sort        any         // 排序
	Remark      any         // 备注
	Status      any         // 菜单状态
	UpdatedAt   *gtime.Time // 更新时间
	CreatedAt   *gtime.Time // 创建时间
	Permissions any         // 菜单权限
}
